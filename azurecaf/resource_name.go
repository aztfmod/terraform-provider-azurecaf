package azurecaf

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceName() *schema.Resource {
	//resourceMapsKeys := make([]string, 0, len(Resources_generated))
	// for k := range Resources {
	// 	resourceMapsKeys = append(resourceMapsKeys, k)
	// }

	return &schema.Resource{
		Create:        resourceNameCreate,
		Read:          schema.Noop,
		Delete:        schema.RemoveFromState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: sliceContainsEmptyString(),
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: sliceContainsEmptyString(),
			},
			"max_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
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
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew: true,
			},
		},
	}
}

// sliceContainsEmptyString check if the slice contains an empty string
func sliceContainsEmptyString() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(v) == 0 {
			es = append(es, fmt.Errorf("emtpy values are not allowed in %s", k))
			return
		}
		return
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

func getNameResult(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	prefixes := d.Get("prefixes").([]string)
	suffixes := d.Get("suffixes").([]string)
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	convention := d.Get("convention").(string)
	desiredMaxLength := d.Get("max_length").(int)

	// Load the regular expression based on the resource type
	var regExFilter string
	var resource ResourceStructure
	var resourceFound bool = false
	if resource, resourceFound = Resources[resourceType]; !resourceFound {
		resource, resourceFound = ResourcesMapping[resourceType]
	}
	if !resourceFound {
		return fmt.Errorf("Invalid resource type %s", resourceType)
	}

	regExFilter = string(resource.RegEx)
	validationRegExPattern := string(resource.ValidationRegExp)
	log.Printf(regExFilter)

	var cafPrefix string
	var randomSuffix string = randSeq(int(resource.MaxLength))

	// configuring the prefix, cafprefix, name, postfix depending on the naming convention
	switch convention {
	case ConventionCafRandom, ConventionCafClassic:
		cafPrefix = resource.CafPrefix
	case ConventionRandom:
		//clear all the field to generate a random
		name = ""
		suffixes = []string{}
	}

	// joning the elements performing first filter to remove non compatible characters based on the resource type
	myRegex, _ := regexp.Compile(regExFilter)
	validationRegEx, _ := regexp.Compile(validationRegExPattern)
	// clear the name first based on the regexp filter of the resource type
	nameList := prefixes
	nameList = append(nameList, []string{cafPrefix, name}...)
	nameList = append(nameList, suffixes...)

	userInputName := strings.Join(nameList, separator)
	userInputName = myRegex.ReplaceAllString(userInputName, "")
	randomSuffix = myRegex.ReplaceAllString(randomSuffix, "")
	// Generate the temporary name based on the concatenation of the values - default case is caf classic
	generatedName := userInputName

	//calculate the max length
	var maxLength int = int(resource.MaxLength)
	if desiredMaxLength > 0 && desiredMaxLength < maxLength {
		maxLength = desiredMaxLength
	}

	//does the generated string contains random chars?
	var containsRandomChar = false
	switch convention {
	case ConventionPassThrough:
		// the naming is already configured
	case ConventionCafClassic:
		// the naming is already configured
	default:
		if len(userInputName) != 0 {
			if len(userInputName) < (maxLength - 1) { // prevent adding a suffix separator as the last character
				containsRandomChar = true
				generatedName = strings.Join([]string{userInputName, randomSuffix}, suffixSeparator)
			} else {
				generatedName = userInputName
			}
		} else {
			containsRandomChar = true
			generatedName = randomSuffix
		}
	}

	// Remove the characters that are not supported in the name based on the regular expression
	filteredGeneratedName := myRegex.ReplaceAllString(generatedName, "")

	var length int = len(filteredGeneratedName)

	if length > maxLength {
		length = maxLength
	}

	result := string(filteredGeneratedName[0:length])
	// making sure the last char is alpha char if we included random string
	if containsRandomChar && len(result) > len(userInputName) {
		randomLastChar := alphagenerator[rand.Intn(len(alphagenerator)-1)]
		resultRune := []rune(result)
		resultRune[len(resultRune)-1] = randomLastChar
		result = string(resultRune)
	}

	if resource.LowerCase {
		result = strings.ToLower(result)
	}

	if !validationRegEx.MatchString(result) {
		return fmt.Errorf("Invalid name for Random CAF naming %s %s Id:%s , the pattern %s doesn't match %s", resource.ResourceTypeName, name, d.Id(), validationRegExPattern, result)
	}

	d.Set("result", result)
	// Set the attribute Id with the value
	//d.SetId("none")
	d.SetId(randSeq(16))
	return nil
}
