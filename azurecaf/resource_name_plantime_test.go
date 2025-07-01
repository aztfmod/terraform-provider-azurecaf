package azurecaf

import (
	"strings"
	"testing"
)

// TestCustomizeDiffLogic validates that plan-time calculation logic works correctly
func TestCustomizeDiffLogic(t *testing.T) {
	// Simple test to validate the core functionality works
	name := "test"
	resourceType := "azurerm_storage_account"
	prefixes := []string{"dev"}
	suffixes := []string{"001"}
	separator := "-"
	randomLength := 0
	randomSeed := int64(0)
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	randomSuffix := randSeq(randomLength, &randomSeed)

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("Error during plan-time calculation: %v", err)
	}

	t.Logf("Generated name: %s", result)
	
	// At minimum, the result should not be empty and should contain the base name
	if result == "" {
		t.Error("Result should not be empty")
	}
	
	if !strings.Contains(result, "test") && !strings.Contains(result, "st") {
		t.Errorf("Result should contain either 'test' or 'st' (slug), got: %s", result)
	}
}

// TestPlanTimeCalculationCore tests that core naming logic works correctly for plan-time calculation
func TestPlanTimeCalculationCore(t *testing.T) {
	// Test the core functionality that CustomizeDiff will use
	name := "plantest"
	resourceType := "azurerm_storage_account"
	prefixes := []string{"env"}
	suffixes := []string{}
	separator := "-"
	randomLength := 0
	randomSeed := int64(0)
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	randomSuffix := randSeq(randomLength, &randomSeed)

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("Error during plan-time calculation: %v", err)
	}

	expectedResult := "envstplantest"
	if result != expectedResult {
		t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
	}

	t.Logf("Plan-time calculated result: %s", result)
}