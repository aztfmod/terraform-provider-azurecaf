package azurecaf

import (
	"testing"
)

func TestAccCafNamingConventionFull_Random(t *testing.T) {
	testCases := []NamingConventionTestCase{
		{
			Name:             "catest",
			Convention:       "random",
			ResourceType:     "st",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
		{
			Name:             "TEST-DEV-AGW-RG",
			Convention:       "random",
			ResourceType:     "azurerm_application_gateway",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
		{
			Name:             "automationtest",
			Convention:       "random",
			ResourceType:     "aaa",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
		{
			Name:             "myacr",
			Convention:       "random",
			ResourceType:     "acr",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
		{
			Name:             "mykeyvault",
			Convention:       "random",
			ResourceType:     "kv",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
		{
			Name:             "myrg",
			Convention:       "random",
			ResourceType:     "rg",
			Prefix:           "utest",
			ExpectedContains: []string{"utest"},
		},
	}

	runMultipleNamingConventionTests(t, testCases)
}
