package caf

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

	return &schema.Resource{
		Create: resourceNamingConventionCreate,
		Read:   resourceNamingConventionRead,
		Delete: resourceNamingConventionDelete,

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
	desiredMaxLength := 100 //d.Get("maxLength").(int)

	// Load the regular expression based on the resource type
	var regExFilter string = ""
	regExFilter = string(Resources[resourceType].RegEx)
	validationRegExPattern := string(Resources[resourceType].ValidationRegExp)
	log.Printf(regExFilter)

	var cafPrefix string = ""
	var randomSuffix string = randSeq(int(Resources[resourceType].MaxLength))

	// configuring the prefix, cafprefix, name, postfix depending on the naming convention
	switch convention {
	case ConventionCafRandom, ConventionCafClassic:
		cafPrefix = Resources[resourceType].CafPrefix
	case ConventionRandom:
		regExFilter = string(alphanumStartletter)
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
	tmpName := strings.Join(nameList, suffixSeparator)
	tmpName = myRegex.ReplaceAllString(tmpName, "")
	randomSuffix = myRegex.ReplaceAllString(randomSuffix, "")
	// Generate the temporary name based on the concatenation of the values - default case is caf classic
	tmpGeneratedName := tmpName

	//calculate the max length
	var maxLength int = int(Resources[resourceType].MaxLength)
	if desiredMaxLength > 0 && desiredMaxLength < maxLength {
		maxLength = desiredMaxLength
	}

	switch convention {
	case ConventionPassThrough:
		tmpGeneratedName = name
	case ConventionCafClassic:
		// the naming is already configured
	default:
		if len(name) != 0 {
			tmpGeneratedName = strings.Join([]string{tmpName, randomSuffix}, suffixSeparator)
		} else {
			tmpGeneratedName = randomSuffix
		}
	}

	// Remove the characters that are not supported in the name based on the regular expression
	filteredTmpGeneratedName := myRegex.ReplaceAllString(tmpGeneratedName, "")

	var length int = len(filteredTmpGeneratedName)

	if length > maxLength {
		length = maxLength
	}

	result := string(filteredTmpGeneratedName[0:length])
	if !validationRegEx.MatchString(result) {
		return fmt.Errorf("Invalid name for Random CAF naming %s %s Id:%s , the pattern %s doesn't match %s", Resources[resourceType].ResourceTypeName, name, d.Id(), validationRegExPattern, result)
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
		b[i] = alphanumgenerator[rand.Intn(len(alphanumgenerator))]
	}
	return string(b)
}
