package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test that data source generates consistent random values for the same inputs
func TestDataSourceAzureCAFName_ConsistentRandom(t *testing.T) {
	// Create a mock resource data with the same parameters
	d1 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"name":          "test",
		"prefixes":      []interface{}{"pre"},
		"suffixes":      []interface{}{"suf"},
		"separator":     "-",
		"resource_type": "",
	})

	d2 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"name":          "test",
		"prefixes":      []interface{}{"pre"},
		"suffixes":      []interface{}{"suf"},
		"separator":     "-",
		"resource_type": "",
	})

	// Call the data source read function for both
	err1 := getNameReadResult(d1, nil)
	if err1 != nil {
		t.Fatalf("First call failed: %v", err1)
	}

	err2 := getNameReadResult(d2, nil)
	if err2 != nil {
		t.Fatalf("Second call failed: %v", err2)
	}

	// The results should be identical since the inputs are the same
	result1 := d1.Get("result").(string)
	result2 := d2.Get("result").(string)

	if result1 != result2 {
		t.Errorf("Expected consistent results, got %s and %s", result1, result2)
	}

	// Also verify the ID is set consistently
	if d1.Id() != d2.Id() {
		t.Errorf("Expected consistent IDs, got %s and %s", d1.Id(), d2.Id())
	}
}

// Test that different inputs produce different random values
func TestDataSourceAzureCAFName_DifferentInputsDifferentRandom(t *testing.T) {
	// Create a mock resource data with different parameters
	d1 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"name":          "test1",
	})

	d2 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"name":          "test2",
	})

	// Call the data source read function for both
	err1 := getNameReadResult(d1, nil)
	if err1 != nil {
		t.Fatalf("First call failed: %v", err1)
	}

	err2 := getNameReadResult(d2, nil)
	if err2 != nil {
		t.Fatalf("Second call failed: %v", err2)
	}

	// The results should be different since the inputs are different
	result1 := d1.Get("result").(string)
	result2 := d2.Get("result").(string)

	if result1 == result2 {
		t.Errorf("Expected different results for different inputs, but both were %s", result1)
	}
}

// Test that explicit random_seed is still respected
func TestDataSourceAzureCAFName_ExplicitSeed(t *testing.T) {
	// Create a mock resource data with explicit seed
	d1 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"random_seed":   12345,
	})

	d2 := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 4,
		"random_seed":   12345,
	})

	// Call the data source read function for both
	err1 := getNameReadResult(d1, nil)
	if err1 != nil {
		t.Fatalf("First call failed: %v", err1)
	}

	err2 := getNameReadResult(d2, nil)
	if err2 != nil {
		t.Fatalf("Second call failed: %v", err2)
	}

	// The results should be identical since the seed is the same
	result1 := d1.Get("result").(string)
	result2 := d2.Get("result").(string)

	if result1 != result2 {
		t.Errorf("Expected consistent results with explicit seed, got %s and %s", result1, result2)
	}
}

// Test that zero random_length returns empty suffix
func TestDataSourceAzureCAFName_ZeroRandomLength(t *testing.T) {
	d := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"random_length": 0,
		"name":          "test",
	})

	err := getNameReadResult(d, nil)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	result := d.Get("result").(string)
	expected := "test"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// Test the deterministic seed generation function
func TestGenerateDeterministicSeed(t *testing.T) {
	// Same inputs should produce same seed
	seed1 := generateDeterministicSeed("test", []string{"pre"}, []string{"suf"}, "-", "rg", 4)
	seed2 := generateDeterministicSeed("test", []string{"pre"}, []string{"suf"}, "-", "rg", 4)
	
	if seed1 != seed2 {
		t.Errorf("Expected same seed for same inputs, got %d and %d", seed1, seed2)
	}
	
	// Different inputs should produce different seeds
	seed3 := generateDeterministicSeed("different", []string{"pre"}, []string{"suf"}, "-", "rg", 4)
	
	if seed1 == seed3 {
		t.Errorf("Expected different seeds for different inputs, but both were %d", seed1)
	}
	
	// Seed should be positive
	if seed1 < 0 {
		t.Errorf("Expected positive seed, got %d", seed1)
	}
}

// Test the randSeqForDataSource function
func TestRandSeqForDataSource(t *testing.T) {
	// Test with same parameters should produce same result
	result1 := randSeqForDataSource(4, nil, "test", []string{"pre"}, []string{"suf"}, "-", "rg")
	result2 := randSeqForDataSource(4, nil, "test", []string{"pre"}, []string{"suf"}, "-", "rg")
	
	if result1 != result2 {
		t.Errorf("Expected same result for same inputs, got %s and %s", result1, result2)
	}
	
	// Test with different parameters should produce different results
	result3 := randSeqForDataSource(4, nil, "different", []string{"pre"}, []string{"suf"}, "-", "rg")
	
	if result1 == result3 {
		t.Errorf("Expected different results for different inputs, but both were %s", result1)
	}
	
	// Test with explicit seed
	seed := int64(12345)
	result4 := randSeqForDataSource(4, &seed, "test", []string{"pre"}, []string{"suf"}, "-", "rg")
	result5 := randSeqForDataSource(4, &seed, "different", []string{"pre"}, []string{"suf"}, "-", "rg")
	
	if result4 != result5 {
		t.Errorf("Expected same result with explicit seed regardless of other params, got %s and %s", result4, result5)
	}
	
	// Test zero length returns empty string
	result6 := randSeqForDataSource(0, nil, "test", []string{"pre"}, []string{"suf"}, "-", "rg")
	if result6 != "" {
		t.Errorf("Expected empty string for zero length, got %s", result6)
	}
	
	// Test that result has correct length
	if len(result1) != 4 {
		t.Errorf("Expected result length 4, got %d", len(result1))
	}
}