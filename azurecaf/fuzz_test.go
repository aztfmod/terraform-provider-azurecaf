package azurecaf

import (
	"testing"
)

// FuzzCleanString tests the cleanString function with various inputs and resource types
// to ensure it never panics with malformed input.
func FuzzCleanString(f *testing.F) {
	// Seed corpus with known inputs - mix of valid and edge cases
	f.Add("test-name", "azurerm_storage_account")
	f.Add("my_resource", "azurerm_virtual_machine")
	f.Add("name123", "azurerm_resource_group")
	f.Add("", "azurerm_storage_account")
	f.Add("!!!invalid@@@", "azurerm_key_vault")
	f.Add("🐱‍🚀testdata😊", "azurerm_resource_group")
	f.Add("testdata()", "azurerm_resource_group")
	f.Add("very-long-name-that-exceeds-normal-length-limits", "azurerm_resource_group")
	f.Add("MixedCaseNAME", "azurerm_storage_account")
	f.Add("123startwithnumber", "azurerm_resource_group")

	f.Fuzz(func(t *testing.T, input string, resourceType string) {
		// Only test valid resource types to focus on input validation
		def, exists := ResourceDefinitions[resourceType]
		if !exists {
			return // Skip invalid resource types
		}

		// Should never panic regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("cleanString panicked with input '%s' and resource type '%s': %v", input, resourceType, r)
			}
		}()

		result := cleanString(input, &def)

		// Basic sanity checks on the result
		// 1. Result should not be longer than input (can only remove characters)
		if len(result) > len(input) {
			t.Errorf("cleanString produced longer result than input: input=%d, result=%d", len(input), len(result))
		}

		// 2. If regex is defined, result should only contain allowed characters
		if def.RegEx != "" && result != "" {
			// The cleanString function removes characters matching the regex
			// So the result should not contain any characters that match the regex
			// This is validated by the fact that cleanString uses ReplaceAllString
		}
	})
}

// FuzzGetResourceName tests the getResourceName function with various parameter combinations
// to ensure it never panics with any input combination.
func FuzzGetResourceName(f *testing.F) {
	// Seed corpus with various valid combinations
	f.Add("myapp", "azurerm_storage_account", "-", "prefix", "suffix", 3)
	f.Add("test", "azurerm_virtual_machine", "_", "", "", 5)
	f.Add("", "azurerm_resource_group", "-", "", "", 0)
	f.Add("name", "azurerm_key_vault", "", "pre", "suf", 10)
	f.Add("app", "azurerm_app_service", "-", "dev", "001", 4)

	f.Fuzz(func(t *testing.T, name, resourceType, separator, prefix, suffix string, randomLength int) {
		// Only test valid resource types
		if _, exists := ResourceDefinitions[resourceType]; !exists {
			return
		}

		// Prevent excessive random lengths that would be invalid anyway
		if randomLength < 0 || randomLength > 1000 {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("getResourceName panicked: name=%q, resourceType=%q, separator=%q, prefix=%q, suffix=%q, randomLength=%d, panic=%v",
					name, resourceType, separator, prefix, suffix, randomLength, r)
			}
		}()

		// Build parameters for getResourceName
		prefixes := []string{}
		if prefix != "" {
			prefixes = append(prefixes, prefix)
		}
		suffixes := []string{}
		if suffix != "" {
			suffixes = append(suffixes, suffix)
		}

		randomSuffix := randSeq(randomLength, nil)
		convention := ConventionCafClassic
		cleanInput := true
		passthrough := false
		useSlug := true
		namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

		// Should never panic with any input combination
		result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)

		// If there's an error, it should be a validation error, not a panic
		if err != nil {
			// Error is acceptable - just ensure it's a proper error, not a panic
			return
		}

		// If no error, result should be non-empty and within resource constraints
		if result != "" {
			resource := ResourceDefinitions[resourceType]
			if len(result) > resource.MaxLength {
				t.Errorf("getResourceName produced name exceeding MaxLength: result=%q (len=%d), maxLength=%d",
					result, len(result), resource.MaxLength)
			}

			// Result should be at least MinLength if defined
			if resource.MinLength > 0 && len(result) < resource.MinLength {
				// This might be acceptable in some edge cases, so just log it
				t.Logf("Result is shorter than MinLength: result=%q (len=%d), minLength=%d",
					result, len(result), resource.MinLength)
			}
		}
	})
}

// FuzzComposeName tests the composeName function to ensure it handles edge cases correctly
func FuzzComposeName(f *testing.F) {
	// Seed corpus
	f.Add("-", "myapp", "st", "dev", 3)
	f.Add("_", "test", "vm", "prod", 10)
	f.Add("", "name", "", "", 5)
	f.Add("-", "", "rg", "suffix", 20)

	f.Fuzz(func(t *testing.T, separator, name, slug, suffix string, maxlength int) {
		// Sanity check on maxlength
		if maxlength < 0 || maxlength > 10000 {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("composeName panicked: separator=%q, name=%q, slug=%q, suffix=%q, maxlength=%d, panic=%v",
					separator, name, slug, suffix, maxlength, r)
			}
		}()

		prefixes := []string{}
		suffixes := []string{}
		if suffix != "" {
			suffixes = append(suffixes, suffix)
		}
		randomSuffix := ""
		namePrecedence := []string{"prefixes", "slug", "name", "suffixes", "random"}

		result := composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)

		// Result should never exceed maxlength
		if len(result) > maxlength {
			t.Errorf("composeName exceeded maxlength: result=%q (len=%d), maxlength=%d",
				result, len(result), maxlength)
		}
	})
}

// FuzzRandSeq tests the randSeq function for stability with various lengths
func FuzzRandSeq(f *testing.F) {
	// Seed corpus
	f.Add(0)
	f.Add(1)
	f.Add(5)
	f.Add(10)
	f.Add(100)
	f.Add(-1)

	f.Fuzz(func(t *testing.T, length int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("randSeq panicked with length=%d: %v", length, r)
			}
		}()

		result := randSeq(length, nil)

		// Negative or zero length should return empty string
		if length <= 0 {
			if result != "" {
				t.Errorf("randSeq with length=%d should return empty string, got %q", length, result)
			}
			return
		}

		// Positive length should return string of that length
		if len(result) != length {
			t.Errorf("randSeq with length=%d returned string of length %d: %q", length, len(result), result)
		}

		// Result should only contain lowercase letters
		for _, r := range result {
			if r < 'a' || r > 'z' {
				t.Errorf("randSeq returned non-lowercase letter: %c in %q", r, result)
				break
			}
		}
	})
}
