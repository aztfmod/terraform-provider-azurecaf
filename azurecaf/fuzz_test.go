package azurecaf

import (
	"strings"
	"testing"
)

func FuzzCleanString(f *testing.F) {
	f.Add("hello-world")
	f.Add("🐱test😊")
	f.Add("")
	f.Add(strings.Repeat("x", 10000))
	f.Add("special!@#$%^&*()")
	f.Add("UPPERCASE")
	f.Add("with spaces and\ttabs")

	resources := []ResourceStructure{
		ResourceDefinitions["azurerm_resource_group"],
		ResourceDefinitions["azurerm_storage_account"],
		ResourceDefinitions["azurerm_application_gateway"],
		{
			ResourceTypeName: "invalid-regex",
			RegEx:            "[",
		},
	}

	f.Fuzz(func(t *testing.T, input string) {
		for i := range resources {
			_ = cleanString(input, &resources[i])
		}
	})
}

func FuzzConcatenateParameters(f *testing.F) {
	f.Add("-", "pre", "prefix2", "name", "suffix", "suffix2")
	f.Add("", "", "", "", "", "")
	f.Add("_", "a", "b", "c", "d", "e")

	f.Fuzz(func(t *testing.T, separator, prefix1, prefix2, name, suffix1, suffix2 string) {
		prefixes := []string{prefix1, prefix2}
		content := []string{name}
		suffixes := []string{suffix1, suffix2}
		_ = concatenateParameters(separator, prefixes, content, suffixes)
	})
}

func FuzzComposeName(f *testing.F) {
	f.Add("-", "pre", "pre2", "name", "slug", "suffix", "suffix2", "rand", 24, false)
	f.Add("", "", "", "", "", "", "", "", 0, false)
	f.Add("_", "alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", 12, true)

	f.Fuzz(func(t *testing.T, separator, prefix1, prefix2, name, slug, suffix1, suffix2, randomSuffix string, maxLength int, errorWhenExceedingMaxLength bool) {
		prefixes := []string{prefix1, prefix2}
		suffixes := []string{suffix1, suffix2}
		namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes", "suffixes", "prefixes"}

		_, _ = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxLength, namePrecedence, errorWhenExceedingMaxLength)
	})
}

func FuzzGetResourceName(f *testing.F) {
	f.Add("seed", "-", "pre", "pre2", "name", "suffix", "suffix2", "rand", true, false, true, false)
	f.Add("", "", "", "", "", "", "", "", false, true, false, true)

	resourceTypes := []string{
		"azurerm_resource_group",
		"azurerm_storage_account",
		"azurerm_application_gateway",
		"azurerm_key_vault",
	}
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	f.Fuzz(func(t *testing.T, resourceHint, separator, prefix1, prefix2, name, suffix1, suffix2, randomSuffix string, cleanInput, passthrough, useSlug, errorWhenExceedingMaxLength bool) {
		resourceType := resourceTypes[0]
		if len(resourceHint) > 0 {
			resourceType = resourceTypes[len(resourceHint)%len(resourceTypes)]
		}

		_, _ = getResourceName(
			resourceType,
			separator,
			[]string{prefix1, prefix2},
			name,
			[]string{suffix1, suffix2},
			randomSuffix,
			ConventionCafClassic,
			cleanInput,
			passthrough,
			useSlug,
			namePrecedence,
			errorWhenExceedingMaxLength,
		)
	})
}

func FuzzNameBuilder(f *testing.F) {
	f.Add(63, "-", "rg", "app", "prod", true, false, true)
	f.Add(0, "", "", "", "", false, false, false)
	f.Add(10, "_", "verylongsegment", "x", "y", false, true, false)

	f.Fuzz(func(t *testing.T, maxLength int, separator, first, second, third string, prependFirst, prependSecond, prependThird bool) {
		builder := NewNameBuilder(maxLength, separator)

		if prependFirst {
			builder.Prepend(first)
		} else {
			builder.Append(first)
		}
		if prependSecond {
			builder.Prepend(second)
		} else {
			builder.Append(second)
		}
		if prependThird {
			builder.Prepend(third)
		} else {
			builder.Append(third)
		}

		_ = builder.GetName()
		_ = builder.GetTrimmedName()
	})
}
