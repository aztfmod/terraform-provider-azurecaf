package azurecaf

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceNameV2 creates and returns the schema for the azurecaf_name resource (version 2).
// This resource generates Azure-compliant resource names following Cloud Adoption Framework
// naming conventions and Azure resource naming requirements.
//
// The resource supports:
//   - Multiple naming conventions (CAF classic, CAF random, passthrough, etc.)
//   - Custom prefixes and suffixes
//   - Random character generation with configurable length
//   - Input sanitization and validation
//   - Multiple resource types in a single configuration
//
// This is an improved version that supersedes the original azurecaf_naming_convention resource.
func resourceNameV2() *schema.Resource {
	// Get all available resource types for validation
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Base name for the resource (will be sanitized according to Azure rules)
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			// List of prefixes to add before the resource name
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			// List of suffixes to add after the resource name
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			// Number of random characters to append to the name
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
		},
	}
}

func resourceNameStateUpgradeV2(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	rawState["use_slug"] = true

	return rawState, nil
}

func resourceName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		Create:        resourceNameCreate,
		Read:          schema.Noop,
		Delete:        schema.RemoveFromState,
		SchemaVersion: 3,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNameV2().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceNameStateUpgradeV2,
				Version: 2,
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceNameImport,
		},

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

func resourceNameCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceNameRead(d, meta)
}

func resourceNameRead(d *schema.ResourceData, meta interface{}) error {
	return getNameResult(d, meta)
}

func resourceNameDelete(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

// resourceNameImport imports an existing Azure resource name into the Terraform state.
// 
// The import ID must be in the format "<resource_type>:<existing_name>", for example:
//   azurerm_storage_account:mystorageaccount123
//
// The function validates the resource type and the existing name against Azure naming requirements,
// sets the resource data fields accordingly (using passthrough mode), and uses the existing name
// as the Terraform resource ID.
func resourceNameImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	// Parse the import ID
	parts := strings.Split(importID, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected '<resource_type>:<existing_name>', got: %s", importID)
	}

	resourceType := parts[0]
	existingName := parts[1]

	// Validate the resource type exists
	resource, err := getResource(resourceType)
	if err != nil {
		return nil, fmt.Errorf("unsupported resource type '%s': %w", resourceType, err)
	}

	// Validate the existing name against Azure naming rules for this resource type
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		return nil, fmt.Errorf("invalid validation regex for resource type '%s': %w", resourceType, err)
	}

	if !validationRegEx.MatchString(existingName) {
		return nil, fmt.Errorf("existing name '%s' does not comply with Azure naming requirements for resource type '%s'. Expected pattern: %s",
			existingName, resourceType, resource.ValidationRegExp)
	}

	// Set the resource data for the imported resource
	// We use passthrough mode to preserve the existing name as-is
	if err := d.Set("name", existingName); err != nil {
		return nil, fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("resource_type", resourceType); err != nil {
		return nil, fmt.Errorf("error setting resource_type: %w", err)
	}
	if err := d.Set("passthrough", true); err != nil {
		return nil, fmt.Errorf("error setting passthrough: %w", err)
	}

	// Set empty slices for prefixes and suffixes since we can't reverse-engineer them
	if err := d.Set("prefixes", []string{}); err != nil {
		return nil, fmt.Errorf("error setting prefixes: %w", err)
	}
	if err := d.Set("suffixes", []string{}); err != nil {
		return nil, fmt.Errorf("error setting suffixes: %w", err)
	}
	if err := d.Set("resource_types", []string{}); err != nil {
		return nil, fmt.Errorf("error setting resource_types: %w", err)
	}

	// Set the result to match the imported name
	if err := d.Set("result", existingName); err != nil {
		return nil, fmt.Errorf("error setting result: %w", err)
	}
	if err := d.Set("results", map[string]string{}); err != nil {
		return nil, fmt.Errorf("error setting results: %w", err)
	}

	// Use the existing name as the Terraform resource ID
	d.SetId(existingName)

	return []*schema.ResourceData{d}, nil
}

func cleanSlice(names []string, resourceDefinition *ResourceStructure) []string {
	for i, name := range names {
		names[i] = cleanString(name, resourceDefinition)
	}
	return names
}

func cleanString(name string, resourceDefinition *ResourceStructure) string {
	myRegex, _ := regexp.Compile(resourceDefinition.RegEx)
	return myRegex.ReplaceAllString(name, "")
}

func concatenateParameters(separator string, parameters ...[]string) string {
	elems := []string{}
	for _, items := range parameters {
		for _, item := range items {
			if item != "" {
				elems = append(elems, []string{item}...)
			}
		}
	}
	return strings.Join(elems, separator)
}

func getResource(resourceType string) (*ResourceStructure, error) {
	if resourceKey, existing := ResourceMaps[resourceType]; existing {
		resourceType = resourceKey
	}
	if resource, resourceFound := ResourceDefinitions[resourceType]; resourceFound {
		return &resource, nil
	}
	return nil, fmt.Errorf("invalid resource type %s", resourceType)
}

// Retrieve the resource slug / shortname based on the resourceType and the selected convention.
func getSlug(resourceType, convention string) string {
	if convention == ConventionCafClassic || convention == ConventionCafRandom {
		if val, ok := ResourceDefinitions[resourceType]; ok {
			return val.CafPrefix
		}
	}
	return ""
}

func trimResourceName(resourceName string, maxLength int) string {
	length := len(resourceName)

	if length > maxLength {
		length = maxLength
	}

	return resourceName[0:length]
}

func convertInterfaceToString(source []interface{}) []string {
	s := make([]string, len(source))
	for i, v := range source {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func composeName(separator string,
	prefixes []string,
	name string,
	slug string,
	suffixes []string,
	randomSuffix string,
	maxlength int,
	namePrecedence []string,
) string {
	contents := []string{}
	currentlength := 0

	for i := 0; i < len(namePrecedence); i++ {
		initialized := 0
		if len(contents) > 0 {
			initialized = len(separator)
		}
		switch c := namePrecedence[i]; c {
		case "name":
			if name != "" {
				if currentlength+len(name)+initialized <= maxlength {
					contents = append(contents, name)
					currentlength = currentlength + len(name) + initialized
				}
			}
		case "slug":
			if slug != "" {
				if currentlength+len(slug)+initialized <= maxlength {
					contents = append([]string{slug}, contents...)
					currentlength = currentlength + len(slug) + initialized
				}
			}
		case "random":
			if randomSuffix != "" {
				if currentlength+len(randomSuffix)+initialized <= maxlength {
					contents = append(contents, randomSuffix)
					currentlength = currentlength + len(randomSuffix) + initialized
				}
			}
		case "suffixes":
			if len(suffixes) > 0 {
				if suffixes[0] != "" {
					if currentlength+len(suffixes[0])+initialized <= maxlength {
						contents = append(contents, suffixes[0])
						currentlength = currentlength + len(suffixes[0]) + initialized
					}
				}
				suffixes = suffixes[1:]
				if len(suffixes) > 0 {
					i--
				}
			}
		case "prefixes":
			if len(prefixes) > 0 {
				if prefixes[len(prefixes)-1] != "" {
					if currentlength+len(prefixes[len(prefixes)-1])+initialized <= maxlength {
						contents = append([]string{prefixes[len(prefixes)-1]}, contents...)
						currentlength = currentlength + len(prefixes[len(prefixes)-1]) + initialized
					}
				}
				prefixes = prefixes[:len(prefixes)-1]
				if len(prefixes) > 0 {
					i--
				}
			}
		}
	}
	content := strings.Join(contents, separator)
	return content
}

func validateResourceType(resourceType string, resourceTypes []string) (bool, error) {
	isEmpty := resourceType == "" && len(resourceTypes) == 0
	if isEmpty {
		return false, fmt.Errorf("resource_type and resource_types parameters are empty, you must specify at least one resource type")
	}
	errorStrings := []string{}
	resourceList := resourceTypes
	if resourceType != "" {
		resourceList = append(resourceList, resourceType)
	}

	for _, resource := range resourceList {
		_, err := getResource(resource)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	if len(errorStrings) > 0 {
		return false, fmt.Errorf("%s", strings.Join(errorStrings, "\n"))
	}
	return true, nil
}

func getResourceName(resourceTypeName string, separator string,
	prefixes []string,
	name string,
	suffixes []string,
	randomSuffix string,
	convention string,
	cleanInput bool,
	passthrough bool,
	useSlug bool,
	namePrecedence []string,
) (string, error) {
	resource, err := getResource(resourceTypeName)
	if err != nil {
		return "", err
	}
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		return "", err
	}

	slug := ""
	if useSlug {
		slug = getSlug(resourceTypeName, convention)
	}

	if cleanInput {
		prefixes = cleanSlice(prefixes, resource)
		suffixes = cleanSlice(suffixes, resource)
		name = cleanString(name, resource)
		separator = cleanString(separator, resource)
		randomSuffix = cleanString(randomSuffix, resource)
	}

	var resourceName string

	if passthrough {
		resourceName = name
	} else {
		resourceName = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, resource.MaxLength, namePrecedence)
	}
	resourceName = trimResourceName(resourceName, resource.MaxLength)

	if resource.LowerCase {
		resourceName = strings.ToLower(resourceName)
	}

	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	return resourceName, nil
}

func getNameResult(d *schema.ResourceData, _ interface{}) error {
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
	resourceTypesRaw, ok := d.Get("resource_types").([]interface{})
	if !ok {
		return fmt.Errorf("resource_types must be an array")
	}
	resourceTypes := convertInterfaceToString(resourceTypesRaw)
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

	// Validate random_length parameter
	if randomLength < 0 {
		return fmt.Errorf("random_length must be non-negative, got: %d", randomLength)
	}

	// Validate against resource type constraints if resource_type is specified
	if resourceType != "" {
		if resource, exists := ResourceDefinitions[resourceType]; exists {
			maxLen := resource.MaxLength
			if randomLength > maxLen {
				return fmt.Errorf("random_length (%d) exceeds maximum length for resource type %s (%d)", randomLength, resourceType, maxLen)
			}
		}
	}

	convention := ConventionCafClassic

	randomSuffix := randSeq(randomLength, &randomSeed)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	isValid, err := validateResourceType(resourceType, resourceTypes)
	if !isValid {
		return err
	}

	if resourceType != "" {
		resourceName, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return err
		}
		if err := d.Set("result", resourceName); err != nil {
			return fmt.Errorf("error setting result: %w", err)
		}
	}
	resourceNames := make(map[string]string, len(resourceTypes))
	for _, resourceTypeName := range resourceTypes {
		var err error
		resourceNames[resourceTypeName], err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return err
		}
	}
	if err := d.Set("results", resourceNames); err != nil {
		return fmt.Errorf("error setting results: %w", err)
	}
	d.SetId(randSeq(16, nil))
	return nil
}
