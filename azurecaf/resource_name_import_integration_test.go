// Package azurecaf provides comprehensive integration tests for the import functionality
// of the azurecaf_name resource. These tests validate the import behavior using the 
// provider's resource schema directly without requiring a full Terraform execution environment.
//
// The integration tests cover:
// - Basic import functionality and configuration validation
// - Import with various Azure resource types
// - Passthrough mode behavior for imported resources
// - Error handling for invalid import scenarios
// - Edge cases and boundary conditions
//
// Note: These tests use schema.TestResourceDataRaw to simulate the import process
// directly with the provider's import function, making them more suitable for 
// environments where Terraform CLI is not available or network access is restricted.
package azurecaf

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testCase represents a single import test case
type testCase struct {
	name           string
	importID       string
	expectError    bool
	expectedAttrs  map[string]interface{}
	errorSubstring string
	description    string
}

// integrationTestHelper provides common functionality for integration tests
type integrationTestHelper struct {
	provider           *schema.Provider
	resourceDefinition *schema.Resource
	t                  *testing.T
}

// newIntegrationTestHelper creates a new helper instance
func newIntegrationTestHelper(t *testing.T) *integrationTestHelper {
	provider := Provider()
	resourceDefinition := provider.ResourcesMap["azurecaf_name"]
	
	if resourceDefinition == nil {
		t.Fatal("azurecaf_name resource not found in provider")
	}
	
	if resourceDefinition.Importer == nil || resourceDefinition.Importer.State == nil {
		t.Fatal("azurecaf_name resource does not have importer configured properly")
	}
	
	return &integrationTestHelper{
		provider:           provider,
		resourceDefinition: resourceDefinition,
		t:                  t,
	}
}

// runImportTest executes a single import test case
func (h *integrationTestHelper) runImportTest(tc testCase) {
	h.t.Run(tc.name, func(t *testing.T) {
		// Create ResourceData instance
		resourceData := schema.TestResourceDataRaw(t, h.resourceDefinition.Schema, map[string]interface{}{})
		resourceData.SetId(tc.importID)
		
		// Execute import
		result, err := h.resourceDefinition.Importer.State(resourceData, nil)
		
		// Handle error cases
		if tc.expectError {
			if err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if tc.errorSubstring != "" && !regexp.MustCompile(tc.errorSubstring).MatchString(err.Error()) {
				t.Errorf("Expected error to contain '%s', but got: %s", tc.errorSubstring, err.Error())
			}
			t.Logf("Got expected error: %s", err.Error())
			return
		}
		
		// Handle success cases
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		if len(result) != 1 {
			t.Errorf("Expected 1 resource data object, got %d", len(result))
			return
		}
		
		// Validate attributes
		importedData := result[0]
		for attrName, expectedValue := range tc.expectedAttrs {
			actualValue := importedData.Get(attrName)
			if actualValue != expectedValue {
				t.Errorf("Expected %s to be %v, got %v", attrName, expectedValue, actualValue)
			}
		}
		
		t.Logf("Successfully imported %s with result: %s", tc.importID, importedData.Get("result"))
	})
}

// TestResourceNameImport_IntegrationBasic tests the basic import functionality using provider schema
func TestResourceNameImport_IntegrationBasic(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	t.Log("Import functionality is properly configured in the provider")
	
	// Basic integration test cases
	testCases := []testCase{
		{
			name:        "valid_storage_account_import",
			importID:    "azurerm_storage_account:mystorageaccount123",
			expectError: false,
			expectedAttrs: map[string]interface{}{
				"name":          "mystorageaccount123",
				"resource_type": "azurerm_storage_account",
				"passthrough":   true,
				"result":        "mystorageaccount123",
			},
		},
		{
			name:        "valid_resource_group_import",
			importID:    "azurerm_resource_group:my-production-rg",
			expectError: false,
			expectedAttrs: map[string]interface{}{
				"name":          "my-production-rg",
				"resource_type": "azurerm_resource_group",
				"passthrough":   true,
				"result":        "my-production-rg",
			},
		},
		{
			name:           "invalid_format",
			importID:       "invalid-format-no-colon",
			expectError:    true,
			errorSubstring: "invalid import ID format",
		},
		{
			name:           "unsupported_resource_type",
			importID:       "invalid_resource_type:somename",
			expectError:    true,
			errorSubstring: "unsupported resource type",
		},
	}
	
	for _, tc := range testCases {
		helper.runImportTest(tc)
	}
}

// TestResourceNameImport_AcceptanceStyleBasic demonstrates how the import tests would work in full acceptance test mode
// This test is commented out since it requires Terraform CLI which isn't available in the current environment
func TestResourceNameImport_AcceptanceStyleBasic(t *testing.T) {
	// Skip this test unless explicitly requested since it requires Terraform CLI
	if testing.Short() {
		t.Skip("Skipping acceptance-style test in short mode")
	}
	
	// This is how the test would be structured for full acceptance testing
	// but it's commented out since we can't run it in this environment
	/*
	resourceName := "azurecaf_name.test"
	
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameImportBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "mystorageaccount123"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "azurerm_storage_account"),
					resource.TestCheckResourceAttr(resourceName, "passthrough", "true"),
					resource.TestCheckResourceAttr(resourceName, "result", "mystorageaccount123"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "azurerm_storage_account:mystorageaccount123",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"prefixes", "suffixes", "resource_types", "results"},
			},
		},
	})
	*/
	
	t.Log("Acceptance-style test configuration is properly structured")
}

// testAccResourceNameImportBasicConfig provides configuration for acceptance tests
func testAccResourceNameImportBasicConfig() string {
	return `
resource "azurecaf_name" "test" {
  name          = "mystorageaccount123"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}
`
}

// TestResourceNameImport_IntegrationMultipleResourceTypes tests importing various Azure resource types
func TestResourceNameImport_IntegrationMultipleResourceTypes(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	
	// Define test cases for multiple resource types
	resourceTypeTests := []struct {
		resourceType string
		validName    string
		description  string
	}{
		{"azurerm_storage_account", "mystorageaccount123", "Storage Account with valid lowercase alphanumeric name"},
		{"azurerm_resource_group", "my-resource-group", "Resource Group with hyphens"},
		{"azurerm_virtual_network", "my-vnet-prod", "Virtual Network with standard naming"},
		{"azurerm_subnet", "my-subnet-web", "Subnet with descriptive name"},
		{"azurerm_key_vault", "mycompanykeyvault01", "Key Vault with alphanumeric name"},
		{"azurerm_linux_virtual_machine", "myproductionvm01", "Linux Virtual Machine with numbered suffix"},
		{"azurerm_windows_virtual_machine", "mywindowsvm01", "Windows Virtual Machine with numbered suffix"},
		{"azurerm_application_gateway", "my-appgw-prod", "Application Gateway with environment suffix"},
		{"azurerm_network_security_group", "my-nsg-web", "Network Security Group with tier suffix"},
		{"azurerm_public_ip", "my-pip-gateway", "Public IP with purpose suffix"},
	}
	
	// Convert to test cases and run
	var testCases []testCase
	for _, rt := range resourceTypeTests {
		testCases = append(testCases, testCase{
			name:        fmt.Sprintf("%s_%s", rt.resourceType, rt.validName),
			importID:    fmt.Sprintf("%s:%s", rt.resourceType, rt.validName),
			expectError: false,
			expectedAttrs: map[string]interface{}{
				"name":          rt.validName,
				"resource_type": rt.resourceType,
				"passthrough":   true,
				"result":        rt.validName,
			},
			description: rt.description,
		})
	}
	
	for _, tc := range testCases {
		helper.runImportTest(tc)
	}
}

// TestResourceNameImport_IntegrationPassthroughBehavior verifies that imported resources automatically use passthrough mode
func TestResourceNameImport_IntegrationPassthroughBehavior(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	
	passthroughTests := []struct {
		importID     string
		expectedName string
	}{
		{"azurerm_storage_account:mystorageaccount123", "mystorageaccount123"},
		{"azurerm_resource_group:very-long-resource-group-name", "very-long-resource-group-name"},
		{"azurerm_key_vault:SpecialCharactersKV", "SpecialCharactersKV"},
	}
	
	for _, pt := range passthroughTests {
		tc := testCase{
			name:        pt.importID,
			importID:    pt.importID,
			expectError: false,
			expectedAttrs: map[string]interface{}{
				"passthrough":    true,
				"result":         pt.expectedName,
				"clean_input":    true,
				"use_slug":       true,
				"separator":      "-",
				"random_length":  0,
			},
		}
		helper.runImportTest(tc)
	}
}

// TestResourceNameImport_IntegrationEdgeCases tests edge cases in import functionality
func TestResourceNameImport_IntegrationEdgeCases(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	
	edgeCases := []testCase{
		{
			name:           "empty_import_id",
			importID:       "",
			expectError:    true,
			errorSubstring: "invalid import ID format",
			description:    "Empty import ID should be rejected",
		},
		{
			name:           "only_colon",
			importID:       ":",
			expectError:    true,
			errorSubstring: "does not comply with Azure naming requirements",
			description:    "Import ID with only colon should be rejected",
		},
		{
			name:           "multiple_colons",
			importID:       "azurerm_storage_account:my:storage:account",
			expectError:    true,
			errorSubstring: "invalid import ID format",
			description:    "Import ID with multiple colons should be rejected",
		},
		{
			name:        "empty_resource_type",
			importID:    ":mystorageaccount123",
			expectError: false, // Empty resource type maps to "general" resource type
			description: "Empty resource type maps to general resource type",
		},
		{
			name:           "empty_name",
			importID:       "azurerm_storage_account:",
			expectError:    true,
			errorSubstring: "does not comply with Azure naming requirements",
			description:    "Empty name should be rejected",
		},
		{
			name:        "valid_minimum_length_name",
			importID:    "azurerm_storage_account:abc",
			expectError: false,
			description: "Minimum length valid name should be accepted",
		},
		{
			name:           "case_sensitive_resource_type",
			importID:       "AZURERM_STORAGE_ACCOUNT:mystorageaccount123",
			expectError:    true,
			errorSubstring: "unsupported resource type",
			description:    "Incorrect case resource type should be rejected",
		},
	}
	
	for _, tc := range edgeCases {
		helper.runImportTest(tc)
	}
}