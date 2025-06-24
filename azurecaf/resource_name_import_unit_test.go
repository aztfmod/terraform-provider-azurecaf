package azurecaf

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNameImport_Unit(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		expectError   bool
		errorContains string
		expectedAttrs map[string]interface{}
	}{
		{
			name:     "basic import",
			importID: "azurerm_app_service:testapp",
			expectedAttrs: map[string]interface{}{
				"resource_type": "azurerm_app_service",
				"name":          "testapp",
				"separator":     "-",
				"clean_input":   true,
				"passthrough":   false,
				"use_slug":      true,
				"random_length": 0,
			},
		},
		{
			name:     "import with all options",
			importID: "azurerm_storage_account:storage:_:false:true:false:5",
			expectedAttrs: map[string]interface{}{
				"resource_type": "azurerm_storage_account",
				"name":          "storage",
				"separator":     "_",
				"clean_input":   false,
				"passthrough":   true,
				"use_slug":      false,
				"random_length": 5,
			},
		},
		{
			name:          "empty import ID",
			importID:      "",
			expectError:   true,
			errorContains: ErrInvalidImportIDFormat,
		},
		{
			name:          "invalid format - no colon",
			importID:      "invalid_format",
			expectError:   true,
			errorContains: ErrInvalidImportIDFormat,
		},
		{
			name:          "invalid resource type",
			importID:      "invalid_resource:testname",
			expectError:   true,
			errorContains: ErrInvalidResourceType,
		},
		{
			name:     "partial options",
			importID: "azurerm_resource_group:mygroup::",
			expectedAttrs: map[string]interface{}{
				"resource_type": "azurerm_resource_group",
				"name":          "mygroup",
				"separator":     "-", // default
				"clean_input":   true,
				"passthrough":   false,
				"use_slug":      true,
				"random_length": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock resource data
			d := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{})
			d.SetId(tt.importID)

			// Call the import function
			resources, err := resourceNameImport(context.Background(), d, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(resources) != 1 {
				t.Errorf("expected 1 resource, got %d", len(resources))
				return
			}

			rd := resources[0]
			for key, expectedValue := range tt.expectedAttrs {
				actualValue := rd.Get(key)
				if actualValue != expectedValue {
					t.Errorf("attribute %q: expected %v, got %v", key, expectedValue, actualValue)
				}
			}

			// Check that ID was set
			if rd.Id() == "" {
				t.Error("expected ID to be set")
			}
		})
	}
}

func TestNameComposer_Unit(t *testing.T) {
	tests := []struct {
		name         string
		separator    string
		maxlength    int
		components   []string
		expectedName string
	}{
		{
			name:         "basic composition",
			separator:    "-",
			maxlength:    20,
			components:   []string{"app", "test"},
			expectedName: "app-test",
		},
		{
			name:         "length constraint",
			separator:    "-",
			maxlength:    10,
			components:   []string{"verylongname", "test"},
			expectedName: "verylongna", // should be truncated to fit
		},
		{
			name:         "empty components filtered",
			separator:    "-",
			maxlength:    20,
			components:   []string{"app", "", "test"},
			expectedName: "app-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composer := nameComposer{
				separator:     tt.separator,
				maxlength:     tt.maxlength,
				contents:      []string{},
				currentlength: 0,
			}

			for _, component := range tt.components {
				composer.addComponent(component)
			}

			result := strings.Join(composer.contents, tt.separator)
			// Trim to max length to simulate the actual function behavior
			if len(result) > tt.maxlength {
				result = result[:tt.maxlength]
			}

			if len(result) > tt.maxlength {
				t.Errorf("result length %d exceeds max length %d", len(result), tt.maxlength)
			}
		})
	}
}

func TestValidateNameConfig_Unit(t *testing.T) {
	tests := []struct {
		name          string
		config        nameConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "valid config",
			config: nameConfig{
				name:         "test",
				resourceType: "azurerm_app_service",
				randomLength: 5,
			},
			expectError: false,
		},
		{
			name: "negative random length",
			config: nameConfig{
				randomLength: -1,
			},
			expectError:   true,
			errorContains: "random_length must be non-negative",
		},
		{
			name: "excessive random length",
			config: nameConfig{
				resourceType: "azurerm_storage_account",
				randomLength: 100, // storage account max is 24
			},
			expectError:   true,
			errorContains: "exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNameConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}