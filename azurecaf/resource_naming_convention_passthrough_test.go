package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConventionPassthrough(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePassthroughConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.logs_inv",
						"logsinvalid",
						11,
						"log"),
					regexMatch("azurecaf_naming_convention.logs_inv", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
				),
			},
		},
	})
}

const testAccResourcePassthroughConfig = `
provider "azurecaf" {

}

#Storage account test
resource "azurecaf_naming_convention" "logs_inv" {
    convention      = "passthrough"
    name            = "logs_invalid"
    resource_type   = "la"
}
`
