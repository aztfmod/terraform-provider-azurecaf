package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConventionFullRandom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRandomConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_st",
						"",
						24,
						""),
					regexMatch("azurecaf_naming_convention.random_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_st2",
						"test",
						24,
						""),
					regexMatch("azurecaf_naming_convention.random_st2", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
				),
			},
		},
	})
}

const testAccResourceRandomConfig = `
provider "azurecaf" {

}

#Storage account test
resource "azurecaf_naming_convention" "random_st" {
    convention      = "random"
    name            = "log"
    resource_type   = "st"
}

resource "azurecaf_naming_convention" "random_st2" {  
	name    = "catest"
	prefix  = "test"
	resource_type    = "st"
	convention  = "random"
  }
  
`
