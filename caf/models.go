package caf

// Convention defines the naming convention to apply. Default is cafclassic. Other options include cafrandom, random and passthrough
type Convention string

const (
	// ConventionCafClassic applies the CAF recommended naming convention
	ConventionCafClassic Convention = "cafclassic"
	// ConventionCafRandom defines the CAF random naming convention
	ConventionCafRandom Convention = "cafrandom"
	// ConventionRandom applies a random naming convention based on the max length of the resource
	ConventionRandom Convention = "random"
	// ConventionPassThrough defines the CAF random naming convention
	ConventionPassThrough Convention = "passthrough"
)

// ResourceType defines the name of the resource
type ResourceType string

const (
	// ResourceTypeRg defines the resource group
	ResourceTypeRg ResourceType = "rg"
	// ResourceTypeSt defines the storage account
	ResourceTypeSt ResourceType = "st"
)

// ResourcePrefix represent the CAF recommended prefix for a resource
type ResourcePrefix string

const (
	// ResourcePrefixRg defines the resource group
	ResourcePrefixRg ResourcePrefix = "rg-"
	// ResourcePrefixSt defines the storage account
	ResourcePrefixSt ResourcePrefix = "st"
)

// ResourceMaxLength represent maximum size for an azure resource
type ResourceMaxLength int

const (
	// ResourceMaxLengthRg defines the max length the resource group
	ResourceMaxLengthRg ResourceMaxLength = 80
	// ResourceMaxLengthSt defines the max length the storage account
	ResourceMaxLengthSt ResourceMaxLength = 24
)
