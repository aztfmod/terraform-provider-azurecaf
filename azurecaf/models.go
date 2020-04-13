package azurecaf

const (
	// ConventionCafClassic applies the CAF recommended naming convention
	ConventionCafClassic string = "cafclassic"
	// ConventionCafRandom defines the CAF random naming convention
	ConventionCafRandom string = "cafrandom"
	// ConventionRandom applies a random naming convention based on the max length of the resource
	ConventionRandom string = "random"
	// ConventionPassThrough defines the CAF random naming convention
	ConventionPassThrough string = "passthrough"
)

const (
	alphanum            string = "[^0-9A-Za-z]"
	alphanumh           string = "/[^0-9A-Za-z,-]/"
	alphanumu           string = "/[^0-9A-Za-z,_]/"
	alphanumhu          string = "/[^0-9A-Za-z,_,-]/"
	alphanumhup         string = "/[^0-9A-Za-z,_,.,-]/"
	alphanumStartletter string = "/\\A[^a-z][^0-9A-Za-z]/"
	unicode             string = `[^-\w\._\(\)]`
)

const (
	suffixSeparator string = "-"
)

// ResourceStructure stores the CafPrefix and the MaxLength of an azure resource
type ResourceStructure struct {
	// Resource type name
	ResourceTypeName string
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string
	// MaxLength attribute define the maximum length of the name
	MinLength int
	// MaxLength attribute define the maximum length of the name
	MaxLength int
	// enforce lowercase
	LowerCase bool
	// Regular expression to apply to the resource type
	RegEx string
	// the Regular expression to validate the generated string
	ValidationRegExp string
}

// Resources currently supported
var Resources = map[string]ResourceStructure{
	"aaa":  {"azure automation account", "aaa", 6, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z,-]{5,49}$"},
	"acr":  {"azure container registry", "acr", 5, 49, true, alphanum, "^[0-9A-Za-z]{5,50}$"},
	"rg":   {"resource group", "rg", 1, 80, false, unicode, `^[-\w\._\(\)]{1,80}$`},
	"st":   {"storage account", "st", 3, 24, true, alphanum, "^[0-9a-z]{3,24}$"},
	"afw":  {"azure firewall", "afw", 1, 80, false, alphanumhup, "^[a-zA-Z][0-9A-Za-z,_,.,-]{0,79}$"},
	"asr":  {"azure site recovery", "asr", 2, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z,-]{1,49}$"},
	"evh":  {"event hub", "evh", 1, 50, false, alphanumh, "^[a-zA-Z][0-9A-Za-z,-]{0,48}[0-9a-zA-Z]{0,1}$"},
	"kv":   {"keyvault", "kv", 3, 24, true, alphanumh, "^[a-zA-Z][0-9A-Za-z,-]{0,22}[0-9a-zA-Z]{0,1}$"},
	"la":   {"loganalytics", "la", 4, 63, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z,-]{3,62}[0-9a-zA-Z]{0,1}$"},
	"nic":  {"network interface card", "nic", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z,_,.,-]{0,79}[0-9a-zA-Z,_]{0,1}$"},
	"nsg":  {"network security group", "nsg", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z,_,.,-]{0,79}[0-9a-zA-Z,_]{0,1}$"},
	"pip":  {"public ip address", "pip", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z,_,.,-]{0,79}[0-9a-zA-Z,_]{0,1}$"},
	"snet": {"virtual network subnet", "snet", 1, 80, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z,_,.,-]{0,79}[0-9a-zA-Z,_]{0,1}$"},
	"vnet": {"virtual network", "vnet", 1, 64, false, alphanumhup, "^[0-9a-zA-Z][0-9A-Za-z,_,.,-]{0,62}[0-9a-zA-Z,_]{0,1}$"},
	"vmw":  {"virtual machine (windows)", "vmw", 1, 15, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z,_,-]{0,13}[0-9a-zA-Z,_]{0,1}$"},
	"vml":  {"virtual machine (linux)", "vml", 1, 64, false, alphanumh, "^[0-9a-zA-Z][0-9A-Za-z,_,-]{0,63}[0-9a-zA-Z,_]{0,1}$"},
	"gen":  {"generic", "vml-", 1, 24, false, alphanum, "^[0-9a-zA-Z]{1,24}$"},
}

// ResourcesMapping enforcing new naming convention
var ResourcesMapping = map[string]ResourceStructure{
	"azurerm_automation_account":              Resources["aaa"],
	"azurerm_container_registry":              Resources["acr"],
	"azurerm_resource_group":                  Resources["rg"],
	"azurerm_storage_account":                 Resources["st"],
	"azurerm_firewall":                        Resources["afw"],
	"azurerm_recovery_services_vault":         Resources["asr"],
	"azurerm_eventhub_namespace":              Resources["evh"],
	"azurerm_key_vault":                       Resources["kv"],
	"azurerm_log_analytics_workspace":         Resources["la"],
	"azurerm_network_interface":               Resources["nic"],
	"azurerm_network_security_group":          Resources["nsg"],
	"azurerm_public_ip":                       Resources["pip"],
	"azurerm_subnet":                          Resources["snet"],
	"azurerm_virtual_network":                 Resources["vnet"],
	"azurerm_windows_virtual_machine_windows": Resources["vmw"],
	"azurerm_windows_virtual_machine_linux":   Resources["linux"],
	"generic":                                 Resources["gen"],
}
