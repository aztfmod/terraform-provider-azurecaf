// Package azurecaf provides data structures and constants for Azure Cloud Adoption Framework
// naming conventions and resource definitions.
//
// This file contains the core constants and enumerations used throughout the provider
// for defining naming conventions, resource types, and validation rules.
package azurecaf

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"math/big"
	"math/rand"
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

// ResourceStructure stores the CafPrefix and the MaxLength of an azure resource
type ResourceStructure struct {
	// Resource type name
	ResourceTypeName string `json:"name"`
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string `json:"slug,omitempty"`
	// MinLength attribute defines the minimum length of the name
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

var (
	alphagenerator = []rune("abcdefghijklmnopqrstuvwxyz")

	// cryptoRandReader is the source of randomness for randomLetter. It is a
	// package-level variable so tests can substitute a failing reader to
	// exercise the error-return branch. Production code MUST NOT reassign it.
	cryptoRandReader io.Reader = cryptorand.Reader
)

// randomLetter returns one character from alphagenerator drawn from
// crypto/rand (cryptographically secure, non-deterministic). It is used for
// non-secret values where unpredictability is preferable to determinism —
// for example, the random ID assigned to resources at apply time.
//
// crypto/rand.Reader only fails when the OS entropy source is broken
// (unreadable /dev/urandom on Linux, failing BCryptGenRandom on Windows),
// which is fatal. We surface the error so callers can return a Terraform
// diagnostic rather than panicking inside the provider process — both
// because tfproviderlint R009 forbids panic() in providers and because a
// diagnostic is the user-visible contract Terraform expects.
func randomLetter() (rune, error) {
	v, err := cryptorand.Int(cryptoRandReader, big.NewInt(int64(len(alphagenerator))))
	if err != nil {
		return 0, fmt.Errorf("azurecaf: crypto/rand.Reader failed: %w", err)
	}
	return alphagenerator[v.Int64()], nil
}

// randSeq generates a random sequence of length characters from alphagenerator.
//
// When seed is nil, characters are drawn from crypto/rand via randomLetter.
// This path has no determinism requirement; using a CSPRNG removes the
// SonarCloud go:S2245 weak-PRNG concern at zero behavioural cost. The error
// is propagated upwards if the OS entropy source fails.
//
// When seed is non-nil, characters are drawn from a math/rand source seeded
// with *seed. This deterministic path is REQUIRED for plan-time name
// visibility (issue #336): the same configuration must yield the same name at
// plan and apply. math/rand is the only stdlib PRNG that can satisfy this, and
// the output is a non-secret Terraform resource name visible in plan and
// state — not a security-sensitive value — so the go:S2245 hotspot at the
// rand.New call below is intentional and accepted. This branch never returns
// a non-nil error.
func randSeq(length int, seed *int64) (string, error) {
	if length <= 0 {
		return "", nil
	}
	b := make([]rune, length)
	if seed == nil {
		for i := range b {
			r, err := randomLetter()
			if err != nil {
				return "", err
			}
			b[i] = r
		}
		return string(b), nil
	}
	rng := rand.New(rand.NewSource(*seed))
	for i := range b {
		b[i] = alphagenerator[rng.Intn(len(alphagenerator))]
	}
	return string(b), nil
}
