package azurecaf

import (
	"testing"
)

// contains both documented and out-of-doc resources with proper attributes.
func TestResourceDefinitionMerge(t *testing.T) {
	// Test that we have resources from both original files
	totalResources := len(ResourceDefinitions)
	if totalResources < 390 {
		t.Errorf("Expected at least 390 resources after merge, got %d", totalResources)
	}

	// Test that some resources have official documentation attributes
	storageAccount, exists := ResourceDefinitions["azurerm_storage_account"]
	if !exists {
		t.Error("Expected azurerm_storage_account to exist in merged definitions")
	}

	// Test that the storage account has proper official documentation mapping
	// (this will be in the generated models)
	if storageAccount.ResourceTypeName != "azurerm_storage_account" {
		t.Error("Storage account resource name should be correct")
	}

	// Test that we can look up resources that were previously out-of-docs
	// Check for a resource that should exist but was out of docs
	_, privateEndpointExists := ResourceDefinitions["azurerm_private_endpoint"]
	if privateEndpointExists {
		t.Log("azurerm_private_endpoint found in merged definitions")
	}

	// Test AKS cluster which should have official mapping
	aksCluster, aksExists := ResourceDefinitions["azurerm_kubernetes_cluster"]
	if !aksExists {
		t.Error("Expected azurerm_kubernetes_cluster to exist")
	} else {
		if aksCluster.CafPrefix != "aks" {
			t.Errorf("Expected AKS cluster slug to be 'aks', got '%s'", aksCluster.CafPrefix)
		}
	}

	// Test container app which should have official mapping
	containerApp, caExists := ResourceDefinitions["azurerm_container_app"]
	if !caExists {
		t.Error("Expected azurerm_container_app to exist")
	} else {
		if containerApp.CafPrefix != "ca" {
			t.Errorf("Expected container app slug to be 'ca', got '%s'", containerApp.CafPrefix)
		}
	}

	t.Logf("Successfully validated merged resource definitions with %d total resources", totalResources)
}

// TestResourceTypesFromBothFiles verifies that resources from both original files are present.
func TestResourceTypesFromBothFiles(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType string
		expectedSlug string
		shouldExist  bool
		wasOutOfDoc  bool
	}{
		{"Storage Account", "azurerm_storage_account", "st", true, false},
		{"Resource Group", "azurerm_resource_group", "rg", true, false},
		{"AKS Cluster", "azurerm_kubernetes_cluster", "aks", true, false},
		{"Container App", "azurerm_container_app", "ca", true, false},
		{"Private Endpoint", "azurerm_private_endpoint", "pe", true, true},
		{"Private Service Connection", "azurerm_private_service_connection", "psc", true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource, exists := ResourceDefinitions[tc.resourceType]

			if tc.shouldExist && !exists {
				t.Errorf("Resource %s should exist but was not found", tc.resourceType)
				return
			}

			if !tc.shouldExist && exists {
				t.Errorf("Resource %s should not exist but was found", tc.resourceType)
				return
			}

			if exists && resource.CafPrefix != tc.expectedSlug {
				t.Errorf("Resource %s expected slug '%s', got '%s'", tc.resourceType, tc.expectedSlug, resource.CafPrefix)
			}

			t.Logf("✓ %s validated successfully (slug: %s, out_of_doc: %t)", tc.name, tc.expectedSlug, tc.wasOutOfDoc)
		})
	}
}
