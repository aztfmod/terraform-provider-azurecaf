package azurecaf

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

var allResourceTypeTestNamePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}

func allResourceTypeNames() []string {
	resourceTypes := make([]string, 0, len(ResourceDefinitions))
	for resourceType := range ResourceDefinitions {
		resourceTypes = append(resourceTypes, resourceType)
	}
	sort.Strings(resourceTypes)
	return resourceTypes
}

func assertGeneratedNameMatchesDefinition(t *testing.T, resourceType string, def ResourceStructure, result string, expectSlug bool) {
	t.Helper()

	validationRegex, err := regexp.Compile(def.ValidationRegExp)
	if err != nil {
		t.Fatalf("failed to compile validation regex for %s: %v", resourceType, err)
	}

	if !validationRegex.MatchString(result) {
		t.Fatalf("generated name %q does not match validation regex %q", result, def.ValidationRegExp)
	}

	if len(result) < def.MinLength {
		t.Fatalf("generated name %q shorter than min length %d", result, def.MinLength)
	}

	if len(result) > def.MaxLength {
		t.Fatalf("generated name %q longer than max length %d", result, def.MaxLength)
	}

	if def.LowerCase && result != strings.ToLower(result) {
		t.Fatalf("generated name %q must be lowercase", result)
	}

	if expectSlug && def.CafPrefix != "" && !strings.Contains(result, def.CafPrefix) {
		t.Fatalf("generated name %q does not contain slug %q", result, def.CafPrefix)
	}
}

func firstCleanableMarker(def ResourceStructure) string {
	cleanRegex, err := regexp.Compile(def.RegEx)
	if err != nil {
		return ""
	}

	for _, candidate := range []string{"!", "@", "#", "$", "%", "?", "🚀", " "} {
		if cleanRegex.MatchString(candidate) {
			return candidate
		}
	}

	return ""
}

// TestAllResourceTypes_NameGeneration validates that every resource type
// in ResourceDefinitions produces valid names under standard conditions.
func TestAllResourceTypes_NameGeneration(t *testing.T) {
	for _, resourceType := range allResourceTypeNames() {
		resourceType := resourceType
		def := ResourceDefinitions[resourceType]

		t.Run(resourceType, func(t *testing.T) {
			t.Parallel()

			result, err := getResourceName(resourceType, "-", []string{"dev"}, "app", []string{"001"}, "", ConventionCafClassic, true, false, true, allResourceTypeTestNamePrecedence, false)
			if err != nil {
				t.Fatalf("failed to generate name: %v", err)
			}

			assertGeneratedNameMatchesDefinition(t, resourceType, def, result, true)
		})
	}
}

// TestAllResourceTypes_MaxLengthTruncation validates truncation behavior.
func TestAllResourceTypes_MaxLengthTruncation(t *testing.T) {
	for _, resourceType := range allResourceTypeNames() {
		resourceType := resourceType
		def := ResourceDefinitions[resourceType]

		t.Run(resourceType, func(t *testing.T) {
			t.Parallel()

			longName := strings.Repeat("a", def.MaxLength+20)
			result, err := getResourceName(resourceType, "-", nil, longName, nil, "", ConventionCafClassic, true, true, false, allResourceTypeTestNamePrecedence, false)
			if err != nil {
				t.Fatalf("failed to truncate long name: %v", err)
			}

			if len(result) != def.MaxLength {
				t.Fatalf("expected truncated name length %d, got %d (%q)", def.MaxLength, len(result), result)
			}

			assertGeneratedNameMatchesDefinition(t, resourceType, def, result, false)
		})
	}
}

// TestAllResourceTypes_MinLengthValidation validates minimum length enforcement.
func TestAllResourceTypes_MinLengthValidation(t *testing.T) {
	for _, resourceType := range allResourceTypeNames() {
		resourceType := resourceType
		def := ResourceDefinitions[resourceType]
		if def.MinLength <= 1 {
			continue
		}

		t.Run(resourceType, func(t *testing.T) {
			t.Parallel()

			result, err := getResourceName(resourceType, "-", nil, "", nil, "", ConventionCafClassic, true, false, true, allResourceTypeTestNamePrecedence, false)
			if err != nil {
				if !strings.Contains(err.Error(), "invalid name") {
					t.Fatalf("expected min length validation error, got: %v", err)
				}
				return
			}

			validationRegex, compileErr := regexp.Compile(def.ValidationRegExp)
			if compileErr != nil {
				t.Fatalf("failed to compile validation regex for %s: %v", resourceType, compileErr)
			}
			if !validationRegex.MatchString(result) {
				t.Fatalf("generated name %q does not match validation regex %q", result, def.ValidationRegExp)
			}
			if len(result) > def.MaxLength {
				t.Fatalf("generated name %q longer than max length %d", result, def.MaxLength)
			}
			if def.LowerCase && result != strings.ToLower(result) {
				t.Fatalf("generated name %q must be lowercase", result)
			}
			if len(result) < def.MinLength && def.CafPrefix != "" && result != def.CafPrefix {
				t.Fatalf("expected short fallback to collapse to slug %q, got %q", def.CafPrefix, result)
			}
		})
	}
}

// TestAllResourceTypes_EmptyName validates behavior with an empty name input.
func TestAllResourceTypes_EmptyName(t *testing.T) {
	for _, resourceType := range allResourceTypeNames() {
		resourceType := resourceType
		def := ResourceDefinitions[resourceType]

		t.Run(resourceType, func(t *testing.T) {
			t.Parallel()

			result, err := getResourceName(resourceType, "-", []string{"dev"}, "", []string{"001"}, "", ConventionCafClassic, true, false, true, allResourceTypeTestNamePrecedence, false)
			if err != nil {
				if !strings.Contains(err.Error(), "invalid name") {
					t.Fatalf("expected empty-name validation error, got: %v", err)
				}
				return
			}

			if result == "" {
				t.Fatal("expected non-empty result when empty name succeeds")
			}

			assertGeneratedNameMatchesDefinition(t, resourceType, def, result, true)
		})
	}
}

// TestAllResourceTypes_SpecialCharacters validates input cleaning across all resource types.
func TestAllResourceTypes_SpecialCharacters(t *testing.T) {
	for _, resourceType := range allResourceTypeNames() {
		resourceType := resourceType
		def := ResourceDefinitions[resourceType]

		t.Run(resourceType, func(t *testing.T) {
			t.Parallel()

			marker := firstCleanableMarker(def)
			prefix := "dev"
			name := "app"
			suffix := "001"
			if marker != "" {
				prefix = "de" + marker + "v"
				name = "ap" + marker + "p"
				suffix = "0" + marker + "01"
			}

			result, err := getResourceName(resourceType, "-", []string{prefix}, name, []string{suffix}, "", ConventionCafClassic, true, false, true, allResourceTypeTestNamePrecedence, false)
			if err != nil {
				t.Fatalf("failed to clean special characters: %v", err)
			}

			assertGeneratedNameMatchesDefinition(t, resourceType, def, result, true)

			if marker != "" && strings.Contains(result, marker) {
				t.Fatalf("generated name %q still contains cleaned marker %q", result, marker)
			}
		})
	}
}
