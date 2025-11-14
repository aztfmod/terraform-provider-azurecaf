package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
)

// The idea of this package it is to check for package completness.
// To update the list of existing resources I did query
// https://registry.terraform.io/v2/provider-versions/7185?include=provider-docs
// them use the jq expression `"azurerm_\(.included[].attributes.title)"`.
// followed by manual cleaning of the non resources doc links.

// ResourceStructure is copied from gen.go.
type ResourceStructure struct {
	ResourceTypeName string `json:"name"`
	CafPrefix        string `json:"slug,omitempty"`
	RegEx            string `json:"regex,omitempty"`
	ValidationRegExp string `json:"validation_regex,omitempty"`
	Scope            string `json:"scope,omitempty"`
	MinLength        int    `json:"min_length"`
	MaxLength        int    `json:"max_length"`
	LowerCase        bool   `json:"lowercase,omitempty"`
	Dashes           bool   `json:"dashes"`
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	sourceDefinitions, err := os.ReadFile(path.Join(wd, "../resourceDefinition.json"))
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
		}
		current = name
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
