package azurecaf

import (
	"regexp"
	"testing"
)

// TestDevCenterGalleryNaming tests the specific requirements for Azure Dev Center Gallery naming
func TestDevCenterGalleryNaming(t *testing.T) {
	resource := ResourceDefinitions["azurerm_dev_center_gallery"]
	exp, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		t.Fatalf("Failed to compile regex for azurerm_dev_center_gallery: %v", err)
	}

	// Test cases based on Azure requirements:
	// - Must be between 1 and 80 characters
	// - Can only include alphanumeric characters, underscores and periods
	// - Cannot start or end with '.' or '_'
	testCases := []struct {
		name     string
		input    string
		expected bool
		reason   string
	}{
		{"single_char", "a", true, "single alphanumeric character should be valid"},
		{"two_chars", "ab", true, "two alphanumeric characters should be valid"},
		{"with_underscore", "a_b", true, "alphanumeric with underscore should be valid"},
		{"with_period", "a.b", true, "alphanumeric with period should be valid"},
		{"complex_valid", "gallery_1.test", true, "complex valid name should be valid"},
		{"start_underscore", "_gallery", false, "name starting with underscore should be invalid"},
		{"end_underscore", "gallery_", false, "name ending with underscore should be invalid"},
		{"start_period", ".gallery", false, "name starting with period should be invalid"},
		{"end_period", "gallery.", false, "name ending with period should be invalid"},
		{"with_dash", "gallery-test", false, "name with dash should be invalid"},
		{"empty", "", false, "empty string should be invalid"},
	}

	// Test 80 character limit (1 start + 78 middle + 1 end)
	maxLength80 := "a"
	for i := 0; i < 78; i++ {
		maxLength80 += "x"
	}
	maxLength80 += "b"
	testCases = append(testCases, struct {
		name     string
		input    string
		expected bool
		reason   string
	}{"max_length_80", maxLength80, true, "exactly 80 characters should be valid"})

	// Test 81 characters (should fail)
	maxLength81 := maxLength80 + "c"
	testCases = append(testCases, struct {
		name     string
		input    string
		expected bool
		reason   string
	}{"exceed_max_length", maxLength81, false, "exceeding 80 characters should be invalid"})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := exp.MatchString(tc.input)
			if result != tc.expected {
				t.Errorf("Test '%s' failed: input '%s' returned %v, expected %v - %s", 
					tc.name, tc.input, result, tc.expected, tc.reason)
			}
		})
	}

	// Verify the resource definition has correct properties
	if resource.MinLength != 1 {
		t.Errorf("Expected MinLength to be 1, got %d", resource.MinLength)
	}
	if resource.MaxLength != 80 {
		t.Errorf("Expected MaxLength to be 80, got %d", resource.MaxLength)
	}
	if resource.Dashes != false {
		t.Errorf("Expected Dashes to be false, got %v", resource.Dashes)
	}
}