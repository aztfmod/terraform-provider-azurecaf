package azurecaf

import (
	"regexp"
	"testing"
)

func TestCompileRegexValidation(t *testing.T) {
	for _, resource := range ResourceDefinitions {
		_, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			t.Logf("Error on the validation regex %s for the resource %s error %v", resource.ValidationRegExp, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
		_, err = regexp.Compile(resource.RegEx)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.RegEx, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
	}
}
