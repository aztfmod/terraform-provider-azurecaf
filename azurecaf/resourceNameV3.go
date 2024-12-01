package azurecaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNameV3() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	return &schema.Resource{
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
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				},
				Optional: true,
				ForceNew: true,
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

func resourceNameStateUpgradeV3(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	randomLength := int(rawState["random_length"].(float64))

	randomSeed := int64(0)
	if rawRandomSeed := rawState["random_seed"]; rawRandomSeed != nil {
		randomSeed = int64(rawRandomSeed.(float64))
	}

	passthrough := rawState["passthrough"].(bool)

	// If a random seed is specified, use that random seed to generate the sequence based on length
	// If passthrough is enabled, generate a new random suffix (knowing that it won't impact existing outputs)
	if randomSeed > 0 || passthrough {
		rawState["random_suffix"] = randSeq(randomLength, &randomSeed)
		return nil, nil
	}

	// Merge the resource type fields in to one array
	namedResources := make(map[string]string)

	resourceType := rawState["resource_type"]
	if resourceType != nil {
		result := rawState["result"].(string)
		namedResources[resourceType.(string)] = result
	}

	resourceTypes := convertInterfaceToString(rawState["resource_types"].([]interface{}))
	results := rawState["results"].(map[string]interface{})
	for _, val := range resourceTypes {
		namedResources[val] = results[val].(string)
	}

	separator := rawState["separator"].(string)
	prefixes := []string{}
	if rawPrefixes := rawState["prefixes"]; rawPrefixes != nil {
		prefixes = convertInterfaceToString(rawPrefixes.([]interface{}))
	}

	name := ""
	if rawName := rawState["name"]; rawName != nil {
		name = rawName.(string)
	}

	suffixes := []string{}
	if rawSuffixes := rawState["suffixes"]; rawSuffixes != nil {
		suffixes = convertInterfaceToString(rawSuffixes.([]interface{}))
	}

	cleanInput := rawState["clean_input"].(bool)
	useSlug := rawState["use_slug"].(bool)

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

outer:
	for resourceTypeName, resourceName := range namedResources {
		resource, err := getResource(resourceTypeName)
		if err != nil {
			return nil, err
		}

		// Decompose the name in to a "best guesss" of the random suffix
		randomSuffix := decomposeName(resource, resourceName, separator, prefixes, name, suffixes, cleanInput, useSlug)

		log.Println(fmt.Sprintf("[INFO] Checking random suffix of %s", randomSuffix))

		// For each generated resource name
		for resourceType, expectedResourceName := range namedResources {
			// Generate the resource name using the current guess
			actualResourceName, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, ConventionCafClassic, cleanInput, passthrough, useSlug, namePrecedence)
			// If the current guess generates an error
			if err != nil {
				log.Println(fmt.Sprintf("[INFO] Got an error when generating a %s name", resourceType))
				// Stop using this guess
				continue outer
			}
			// If the current guess doesn't produce the expected value
			if !strings.EqualFold(expectedResourceName, actualResourceName) {
				log.Println(fmt.Sprintf("[INFO] Got an mismatch when generating a %s name", resourceType))
				// Stop using this guess
				continue outer
			}
		}

		rawState["random_suffix"] = randomSuffix
		return rawState, nil
	}

	return nil, fmt.Errorf("Unable to identify a consistent random suffix")
}

func decomposeName(resourceDefinition *ResourceStructure, resourceName string, separator string, prefixes []string, name string, suffixes []string, cleanInput bool, useSlug bool) string {
	if cleanInput {
		prefixes = cleanSlice(prefixes, resourceDefinition)
		name = cleanString(name, resourceDefinition)
		suffixes = cleanSlice(suffixes, resourceDefinition)
		separator = cleanString(separator, resourceDefinition)
	}

	randomSuffix := resourceName

	// For each prefix
	for _, val := range prefixes {
		// Remove each prefix (and separator if it exists)
		randomSuffix = strings.TrimPrefix(randomSuffix, val)
		randomSuffix = strings.TrimPrefix(randomSuffix, separator)
	}

	// If the CAF prefx is included
	if useSlug {
		// Remove the CAF prefix (and separator if it exists)
		randomSuffix = strings.TrimPrefix(randomSuffix, resourceDefinition.CafPrefix)
		randomSuffix = strings.TrimPrefix(randomSuffix, separator)
	}

	// If a name is included
	if len(name) > 0 {
		// Remove the name (and separator if it exists)
		randomSuffix = strings.TrimPrefix(randomSuffix, name)
		randomSuffix = strings.TrimPrefix(randomSuffix, separator)
	}

	// For each suffix (in reverse order)
	for idx := len(suffixes) - 1; idx >= 0; idx-- {
		// Remove each suffix (and separator if it exists)
		randomSuffix = strings.TrimSuffix(randomSuffix, suffixes[idx])
		randomSuffix = strings.TrimSuffix(randomSuffix, separator)
	}

	return randomSuffix
}
