package azurecaf

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const errResourceNotFound = "azurecaf_name resource not found"

func setData(prefixes []string, name string, suffixes []string, cleanInput bool) *schema.ResourceData {
	data := &schema.ResourceData{}
	data.Set("name", name)
	data.Set("prefixes", prefixes)
	data.Set("suffixes", suffixes)
	data.Set("clean_input", cleanInput)
	return data
}

func TestCleanInput_no_changes(t *testing.T) {
	data := "testdata"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if data != result {
		t.Errorf("Expected %s but received %s", data, result)
	}
}

func TestCleanInput_remove_always(t *testing.T) {
	data := "🐱‍🚀testdata😊"
	expected := "testdata"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestCleanInput_not_remove_special_allowed_chars(t *testing.T) {
	data := "testdata()"
	expected := "testdata()"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestCleanStringInvalidRegex(t *testing.T) {
	// Inject an invalid regex pattern to verify cleanString returns the input unchanged (no panic)
	resource := ResourceDefinitions["azurerm_resource_group"]
	invalidResource := resource
	invalidResource.RegEx = "[" // Invalid regex
	input := "my-test-name"
	result := cleanString(input, &invalidResource)
	if result != input {
		t.Errorf("Expected cleanString to return input unchanged for invalid regex, got %q", result)
	}
}

func TestCleanSplice_no_changes(t *testing.T) {
	data := []string{"testdata", "test", "data"}
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanSlice(data, &resource)
	for i := range data {
		if data[i] != result[i] {
			t.Errorf("Expected %s but received %s", data[i], result[i])
		}
	}
}

func TestConcatenateParameters_azurerm_public_ip_prefix(t *testing.T) {
	prefixes := []string{"pre"}
	suffixes := []string{"suf"}
	content := []string{"name", "ip"}
	separator := "-"
	expected := "pre-name-ip-suf"
	result := concatenateParameters(separator, prefixes, content, suffixes)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestGetSlug(t *testing.T) {
	resourceType := "azurerm_resource_group"
	convention := ConventionCafClassic
	result := getSlug(resourceType, convention)
	expected := "rg"
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestGetSlug_unknown(t *testing.T) {
	resourceType := "azurerm_does_not_exist"
	convention := ConventionCafClassic
	result := getSlug(resourceType, convention)
	expected := ""
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestAccResourceName_CafClassic(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal(errResourceNotFound)
	}

	// Test case 1: Resource Group
	t.Run("ClassicResourceGroup", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "myrg",
			"resource_type": "azurerm_resource_group",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1-pr2-rg-myrg-"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_resource_group"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 2: Container App with invalid name cleaning
	t.Run("ContainerAppInvalidCleaning", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_ca_name",
			"resource_type": "azurerm_container_app",
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "ca-myinvalidcaname"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_app"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 3: Passthrough mode
	t.Run("PassthroughMode", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "passthRough",
			"resource_type": "azurerm_container_registry",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
			"passthrough":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expected := "passthrough" // passthrough mode normalizes case
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}
	})

	// Test case 4: Container App Environment
	t.Run("ContainerAppEnvironment", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_cae_name",
			"resource_type": "azurerm_container_app_environment",
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "cae-myinvalidcaename"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_app_environment"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 5: Container Registry
	t.Run("ContainerRegistry", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_acr_name",
			"resource_type": "azurerm_container_registry",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1pr2crmyinvalidacrname"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_registry"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	t.Log("CAF Classic naming tests completed successfully")
}

func TestAccResourceName_CafClassicRSV(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal(errResourceNotFound)
	}

	// Test Recovery Services Vault
	t.Run("RecoveryServicesVault", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "test",
			"resource_type": "azurerm_recovery_services_vault",
			"prefixes":      []interface{}{"pr1"},
			"suffixes":      []interface{}{"su1"},
			"random_length": 2,
			"random_seed":   1,
			"clean_input":   true,
			"passthrough":   false,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1-rsv-test-"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Should end with su1 suffix
		if !strings.HasSuffix(result, "-su1") {
			t.Errorf("Expected result to end with '-su1', got '%s'", result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_recovery_services_vault"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	t.Log("CAF Classic RSV naming test completed successfully")
}

func TestComposeName(t *testing.T) {
	namePrecedence := []string{"name", "random", "slug", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name, _ := composeName("-", prefixes, "name", "slug", suffixes, "rd", 21, namePrecedence, false)
	expected := "a-b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrect(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name, _ := composeName("-", prefixes, "name", "slug", suffixes, "rd", 19, namePrecedence, false)
	expected := "b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutMaxLength(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{}
	suffixes := []string{}
	name, _ := composeName("-", prefixes, "aaaaaaaaaa", "bla", suffixes, "", 10, namePrecedence, false)
	expected := "aaaaaaaaaa"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrectSuffixes(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name, _ := composeName("-", prefixes, "name", "slug", suffixes, "rd", 15, namePrecedence, false)
	expected := "slug-name-rd-c"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeEmptyStringArray(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"", "b"}
	suffixes := []string{"", "d"}
	name, _ := composeName("-", prefixes, "", "", suffixes, "", 15, namePrecedence, false)
	expected := "b-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeName_ErrorWhenExceedingMaxLength(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"prefix"}
	suffixes := []string{"suffix"}
	_, err := composeName("-", prefixes, "verylongname", "", suffixes, "", 10, namePrecedence, true)
	if err == nil {
		t.Errorf("expected error when name exceeds max length, got nil")
	}
}

func TestResourceName_ErrorWhenExceedingMaxLength(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]

	resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
		"name":                            "verylongnamethatwillexceedmaxlength",
		"prefixes":                        []interface{}{"prefix1", "prefix2"},
		"suffixes":                        []interface{}{"suffix1", "suffix2"},
		"resource_type":                   "azurerm_storage_account",
		"use_slug":                        true,
		"clean_input":                     true,
		"separator":                       "-",
		"error_when_exceeding_max_length": true,
	})

	err := nameResource.Create(resourceData, nil)
	if err == nil {
		t.Errorf("expected error when name exceeds max length, got nil")
	}
	expectedPattern := regexp.MustCompile(`exceeds maximum length of \d+ by \d+ characters`)
	if !expectedPattern.MatchString(err.Error()) {
		t.Errorf("error %q does not match pattern %q", err.Error(), expectedPattern.String())
	}
}

func TestValidResourceType_validParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resourceTypes := []string{"azurerm_container_registry", "azurerm_storage_account"}
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if !isValid {
		t.Logf("resource types considered invalid while input parameters are valid")
		t.Fail()
	}
	if err != nil {
		t.Logf("resource validation generated an unexpected error %s", err.Error())
		t.Fail()
	}
}
func TestValidResourceType_invalidParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resourceTypes := []string{"azurerm_not_supported", "azurerm_storage_account"}
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if isValid {
		t.Logf("resource types considered valid while input parameters are invalid")
		t.Fail()
	}
	if err == nil {
		t.Logf("resource validation did generate an error while the input is invalid")
		t.Fail()
	}
}

func TestGetResourceNameValid(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, true, namePrecedence, false)
	expected := "a-b-rg-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidRsv(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_recovery_services_vault", "-", []string{"a", "b"}, "test", nil, "1234", "cafclassic", true, false, true, namePrecedence, false)
	expected := "a-b-rsv-test-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidNoSlug(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, false, namePrecedence, false)
	expected := "a-b-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameInvalidResourceType(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_invalid", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, true, namePrecedence, false)
	expected := "a-b-rg-myrg-1234"

	if err == nil {
		t.Logf("Expected a validation error, got nil")
		t.Fail()
	}
	if expected == resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

func TestGetResourceNamePassthrough(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, _ := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, true, true, namePrecedence, false)
	expected := "myrg"

	if expected != resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

func testResourceNameStateDataV2() map[string]interface{} {
	return map[string]interface{}{}
}

func testResourceNameStateDataV3() map[string]interface{} {
	return map[string]interface{}{
		"use_slug": true,
	}
}

func TestResourceExampleInstanceStateUpgradeV2(t *testing.T) {
	expected := testResourceNameStateDataV3()
	actual, err := resourceNameStateUpgradeV2(context.Background(), testResourceNameStateDataV2(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

const testAccResourceNameCafClassicConfig = `


# Resource Group
resource "azurecaf_name" "classic_rg" {
    name            = "myrg"
	resource_type   = "azurerm_resource_group"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_ca_invalid" {
    name            = "my_invalid_ca_name"
	resource_type   = "azurerm_container_app"
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_cae_invalid" {
    name            = "my_invalid_cae_name"
	resource_type   = "azurerm_container_app_environment"
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_acr_invalid" {
    name            = "my_invalid_acr_name"
	resource_type   = "azurerm_container_registry"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "passthrough" {
    name            = "passthRough"
	resource_type   = "azurerm_container_registry"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
	passthrough     = true
}


resource "azurecaf_name" "apim" {
	name = "apim"
	resource_type = "azurerm_api_management_service"
	prefixes = ["vsic"]
	random_length = 0
	clean_input = true
	passthrough = false
   }
`

const testAccResourceNameCafClassicConfigRsv = `


# Resource Group

resource "azurecaf_name" "rsv" {
    name            = "test"
	resource_type   = "azurerm_recovery_services_vault"
	prefixes        = ["pr1"]
	suffixes        = ["su1"]
	random_length   = 2
	random_seed     = 1
	clean_input     = true
	passthrough     = false
}
`

// TestResourceNameHasCustomizeDiff verifies that CustomizeDiff is registered
// on the azurecaf_name resource, which is required for plan-time visibility.
func TestResourceNameHasCustomizeDiff(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal(errResourceNotFound)
	}
	if nameResource.CustomizeDiff == nil {
		t.Fatal("azurecaf_name resource must have CustomizeDiff for plan-time computation")
	}
}

// TestPlanApplyConsistency verifies that the same random_seed produces the same
// result across multiple invocations, which is critical for plan-apply consistency.
func TestPlanApplyConsistency(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]

	input := map[string]interface{}{
		"name":          "myapp",
		"resource_type": "azurerm_resource_group",
		"prefixes":      []interface{}{"dev"},
		"suffixes":      []interface{}{"001"},
		"random_seed":   42,
		"random_length": 5,
		"clean_input":   true,
		"use_slug":      true,
	}

	// Call Create twice with the same seed
	rd1 := schema.TestResourceDataRaw(t, nameResource.Schema, input)
	if err := nameResource.Create(rd1, nil); err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	result1 := rd1.Get("result").(string)

	rd2 := schema.TestResourceDataRaw(t, nameResource.Schema, input)
	if err := nameResource.Create(rd2, nil); err != nil {
		t.Fatalf("Second call failed: %v", err)
	}
	result2 := rd2.Get("result").(string)

	if result1 != result2 {
		t.Errorf("Same seed must produce same result: first=%q second=%q", result1, result2)
	}
	if result1 == "" {
		t.Error("Result must not be empty")
	}
}

// TestPlanTimeMultipleResourceTypes verifies that the results map is correctly
// populated for multiple resource types, which is the main use case for the
// azurecaf_name resource (vs data source).
func TestPlanTimeMultipleResourceTypes(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]

	resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
		"name":           "myapp",
		"resource_type":  "azurerm_resource_group",
		"resource_types": []interface{}{"azurerm_storage_account", "azurerm_key_vault"},
		"prefixes":       []interface{}{"dev"},
		"random_seed":    100,
		"random_length":  3,
		"clean_input":    true,
		"use_slug":       true,
	})

	if err := nameResource.Create(resourceData, nil); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check primary result
	result := resourceData.Get("result").(string)
	if !strings.Contains(result, "rg") {
		t.Errorf("Expected result to contain slug 'rg', got %q", result)
	}

	// Check results map
	results := resourceData.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 entries in results map, got %d", len(results))
	}

	stResult, ok := results["azurerm_storage_account"]
	if !ok {
		t.Fatal("Expected azurerm_storage_account in results map")
	}
	if stResult == "" {
		t.Error("azurerm_storage_account result must not be empty")
	}

	kvResult, ok := results["azurerm_key_vault"]
	if !ok {
		t.Fatal("Expected azurerm_key_vault in results map")
	}
	if kvResult == "" {
		t.Error("azurerm_key_vault result must not be empty")
	}
}

// TestDeterministicWithoutRandom verifies that when random_length is 0,
// results are fully deterministic regardless of seed.
func TestDeterministicWithoutRandom(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]

	input := map[string]interface{}{
		"name":          "myapp",
		"resource_type": "azurerm_resource_group",
		"prefixes":      []interface{}{"dev"},
		"suffixes":      []interface{}{"001"},
		"random_length": 0,
		"clean_input":   true,
		"use_slug":      true,
	}

	rd1 := schema.TestResourceDataRaw(t, nameResource.Schema, input)
	if err := nameResource.Create(rd1, nil); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	result := rd1.Get("result").(string)

	expected := "dev-rg-myapp-001"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestRandSeqDeterminism verifies that randSeq with the same seed always
// produces identical output, which is the foundation of plan-apply consistency.
func TestRandSeqDeterminism(t *testing.T) {
	seed1 := int64(12345)
	seed2 := int64(12345)

	r1 := randSeq(8, &seed1)
	r2 := randSeq(8, &seed2)

	if r1 != r2 {
		t.Errorf("Same seed must produce same sequence: %q vs %q", r1, r2)
	}
	if len(r1) != 8 {
		t.Errorf("Expected length 8, got %d", len(r1))
	}

	// Different seed → different result
	seed3 := int64(99999)
	r3 := randSeq(8, &seed3)
	if r1 == r3 {
		t.Errorf("Different seeds should (almost certainly) produce different sequences")
	}
}

// TestRandSeqZeroLength verifies empty string for zero or negative length.
func TestRandSeqZeroLength(t *testing.T) {
	seed := int64(1)
	if r := randSeq(0, &seed); r != "" {
		t.Errorf("Expected empty string for length 0, got %q", r)
	}
	if r := randSeq(-1, &seed); r != "" {
		t.Errorf("Expected empty string for negative length, got %q", r)
	}
}

// TestComputeNamesMatchesGetNameResult verifies the refactored computeNames
// function produces identical output to the Create path (getNameResult).
func TestComputeNamesMatchesGetNameResult(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]

	input := map[string]interface{}{
		"name":           "myapp",
		"resource_type":  "azurerm_resource_group",
		"resource_types": []interface{}{"azurerm_storage_account", "azurerm_key_vault"},
		"prefixes":       []interface{}{"dev"},
		"suffixes":       []interface{}{"001"},
		"random_seed":    42,
		"random_length":  4,
		"clean_input":    true,
		"use_slug":       true,
	}

	// Path 1: via Create (getNameResult)
	rd := schema.TestResourceDataRaw(t, nameResource.Schema, input)
	if err := nameResource.Create(rd, nil); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	createResult := rd.Get("result").(string)
	createResults := rd.Get("results").(map[string]interface{})

	// Path 2: via computeNames directly
	p := namingParams{
		name:         "myapp",
		resourceType: "azurerm_resource_group",
		resourceTypes: []string{"azurerm_storage_account", "azurerm_key_vault"},
		prefixes:     []string{"dev"},
		suffixes:     []string{"001"},
		randomSeed:   42,
		randomLength: 4,
		cleanInput:   true,
		useSlug:      true,
		separator:    "-",
	}
	computeResult, computeResults, err := computeNames(p)
	if err != nil {
		t.Fatalf("computeNames failed: %v", err)
	}

	if createResult != computeResult {
		t.Errorf("result mismatch: Create=%q computeNames=%q", createResult, computeResult)
	}
	for k, v := range createResults {
		if computeResults[k] != v.(string) {
			t.Errorf("results[%s] mismatch: Create=%q computeNames=%q", k, v, computeResults[k])
		}
	}
}

// TestComputeNamesOnlyResourceTypes verifies computeNames when only
// resource_types is set (no resource_type), result should be empty.
func TestComputeNamesOnlyResourceTypes(t *testing.T) {
	p := namingParams{
		name:          "myapp",
		resourceTypes: []string{"azurerm_storage_account"},
		prefixes:      []string{"dev"},
		randomSeed:    42,
		randomLength:  3,
		cleanInput:    true,
		useSlug:       true,
		separator:     "-",
	}
	result, results, err := computeNames(p)
	if err != nil {
		t.Fatalf("computeNames failed: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty result when resource_type is not set, got %q", result)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 entry in results, got %d", len(results))
	}
	if results["azurerm_storage_account"] == "" {
		t.Error("azurerm_storage_account result must not be empty")
	}
}

// TestComputeNamesInvalidResourceType verifies computeNames returns an error
// for an invalid resource type.
func TestComputeNamesInvalidResourceType(t *testing.T) {
	p := namingParams{
		name:         "myapp",
		resourceType: "azurerm_nonexistent",
		cleanInput:   true,
		useSlug:      true,
		separator:    "-",
	}
	_, _, err := computeNames(p)
	if err == nil {
		t.Error("Expected error for invalid resource type, got nil")
	}
}

// TestComposeName_ErrorWhenExceedingMaxLength_Success verifies that composeName
// returns the full name (not trimmed) when errorWhenExceedingMaxLength is true
// and the name fits within the limit.
func TestComposeName_ErrorWhenExceedingMaxLength_Success(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, err := composeName("-", []string{"a"}, "b", "", []string{"c"}, "", 100, namePrecedence, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "a-b-c"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestComputeNamesInvalidResourceTypes verifies computeNames returns an error
// when one of the resource_types entries is invalid.
func TestComputeNamesInvalidResourceTypes(t *testing.T) {
	p := namingParams{
		resourceType:  "azurerm_resource_group",
		resourceTypes: []string{"azurerm_nonexistent_type"},
		name:          "myapp",
		cleanInput:    true,
		useSlug:       true,
		separator:     "-",
	}
	_, _, err := computeNames(p)
	if err == nil {
		t.Error("Expected error for invalid resource_types entry, got nil")
	}
}

// TestDataNameReadError verifies that the data source returns a diagnostic
// error when given an invalid resource type.
func TestDataNameReadError(t *testing.T) {
	provider := Provider()
	dataSource := provider.DataSourcesMap["azurecaf_name"]
	if dataSource == nil {
		t.Fatal("azurecaf_name data source not found")
	}

	// Use error_when_exceeding_max_length to trigger a validation error:
	// a very long name for a short-max-length resource
	rd := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
		"name":                            "averyveryveryverylongnamethatwillexceedmaxlength",
		"resource_type":                   "azurerm_storage_account",
		"prefixes":                        []interface{}{"prefix1", "prefix2", "prefix3"},
		"suffixes":                        []interface{}{"suffix1", "suffix2", "suffix3"},
		"clean_input":                     true,
		"use_slug":                        true,
		"error_when_exceeding_max_length": true,
	})

	diags := dataSource.ReadContext(context.Background(), rd, nil)
	if !diags.HasError() {
		t.Error("Expected error diagnostic when name exceeds max length, got none")
	}
}

// TestNamingConventionWithMaxLength covers the desiredMaxLength override branch
// in getResult where max_length < resource.MaxLength.
func TestNamingConventionWithMaxLength(t *testing.T) {
	provider := Provider()
	ncResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if ncResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	rd := schema.TestResourceDataRaw(t, ncResource.Schema, map[string]interface{}{
		"name":          "myrg",
		"convention":    "cafclassic",
		"resource_type": "rg",
		"max_length":    15,
	})

	if err := ncResource.Create(rd, nil); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	result := rd.Get("result").(string)
	if len(result) > 15 {
		t.Errorf("Expected result length <= 15, got %d (%q)", len(result), result)
	}
}

// TestAccResourceName_PlanTimeVisibility is an acceptance test that verifies
// CustomizeDiff populates result and results during the plan phase, so they
// are visible in terraform plan output (not "known after apply").
func TestAccResourceName_PlanTimeVisibility(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "azurecaf_name" "plan_test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  random_seed   = 42
  clean_input   = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("azurecaf_name.plan_test", "result", regexp.MustCompile(`^dev-rg-myapp-[a-z]{4}-001$`)),
				),
			},
		},
	})
}

// TestAccResourceName_PlanTimeMultipleTypes is an acceptance test that verifies
// CustomizeDiff populates both result and results for multiple resource types.
func TestAccResourceName_PlanTimeMultipleTypes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "azurecaf_name" "multi_test" {
  name           = "myapp"
  resource_type  = "azurerm_resource_group"
  resource_types = ["azurerm_storage_account", "azurerm_key_vault"]
  prefixes       = ["dev"]
  suffixes       = ["001"]
  random_length  = 3
  random_seed    = 100
  clean_input    = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("azurecaf_name.multi_test", "result", regexp.MustCompile(`^dev-rg-myapp-.+-001$`)),
					resource.TestCheckResourceAttrSet("azurecaf_name.multi_test", "results.azurerm_storage_account"),
					resource.TestCheckResourceAttrSet("azurecaf_name.multi_test", "results.azurerm_key_vault"),
				),
			},
		},
	})
}

// TestAccResourceName_PlanTimeAutoSeed is an acceptance test that verifies
// that when random_seed is not set but random_length > 0, the plan-apply
// cycle succeeds (falls back to apply-time computation).
func TestAccResourceName_PlanTimeAutoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "azurecaf_name" "auto_seed_test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  random_length = 5
  clean_input   = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					// Result should still be computed (at apply time) and non-empty
					resource.TestMatchResourceAttr("azurecaf_name.auto_seed_test", "result", regexp.MustCompile(`^dev-rg-myapp-[a-z]{5}$`)),
				),
			},
		},
	})
}
