package azurecaf

import (
	"testing"
)

// BenchmarkCleanString measures the performance of input sanitization.
func BenchmarkCleanString(b *testing.B) {
	resource := ResourceDefinitions["azurerm_storage_account"]
	input := "My-App_Name!@#$%^&*() with spaces"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanString(input, &resource)
	}
}

// BenchmarkCleanString_Long measures cleaning performance with long input.
func BenchmarkCleanString_Long(b *testing.B) {
	resource := ResourceDefinitions["azurerm_storage_account"]
	input := "a-very-long-application-name-that-contains-many-special-characters!@#$and-should-be-cleaned-properly-by-the-provider"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanString(input, &resource)
	}
}

// BenchmarkConcatenateParameters measures name concatenation performance.
func BenchmarkConcatenateParameters(b *testing.B) {
	prefixes := []string{"dev", "eastus"}
	content := []string{"myapp"}
	suffixes := []string{"001", "prod"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		concatenateParameters("-", prefixes, content, suffixes)
	}
}

// BenchmarkComposeName measures the full name composition logic.
func BenchmarkComposeName(b *testing.B) {
	resource := ResourceDefinitions["azurerm_storage_account"]
	precedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composeName("-",
			[]string{"dev"},
			"myapp",
			resource.CafPrefix,
			[]string{"001"},
			"",
			resource.MaxLength,
			precedence,
			false,
		)
	}
}

// BenchmarkGetResourceName measures end-to-end name generation.
func BenchmarkGetResourceName(b *testing.B) {
	precedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getResourceName(
			"azurerm_storage_account",
			"-",
			[]string{"dev"},
			"myapp",
			[]string{"001"},
			"",
			ConventionCafClassic,
			true,
			false,
			true,
			precedence,
			false,
		)
	}
}

// BenchmarkGetResourceName_AllTypes measures generation across many resource types.
func BenchmarkGetResourceName_AllTypes(b *testing.B) {
	precedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	types := allResourceTypeNames()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt := types[i%len(types)]
		getResourceName(
			rt,
			"-",
			[]string{"dev"},
			"app",
			[]string{"001"},
			"",
			ConventionCafClassic,
			true,
			false,
			true,
			precedence,
			false,
		)
	}
}

// BenchmarkGetResource measures resource definition lookup.
func BenchmarkGetResource(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getResource("azurerm_storage_account")
	}
}

// BenchmarkGetSlug measures slug retrieval.
func BenchmarkGetSlug(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getSlug("azurerm_storage_account", ConventionCafClassic)
	}
}
