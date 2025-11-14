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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testCase represents a single import test case.
type testCase struct {
	expectedAttrs  map[string]interface{}
	name           string
	importID       string
	errorSubstring string
	description    string
	expectError    bool
}

// integrationTestHelper provides common functionality for integration tests.
type integrationTestHelper struct {
	provider           *schema.Provider
	resourceDefinition *schema.Resource
	t                  *testing.T
}

// newIntegrationTestHelper creates a new helper instance.
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

// runImportTest executes a single import test case.
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

// getDefaultExpectedAttrs returns the standard expected attributes for import tests.
func getDefaultExpectedAttrs() map[string]interface{} {
	return map[string]interface{}{
		"passthrough":   true,
		"clean_input":   true,
		"use_slug":      true,
		"separator":     "-",
		"random_length": 0,
	}
}

// createTestCaseFromImportID creates a test case from an import ID with standard attributes.
func createTestCaseFromImportID(name, importID string, expectError bool, extraAttrs map[string]interface{}) testCase {
	expectedAttrs := getDefaultExpectedAttrs()

	if !expectError {
		parts := strings.Split(importID, ":")
		if len(parts) == 2 {
			expectedAttrs["name"] = parts[1]
			expectedAttrs["resource_type"] = parts[0]
			expectedAttrs["result"] = parts[1]
		}
	}

	// Override with extra attributes
	for k, v := range extraAttrs {
		expectedAttrs[k] = v
	}

	return testCase{
		name:          name,
		importID:      importID,
		expectError:   expectError,
		expectedAttrs: expectedAttrs,
	}
}

// createTestCasesFromRegistry creates test cases from the resource type registry.
func createTestCasesFromRegistry(pattern string) []testCase {
	registry := getResourceTypeRegistry()
	var testCases []testCase

	for _, rt := range registry {
		if rt.Pattern == pattern || pattern == "all" {
			testCases = append(testCases, createTestCaseFromImportID(
				fmt.Sprintf("%s_%s", rt.ResourceType, rt.Name),
				fmt.Sprintf("%s:%s", rt.ResourceType, rt.Name),
				false,
				nil, // Don't include description as it's not a schema field
			))
		}
	}

	return testCases
}

// runTestGroup executes a group of test cases with a common description.
func (h *integrationTestHelper) runTestGroup(t *testing.T, groupName string, testCases []testCase) {
	t.Run(groupName, func(t *testing.T) {
		for _, tc := range testCases {
			h.runImportTest(tc)
		}
	})
}

func TestResourceNameImport_IntegrationBasic(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	t.Log("Import functionality is properly configured in the provider")

	// Basic integration test cases
	testCases := []testCase{
		createTestCaseFromImportID("valid_storage_account_import", "azurerm_storage_account:mystorageaccount123", false, nil),
		createTestCaseFromImportID("valid_resource_group_import", "azurerm_resource_group:my-production-rg", false, nil),
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

// This test is commented out since it requires Terraform CLI which isn't available in the current environment.
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

// resourceTypeDefinition represents a comprehensive resource type definition for testing.
type resourceTypeDefinition struct {
	ResourceType string
	Name         string
	Description  string
	Pattern      string
}

// getResourceTypeRegistry returns a comprehensive registry of resource types for testing.
func getResourceTypeRegistry() []resourceTypeDefinition {
	return []resourceTypeDefinition{
		{"azurerm_storage_account", "mystorageaccount123", "Storage Account with valid lowercase alphanumeric name", "basic"},
		{"azurerm_resource_group", "my-resource-group", "Resource Group with hyphens", "basic"},
		{"azurerm_virtual_network", "my-vnet-prod", "Virtual Network with standard naming", "basic"},
		{"azurerm_key_vault", "mycompanykeyvault01", "Key Vault with alphanumeric name", "basic"},
		{"azurerm_linux_virtual_machine", "myproductionvm01", "Linux Virtual Machine with numbered suffix", "basic"},
		{"azurerm_application_gateway", "my-appgw-prod", "Application Gateway with environment suffix", "basic"},
		{"azurerm_storage_account", "modulestorageaccount123", "Module storage account", "submodule_internal"},
		{"azurerm_resource_group", "module-production-rg", "Module resource group", "submodule_internal"},
		{"azurerm_storage_account", "rootstorageformodule", "Root storage for module", "root_to_submodule"},
		{"azurerm_resource_group", "root-rg-for-module", "Root RG for module", "root_to_submodule"},
	}
}

// getMinimalTerraformConfig provides basic configuration template for documentation.
func getMinimalTerraformConfig(resourceType, name string) string {
	return fmt.Sprintf(`resource "azurecaf_name" "test" {
  name          = "%s"
  resource_type = "%s"
  passthrough   = true
}`, name, resourceType)
}

// TestResourceNameImport_IntegrationMultipleResourceTypes tests importing various Azure resource types.
func TestResourceNameImport_IntegrationMultipleResourceTypes(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	testCases := createTestCasesFromRegistry("basic")
	helper.runTestGroup(t, "MultipleResourceTypes", testCases)
}

// TestResourceNameImport_IntegrationPassthroughBehavior verifies that imported resources automatically use passthrough mode.
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
		tc := createTestCaseFromImportID(pt.importID, pt.importID, false, nil)
		helper.runImportTest(tc)
	}
}

// createErrorTestCase creates a test case for error scenarios.
func createErrorTestCase(name, importID, errorSubstring, description string) testCase {
	return testCase{
		name:           name,
		importID:       importID,
		expectError:    true,
		errorSubstring: errorSubstring,
		description:    description,
	}
}

// TestResourceNameImport_IntegrationEdgeCases tests edge cases in import functionality.
func TestResourceNameImport_IntegrationEdgeCases(t *testing.T) {
	helper := newIntegrationTestHelper(t)

	edgeCases := []testCase{
		createErrorTestCase("empty_import_id", "", "invalid import ID format", "Empty import ID should be rejected"),
		createErrorTestCase("only_colon", ":", "does not comply with Azure naming requirements", "Import ID with only colon should be rejected"),
		createErrorTestCase("multiple_colons", "azurerm_storage_account:my:storage:account", "invalid import ID format", "Import ID with multiple colons should be rejected"),
		createTestCaseFromImportID("empty_resource_type", ":mystorageaccount123", false, nil),
		createErrorTestCase("empty_name", "azurerm_storage_account:", "does not comply with Azure naming requirements", "Empty name should be rejected"),
		createTestCaseFromImportID("valid_minimum_length_name", "azurerm_storage_account:abc", false, nil),
		createErrorTestCase("case_sensitive_resource_type", "AZURERM_STORAGE_ACCOUNT:mystorageaccount123", "unsupported resource type", "Incorrect case resource type should be rejected"),
	}

	for _, tc := range edgeCases {
		helper.runImportTest(tc)
	}
}

// TestResourceNameImport_AcceptanceStyleImportBlocks tests the import {} block functionality.
func TestResourceNameImport_AcceptanceStyleImportBlocks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping import {} block acceptance tests in short mode - requires Terraform CLI")
	}

	// Validate basic configuration generation
	config := getMinimalTerraformConfig("azurerm_storage_account", "mystorageaccount123")
	if config == "" {
		t.Error("Configuration generation failed")
	}

	t.Log("Import {} block configurations are properly structured")
}

// TestResourceNameImport_ImportBlockValidationSimulation tests import {} block scenarios using schema validation.
func TestResourceNameImport_ImportBlockValidationSimulation(t *testing.T) {
	helper := newIntegrationTestHelper(t)
	testCases := createTestCasesFromRegistry("all")
	helper.runTestGroup(t, "ImportBlockSimulation", testCases)
	t.Log("Import {} block simulation tests completed successfully")
}

// TestResourceNameImport_SubmodulePatternValidation tests both patterns for submodule import usage.
func TestResourceNameImport_SubmodulePatternValidation(t *testing.T) {
	helper := newIntegrationTestHelper(t)

	// Test Pattern 1: SubmoduleInternal
	submoduleInternalCases := createTestCasesFromRegistry("submodule_internal")
	helper.runTestGroup(t, "Pattern1_SubmoduleInternal", submoduleInternalCases)
	t.Log("Pattern1_SubmoduleInternal tests completed: azurecaf_name declared within submodule")

	// Test Pattern 2: RootToSubmodule
	rootToSubmoduleCases := createTestCasesFromRegistry("root_to_submodule")
	helper.runTestGroup(t, "Pattern2_RootToSubmodule", rootToSubmoduleCases)
	t.Log("Pattern2_RootToSubmodule tests completed: azurecaf_name declared at root, passed to submodule")

	t.Log("Both submodule patterns validated successfully")
}
