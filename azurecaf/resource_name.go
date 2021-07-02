package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/google/martian/v3/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	log.Debugf("Initialized provider")

	return &schema.Resource{
		Create: resourceNameCreate,
		Update: resourceNameUpdate,
		Read:   resourceNameRead,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{

			State: importState,
		},
		SchemaVersion: 4,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNameV2().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceNameStateUpgradeV2,
				Version: 2,
			},
			{
				Type:    resourceNameV3().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceNameStateUpgradeV3,
				Version: 3,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: false,
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: false,
			},
			"random_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"separator": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "-",
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     false,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				},
				Optional: true,
				ForceNew: false,
			},
			"random_seed": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"use_slug": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
		CustomizeDiff: getDifference,
	}
}

func importState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("result", d.Id())
	d.SetId("imported")
	return []*schema.ResourceData{d}, nil
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
	convention := ConventionCafClassic
	randomSuffix := randSeq(int(randomLength), &randomSeed)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, id, err :=
		getData(resourceType, resourceTypes, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	previousResults := d.Get("results")
	previousResult := d.Get("result")
	log.Debugf("%v, %v", previousResults, previousResult)
	d.SetNew("result", result)
	switch v := previousResults.(type) {
	default:
		d.SetNew("results", results)
	case map[string]interface{}:
		for k, vv := range v {
			if _, ok := results[k]; !ok || results[k] != vv {
				d.SetNew("results", results)
			}
		}
	}
	if err != nil {
		return err
	}
	log.Debugf("%v,%v,%v", result, results, id)
	return nil
}

func resourceNameCreate(d *schema.ResourceData, meta interface{}) error {
	return getNameResult(d, meta)
}

func resourceNameUpdate(d *schema.ResourceData, meta interface{}) error {
	return getNameResult(d, meta)
}

func resourceNameRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func convertInterfaceToString(source []interface{}) []string {
	s := make([]string, len(source))
	for i, v := range source {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func getNameResult(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Get Name result, %v", d)
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

	convention := ConventionCafClassic

	randomSuffix := randSeq(int(randomLength), &randomSeed)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, id, err :=
		getData(resourceType, resourceTypes, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return err
	}
	if len(result) > 0 {
		d.Set("result", result)
	}
	if len(results) > 0 {
		d.Set("results", results)
	}
	d.SetId(id)
	return nil
}

func getData(resourceType string, resourceTypes []string, separator string, prefixes []string, name string, suffixes []string, randomSuffix string, convention string, cleanInput bool, passthrough bool, useSlug bool, namePrecedence []string) (result string, results map[string]string, id string, err error) {
	isValid, err := validateResourceType(resourceType, resourceTypes)
	if !isValid {
		return
	}
	if results == nil {
		results = make(map[string]string)
	}
	ids := []string{}
	if len(resourceType) > 0 {
		result, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		results[resourceType] = result
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceType, result))
	}

	for _, resourceTypeName := range resourceTypes {
		results[resourceTypeName], err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceTypeName, results[resourceTypeName]))
	}
	id = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))
	return
}
