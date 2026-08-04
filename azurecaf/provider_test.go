package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"azurecaf": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	provider := Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, exists := provider.ResourcesMap["azurecaf_naming_convention"]; exists {
		t.Error("azurecaf_naming_convention must not be registered in v2")
	}
	if _, exists := provider.ResourcesMap["azurecaf_name"]; !exists {
		t.Error("azurecaf_name resource is not registered")
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	return nil
}
