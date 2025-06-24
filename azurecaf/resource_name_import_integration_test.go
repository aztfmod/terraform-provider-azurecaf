package azurecaf

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestResourceNameImport_IntegrationBasic tests basic import functionality for azurecaf_name resource
func TestResourceNameImport_IntegrationBasic(t *testing.T) {
	resourceName := "azurecaf_name.test"
	
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameImportBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testapp"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "azurerm_app_service"),
					resource.TestCheckResourceAttr(resourceName, "separator", "-"),
					resource.TestCheckResourceAttr(resourceName, "clean_input", "true"),
					resource.TestCheckResourceAttr(resourceName, "passthrough", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_slug", "true"),
					resource.TestCheckResourceAttr(resourceName, "random_length", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "result"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccResourceNameImportStateIdFunc,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"result", "results"},
			},
		},
	})
}

// TestResourceNameImport_IntegrationWithOptions tests import with various configuration options
func TestResourceNameImport_IntegrationWithOptions(t *testing.T) {
	resourceName := "azurecaf_name.test_options"
	
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameImportOptionsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "storage"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "azurerm_storage_account"),
					resource.TestCheckResourceAttr(resourceName, "separator", "_"),
					resource.TestCheckResourceAttr(resourceName, "clean_input", "false"),
					resource.TestCheckResourceAttr(resourceName, "passthrough", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_slug", "false"),
					resource.TestCheckResourceAttr(resourceName, "random_length", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccResourceNameImportWithOptionsStateIdFunc,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"result", "results"},
			},
		},
	})
}

// TestResourceNameImport_ErrorCases tests error handling during import
func TestResourceNameImport_ErrorCases(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testAccResourceNameImportErrorConfig,
				ResourceName:  "azurecaf_name.test_error",
				ImportState:   true,
				ImportStateId: "invalid_resource_type:testname",
				ExpectError:   regexp.MustCompile(ErrInvalidResourceType),
			},
			{
				Config:        testAccResourceNameImportErrorConfig,
				ResourceName:  "azurecaf_name.test_error",
				ImportState:   true,
				ImportStateId: "invalid_format",
				ExpectError:   regexp.MustCompile(ErrInvalidImportIDFormat),
			},
			{
				Config:        testAccResourceNameImportErrorConfig,
				ResourceName:  "azurecaf_name.test_error",
				ImportState:   true,
				ImportStateId: "",
				ExpectError:   regexp.MustCompile(ErrInvalidImportIDFormat),
			},
		},
	})
}

// TestResourceNameImport_AcceptanceStyleBasic tests acceptance-style import with basic configuration
func TestResourceNameImport_AcceptanceStyleBasic(t *testing.T) {
	resourceName := "azurecaf_name.acceptance_test"
	
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameImportAcceptanceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "webapp"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "azurerm_app_service"),
					testAccCheckResourceNameExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccResourceNameImportAcceptanceStateIdFunc,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"result", "results"},
			},
		},
	})
}

// Helper functions for import state ID generation
func testAccResourceNameImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["azurecaf_name.test"]
	if !ok {
		return "", fmt.Errorf("resource not found: azurecaf_name.test")
	}
	
	resourceType := rs.Primary.Attributes["resource_type"]
	name := rs.Primary.Attributes["name"]
	
	if resourceType == "" || name == "" {
		return "", fmt.Errorf("resource_type or name not set")
	}
	
	return fmt.Sprintf("%s:%s", resourceType, name), nil
}

func testAccResourceNameImportWithOptionsStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["azurecaf_name.test_options"]
	if !ok {
		return "", fmt.Errorf("resource not found: azurecaf_name.test_options")
	}
	
	resourceType := rs.Primary.Attributes["resource_type"]
	name := rs.Primary.Attributes["name"]
	separator := rs.Primary.Attributes["separator"]
	cleanInput := rs.Primary.Attributes["clean_input"]
	passthrough := rs.Primary.Attributes["passthrough"]
	useSlug := rs.Primary.Attributes["use_slug"]
	randomLength := rs.Primary.Attributes["random_length"]
	
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", 
		resourceType, name, separator, cleanInput, passthrough, useSlug, randomLength), nil
}

func testAccResourceNameImportAcceptanceStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["azurecaf_name.acceptance_test"]
	if !ok {
		return "", fmt.Errorf("resource not found: azurecaf_name.acceptance_test")
	}
	
	resourceType := rs.Primary.Attributes["resource_type"]
	name := rs.Primary.Attributes["name"]
	
	return fmt.Sprintf("%s:%s", resourceType, name), nil
}

// Helper function to check if resource exists
func testAccCheckResourceNameExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}
		
		return nil
	}
}

// Test configurations
const testAccResourceNameImportBasicConfig = `
resource "azurecaf_name" "test" {
  name          = "testapp"
  resource_type = "azurerm_app_service"
  separator     = "-"
  clean_input   = true
  passthrough   = false
  use_slug      = true
  random_length = 0
}
`

const testAccResourceNameImportOptionsConfig = `
resource "azurecaf_name" "test_options" {
  name          = "storage"
  resource_type = "azurerm_storage_account"
  separator     = "_"
  clean_input   = false
  passthrough   = true
  use_slug      = false
  random_length = 5
}
`

const testAccResourceNameImportErrorConfig = `
resource "azurecaf_name" "test_error" {
  name          = "errortest"
  resource_type = "azurerm_resource_group"
}
`

const testAccResourceNameImportAcceptanceConfig = `
resource "azurecaf_name" "acceptance_test" {
  name          = "webapp"
  resource_type = "azurerm_app_service"
}
`