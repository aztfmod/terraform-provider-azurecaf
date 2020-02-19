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
	return &schema.Resource{
		Create: resourceNamingConventionCreate,
		Read:   resourceNamingConventionRead,
		Update: resourceNamingConventionUpdate,
		Delete: resourceNamingConventionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"convention": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ConventionCafRandom,
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
			},
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					ResourceTypeAaa,
					ResourceTypeRg,
					ResourceTypeSt,
				}, false),
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
	log.Printf(regExFilter)

	var suffixRandom string = ""
	var cafPrefix string = ""

	if convention == "cafrandom" {
		suffixRandom = "-" + randSeq(int(Resources[resourceType].MaxLength))
		cafPrefix = Resources[resourceType].CafPrefix
	} else if convention == "cafclassic" {
		cafPrefix = Resources[resourceType].CafPrefix
	} else if convention == "random" {
		suffixRandom = randSeq(int(Resources[resourceType].MaxLength - 1))
		regExFilter = string(alphanumStartletter)
	}

	// Generate the temporary name based on the concatenation of the values
	tmpGeneratedName := fmt.Sprintf("%s%s%s%s", prefix, cafPrefix, name, suffixRandom)

	myRegex, _ := regexp.Compile(regExFilter)
	// Remove the characters that are not supported in the name based on the regular expression
	filteredTmpGeneratedName := myRegex.ReplaceAllString(tmpGeneratedName, "")

	var maxLength int = 0
	maxLength = len(filteredTmpGeneratedName)

	if maxLength > int(Resources[resourceType].MaxLength) {
		maxLength = int(Resources[resourceType].MaxLength)
	}

	d.Set("result", filteredTmpGeneratedName[0:maxLength])

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
