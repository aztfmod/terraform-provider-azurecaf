package azurecaf

import (
	"testing"
)

func TestNamingConventionDeprecation(t *testing.T) {
	provider := Provider()
	resource := provider.ResourcesMap["azurecaf_naming_convention"]

	if resource.DeprecationMessage == "" {
		t.Error("azurecaf_naming_convention should have a deprecation message")
	}

	expectedMessage := "The azurecaf_naming_convention resource is deprecated. Please use azurecaf_name instead, which provides more flexibility and supports a broader range of Azure resources."
	if resource.DeprecationMessage != expectedMessage {
		t.Errorf("Expected deprecation message '%s', got '%s'", expectedMessage, resource.DeprecationMessage)
	}
}
