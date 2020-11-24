package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
)

// The idea of this package it is to check for package completness
// To update the list of existing resources I did query
// https://registry.terraform.io/v2/provider-versions/7185?include=provider-docs
// them use the jq espression `"azurerm_\(.included[].attributes.title)"`
// followed by manual cleaning of the non resources doc links

// ResourceStructure resource definition structure
// Copied from gen.go
type ResourceStructure struct {
	// Resource type name
	ResourceTypeName string `json:"name"`
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string `json:"slug,omitempty"`
	// MaxLength attribute define the maximum length of the name
	MinLength int `json:"min_length"`
	// MaxLength attribute define the maximum length of the name
	MaxLength int `json:"max_length"`
	// enforce lowercase
	LowerCase bool `json:"lowercase,omitempty"`
	// Regular expression to apply to the resource type
	RegEx string `json:"regex,omitempty"`
	// the Regular expression to validate the generated string
	ValidationRegExp string `json:"validation_regex,omitempty"`
	// can the resource include dashes
	Dashes bool `json:"dashes"`
	// The scope of this name where it needs to be unique
	Scope string `json:"scope,omitempty"`
}

func main() {
	wd, err := os.Getwd()
	sourceDefinitions, err := ioutil.ReadFile(path.Join(wd, "../resourceDefinition.json"))
	if err != nil {
		log.Fatal(err)
	}
	s, err := readLines(path.Join(wd, "/existing_tf_resources.txt"))
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(s)
	var data []ResourceStructure
	err = json.Unmarshal(sourceDefinitions, &data)
	if err != nil {
		log.Fatal(err)
	}
	implemented := make(map[string]bool)
	for _, name := range s {
		_, found := findByName(data, name)
		implemented[name] = found
	}
	fmt.Println("|resource | status |")
	fmt.Println("|---|---|")
	current := ""
	for _, name := range s {
		if name == current {
			continue
		} else {
			current = name
		}
		status := "❌"
		if implemented[name] {
			status = "✔"
		}
		fmt.Printf("|%s | %s |\n", name, status)
	}

}

func findByName(slice []ResourceStructure, name string) (int, bool) {
	for i, item := range slice {
		if item.ResourceTypeName == name {
			return i, true
		}
	}
	return -1, false
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
