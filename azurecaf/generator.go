package azurecaf

import (
	"math/rand"
)

var (
	alphagenerator = []rune("abcdefghijklmnopqrstuvwxyz")
)

// Generate a random value to add to the resource names
func randSeq(length int, seed int64) string {
	if length == 0 {
		return ""
	}
	// initialize random seed
	rand.Seed(seed)
	// generate at least one random character
	b := make([]rune, length)
	for i := range b {
		// We need the random generated string to start with a letter
		b[i] = alphagenerator[rand.Intn(len(alphagenerator)-1)]
	}
	return string(b)
}
