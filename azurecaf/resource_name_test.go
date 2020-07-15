package azurecaf

import (
	"fmt"
	"testing"

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

func TestHello(t *testing.T) {
	data := setData([]string{}, "", []string{}, true)
	fmt.Println(data)
}
