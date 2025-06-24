package azurecaf

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Constants for error messages to avoid duplication
const (
	ErrInvalidImportIDFormat = "invalid import ID format"
	ErrResourceTypeRequired  = "resource_type is required for import"
	ErrInvalidResourceType   = "invalid resource type"
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

func resourceNameStateUpgradeV2(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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
		Importer: &schema.ResourceImporter{
			StateContext: resourceNameImport,
		},
		SchemaVersion: 3,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNameV2().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceNameStateUpgradeV2,
				Version: 2,
			},
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

// resourceNameImport handles the import of azurecaf_name resources
func resourceNameImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if importID == "" {
		return nil, fmt.Errorf("%s: empty import ID", ErrInvalidImportIDFormat)
	}

	// Parse import ID format: resourceType:name:separator:clean_input:passthrough:use_slug:random_length
	parts := strings.Split(importID, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("%s: expected format 'resourceType:name[:separator:clean_input:passthrough:use_slug:random_length]'", ErrInvalidImportIDFormat)
	}

	resourceType := parts[0]
	name := parts[1]

	// Validate resource type
	if _, exists := ResourceDefinitions[resourceType]; !exists {
		return nil, fmt.Errorf("%s: %s", ErrInvalidResourceType, resourceType)
	}

	// Set basic required fields
	d.Set("resource_type", resourceType)
	d.Set("name", name)

	// Set optional fields with defaults
	separator := "-"
	cleanInput := true
	passthrough := false
	useSlug := true
	randomLength := 0

	// Parse optional fields if provided
	if len(parts) > 2 && parts[2] != "" {
		separator = parts[2]
	}
	if len(parts) > 3 && parts[3] != "" {
		if val, err := strconv.ParseBool(parts[3]); err == nil {
			cleanInput = val
		}
	}
	if len(parts) > 4 && parts[4] != "" {
		if val, err := strconv.ParseBool(parts[4]); err == nil {
			passthrough = val
		}
	}
	if len(parts) > 5 && parts[5] != "" {
		if val, err := strconv.ParseBool(parts[5]); err == nil {
			useSlug = val
		}
	}
	if len(parts) > 6 && parts[6] != "" {
		if val, err := strconv.Atoi(parts[6]); err == nil && val >= 0 {
			randomLength = val
		}
	}

	// Set the parsed values
	d.Set("separator", separator)
	d.Set("clean_input", cleanInput)
	d.Set("passthrough", passthrough)
	d.Set("use_slug", useSlug)
	d.Set("random_length", randomLength)

	// Generate a unique ID for the imported resource
	d.SetId(randSeq(16, nil))

	return []*schema.ResourceData{d}, nil
}

func resourceNameCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceNameRead(d, meta)
}

func resourceNameRead(d *schema.ResourceData, meta interface{}) error {
	return getNameResult(d, meta)
}

func resourceNameDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
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
			if len(item) > 0 {
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

// Retrieve the resource slug / shortname based on the resourceType and the selected convention
func getSlug(resourceType string, convention string) string {
	if convention == ConventionCafClassic || convention == ConventionCafRandom {
		if val, ok := ResourceDefinitions[resourceType]; ok {
			return val.CafPrefix
		}
	}
	return ""
}

func trimResourceName(resourceName string, maxLength int) string {
	var length int = len(resourceName)

	if length > maxLength {
		length = maxLength
	}

	return string(resourceName[0:length])
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
	namePrecedence []string) string {
	composer := nameComposer{
		separator:     separator,
		maxlength:     maxlength,
		contents:      []string{},
		currentlength: 0,
	}

	for i := 0; i < len(namePrecedence); i++ {
		switch namePrecedence[i] {
		case "name":
			composer.addComponent(name)
		case "slug":
			composer.prependComponent(slug)
		case "random":
			composer.addComponent(randomSuffix)
		case "suffixes":
			if len(suffixes) > 0 {
				composer.addComponent(suffixes[0])
				suffixes = suffixes[1:]
				if len(suffixes) > 0 {
					i--
				}
			}
		case "prefixes":
			if len(prefixes) > 0 {
				composer.prependComponent(prefixes[len(prefixes)-1])
				prefixes = prefixes[:len(prefixes)-1]
				if len(prefixes) > 0 {
					i--
				}
			}
		}
	}
	return strings.Join(composer.contents, separator)
}

// nameComposer helps build names while respecting length constraints
type nameComposer struct {
	separator     string
	maxlength     int
	contents      []string
	currentlength int
}

// addComponent adds a component to the end if it fits within length constraints
func (nc *nameComposer) addComponent(component string) {
	if len(component) > 0 && nc.canAddComponent(component) {
		nc.contents = append(nc.contents, component)
		nc.updateLength(component)
	}
}

// prependComponent adds a component to the beginning if it fits within length constraints
func (nc *nameComposer) prependComponent(component string) {
	if len(component) > 0 && nc.canAddComponent(component) {
		nc.contents = append([]string{component}, nc.contents...)
		nc.updateLength(component)
	}
}

// canAddComponent checks if a component can be added without exceeding max length
func (nc *nameComposer) canAddComponent(component string) bool {
	separatorLength := 0
	if len(nc.contents) > 0 {
		separatorLength = len(nc.separator)
	}
	return nc.currentlength+len(component)+separatorLength <= nc.maxlength
}

// updateLength updates the current length after adding a component
func (nc *nameComposer) updateLength(component string) {
	separatorLength := 0
	if len(nc.contents) > 1 { // Only count separator if we have more than one component
		separatorLength = len(nc.separator)
	}
	nc.currentlength += len(component) + separatorLength
}

func validateResourceType(resourceType string, resourceTypes []string) (bool, error) {
	isEmpty := len(resourceType) == 0 && len(resourceTypes) == 0
	if isEmpty {
		return false, fmt.Errorf("resource_type and resource_types parameters are empty, you must specify at least one resource type")
	}
	errorStrings := []string{}
	resourceList := resourceTypes
	if len(resourceType) > 0 {
		resourceList = append(resourceList, resourceType)
	}

	for _, resource := range resourceList {
		_, err := getResource(resource)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	if len(errorStrings) > 0 {
		return false, fmt.Errorf(strings.Join(errorStrings, "\n"))
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
	namePrecedence []string) (string, error) {

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

func getNameResult(d *schema.ResourceData, meta interface{}) error {
	config := extractConfigFromResourceData(d)
	
	if err := validateNameConfig(config); err != nil {
		return err
	}

	convention := ConventionCafClassic
	randomSuffix := randSeq(config.randomLength, &config.randomSeed)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	isValid, err := validateResourceType(config.resourceType, config.resourceTypes)
	if !isValid {
		return err
	}

	// Handle single resource type
	if len(config.resourceType) > 0 {
		resourceName, err := getResourceName(config.resourceType, config.separator, config.prefixes, config.name, config.suffixes, randomSuffix, convention, config.cleanInput, config.passthrough, config.useSlug, namePrecedence)
		if err != nil {
			return err
		}
		d.Set("result", resourceName)
	}

	// Handle multiple resource types
	resourceNames := make(map[string]string, len(config.resourceTypes))
	for _, resourceTypeName := range config.resourceTypes {
		var err error
		resourceNames[resourceTypeName], err = getResourceName(resourceTypeName, config.separator, config.prefixes, config.name, config.suffixes, randomSuffix, convention, config.cleanInput, config.passthrough, config.useSlug, namePrecedence)
		if err != nil {
			return err
		}
	}
	d.Set("results", resourceNames)
	d.SetId(randSeq(16, nil))
	return nil
}

// nameConfig holds all configuration parameters for name generation
type nameConfig struct {
	name          string
	prefixes      []string
	suffixes      []string
	separator     string
	resourceType  string
	resourceTypes []string
	cleanInput    bool
	passthrough   bool
	useSlug       bool
	randomLength  int
	randomSeed    int64
}

// extractConfigFromResourceData extracts configuration from Terraform resource data
func extractConfigFromResourceData(d *schema.ResourceData) nameConfig {
	return nameConfig{
		name:          d.Get("name").(string),
		prefixes:      convertInterfaceToString(d.Get("prefixes").([]interface{})),
		suffixes:      convertInterfaceToString(d.Get("suffixes").([]interface{})),
		separator:     d.Get("separator").(string),
		resourceType:  d.Get("resource_type").(string),
		resourceTypes: convertInterfaceToString(d.Get("resource_types").([]interface{})),
		cleanInput:    d.Get("clean_input").(bool),
		passthrough:   d.Get("passthrough").(bool),
		useSlug:       d.Get("use_slug").(bool),
		randomLength:  d.Get("random_length").(int),
		randomSeed:    int64(d.Get("random_seed").(int)),
	}
}

// validateNameConfig validates the configuration parameters
func validateNameConfig(config nameConfig) error {
	// Validate random_length parameter
	if config.randomLength < 0 {
		return fmt.Errorf("random_length must be non-negative, got: %d", config.randomLength)
	}

	// Validate against resource type constraints if resource_type is specified
	if config.resourceType != "" {
		if resource, exists := ResourceDefinitions[config.resourceType]; exists {
			maxLen := resource.MaxLength
			if config.randomLength > maxLen {
				return fmt.Errorf("random_length (%d) exceeds maximum length for resource type %s (%d)", config.randomLength, config.resourceType, maxLen)
			}
		}
	}

	return nil
}
