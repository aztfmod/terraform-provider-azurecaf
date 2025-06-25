#!/usr/bin/env python3
"""
Script to merge resourceDefinition.json and resourceDefinition_out_of_docs.json
according to the requirements:
1. Combine both files into a single resourceDefinition.json
2. Add out_of_doc: true to resources that are not in official documentation
3. Add official documentation attributes for all resources
"""

import json
import sys
import os
from collections import OrderedDict

def load_json_file(filepath):
    """Load and parse a JSON file."""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading {filepath}: {e}")
        sys.exit(1)

def save_json_file(data, filepath):
    """Save data to a JSON file with proper formatting."""
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=4, ensure_ascii=False)
    except Exception as e:
        print(f"Error saving {filepath}: {e}")
        sys.exit(1)

def get_official_doc_mapping():
    """
    Return the official documentation mapping based on the Azure CAF documentation.
    This mapping is based on the official table from:
    https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
    
    The mapping contains known resources that are officially documented.
    """
    
    # This is a subset of known mappings based on the official documentation
    # We'll start with the examples given in the issue
    official_mapping = {
        # Container services
        "azurerm_kubernetes_cluster": {
            "resource": "AKS cluster",
            "resource_provider_namespace": "Microsoft.ContainerService/managedClusters",
            "slug": "aks"
        },
        "azurerm_kubernetes_cluster_node_pool": {
            "resource": "AKS user node pool", 
            "resource_provider_namespace": "Microsoft.ContainerService/managedClusters/agentPools",
            "slug": "np"
        },
        "azurerm_container_app": {
            "resource": "Container apps",
            "resource_provider_namespace": "Microsoft.App/containerApps", 
            "slug": "ca"
        },
        "azurerm_container_app_environment": {
            "resource": "Container apps environment",
            "resource_provider_namespace": "Microsoft.App/managedEnvironments",
            "slug": "cae"
        },
        
        # Common Azure resources that are definitely in official docs
        "azurerm_storage_account": {
            "resource": "Storage account",
            "resource_provider_namespace": "Microsoft.Storage/storageAccounts",
            "slug": "st"
        },
        "azurerm_resource_group": {
            "resource": "Resource group", 
            "resource_provider_namespace": "Microsoft.Resources/resourceGroups",
            "slug": "rg"
        },
        "azurerm_virtual_machine": {
            "resource": "Virtual machine",
            "resource_provider_namespace": "Microsoft.Compute/virtualMachines", 
            "slug": "vm"
        },
        "azurerm_key_vault": {
            "resource": "Key Vault",
            "resource_provider_namespace": "Microsoft.KeyVault/vaults",
            "slug": "kv"
        },
        "azurerm_app_service": {
            "resource": "App Service",
            "resource_provider_namespace": "Microsoft.Web/sites", 
            "slug": "app"
        },
        "azurerm_virtual_network": {
            "resource": "Virtual network",
            "resource_provider_namespace": "Microsoft.Network/virtualNetworks",
            "slug": "vnet"
        },
        "azurerm_subnet": {
            "resource": "Subnet",
            "resource_provider_namespace": "Microsoft.Network/virtualNetworks/subnets", 
            "slug": "snet"
        }
    }
    
    return official_mapping

def merge_resource_definitions(main_file, out_of_docs_file, output_file):
    """
    Merge the two resource definition files according to requirements.
    """
    
    print("Loading resource definition files...")
    main_resources = load_json_file(main_file)
    out_of_docs_resources = load_json_file(out_of_docs_file)
    
    print(f"Loaded {len(main_resources)} resources from main file")
    print(f"Loaded {len(out_of_docs_resources)} resources from out of docs file")
    
    # Get official documentation mapping
    official_mapping = get_official_doc_mapping()
    
    # Create a combined list
    combined_resources = []
    
    # Process main resources (these are considered to be in official docs)
    for resource in main_resources:
        resource_name = resource.get("name", "")
        
        # Add official doc attributes if available
        if resource_name in official_mapping:
            mapping = official_mapping[resource_name]
            resource["resource"] = mapping["resource"]
            resource["resource_provider_namespace"] = mapping["resource_provider_namespace"]
            # Keep existing slug but could update if needed
            if "slug" not in resource:
                resource["slug"] = mapping["slug"]
        else:
            # For resources not in our mapping, we'll assume they are in official docs
            # and add placeholder values that can be updated later
            if "resource" not in resource:
                resource["resource"] = f"Azure {resource_name.replace('azurerm_', '').replace('_', ' ').title()}"
            if "resource_provider_namespace" not in resource:
                resource["resource_provider_namespace"] = "Unknown"
        
        combined_resources.append(resource)
    
    # Process out-of-docs resources (these get out_of_doc: true)
    for resource in out_of_docs_resources:
        resource_name = resource.get("name", "")
        
        # Add the out_of_doc flag
        resource["out_of_doc"] = True
        
        # Add official doc attributes if available (some out_of_doc resources might have been added to official docs)
        if resource_name in official_mapping:
            mapping = official_mapping[resource_name]
            resource["resource"] = mapping["resource"]
            resource["resource_provider_namespace"] = mapping["resource_provider_namespace"]
            # Keep existing slug
            if "slug" not in resource:
                resource["slug"] = mapping["slug"]
        else:
            # Add placeholder values for consistency
            if "resource" not in resource:
                resource["resource"] = f"Azure {resource_name.replace('azurerm_', '').replace('_', ' ').title()}"
            if "resource_provider_namespace" not in resource:
                resource["resource_provider_namespace"] = "Unknown"
        
        combined_resources.append(resource)
    
    # Sort by resource name for consistency
    combined_resources.sort(key=lambda x: x.get("name", ""))
    
    print(f"Combined {len(combined_resources)} resources total")
    print(f"Resources marked as out_of_doc: {len(out_of_docs_resources)}")
    
    # Save the combined file
    save_json_file(combined_resources, output_file)
    print(f"Saved combined resources to {output_file}")
    
    return True

def main():
    if len(sys.argv) != 4:
        print("Usage: python3 merge_resource_definitions.py <main_file> <out_of_docs_file> <output_file>")
        sys.exit(1)
    
    main_file = sys.argv[1]
    out_of_docs_file = sys.argv[2] 
    output_file = sys.argv[3]
    
    # Verify input files exist
    if not os.path.exists(main_file):
        print(f"Error: {main_file} does not exist")
        sys.exit(1)
    
    if not os.path.exists(out_of_docs_file):
        print(f"Error: {out_of_docs_file} does not exist")
        sys.exit(1)
    
    # Perform the merge
    success = merge_resource_definitions(main_file, out_of_docs_file, output_file)
    
    if success:
        print("Resource definition files merged successfully!")
    else:
        print("Failed to merge resource definition files")
        sys.exit(1)

if __name__ == "__main__":
    main()