package azurecaf

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var updateGolden = flag.Bool("update-golden", false, "update golden files")

func TestGolden_AllResources(t *testing.T) {
	goldenFile := filepath.Join("testdata", "golden", "all_resources.json")
	results := make(map[string]string, len(ResourceDefinitions))

	resourceTypes := sortedResourceTypes(ResourceDefinitions)
	for _, resourceType := range resourceTypes {
		results[resourceType] = generateAzureCAFName(t, resourceType, map[string]interface{}{})
	}

	assertGoldenStringMap(t, goldenFile, results)
}

func TestGolden_NamingConventions(t *testing.T) {
	goldenFile := filepath.Join("testdata", "golden", "naming_conventions.json")

	results := map[string]string{
		ConventionCafClassic:  generateLegacyConventionName(t, ConventionCafClassic, "rg", "dev", "myapp", "001", 24, 12345),
		ConventionCafRandom:   generateLegacyConventionName(t, ConventionCafRandom, "rg", "dev", "myapp", "001", 24, 12345),
		ConventionRandom:      generateLegacyConventionName(t, ConventionRandom, "rg", "dev", "myapp", "001", 24, 12345),
		ConventionPassThrough: generateLegacyConventionName(t, ConventionPassThrough, "rg", "dev", "myapp", "001", 24, 12345),
	}

	assertGoldenStringMap(t, goldenFile, results)
}

func TestGolden_EdgeCases(t *testing.T) {
	goldenFile := filepath.Join("testdata", "golden", "edge_cases.json")

	results := map[string]string{
		"empty_name": generateAzureCAFName(t, "azurerm_resource_group", map[string]interface{}{
			"name": "",
		}),
		"very_long_name": generateAzureCAFName(t, "azurerm_resource_group", map[string]interface{}{
			"name": strings.Repeat("v", 84),
		}),
		"special_chars": generateAzureCAFName(t, "azurerm_container_app", map[string]interface{}{
			"name":     "my$app_!cool",
			"prefixes": []interface{}{"dev!"},
			"suffixes": []interface{}{"00#1"},
		}),
	}

	assertGoldenStringMap(t, goldenFile, results)
}

func generateAzureCAFName(t *testing.T, resourceType string, overrides map[string]interface{}) string {
	t.Helper()

	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	testData := map[string]interface{}{
		"name":          "myapp",
		"resource_type": resourceType,
		"prefixes":      []interface{}{"dev"},
		"suffixes":      []interface{}{"001"},
		"separator":     "-",
		"random_length": 0,
		"random_seed":   12345,
		"clean_input":   true,
		"use_slug":      true,
	}
	for key, value := range overrides {
		testData[key] = value
	}

	resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, testData)
	if err := nameResource.Create(resourceData, nil); err != nil {
		t.Fatalf("failed to generate name for %s: %v", resourceType, err)
	}

	result, ok := resourceData.Get("result").(string)
	if !ok || result == "" {
		t.Fatalf("empty result for %s", resourceType)
	}

	return result
}

func generateLegacyConventionName(t *testing.T, convention, resourceType, prefix, name, postfix string, desiredMaxLength int, seed int64) string {
	t.Helper()

	resource, ok := Resources[resourceType]
	if !ok {
		resource, ok = ResourcesMapping[resourceType]
	}
	if !ok {
		t.Fatalf("legacy resource type %s not found", resourceType)
	}

	myRegex, err := regexp.Compile(resource.RegEx)
	if err != nil {
		t.Fatalf("invalid regex for %s: %v", resourceType, err)
	}
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		t.Fatalf("invalid validation regex for %s: %v", resourceType, err)
	}

	cafPrefix := ""
	rng := rand.New(rand.NewSource(seed))
	randomSuffix := randSeqDeterministic(resource.MaxLength, rng)

	switch convention {
	case ConventionCafRandom, ConventionCafClassic:
		cafPrefix = resource.CafPrefix
	case ConventionRandom:
		name = ""
		postfix = ""
	}

	nameList := []string{}
	for _, s := range []string{prefix, cafPrefix, name, postfix} {
		if strings.TrimSpace(s) != "" {
			nameList = append(nameList, s)
		}
	}

	userInputName := strings.Join(nameList, suffixSeparator)
	userInputName = myRegex.ReplaceAllString(userInputName, "")
	randomSuffix = myRegex.ReplaceAllString(randomSuffix, "")

	generatedName := userInputName
	maxLength := resource.MaxLength
	if desiredMaxLength > 0 && desiredMaxLength < maxLength {
		maxLength = desiredMaxLength
	}

	containsRandomChar := false
	switch convention {
	case ConventionPassThrough, ConventionCafClassic:
	default:
		if len(userInputName) != 0 {
			if len(userInputName) < (maxLength - 1) {
				containsRandomChar = true
				generatedName = strings.Join([]string{userInputName, randomSuffix}, suffixSeparator)
			} else {
				generatedName = userInputName
			}
		} else {
			containsRandomChar = true
			generatedName = randomSuffix
		}
	}

	filteredGeneratedName := myRegex.ReplaceAllString(generatedName, "")
	length := len(filteredGeneratedName)
	if length > maxLength {
		length = maxLength
	}

	result := filteredGeneratedName[:length]
	if containsRandomChar && len(result) > len(userInputName) {
		randomLastChar := alphagenerator[rng.Intn(len(alphagenerator)-1)]
		resultRune := []rune(result)
		resultRune[len(resultRune)-1] = randomLastChar
		result = string(resultRune)
	}

	if resource.LowerCase {
		result = strings.ToLower(result)
	}

	if !validationRegEx.MatchString(result) {
		t.Fatalf("generated name %q does not match validation regex %s for %s", result, resource.ValidationRegExp, resourceType)
	}

	return result
}

func assertGoldenStringMap(t *testing.T, goldenFile string, got map[string]string) {
	t.Helper()

	if *updateGolden {
		writeGoldenFile(t, goldenFile, got)
		return
	}

	want := readGoldenStringMap(t, goldenFile)
	compareGoldenStringMaps(t, goldenFile, want, got)
}

func writeGoldenFile(t *testing.T, goldenFile string, content map[string]string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(goldenFile), 0o755); err != nil {
		t.Fatalf("failed to create golden directory for %s: %v", goldenFile, err)
	}

	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal golden content for %s: %v", goldenFile, err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(goldenFile, data, 0o644); err != nil {
		t.Fatalf("failed to write golden file %s: %v", goldenFile, err)
	}
}

func readGoldenStringMap(t *testing.T, goldenFile string) map[string]string {
	t.Helper()

	data, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("golden file %s not found, rerun with -update-golden: %v", goldenFile, err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal golden file %s: %v", goldenFile, err)
	}

	return result
}

func compareGoldenStringMaps(t *testing.T, goldenFile string, want, got map[string]string) {
	t.Helper()

	missing := []string{}
	extra := []string{}
	mismatched := []string{}

	for key, expected := range want {
		actual, ok := got[key]
		if !ok {
			missing = append(missing, key)
			continue
		}
		if actual != expected {
			mismatched = append(mismatched, fmt.Sprintf("%s: got %q want %q", key, actual, expected))
		}
	}

	for key := range got {
		if _, ok := want[key]; !ok {
			extra = append(extra, key)
		}
	}

	sort.Strings(missing)
	sort.Strings(extra)
	sort.Strings(mismatched)

	if len(missing) == 0 && len(extra) == 0 && len(mismatched) == 0 {
		return
	}

	for _, key := range missing {
		t.Errorf("%s missing key %q", goldenFile, key)
	}
	for _, key := range extra {
		t.Errorf("%s has new key %q, rerun with -update-golden", goldenFile, key)
	}
	for _, mismatch := range mismatched {
		t.Error(mismatch)
	}
}

func sortedResourceTypes(definitions map[string]ResourceStructure) []string {
	resourceTypes := make([]string, 0, len(definitions))
	for resourceType := range definitions {
		resourceTypes = append(resourceTypes, resourceType)
	}
	sort.Strings(resourceTypes)
	return resourceTypes
}

func randSeqDeterministic(length int, rng *rand.Rand) string {
	if length <= 0 {
		return ""
	}

	b := make([]rune, length)
	for i := range b {
		b[i] = alphagenerator[rng.Intn(len(alphagenerator)-1)]
	}
	return string(b)
}
