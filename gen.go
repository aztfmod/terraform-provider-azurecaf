// Code Generation Tool for Azure CAF Provider
//
// This tool generates Go code from JSON resource definitions to ensure the provider
// stays current with Azure resource naming requirements and Cloud Adoption Framework
// guidelines.
//
// The generator reads resourceDefinition.json and creates models_generated.go with:
//   - Resource type constants and mappings
//   - Validation rules and constraints
//   - Naming convention logic
//   - Resource slug mappings
//
// Usage: go generate (automatically runs this file via go:generate directive in main.go)

//go:build ignore
// +build ignore

// This program generates models_generated.go. It can be invoked by running
// go generate from the project root.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"text/template"
	"time"
)

// ResourceStructure defines the schema for Azure resource naming requirements
// as specified in the resourceDefinition.json file.
//
// Each resource type has specific constraints that must be enforced when
// generating compliant names for Azure resources.
type ResourceStructure struct {
	// ResourceTypeName is the full Terraform resource type name (e.g., "azurerm_storage_account")
	ResourceTypeName string `json:"name"`

	// CafPrefix is the Cloud Adoption Framework abbreviation for this resource type (e.g., "st" for storage account)
	// This slug is used as a prefix in generated names to indicate resource type
	CafPrefix string `json:"slug,omitempty"`

	// MinLength defines the minimum allowed length for the resource name
	MinLength int `json:"min_length"`

	// MaxLength defines the maximum allowed length for the resource name
	MaxLength int `json:"max_length"`

	// LowerCase indicates whether the resource name must be entirely lowercase
	LowerCase bool `json:"lowercase,omitempty"`

	// RegEx is the cleaning regex pattern used to remove invalid characters from input names
	// Characters matching this pattern will be stripped from the name
	RegEx string `json:"regex,omitempty"`

	// ValidationRegExp is the validation regex that the final generated name must match
	// This ensures the generated name complies with Azure's naming requirements
	ValidationRegExp string `json:"validation_regex,omitempty"`

	// Dashes indicates whether the resource type allows dash characters in names
	Dashes bool `json:"dashes"`

	// Scope defines where the resource name must be unique (e.g., "global", "resourceGroup", "parent")
	Scope string `json:"scope,omitempty"`

	// OutOfDoc indicates whether this resource is not present in the official Azure CAF documentation
	OutOfDoc bool `json:"out_of_doc,omitempty"`

	// Resource is the official resource name from Azure CAF documentation
	Resource string `json:"resource,omitempty"`

	// ResourceProviderNamespace is the Azure resource provider namespace from official documentation
	ResourceProviderNamespace string `json:"resource_provider_namespace,omitempty"`
}

// templateData holds the data structure passed to the Go template for code generation
type templateData struct {
	GeneratedTime      time.Time           // Timestamp when the code was generated
	ResourceStructures []ResourceStructure // All resource definitions from JSON
	SlugMap            map[string]string   // Mapping of CAF prefixes to resource types
}

// main is the entry point for the code generator.
// It performs the following steps:
//  1. Reads resource definitions from resourceDefinition.json
//  2. Loads and parses Go templates from the templates/ directory
//  3. Processes the resource data to create mappings and deduplicate entries
//  4. Generates models_generated.go with all resource definitions and validation logic
func main() {
	// Get the current working directory to locate input files
	wd, err := os.Getwd()
	if err != nil {
		log.Panicln("No directory found")
	}

	fmt.Println() // Add spacing for readability

	// Load all template files from the templates directory
	files, err := ioutil.ReadDir(path.Join(wd, "templates"))
	if err != nil {
		log.Fatal(err)
	}

	// Build list of template file paths
	var fileNames = make([]string, len(files))
	for i, file := range files {
		fileNames[i] = path.Join(wd, "templates", file.Name())
	}

	// Parse all templates and register custom functions
	parsedTemplate, err := template.New("templates").Funcs(template.FuncMap{
		// Terraform does not yet support lookahead in their regex function,
		// so we need to clean regex patterns for compatibility
		"cleanRegex": func(dirtyString string) string {
			// Remove lookahead patterns that Terraform cannot handle
			var re = regexp.MustCompile(`(?m)\(\?\=.{\d+,\d+}\$\)|\(\?\!\..*--\)`)
			return re.ReplaceAllString(dirtyString, "")
		},
	}).ParseFiles(fileNames...)
	if err != nil {
		log.Fatal(err)
	}

	// Read the combined resource definitions from JSON file
	// This file now contains both documented and undocumented resources
	sourceDefinitions, err := ioutil.ReadFile(path.Join(wd, "resourceDefinition.json"))
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON resource definitions into Go structs
	var uniqueData []ResourceStructure
	err = json.Unmarshal(sourceDefinitions, &uniqueData)
	if err != nil {
		log.Fatal(err)
	}

	// Sort by resource type name for consistent output
	sort.SliceStable(uniqueData, func(i, j int) bool {
		return uniqueData[i].ResourceTypeName < uniqueData[j].ResourceTypeName
	})

	// Build a mapping of CAF prefixes (slugs) to resource types
	// This allows reverse lookup from slug to resource type name
	slugMap := make(map[string]string)
	for _, res := range uniqueData {
		if _, exists := slugMap[res.CafPrefix]; !exists {
			slugMap[res.CafPrefix] = res.ResourceTypeName
		}
	}

	// Generate the Go source file using the parsed template
	modelsFile, err := os.OpenFile(path.Join(wd, "azurecaf/models_generated.go"), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer modelsFile.Close()

	// Execute the template with our processed data
	err = parsedTemplate.ExecuteTemplate(modelsFile, "model.tmpl", templateData{
		GeneratedTime:      time.Now(),
		ResourceStructures: uniqueData,
		SlugMap:            slugMap,
	})

	if err != nil {
		log.Fatalf("execution failed: %s", err)
	}
	log.Println("File generated")
}
