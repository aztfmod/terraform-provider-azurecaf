package azurecaf

import (
	"testing"
)

// Test the import functionality with specific unit tests
func TestResourceNameImportUnit(t *testing.T) {
	// Create a resource instance
	r := resourceName()

	tests := []struct {
		name                 string
		importID             string
		expectError          bool
		expectedResourceType string
		expectedName         string
	}{
		{
			name:                 "valid storage account import",
			importID:             "azurerm_storage_account:mystorageaccount123",
			expectError:          false,
			expectedResourceType: "azurerm_storage_account",
			expectedName:         "mystorageaccount123",
		},
		{
			name:                 "valid resource group import",
			importID:             "azurerm_resource_group:my-resource-group",
			expectError:          false,
			expectedResourceType: "azurerm_resource_group",
			expectedName:         "my-resource-group",
		},
		{
			name:        "invalid import ID format",
			importID:    "invalid-format",
			expectError: true,
		},
		{
			name:        "unsupported resource type",
			importID:    "invalid_resource_type:somename",
			expectError: true,
		},
		{
			name:        "invalid name for resource type",
			importID:    "azurerm_storage_account:Invalid-Storage-Account-Name!",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new ResourceData instance for each test
			d := r.TestResourceData()
			d.SetId(tt.importID)

			// Call the import function
			result, err := resourceNameImport(d, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != 1 {
				t.Errorf("Expected 1 resource data object, got %d", len(result))
				return
			}

			// Verify the imported data
			importedData := result[0]

			if importedData.Get("resource_type").(string) != tt.expectedResourceType {
				t.Errorf("Expected resource_type %s, got %s",
					tt.expectedResourceType, importedData.Get("resource_type").(string))
			}

			if importedData.Get("name").(string) != tt.expectedName {
				t.Errorf("Expected name %s, got %s",
					tt.expectedName, importedData.Get("name").(string))
			}

			if importedData.Get("passthrough").(bool) != true {
				t.Errorf("Expected passthrough to be true, got %v", importedData.Get("passthrough"))
			}

			if importedData.Get("result").(string) != tt.expectedName {
				t.Errorf("Expected result %s, got %s",
					tt.expectedName, importedData.Get("result").(string))
			}
		})
	}
}
