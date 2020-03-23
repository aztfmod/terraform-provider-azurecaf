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

func resourceNamingConvention() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(Resources))
	for k := range Resources {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	for k := range ResourcesMapping {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		Create:        resourceNamingConventionCreate,
		Read:          schema.Noop,
		Delete:        schema.RemoveFromState,
		SchemaVersion: 2,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"convention": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ConventionCafRandom,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ConventionCafClassic,
					ConventionCafRandom,
					ConventionRandom,
					ConventionPassThrough,
				}, false),
			},
			"prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"postfix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"max_length": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceNamingConventionCreate(d *schema.ResourceData, m interface{}) error {

	// Set the attribute Id with the value
	d.SetId(randSeq(16))

	return resourceNamingConventionRead(d, m)
}

func resourceNamingConventionRead(d *schema.ResourceData, m interface{}) error {
	return getResult(d, m)
}

func resourceNamingConventionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getResult(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	prefix := d.Get("prefix").(string)
	postfix := d.Get("postfix").(string)
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
		postfix = ""
	}

	// joning the elements performing first filter to remove non compatible characters based on the resource type
	myRegex, _ := regexp.Compile(regExFilter)
	validationRegEx, _ := regexp.Compile(validationRegExPattern)
	// clear the name first based on the regexp filter of the resource type
	nameList := []string{}
	for _, s := range []string{prefix, cafPrefix, name, postfix} {
		if strings.TrimSpace(s) != "" {
			nameList = append(nameList, s)
		}
	}
	userInputName := strings.Join(nameList, suffixSeparator)
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
		generatedName = name
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

	return nil
}

var (
	alphanumgenerator = []rune("01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphagenerator    = []rune("abcdefghijklmnopqrstuvwxyz")
)

// Generate a random value to add to the resource names
func randSeq(n int) string {
	// generate at least one random character
	b := make([]rune, n)
	for i := range b {
		b[i] = alphanumgenerator[rand.Intn(len(alphanumgenerator)-1)]
	}
	return string(b)
}
