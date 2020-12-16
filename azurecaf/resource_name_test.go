package azurecaf

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_name.classic_rg",
						"pr1-pr2-rg-myrg-",
						29,
						"pr1-pr2"),
					regexMatch("azurecaf_name.classic_rg", regexp.MustCompile(ResourceDefinitions["azurerm_resource_group"].ValidationRegExp), 1),
				),
			},
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_name.classic_acr_invalid",
						"pr1pr2crmyinvalidacrname",
						35,
						"pr1pr2"),
					regexMatch("azurecaf_name.classic_acr_invalid", regexp.MustCompile(ResourceDefinitions["azurerm_container_registry"].ValidationRegExp), 1),
				),
			},
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_name.passthrough",
						"passthrough",
						11,
						""),
					regexMatch("azurecaf_name.passthrough", regexp.MustCompile(ResourceDefinitions["azurerm_container_registry"].ValidationRegExp), 1),
				),
			},
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_name.apim",
						"vsic-apim-apim",
						14,
						"vsic"),
					regexMatch("azurecaf_name.apim", regexp.MustCompile(ResourceDefinitions["azurerm_api_management_service"].ValidationRegExp), 1),
				),
			},
		},
	})
}

func TestAccResourceNameRsv_CafClassic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameCafClassicConfigRsv,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.rsv",
						"pr1-rsv-test-gm-su1",
						19,
						""),
					regexMatch("azurecaf_name.rsv", regexp.MustCompile(ResourceDefinitions["azurerm_recovery_services_vault"].ValidationRegExp), 1),
				),
			},
		},
	})
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
	actual, err := resourceNameStateUpgradeV2(testResourceNameStateDataV2(), nil)
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
