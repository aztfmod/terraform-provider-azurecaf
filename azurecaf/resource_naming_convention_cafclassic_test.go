package azurecaf

import (
	"testing"
)

func TestAccCafNamingConvention_Classic(t *testing.T) {
	testCases := []NamingConventionTestCase{
		{
			Name:             "log",
			Convention:       "cafclassic",
			ResourceType:     "st",
			ExpectedContains: []string{"st"},
		},
		{
			Name:             "automation",
			Convention:       "cafclassic",
			ResourceType:     "aaa",
			ExpectedContains: []string{"aaa"},
		},
		{
			Name:             "registry",
			Convention:       "cafclassic",
			ResourceType:     "acr",
			ExpectedContains: []string{"acr"},
		},
		{
			Name:             "myrg",
			Convention:       "cafclassic",
			ResourceType:     "rg",
			ExpectedContains: []string{"rg"},
		},
		{
			Name:             "passepartout",
			Convention:       "cafclassic",
			ResourceType:     "kv",
			ExpectedContains: []string{"kv"},
		},
		{
			Name:             "fire",
			Convention:       "cafclassic",
			ResourceType:     "afw",
			ExpectedContains: []string{"afw"},
		},
		{
			Name:             "recov",
			Convention:       "cafclassic",
			ResourceType:     "asr",
			ExpectedContains: []string{"asr"},
		},
		{
			Name:             "hub",
			Convention:       "cafclassic",
			ResourceType:     "evh",
			ExpectedContains: []string{"evh"},
		},
		{
			Name:             "kubedemo",
			Convention:       "cafclassic",
			ResourceType:     "aks",
			ExpectedContains: []string{"aks"},
		},
		{
			Name:             "kubedemodns",
			Convention:       "cafclassic",
			ResourceType:     "aksdns",
			ExpectedContains: []string{"aksdns"},
		},
	}

	runMultipleNamingConventionTests(t, testCases)
}
