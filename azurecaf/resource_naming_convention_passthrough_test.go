package azurecaf

import (
	"testing"
)

func TestAccCafNamingConvention_Passthrough(t *testing.T) {
	testCases := []NamingConventionTestCase{
		{
			Name:         "logs_invalid",
			Convention:   "passthrough",
			ResourceType: "la",
			// For passthrough, we expect cleaned input
			ExpectedContains: []string{"logsinvalid"}, // underscores cleaned
		},
		{
			Name:             "TEST-DEV-AGW-RG",
			Convention:       "passthrough",
			ResourceType:     "azurerm_application_gateway",
			ExpectedContains: []string{"TEST-DEV-AGW-RG"},
		},
		{
			Name:             "myacr",
			Convention:       "passthrough",
			ResourceType:     "acr",
			ExpectedContains: []string{"myacr"},
		},
		{
			Name:             "myrg",
			Convention:       "passthrough",
			ResourceType:     "rg",
			ExpectedContains: []string{"myrg"},
		},
		{
			Name:             "mykeyvault",
			Convention:       "passthrough",
			ResourceType:     "kv",
			ExpectedContains: []string{"mykeyvault"},
		},
	}

	runMultipleNamingConventionTests(t, testCases)
}
