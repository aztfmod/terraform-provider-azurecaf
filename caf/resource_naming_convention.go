package caf

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	getResult(d, m)
	return nil
}

func resourceNamingConventionUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceNamingConventionRead(d, m)
}

func resourceNamingConventionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getResult(d *schema.ResourceData, m interface{}) error {

	// time.Sleep(30 * time.Second)
	name := d.Get("name").(string)
	prefix := d.Get("prefix").(string)
	resourceType := d.Get("resource_type").(string)
	convention := d.Get("convention").(string)

	// Load the regular expression based on the resource type
	var regExFilter string = ""
	regExFilter = string(Resources[resourceType].RegEx)
	validationRegExPattern := string(Resources[resourceType].ValidationRegExp)
	log.Printf(regExFilter)

	var suffixSeparator string = ""
	var cafPrefix string = ""
	var randomSuffix string = randSeq(int(Resources[resourceType].MaxLength))
	if convention == "cafrandom" {
		suffixSeparator = "-"
		cafPrefix = Resources[resourceType].CafPrefix
	} else if convention == "cafclassic" {
		cafPrefix = Resources[resourceType].CafPrefix
	} else if convention == "random" {
		regExFilter = string(alphanumStartletter)
	}

	myRegex, _ := regexp.Compile(regExFilter)
	validationRegEx, _ := regexp.Compile(validationRegExPattern)
	// clear the name first based on the regexp filter of the resource type
	tmpName := fmt.Sprintf("%s%s%s%s", prefix, cafPrefix, name, suffixSeparator)
	tmpName = myRegex.ReplaceAllString(tmpName, "")
	// Generate the temporary name based on the concatenation of the values
	tmpGeneratedName := fmt.Sprintf("%s%s", tmpName, randomSuffix)

	// Remove the characters that are not supported in the name based on the regular expression
	filteredTmpGeneratedName := myRegex.ReplaceAllString(tmpGeneratedName, "")

	var maxLength int = 0
	maxLength = len(filteredTmpGeneratedName)

	if maxLength > int(Resources[resourceType].MaxLength) {
		maxLength = int(Resources[resourceType].MaxLength)
	}

	result := string(filteredTmpGeneratedName[0:maxLength])
	if !validationRegEx.MatchString(result) {
		return fmt.Errorf("The pattern %s doesn't match %s", validationRegExPattern, result)
	}

	d.Set("result", result)

	return nil
}

var letters = []rune("01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generate a random value to add to the resource names
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
