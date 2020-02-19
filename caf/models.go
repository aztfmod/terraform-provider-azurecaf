package caf

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
	// ResourceTypeAaa defines the azure automation account
	ResourceTypeAaa string = "aaa"
	// ResourceTypeAcr defines the azure container registry
	ResourceTypeAcr string = "acr"
	// ResourceTypeRg defines the resource group
	ResourceTypeRg string = "rg"
	// ResourceTypeSt defines the storage account
	ResourceTypeSt string = "st"
)

const (
	alphanum            string = "[^0-9A-Za-z]"
	alphanumh           string = "/[^0-9A-Za-z,-]/"
	alphanumu           string = "/[^0-9A-Za-z,_]/"
	alphanumhu          string = "/[^0-9A-Za-z,_,-]/"
	alphanumhup         string = "/[^0-9A-Za-z,_,.,-]/"
	alphanumStartletter string = "/\\A[^a-z][^0-9A-Za-z]/"
)

// ResourceStructure stores the CafPrefix and the MaxLength of an azure resource
type ResourceStructure struct {
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string
	// MaxLength attribute define the maximum length of the name
	MaxLength int
	// Regular expression to apply to the resource type
	RegEx string
}

// Resources currently supported
var Resources = map[string]ResourceStructure{
	"aaa": {"aaa-", 50, alphanumh},
	"acr": {"acr-", 49, alphanum},
	"rg":  {"rg-", 80, alphanumhup},
	"st":  {"st-", 24, alphanum},
}
