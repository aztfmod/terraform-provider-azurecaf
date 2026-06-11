package azurecaf

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These tests close the coverage gaps introduced by the plan-time visibility
// feature (#336/#437) and the panic→error refactor. Each test targets a
// specific uncovered branch identified by `go tool cover -func`.

// TestResourceNameImportUnsupportedResourceType exercises the error branch in
// resourceNameImport (lines 307-309) where getResource returns a not-found
// error for a syntactically valid but unknown resource type.
func TestResourceNameImportUnsupportedResourceType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{})
	rd.SetId("azurerm_does_not_exist_xyz:somename123")

	_, err := resourceNameImport(rd, nil)
	if err == nil {
		t.Fatal("expected error for unsupported resource type in import ID")
	}
	if !strings.Contains(err.Error(), "unsupported resource type") {
		t.Fatalf("expected 'unsupported resource type' in error, got: %v", err)
	}
}

// TestComputeNamesRandomLengthExceedsResourceMax covers lines 428-430 in
// computeNames where random_length exceeds the resource type's MaxLength.
func TestComputeNamesRandomLengthExceedsResourceMax(t *testing.T) {
	// azurerm_storage_account has MaxLength=24.
	p := namingParams{
		resourceType: "azurerm_storage_account",
		randomLength: 999,
	}
	_, _, err := computeNames(p)
	if err == nil {
		t.Fatal("expected error when random_length exceeds resource MaxLength")
	}
	if !strings.Contains(err.Error(), "exceeds maximum length") {
		t.Fatalf("expected 'exceeds maximum length' in error, got: %v", err)
	}
}

// TestComputeNamesPropagatesRandSeqError covers lines 406-408 in computeNames
// where the unseeded randSeq call fails because crypto/rand is broken.
func TestComputeNamesPropagatesRandSeqError(t *testing.T) {
	restore := failingReader()
	defer restore()

	p := namingParams{
		resourceType: "azurerm_resource_group",
		randomLength: 4,
		// randomSeedSet=false → randSeq draws from crypto/rand → fails
	}
	_, _, err := computeNames(p)
	if err == nil {
		t.Fatal("expected error when crypto/rand fails during random suffix generation")
	}
	if !strings.Contains(err.Error(), "failed to generate random suffix") {
		t.Fatalf("expected wrapped randSeq error, got: %v", err)
	}
}

// TestGetNameResultPropagatesRandSeqError covers lines 713-715 in getNameResult
// where randSeq(16, nil) fails for the SetId call.
func TestGetNameResultPropagatesRandSeqError(t *testing.T) {
	restore := failingReader()
	defer restore()

	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_resource_group",
		"random_seed":   42, // seeded path so computeNames does not error
	})

	err := getNameResult(rd, nil)
	if err == nil {
		t.Fatal("expected error when crypto/rand fails during SetId generation")
	}
	if !strings.Contains(err.Error(), "failed to generate resource id") {
		t.Fatalf("expected wrapped SetId error, got: %v", err)
	}
}

// TestResourceNameCustomizeDiffNewResourceDeterministic exercises the happy
// path of resourceNameCustomizeDiff with random_seed set, so computeNames is
// run and SetNew is called for both result and results (lines 470-484).
func TestResourceNameCustomizeDiffNewResourceDeterministic(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":           "myapp",
		"resource_type":  "azurerm_resource_group",
		"resource_types": []interface{}{"azurerm_resource_group", "azurerm_storage_account"},
		"random_seed":    42,
		"random_length":  4,
	})

	diff, err := res.Diff(ctx, &terraform.InstanceState{}, config, nil)
	if err != nil {
		t.Fatalf("Diff returned unexpected error: %v", err)
	}
	if diff == nil {
		t.Fatal("expected non-nil diff")
	}
	if diff.Attributes["result"] == nil || diff.Attributes["result"].New == "" {
		t.Errorf("expected result to be set at plan time, got: %+v", diff.Attributes["result"])
	}
}

// TestResourceNameCustomizeDiffSkipsWhenRandomLengthWithoutSeed covers lines
// 459-468: when random_length > 0 but no random_seed is given, validations
// still run but result/results stay "(known after apply)".
func TestResourceNameCustomizeDiffSkipsWhenRandomLengthWithoutSeed(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":          "myapp",
		"resource_type": "azurerm_resource_group",
		"random_length": 4,
		// random_seed deliberately omitted
	})

	diff, err := res.Diff(ctx, &terraform.InstanceState{}, config, nil)
	if err != nil {
		t.Fatalf("Diff returned unexpected error: %v", err)
	}
	// result attribute should not be populated at plan time
	if diff != nil && diff.Attributes["result"] != nil && diff.Attributes["result"].New != "" {
		t.Errorf("expected result to remain unset when random_length is set without random_seed, got: %q",
			diff.Attributes["result"].New)
	}
}

// TestResourceNameCustomizeDiffSurfacesValidationErrors covers the validation
// error path in CustomizeDiff (e.g., random_length exceeds the resource
// MaxLength), ensuring computeNames errors propagate to plan output.
func TestResourceNameCustomizeDiffSurfacesValidationErrors(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":          "myapp",
		"resource_type": "azurerm_storage_account",
		"random_length": 999, // exceeds storage account MaxLength (24)
	})

	_, err := res.Diff(ctx, &terraform.InstanceState{}, config, nil)
	if err == nil {
		t.Fatal("expected validation error from CustomizeDiff for oversized random_length")
	}
	if !strings.Contains(err.Error(), "exceeds maximum length") {
		t.Fatalf("expected 'exceeds maximum length' in error, got: %v", err)
	}
}

// TestResourceNameCustomizeDiffExistingResourceNoChange covers lines 441-456
// where d.Id() != "" and no relevant attribute has changed, so the function
// returns early without recomputing.
func TestResourceNameCustomizeDiffExistingResourceNoChange(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	// State with non-empty ID and existing attribute values.
	state := &terraform.InstanceState{
		ID: "existing-resource-id",
		Attributes: map[string]string{
			"name":                            "myapp",
			"resource_type":                   "azurerm_resource_group",
			"random_length":                   "0",
			"random_seed":                     "0",
			"separator":                       "-",
			"clean_input":                     "true",
			"passthrough":                     "false",
			"use_slug":                        "true",
			"error_when_exceeding_max_length": "false",
			"result":                          "rg-myapp",
		},
	}
	// Config matches state exactly → no changes detected.
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":          "myapp",
		"resource_type": "azurerm_resource_group",
	})

	_, err := res.Diff(ctx, state, config, nil)
	if err != nil {
		t.Fatalf("Diff returned unexpected error for no-change existing resource: %v", err)
	}
}

// TestComputeNamesResourceTypesLoopError covers lines 428-430 in computeNames:
// when the resource_types loop calls getResourceName and that call fails (e.g.
// because the produced name exceeds the MaxLength for one of the types in the
// list and error_when_exceeding_max_length is true).
func TestComputeNamesResourceTypesLoopError(t *testing.T) {
	p := namingParams{
		// resourceType empty → first getResourceName block is skipped,
		// loop is the only path that calls getResourceName.
		resourceType:                "",
		resourceTypes:               []string{"azurerm_storage_account"},
		name:                        "thisnameistoolongforstorageaccount",
		separator:                   "-",
		cleanInput:                  false,
		passthrough:                 false,
		useSlug:                     true,
		errorWhenExceedingMaxLength: true,
	}
	_, _, err := computeNames(p)
	if err == nil {
		t.Fatal("expected error when resource_types loop generates an over-length name")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "exceed") &&
		!strings.Contains(strings.ToLower(err.Error()), "max") &&
		!strings.Contains(strings.ToLower(err.Error()), "length") {
		t.Fatalf("expected length-exceeded error from resource_types loop, got: %v", err)
	}
}

// TestResourceNameCustomizeDiffMainPathError covers lines 471-473 in
// resourceNameCustomizeDiff: with random_seed set (so the early
// random-length-without-seed branch is bypassed), computeNames still fails
// because the produced name exceeds the MaxLength for the resource type.
func TestResourceNameCustomizeDiffMainPathError(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":                            "thisnameistoolongforstorageaccount",
		"resource_type":                   "azurerm_storage_account",
		"random_seed":                     1, // bypasses the early random-length-without-seed branch
		"error_when_exceeding_max_length": true,
	})

	_, err := res.Diff(ctx, &terraform.InstanceState{}, config, nil)
	if err == nil {
		t.Fatal("expected validation error from CustomizeDiff main path when name exceeds MaxLength")
	}
}

// TestResourceNameCustomizeDiffExistingResourceWithChange covers lines 449-451
// in resourceNameCustomizeDiff: when d.Id() != "" and at least one tracked
// attribute changes, needsRecompute becomes true and the function proceeds to
// recompute the names instead of returning early.
func TestResourceNameCustomizeDiffExistingResourceWithChange(t *testing.T) {
	res := resourceName()
	ctx := context.Background()

	// Existing resource state: name="oldname".
	state := &terraform.InstanceState{
		ID: "existing-resource-id",
		Attributes: map[string]string{
			"name":          "oldname",
			"resource_type": "azurerm_resource_group",
			"random_seed":   "0",
			"random_length": "0",
			"result":        "rg-oldname",
		},
	}
	// New config has a different name → d.HasChange("name") returns true →
	// needsRecompute=true; break.
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":          "newname",
		"resource_type": "azurerm_resource_group",
	})

	_, err := res.Diff(ctx, state, config, nil)
	if err != nil {
		t.Fatalf("Diff returned unexpected error: %v", err)
	}
}

// TestResourceNameImportInvalidValidationRegex covers lines 307-309 in
// resourceNameImport: when the resource registered for the requested type has
// an invalid ValidationRegExp, regexp.Compile returns an error which is
// wrapped and surfaced. We inject a bad resource into ResourceDefinitions for
// the duration of the test and restore the original state afterward.
func TestResourceNameImportInvalidValidationRegex(t *testing.T) {
	const badType = "azurerm_test_badregex_import"

	orig, existed := ResourceDefinitions[badType]
	ResourceDefinitions[badType] = ResourceStructure{
		ResourceTypeName: badType,
		CafPrefix:        "bad",
		MaxLength:        50,
		MinLength:        1,
		RegEx:            ".*",
		ValidationRegExp: "(unclosed", // invalid regex
		LowerCase:        false,
		Dashes:           true,
		Scope:            "resourceGroup",
	}
	defer func() {
		if existed {
			ResourceDefinitions[badType] = orig
		} else {
			delete(ResourceDefinitions, badType)
		}
	}()

	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{})
	rd.SetId(badType + ":somename")

	_, err := resourceNameImport(rd, nil)
	if err == nil {
		t.Fatal("expected error when ValidationRegExp fails to compile")
	}
	if !strings.Contains(err.Error(), "invalid validation regex") {
		t.Fatalf("expected wrapped regex compile error, got: %v", err)
	}
}

