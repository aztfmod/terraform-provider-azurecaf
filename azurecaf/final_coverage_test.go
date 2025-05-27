package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test getNameResult resource_type constraint with excessive random_length
func TestGetNameResultExcessiveRandomLength(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_storage_account", // max length 24
		"random_length": 25,                        // exceeds max length
	})

	err := getNameResult(rd, nil)
	if err == nil {
		t.Error("Expected error for exceeding max length but got none")
	}
}
