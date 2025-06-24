// Package azurecaf provides data structures and constants for Azure Cloud Adoption Framework
// naming conventions and resource definitions.
//
// This file contains the core constants and enumerations used throughout the provider
// for defining naming conventions, resource types, and validation rules.
package azurecaf

import (
	"math/rand"
	"time"
)

// Naming convention constants define the different methodologies supported by the provider
// for generating Azure resource names that comply with CAF guidelines.
const (
	// ConventionCafClassic applies the CAF recommended naming convention
	// Format: [prefix]-[resource-slug]-[name]-[suffix]
	ConventionCafClassic string = "cafclassic"

	// ConventionCafRandom defines the CAF random naming convention
	// Fills remaining space with random characters up to maximum length
	ConventionCafRandom string = "cafrandom"

	// ConventionRandom applies a random naming convention based on the max length of the resource
	// Generates completely random names within Azure resource constraints
	ConventionRandom string = "random"

	// ConventionPassThrough validates existing names without modification
	// Used for checking compliance of pre-existing resource names
	ConventionPassThrough string = "passthrough"
)

// Regular expression patterns for cleaning input strings to ensure compliance
// with Azure naming requirements. These patterns define which characters to remove
// from input names for different Azure resource types.
const (
	// alphanum removes all characters except alphanumeric (a-z, A-Z, 0-9)
	alphanum string = "[^0-9A-Za-z]"

	// alphanumh allows alphanumeric characters and hyphens
	alphanumh string = "[^0-9A-Za-z-]"

	// alphanumu allows alphanumeric characters and underscores
	alphanumu string = "[^0-9A-Za-z_]"

	// alphanumhu allows alphanumeric characters, hyphens, and underscores
	alphanumhu string = "[^0-9A-Za-z_-]"

	// alphanumhup allows alphanumeric characters, hyphens, underscores, and periods
	alphanumhup string = "[^0-9A-Za-z_.-]"

	// unicode allows extended unicode characters, word characters, and common punctuation
	unicode string = `[^-\w\._\(\)]`

	// invappi defines invalid characters for Application Insights resources
	invappi string = "[%&\\?/]"

	// invsqldb defines invalid characters for SQL Database resources
	invsqldb string = "[<>*%&:\\/?]"

	//Need to find a way to filter beginning and end of string
	//alphanumstartletter string = "\\A[^a-z][^0-9A-Za-z]"
)

const (
	suffixSeparator string = "-"
)

// ResourceStructure stores the CafPrefix and the MaxLength of an azure resource
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
	ValidationRegExp string `json:"validatation_regex,omitempty"`
	// can the resource include dashes
	Dashes bool `json:"dashes"`
	// The scope of this name where it needs to be unique
	Scope string `json:"scope,omitempty"`
}

var (
	alphagenerator = []rune("abcdefghijklmnopqrstuvwxyz")
)

// Generate a random value to add to the resource names
func randSeq(length int, seed *int64) string {
	// Handle invalid input: negative or zero length
	if length <= 0 {
		return ""
	}
	// initialize random seed
	if seed == nil || *seed == 0 {
		value := time.Now().UnixNano()
		seed = &value
	}
	rand.Seed(*seed)
	// generate at least one random character
	b := make([]rune, length)
	for i := range b {
		// We need the random generated string to start with a letter
		b[i] = alphagenerator[rand.Intn(len(alphagenerator)-1)]
	}
	return string(b)
}

// Resources currently supported
var Resources = map[string]ResourceStructure{
	"aaa":    {"azure automation account", "aaa", 6, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z-]{5,49}$", true, "resourceGroup"},
	"ac":     {"azure container app", "ac", 1, 32, true, alphanumh, "^[a-z0-9][a-z0-9-]{0,30}[a-z0-9]$", true, "resourceGroup"},
	"ace":    {"azure container app environment", "ace", 1, 60, false, alphanumh, "^[0-9A-Za-z][0-9A-Za-z-]{0,58}[0-9a-zA-Z]$", true, "resourceGroup"},
	"acr":    {"azure container registry", "acr", 5, 50, true, alphanum, "^[0-9A-Za-z]{5,50}$", true, "resourceGroup"},
	"afw":    {"azure firewall", "afw", 1, 80, false, alphanumhup, "^[a-zA-Z][0-9A-Za-z_.-]{0,79}$", true, "resourceGroup"},
	"agw":    {"application gateway", "agw", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"aks":    {"azure kubernetes service", "aks", 1, 63, false, alphanumhu, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,61}[0-9a-zA-Z]$", true, "resourceGroup"},
	"aksdns": {"aksdns prefix", "aksdns", 3, 45, false, alphanumh, "^[a-zA-Z][0-9A-Za-z-]{0,43}[0-9a-zA-Z]$", true, "resourceGroup"},
	"aksnpl": {"aks node pool for Linux", "aksnpl", 2, 12, true, alphanum, "^[a-zA-Z][0-9a-z]{0,11}$", true, "resourceGroup"},
	"aksnpw": {"aks node pool for Windows", "aksnpw", 2, 6, true, alphanum, "^[a-zA-Z][0-9a-z]{0,5}$", true, "resourceGroup"},
	"apim":   {"api management", "apim", 1, 50, false, alphanum, "^[a-zA-Z][0-9A-Za-z]{0,49}$", true, "resourceGroup"},
	"app":    {"web app", "app", 2, 60, false, alphanumh, "^[0-9A-Za-z][0-9A-Za-z-]{0,58}[0-9a-zA-Z]$", true, "resourceGroup"},
	"appi":   {"application insights", "appi", 1, 260, false, invappi, "^[^%&\\?/. ][^%&\\?/]{0,258}[^%&\\?/. ]$", true, "resourceGroup"},
	"ase":    {"app service environment", "ase", 2, 36, false, alphanumh, "^[0-9A-Za-z-]{2,36}$", true, "resourceGroup"},
	"asr":    {"azure site recovery", "asr", 2, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z-]{1,49}$", true, "resourceGroup"},
	"dcr":    {"data collection rule", "dcr", 3, 44, false, alphanumhup, "^[a-zA-Z0-9][a-zA-Z0-9-]{1,42}[a-zA-Z0-9]$", true, "resourceGroup"},
	"evh":    {"event hub", "evh", 1, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z-]{0,48}[0-9a-zA-Z]$", true, "resourceGroup"},
	"gen":    {"generic", "gen", 1, 24, false, alphanum, "^[0-9a-zA-Z]{1,24}$", true, "resourceGroup"},
	"kv":     {"keyvault", "kv", 3, 24, true, alphanumh, "^[a-zA-Z][0-9A-Za-z-]{0,22}[0-9a-zA-Z]$", true, "resourceGroup"},
	"la":     {"loganalytics", "la", 4, 63, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z-]{3,61}[0-9a-zA-Z]$", true, "resourceGroup"},
	"las":    {"log analytics solution", "las", 4, 63, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z-]{2,61}[0-9a-zA-Z]$", true, "parent"},
	"laqp":   {"log analytics query pack", "laqp", 4, 63, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z-]{2,61}[0-9a-zA-Z]$", true, "parent"},
	"nic":    {"network interface card", "nic", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"nsg":    {"network security group", "nsg", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"pip":    {"public ip address", "pip", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"plan":   {"app service plan", "plan", 1, 40, false, alphanumh, "^[0-9A-Za-z-]{1,40}$", true, "resourceGroup"},
	"rg":     {"resource group", "rg", 1, 80, false, unicode, `^[-\w\._\(\)]{1,80}$`, true, "resourceGroup"},
	"snet":   {"virtual network subnet", "snet", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"sql":    {"azure sql db server", "sql", 1, 63, true, alphanumh, "^[0-9a-z][0-9a-z-]{0,61}[0-9a-z]$", true, "resourceGroup"},
	"sqldb":  {"azure sql db", "sqldb", 1, 128, false, invsqldb, "^[^<>*%&:\\/?. ][^<>*%&:\\/?]{0,126}[^<>*%&:\\/?. ]$", true, "resourceGroup"},
	"st":     {"storage account", "st", 3, 24, true, alphanum, "^[0-9a-z]{3,24}$", true, "resourceGroup"},
	"vml":    {"virtual machine (linux)", "vml", 1, 64, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z_-]{0,62}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"vmw":    {"virtual machine (windows)", "vmw", 1, 15, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z_-]{0,13}[0-9a-zA-Z_]$", true, "resourceGroup"},
	"vnet":   {"virtual network", "vnet", 2, 64, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z_.-]{0,62}[0-9a-zA-Z_]$", true, "resourceGroup"},
}

// ResourcesMapping enforcing new naming convention
var ResourcesMapping = map[string]ResourceStructure{
	"azurerm_automation_account":              Resources["aaa"],
	"azurerm_container_app":                   Resources["ac"],
	"azurerm_container_app_environment":       Resources["ace"],
	"azurerm_container_registry":              Resources["acr"],
	"azurerm_firewall":                        Resources["afw"],
	"azurerm_application_gateway":             Resources["agw"],
	"azurerm_api_management":                  Resources["apim"],
	"azurerm_app_service":                     Resources["app"],
	"azurerm_application_insights":            Resources["appi"],
	"azurerm_app_service_environment":         Resources["ase"],
	"azurerm_recovery_services_vault":         Resources["asr"],
	"azurerm_eventhub_namespace":              Resources["evh"],
	"generic":                                 Resources["gen"],
	"azurerm_key_vault":                       Resources["kv"],
	"azurerm_kubernetes_cluster":              Resources["aks"],
	"aks_dns_prefix":                          Resources["aksdns"],
	"aks_node_pool_linux":                     Resources["aksnpl"],
	"aks_node_pool_windows":                   Resources["aksnpw"],
	"azurerm_log_analytics_workspace":         Resources["la"],
	"azurerm_log_analytics_solution":          Resources["las"],
	"azurerm_log_analytics_query_pack":        Resources["laqp"],
	"azurerm_monitor_data_collection_rule":    Resources["dcr"],
	"azurerm_network_interface":               Resources["nic"],
	"azurerm_network_security_group":          Resources["nsg"],
	"azurerm_public_ip":                       Resources["pip"],
	"azurerm_app_service_plan":                Resources["plan"],
	"azurerm_service_plan":                    Resources["plan"],
	"azurerm_resource_group":                  Resources["rg"],
	"azurerm_subnet":                          Resources["snet"],
	"azurerm_sql_server":                      Resources["sql"],
	"azurerm_sql_database":                    Resources["sqldb"],
	"azurerm_storage_account":                 Resources["st"],
	"azurerm_windows_virtual_machine_linux":   Resources["vml"],
	"azurerm_windows_virtual_machine_windows": Resources["vmw"],
	"azurerm_virtual_network":                 Resources["vnet"],
}
