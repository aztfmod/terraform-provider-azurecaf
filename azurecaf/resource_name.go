package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNameCreate,
		UpdateContext: resourceNameUpdate,
		ReadContext:   resourceNameRead,
		Delete:        schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 4,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    schemas.V2().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.ResourceNameStateUpgradeV2,
				Version: 2,
			},
			{
				Type:    schemas.V3().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.ResourceNameStateUpgradeV3,
				Version: 3,
			},
		},
		Schema:        schemas.V4_Schema(),
		CustomizeDiff: getDifference,
	}
}

func getDifference(context context.Context, d *schema.ResourceDiff, resource interface{}) error {
	name := d.Get("name").(string)
	prefixes := convertInterfaceToString(d.Get("prefixes").([]interface{}))
	suffixes := convertInterfaceToString(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	resourceTypes := convertInterfaceToString(d.Get("resource_types").([]interface{}))
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))
	randomString := d.Get("random_string").(string)
	randomSuffix := randSeq(int(randomLength), randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
		d.SetNew("random_string", randomSuffix)
	}
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, _, err :=
		getData(resourceType, resourceTypes, separator,
			prefixes, name, suffixes, randomSuffix,
			cleanInput, passthrough, useSlug, namePrecedence)
	if !d.GetRawState().IsNull() {
		d.SetNew("result", result)
		d.SetNew("results", results)
	}
	if err != nil {
		return err
	}
	return nil
}

func resourceNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func convertInterfaceToString(source []interface{}) []string {
	s := make([]string, len(source))
	for i, v := range source {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func getNameResult(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	prefixes := convertInterfaceToString(d.Get("prefixes").([]interface{}))
	suffixes := convertInterfaceToString(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	resourceTypes := convertInterfaceToString(d.Get("resource_types").([]interface{}))
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))

	if randomSeed == 0 {
		randomSeed = time.Now().UnixMicro()
		d.Set("random_seed", randomSeed)
	}
	randomString := d.Get("random_string").(string)
	randomSuffix := randSeq(int(randomLength), randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
		d.Set("random_string", randomSuffix)
	}

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, id, err :=
		getData(resourceType, resourceTypes, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) > 0 {
		d.Set("result", result)
	}
	if len(results) > 0 {
		d.Set("results", results)
	}
	d.SetId(id)
	return diags
}

func getData(resourceType string, resourceTypes []string, separator string, prefixes []string, name string, suffixes []string, randomSuffix string, cleanInput bool, passthrough bool, useSlug bool, namePrecedence []string) (result string, results map[string]string, id string, err error) {
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if !isValid {
		return
	}
	if results == nil {
		results = make(map[string]string)
	}
	ids := []string{}
	if len(resourceType) > 0 {
		result, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		results[resourceType] = result
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceType, result))
	}

	for _, resourceTypeName := range resourceTypes {
		results[resourceTypeName], err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceTypeName, results[resourceTypeName]))
	}
	id = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))
	return
}
