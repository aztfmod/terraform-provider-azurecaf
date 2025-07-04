package azurecaf

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setData(prefixes []string, name string, suffixes []string, cleanInput bool) *schema.ResourceData {
	data := &schema.ResourceData{}
	data.Set("name", name)
	data.Set("prefixes", prefixes)
	data.Set("suffixes", suffixes)
	data.Set("clean_input", cleanInput)
	return data
}

func TestCleanInput_no_changes(t *testing.T) {
	data := "testdata"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if data != result {
		t.Errorf("Expected %s but received %s", data, result)
	}
}

func TestCleanInput_remove_always(t *testing.T) {
	data := "🐱‍🚀testdata😊"
	expected := "testdata"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestCleanInput_not_remove_special_allowed_chars(t *testing.T) {
	data := "testdata()"
	expected := "testdata()"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestCleanSplice_no_changes(t *testing.T) {
	data := []string{"testdata", "test", "data"}
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanSlice(data, &resource)
	for i := range data {
		if data[i] != result[i] {
			t.Errorf("Expected %s but received %s", data[i], result[i])
		}
	}
}

func TestConcatenateParameters_azurerm_public_ip_prefix(t *testing.T) {
	prefixes := []string{"pre"}
	suffixes := []string{"suf"}
	content := []string{"name", "ip"}
	separator := "-"
	expected := "pre-name-ip-suf"
	result := concatenateParameters(separator, prefixes, content, suffixes)
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestGetSlug(t *testing.T) {
	resourceType := "azurerm_resource_group"
	convention := ConventionCafClassic
	result := getSlug(resourceType, convention)
	expected := "rg"
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestGetSlug_unknown(t *testing.T) {
	resourceType := "azurerm_does_not_exist"
	convention := ConventionCafClassic
	result := getSlug(resourceType, convention)
	expected := ""
	if result != expected {
		t.Errorf("Expected %s but received %s", expected, result)
	}
}

func TestAccResourceName_CafClassic(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	// Test case 1: Resource Group
	t.Run("ClassicResourceGroup", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "myrg",
			"resource_type": "azurerm_resource_group",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1-pr2-rg-myrg-"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_resource_group"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 2: Container App with invalid name cleaning
	t.Run("ContainerAppInvalidCleaning", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_ca_name",
			"resource_type": "azurerm_container_app",
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "ca-myinvalidcaname"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_app"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 3: Passthrough mode
	t.Run("PassthroughMode", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "passthRough",
			"resource_type": "azurerm_container_registry",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
			"passthrough":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expected := "passthrough" // passthrough mode normalizes case
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}
	})

	// Test case 4: Container App Environment
	t.Run("ContainerAppEnvironment", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_cae_name",
			"resource_type": "azurerm_container_app_environment",
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "cae-myinvalidcaename"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_app_environment"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	// Test case 5: Container Registry
	t.Run("ContainerRegistry", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "my_invalid_acr_name",
			"resource_type": "azurerm_container_registry",
			"prefixes":      []interface{}{"pr1", "pr2"},
			"suffixes":      []interface{}{"su1", "su2"},
			"random_seed":   1,
			"random_length": 5,
			"clean_input":   true,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1pr2crmyinvalidacrname"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_container_registry"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	t.Log("CAF Classic naming tests completed successfully")
}

func TestAccResourceName_CafClassicRSV(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	// Test Recovery Services Vault
	t.Run("RecoveryServicesVault", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "test",
			"resource_type": "azurerm_recovery_services_vault",
			"prefixes":      []interface{}{"pr1"},
			"suffixes":      []interface{}{"su1"},
			"random_length": 2,
			"random_seed":   1,
			"clean_input":   true,
			"passthrough":   false,
		})

		err := nameResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expectedSubstring := "pr1-rsv-test-"
		if !strings.Contains(result, expectedSubstring) {
			t.Errorf("Expected result to contain '%s', got '%s'", expectedSubstring, result)
		}

		// Should end with su1 suffix
		if !strings.HasSuffix(result, "-su1") {
			t.Errorf("Expected result to end with '-su1', got '%s'", result)
		}

		// Validate against Azure naming requirements
		if !regexp.MustCompile(ResourceDefinitions["azurerm_recovery_services_vault"].ValidationRegExp).MatchString(result) {
			t.Errorf("Result '%s' does not match Azure naming requirements", result)
		}
	})

	t.Log("CAF Classic RSV naming test completed successfully")
}

func TestComposeName(t *testing.T) {
	namePrecedence := []string{"name", "random", "slug", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 21, namePrecedence)
	expected := "a-b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrect(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 19, namePrecedence)
	expected := "b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutMaxLength(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{}
	suffixes := []string{}
	name := composeName("-", prefixes, "aaaaaaaaaa", "bla", suffixes, "", 10, namePrecedence)
	expected := "aaaaaaaaaa"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrectSuffixes(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 15, namePrecedence)
	expected := "slug-name-rd-c"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeEmptyStringArray(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"", "b"}
	suffixes := []string{"", "d"}
	name := composeName("-", prefixes, "", "", suffixes, "", 15, namePrecedence)
	expected := "b-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestValidResourceType_validParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resourceTypes := []string{"azurerm_container_registry", "azurerm_storage_account"}
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if !isValid {
		t.Logf("resource types considered invalid while input parameters are valid")
		t.Fail()
	}
	if err != nil {
		t.Logf("resource validation generated an unexpected error %s", err.Error())
		t.Fail()
	}
}
func TestValidResourceType_invalidParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resourceTypes := []string{"azurerm_not_supported", "azurerm_storage_account"}
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if isValid {
		t.Logf("resource types considered valid while input parameters are invalid")
		t.Fail()
	}
	if err == nil {
		t.Logf("resource validation did generate an error while the input is invalid")
		t.Fail()
	}
}

func TestGetResourceNameValid(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, true, namePrecedence)
	expected := "a-b-rg-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidRsv(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_recovery_services_vault", "-", []string{"a", "b"}, "test", nil, "1234", "cafclassic", true, false, true, namePrecedence)
	expected := "a-b-rsv-test-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidNoSlug(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, false, namePrecedence)
	expected := "a-b-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameInvalidResourceType(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_invalid", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, false, true, namePrecedence)
	expected := "a-b-rg-myrg-1234"

	if err == nil {
		t.Logf("Expected a validation error, got nil")
		t.Fail()
	}
	if expected == resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

func TestGetResourceNamePassthrough(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, _ := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", "cafclassic", true, true, true, namePrecedence)
	expected := "myrg"

	if expected != resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

func testResourceNameStateDataV2() map[string]interface{} {
	return map[string]interface{}{}
}

func testResourceNameStateDataV3() map[string]interface{} {
	return map[string]interface{}{
		"use_slug": true,
	}
}

func TestResourceExampleInstanceStateUpgradeV2(t *testing.T) {
	expected := testResourceNameStateDataV3()
	actual, err := resourceNameStateUpgradeV2(context.Background(), testResourceNameStateDataV2(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

const testAccResourceNameCafClassicConfig = `


# Resource Group
resource "azurecaf_name" "classic_rg" {
    name            = "myrg"
	resource_type   = "azurerm_resource_group"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_ca_invalid" {
    name            = "my_invalid_ca_name"
	resource_type   = "azurerm_container_app"
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_cae_invalid" {
    name            = "my_invalid_cae_name"
	resource_type   = "azurerm_container_app_environment"
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "classic_acr_invalid" {
    name            = "my_invalid_acr_name"
	resource_type   = "azurerm_container_registry"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
}

resource "azurecaf_name" "passthrough" {
    name            = "passthRough"
	resource_type   = "azurerm_container_registry"
	prefixes        = ["pr1", "pr2"]
	suffixes        = ["su1", "su2"]
	random_seed     = 1
	random_length   = 5
	clean_input     = true
	passthrough     = true
}


resource "azurecaf_name" "apim" {
	name = "apim"
	resource_type = "azurerm_api_management_service"
	prefixes = ["vsic"]
	random_length = 0
	clean_input = true
	passthrough = false
   }
`

const testAccResourceNameCafClassicConfigRsv = `


# Resource Group

resource "azurecaf_name" "rsv" {
    name            = "test"
	resource_type   = "azurerm_recovery_services_vault"
	prefixes        = ["pr1"]
	suffixes        = ["su1"]
	random_length   = 2
	random_seed     = 1
	clean_input     = true
	passthrough     = false
}
`
