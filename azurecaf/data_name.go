package azurecaf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// dataName creates and returns the schema for the azurecaf_name data source.
//
// This data source provides the same naming functionality as the azurecaf_name resource
// but is evaluated during the plan phase, making the generated name visible before
// resource creation. This is the recommended approach for most use cases.
//
// Key benefits of using the data source over the resource:
//   - Names are generated during terraform plan, providing early visibility
//   - No state management required (data sources are read-only)
//   - Better for single resource name generation
//   - Integrates naturally with Terraform's data flow
//
// Use the resource version when you need to generate multiple related names
// using the resource_types parameter.
func dataName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		ReadContext: dataNameRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"random_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"separator": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "-",
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
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
	}
}

func dataNameRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	if err := getNameReadResult(d); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func getNameReadResult(d *schema.ResourceData) error {
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}
	prefixesRaw, ok := d.Get("prefixes").([]interface{})
	if !ok {
		return fmt.Errorf("prefixes must be an array")
	}
	prefixes := convertInterfaceToString(prefixesRaw)
	suffixesRaw, ok := d.Get("suffixes").([]interface{})
	if !ok {
		return fmt.Errorf("suffixes must be an array")
	}
	suffixes := convertInterfaceToString(suffixesRaw)
	separator, ok := d.Get("separator").(string)
	if !ok {
		return fmt.Errorf("separator must be a string")
	}
	resourceType, ok := d.Get("resource_type").(string)
	if !ok {
		return fmt.Errorf("resource_type must be a string")
	}
	cleanInput, ok := d.Get("clean_input").(bool)
	if !ok {
		return fmt.Errorf("clean_input must be a boolean")
	}
	passthrough, ok := d.Get("passthrough").(bool)
	if !ok {
		return fmt.Errorf("passthrough must be a boolean")
	}
	useSlug, ok := d.Get("use_slug").(bool)
	if !ok {
		return fmt.Errorf("use_slug must be a boolean")
	}
	randomLength, ok := d.Get("random_length").(int)
	if !ok {
		return fmt.Errorf("random_length must be an integer")
	}
	randomSeedInt, ok := d.Get("random_seed").(int)
	if !ok {
		return fmt.Errorf("random_seed must be an integer")
	}
	randomSeed := int64(randomSeedInt)

	convention := ConventionCafClassic

	randomSuffix := randSeq(randomLength, &randomSeed)

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	resourceName, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return err
	}
	if err := d.Set("result", resourceName); err != nil {
		return fmt.Errorf("error setting result: %w", err)
	}

	d.SetId(resourceName)
	return nil
}
