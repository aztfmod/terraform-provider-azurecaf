package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestE2EBasic(t *testing.T) {
	// Build the provider first
	fmt.Println("Building terraform-provider-azurecaf...")
	cmd := exec.Command("make", "build")
	cmd.Dir = ".."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build provider: %v", err)
	}

	// Create a temporary directory for our test
	testDir, err := os.MkdirTemp("", "azurecaf-e2e-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create a simple terraform configuration
	tfConfig := `
terraform {
  required_providers {
    azurecaf = {
      source = "registry.terraform.io/aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {}

resource "azurecaf_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["test"]
}

output "result" {
  value = azurecaf_name.test.result
}
`

	// Write the terraform configuration
	configPath := filepath.Join(testDir, "main.tf")
	if err := os.WriteFile(configPath, []byte(tfConfig), 0644); err != nil {
		t.Fatalf("Failed to write terraform config: %v", err)
	}

	// Create terraform.rc for local provider override
	providerPath, _ := filepath.Abs("../terraform-provider-azurecaf")
	overrideConfig := fmt.Sprintf(`
provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}
`, filepath.Dir(providerPath))

	rcPath := filepath.Join(testDir, "terraform.rc")
	if err := os.WriteFile(rcPath, []byte(overrideConfig), 0644); err != nil {
		t.Fatalf("Failed to write terraform.rc: %v", err)
	}

	// Run terraform init
	fmt.Println("Running terraform init...")
	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = testDir
	initCmd.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+rcPath)
	if output, err := initCmd.CombinedOutput(); err != nil {
		t.Logf("Terraform init output: %s", output)
		// Init might fail with dev_overrides, that's expected
	}

	// Run terraform plan
	fmt.Println("Running terraform plan...")
	planCmd := exec.Command("terraform", "plan")
	planCmd.Dir = testDir
	planCmd.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+rcPath)
	output, err := planCmd.CombinedOutput()
	
	fmt.Printf("Terraform plan output:\n%s\n", output)
	
	if err != nil {
		t.Fatalf("Terraform plan failed: %v", err)
	}

	// Check if the output contains expected content
	outputStr := string(output)
	if !containsAll(outputStr, []string{"azurecaf_name.test", "will be created"}) {
		t.Fatalf("Terraform plan output doesn't contain expected content")
	}

	fmt.Println("✅ E2E test passed!")
}

func containsAll(text string, substrings []string) bool {
	for _, substr := range substrings {
		if !contains(text, substr) {
			return false
		}
	}
	return true
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && indexOf(text, substr) >= 0
}

func indexOf(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
