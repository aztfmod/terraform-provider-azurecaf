package azurecaf

import (
	"testing"
)

func TestAccCafNamingConventionCaf_Random(t *testing.T) {
	testCases := []NamingConventionTestCase{
		{
			Name:             "log",
			Convention:       "cafrandom",
			ResourceType:     "st",
			Prefix:           "rdmi",
			ExpectedContains: []string{"rdmi", "st"},
		},
		{
			Name:             "rg",
			Convention:       "cafrandom",
			ResourceType:     "rg",
			Prefix:           "caf-dev-",
			ExpectedContains: []string{"caf-dev-", "rg"},
		},
		{
			Name:             "mykeyvault",
			Convention:       "cafrandom",
			ResourceType:     "kv",
			Prefix:           "rdmi",
			ExpectedContains: []string{"rdmi", "kv"},
		},
		{
			Name:             "myacr",
			Convention:       "cafrandom",
			ResourceType:     "acr",
			Prefix:           "rdmi",
			ExpectedContains: []string{"rdmi", "acr"},
		},
		{
			Name:             "aks",
			Convention:       "cafrandom",
			ResourceType:     "aks",
			Prefix:           "rdmi",
			ExpectedContains: []string{"rdmi", "aks"},
		},
		{
			Name:             "nsg",
			Convention:       "cafrandom",
			ResourceType:     "nsg",
			Prefix:           "rdmi",
			ExpectedContains: []string{"rdmi", "nsg"},
		},
	}

	runMultipleNamingConventionTests(t, testCases)
}
