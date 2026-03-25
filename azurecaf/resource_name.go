package azurecaf

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		CustomizeDiff: customdiff.All(resourceNameCustomizeDiff),
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
			"error_when_exceeding_max_length": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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

func resourceNameDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceNameImport handles importing existing resource names.
// Import ID format: <resource_type>:<existing_name>
// Example: azurerm_storage_account:mystorageaccount123
func resourceNameImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	d.Set("name", existingName)
	d.Set("resource_type", resourceType)
	d.Set("passthrough", true)

	// Set empty slices for prefixes and suffixes since we can't reverse-engineer them
	d.Set("prefixes", []string{})
	d.Set("suffixes", []string{})
	d.Set("resource_types", []string{})

	// Set the result to match the imported name
	d.Set("result", existingName)
	d.Set("results", map[string]string{})

	// Use the existing name as the Terraform resource ID
	d.SetId(existingName)

	return []*schema.ResourceData{d}, nil
}

// schemaGetter is implemented by both schema.ResourceData and schema.ResourceDiff,
// allowing shared parameter extraction logic.
type schemaGetter interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

// namingParams holds the extracted input parameters for name computation.
type namingParams struct {
	name                        string
	prefixes                    []string
	suffixes                    []string
	separator                   string
	resourceType                string
	resourceTypes               []string
	cleanInput                  bool
	passthrough                 bool
	useSlug                     bool
	randomLength                int
	randomSeed                  int64
	randomSeedSet               bool
	errorWhenExceedingMaxLength bool
}

// extractNamingParams reads naming input parameters from a schema getter.
func extractNamingParams(d schemaGetter) namingParams {
	seedVal, seedSet := d.GetOk("random_seed")
	var seed int64
	if seedSet {
		seed = int64(seedVal.(int))
	}
	return namingParams{
		name:                        d.Get("name").(string),
		prefixes:                    convertInterfaceToString(d.Get("prefixes").([]interface{})),
		suffixes:                    convertInterfaceToString(d.Get("suffixes").([]interface{})),
		separator:                   d.Get("separator").(string),
		resourceType:                d.Get("resource_type").(string),
		resourceTypes:               convertInterfaceToString(d.Get("resource_types").([]interface{})),
		cleanInput:                  d.Get("clean_input").(bool),
		passthrough:                 d.Get("passthrough").(bool),
		useSlug:                     d.Get("use_slug").(bool),
		randomLength:                d.Get("random_length").(int),
		randomSeed:                  seed,
		randomSeedSet:               seedSet,
		errorWhenExceedingMaxLength: d.Get("error_when_exceeding_max_length").(bool),
	}
}

// computeNames generates the result and results map from the given parameters.
func computeNames(p namingParams) (string, map[string]string, error) {
	// Validate random_length parameter
	if p.randomLength < 0 {
		return "", nil, fmt.Errorf("random_length must be non-negative, got: %d", p.randomLength)
	}

	// Validate against resource type constraints if resource_type is specified
	if p.resourceType != "" {
		if resource, exists := ResourceDefinitions[p.resourceType]; exists {
			if p.randomLength > resource.MaxLength {
				return "", nil, fmt.Errorf("random_length (%d) exceeds maximum length for resource type %s (%d)", p.randomLength, p.resourceType, resource.MaxLength)
			}
		}
	}

	var seedPtr *int64
	if p.randomSeedSet {
		seedPtr = &p.randomSeed
	}
	randomSuffix := randSeq(p.randomLength, seedPtr)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	convention := ConventionCafClassic

	isValid, err := validateResourceType(p.resourceType, p.resourceTypes)
	if !isValid {
		return "", nil, err
	}

	var result string
	if len(p.resourceType) > 0 {
		result, err = getResourceName(p.resourceType, p.separator, p.prefixes, p.name, p.suffixes, randomSuffix, convention, p.cleanInput, p.passthrough, p.useSlug, namePrecedence, p.errorWhenExceedingMaxLength)
		if err != nil {
			return "", nil, err
		}
	}

	resourceNames := make(map[string]string, len(p.resourceTypes))
	for _, resourceTypeName := range p.resourceTypes {
		resourceNames[resourceTypeName], err = getResourceName(resourceTypeName, p.separator, p.prefixes, p.name, p.suffixes, randomSuffix, convention, p.cleanInput, p.passthrough, p.useSlug, namePrecedence, p.errorWhenExceedingMaxLength)
		if err != nil {
			return "", nil, err
		}
	}

	return result, resourceNames, nil
}

// resourceNameCustomizeDiff computes naming values during the plan phase so that
// users can see the actual resource names in terraform plan output instead of
// "(known after apply)". This runs during plan for new or replaced resources.
func resourceNameCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// For existing resources with no relevant input changes, values are already in state.
	if d.Id() != "" {
		needsRecompute := false
		for _, attr := range []string{
			"name", "prefixes", "suffixes", "separator",
			"resource_type", "resource_types", "clean_input",
			"passthrough", "use_slug", "random_length",
			"random_seed", "error_when_exceeding_max_length",
		} {
			if d.HasChange(attr) {
				needsRecompute = true
				break
			}
		}
		if !needsRecompute {
			return nil
		}
	}

	p := extractNamingParams(d)

	// When random_length > 0 but no random_seed is explicitly provided, we cannot
	// compute deterministic names at plan time (each CustomizeDiff call would
	// generate a different seed). Still run validations so errors surface during
	// plan, but skip setting result/results (they remain "known after apply").
	if p.randomLength > 0 && !p.randomSeedSet {
		_, _, err := computeNames(p)
		return err
	}

	result, resourceNames, err := computeNames(p)
	if err != nil {
		return err
	}

	if len(p.resourceType) > 0 {
		if err := d.SetNew("result", result); err != nil {
			return fmt.Errorf("failed to set result: %w", err)
		}
	}
	if err := d.SetNew("results", resourceNames); err != nil {
		return fmt.Errorf("failed to set results: %w", err)
	}

	return nil
}

func cleanSlice(names []string, resourceDefinition *ResourceStructure) []string {
	for i, name := range names {
		names[i] = cleanString(name, resourceDefinition)
	}
	return names
}

func cleanString(name string, resourceDefinition *ResourceStructure) string {
	myRegex, err := regexp.Compile(resourceDefinition.RegEx)
	if err != nil {
		log.Printf("[WARN] invalid regex pattern %q for resource %s: %v", resourceDefinition.RegEx, resourceDefinition.ResourceTypeName, err)
		return name
	}
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
	namePrecedence []string,
	errorWhenExceedingMaxLength bool) (string, error) {
	nameBuilder := NewNameBuilder(maxlength, separator)

	for i := 0; i < len(namePrecedence); i++ {
		switch c := namePrecedence[i]; c {
		case "name":
			if len(name) > 0 {
				nameBuilder.Append(name)
			}
		case "slug":
			if len(slug) > 0 {
				nameBuilder.Prepend(slug)
			}
		case "random":
			if len(randomSuffix) > 0 {
				nameBuilder.Append(randomSuffix)
			}
		case "suffixes":
			if len(suffixes) > 0 {
				if len(suffixes[0]) > 0 {
					nameBuilder.Append(suffixes[0])
				}
				suffixes = suffixes[1:]
				if len(suffixes) > 0 {
					i--
				}
			}
		case "prefixes":
			if len(prefixes) > 0 {
				if len(prefixes[len(prefixes)-1]) > 0 {
					nameBuilder.Prepend(prefixes[len(prefixes)-1])
				}
				prefixes = prefixes[:len(prefixes)-1]
				if len(prefixes) > 0 {
					i--
				}
			}
		}
	}
	if errorWhenExceedingMaxLength {
		content := nameBuilder.GetName()
		contentLength := len(content)
		if contentLength > maxlength {
			return "", fmt.Errorf("composed name '%s' exceeds maximum length of %d by %d characters", content, maxlength, contentLength-maxlength)
		}
		return content, nil
	}
	content := nameBuilder.GetTrimmedName()
	return content, nil
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
	errorWhenExceedingMaxLength bool) (string, error) {

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
		resourceName, err = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, resource.MaxLength, namePrecedence, errorWhenExceedingMaxLength)
		if err != nil {
			return "", err
		}
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
	p := extractNamingParams(d)

	result, resourceNames, err := computeNames(p)
	if err != nil {
		return err
	}

	if len(p.resourceType) > 0 {
		if err := d.Set("result", result); err != nil {
			return fmt.Errorf("failed to set result: %w", err)
		}
	}
	if err := d.Set("results", resourceNames); err != nil {
		return fmt.Errorf("failed to set results: %w", err)
	}
	d.SetId(randSeq(16, nil))
	return nil
}
