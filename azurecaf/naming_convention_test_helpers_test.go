package azurecaf

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NamingConventionTestCase represents a test case for naming conventions
type NamingConventionTestCase struct {
	Name                string
	Convention          string
	ResourceType        string
	Prefix              string
	Suffix              string
	ExpectedContains    []string
	ExpectedNotContains []string
}

// runNamingConventionTest runs a single naming convention test case
func runNamingConventionTest(t *testing.T, testCase NamingConventionTestCase) {
	provider := Provider()
	namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if namingConventionResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	// Build test data
	testData := map[string]interface{}{
		"convention":    testCase.Convention,
		"name":          testCase.Name,
		"resource_type": testCase.ResourceType,
	}

	if testCase.Prefix != "" {
		testData["prefix"] = testCase.Prefix
	}

	if testCase.Suffix != "" {
		testData["suffix"] = testCase.Suffix
	}

	resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, testData)

	err := namingConventionResource.Create(resourceData, nil)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	result := resourceData.Get("result").(string)
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Validate expected content
	for _, expected := range testCase.ExpectedContains {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected result to contain '%s', got '%s'", expected, result)
		}
	}

	// Validate content that should not be present
	for _, notExpected := range testCase.ExpectedNotContains {
		if strings.Contains(result, notExpected) {
			t.Errorf("Expected result to NOT contain '%s', got '%s'", notExpected, result)
		}
	}

	// Validate against Azure naming requirements if Resources map exists
	if resource, exists := Resources[testCase.ResourceType]; exists && resource.ValidationRegExp != "" {
		if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements for %s", result, testCase.ResourceType)
		}
	}

	t.Logf("Test case '%s' generated: %s", testCase.ResourceType, result)
}

// runMultipleNamingConventionTests runs multiple test cases in parallel
func runMultipleNamingConventionTests(t *testing.T, testCases []NamingConventionTestCase) {
	for _, testCase := range testCases {
		testCase := testCase // capture loop variable
		t.Run(testCase.ResourceType, func(t *testing.T) {
			runNamingConventionTest(t, testCase)
		})
	}
}

// createBasicTestCase creates a basic test case with common defaults
func createBasicTestCase(resourceType, convention, name string) NamingConventionTestCase {
	return NamingConventionTestCase{
		Name:             name,
		Convention:       convention,
		ResourceType:     resourceType,
		ExpectedContains: []string{}, // Will be populated by specific tests
	}
}

// validateNamingConventionResult performs common validation on naming convention results
func validateNamingConventionResult(t *testing.T, result, resourceType, convention string, shouldContain []string) {
	if result == "" {
		t.Error("Expected non-empty result")
		return
	}

	// Validate expected content
	for _, expected := range shouldContain {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected result to contain '%s', got '%s'", expected, result)
		}
	}

	// Validate against Azure naming requirements
	if resource, exists := Resources[resourceType]; exists && resource.ValidationRegExp != "" {
		if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements for %s", result, resourceType)
		}
	}

	t.Logf("Convention '%s' for '%s' generated: %s", convention, resourceType, result)
}
