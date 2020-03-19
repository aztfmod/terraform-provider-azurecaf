package caf

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCafNamingConvention(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCafConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingRandomStorageAccount(
						"caf_naming_convention.nc",
						"log",
						Resources["st"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.nc", regexp.MustCompile("^[0-9A-Za-z]{3,24}$"), 1),
				),
			},
		},
	})
}

func testAccCafNamingRandomStorageAccount(id string, name string, expectedLength int, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attrs := rs.Primary.Attributes

		result := attrs["result"]
		if len(result) != expectedLength {
			return fmt.Errorf("got %s %d result items; want %d", result, len(result), expectedLength)
		}
		if !strings.HasPrefix(result, prefix) {
			return fmt.Errorf("got %s which doesn't start with %s", result, prefix)
		}
		if !strings.Contains(result, name) {
			return fmt.Errorf("got %s which doesn't contain the name %s", result, name)
		}
		return nil
	}
}

func regexMatch(id string, exp *regexp.Regexp, requiredMatches int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		result := rs.Primary.Attributes["result"]

		if matches := exp.FindAllStringSubmatchIndex(result, -1); len(matches) != requiredMatches {
			return fmt.Errorf("result string is %s; did not match %s, got %d", result, exp, len(matches))
		}

		return nil
	}
}

const testAccResourceCafConfig = `
resource "caf_naming_convention" "nc" {
    convention      = "cafrandom"
    name            = "log"
    prefix          = "rdmi-"
    resource_type   = "st"
}
`
