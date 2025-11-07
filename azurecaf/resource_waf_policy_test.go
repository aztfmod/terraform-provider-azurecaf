package azurecaf

import (
	"regexp"
	"strings"
	"testing"
)


// TestWebApplicationFirewallPolicy_Hyphens validates that the WAF policy resource
// correctly allows hyphens in names (issue fix)
func TestWebApplicationFirewallPolicy_Hyphens(t *testing.T) {
resource := ResourceDefinitions["azurerm_web_application_firewall_policy"]

// Test that dashes are enabled
if !resource.Dashes {
t.Errorf("Expected dashes to be enabled for azurerm_web_application_firewall_policy, got disabled")
}

// Test that the slug is "waf" (CAF best practice)
if resource.CafPrefix != "waf" {
t.Errorf("Expected slug to be 'waf', got '%s'", resource.CafPrefix)
}

// Compile the validation regex
validationRegex, err := regexp.Compile(resource.ValidationRegExp)
if err != nil {
t.Fatalf("Failed to compile validation regex: %v", err)
}

// Compile the cleaning regex
cleaningRegex, err := regexp.Compile(resource.RegEx)
if err != nil {
t.Fatalf("Failed to compile cleaning regex: %v", err)
}

// Test cases for valid names with hyphens, periods, and underscores
validTestCases := []string{
"waf-policy-name",         // hyphens should be allowed
"waf.policy.name",         // periods should be allowed
"waf_policy_name",         // underscores should be allowed
"waf-policy_name.test",    // combination should be allowed
"a-b",                     // minimum valid length with hyphen
"MyWAF-Policy_2024.v1",   // mixed case with special chars
}

for _, testCase := range validTestCases {
// Test validation regex
if !validationRegex.MatchString(testCase) {
t.Errorf("Valid name '%s' should match validation regex but didn't", testCase)
}

// Test that cleaning regex doesn't remove valid characters
cleaned := cleaningRegex.ReplaceAllString(testCase, "")
if cleaned != testCase {
t.Errorf("Cleaning regex incorrectly modified '%s' to '%s'", testCase, cleaned)
}
}

// Test cases for invalid names that should be rejected
invalidTestCases := []string{
"-waf-policy",     // starts with hyphen
"waf-policy-",     // ends with hyphen
"waf policy",      // contains space
"waf@policy",      // contains @
"waf#policy",      // contains #
".waf",            // starts with period
"waf.",            // ends with period
}

for _, testCase := range invalidTestCases {
if validationRegex.MatchString(testCase) {
t.Errorf("Invalid name '%s' should not match validation regex but did", testCase)
}
}

// Test that the cleaning regex removes invalid characters but keeps hyphens, periods, underscores
testWithInvalidChars := "waf-policy@name#test_2024.v1"
cleaned := cleaningRegex.ReplaceAllString(testWithInvalidChars, "")
expected := "waf-policynametest_2024.v1" // @ and # are removed, but hyphens, periods, underscores remain
if cleaned != expected {
t.Errorf("Cleaning regex produced '%s', expected '%s'", cleaned, expected)
}
}

// TestWebApplicationFirewallPolicy_LengthConstraints validates length constraints
func TestWebApplicationFirewallPolicy_LengthConstraints(t *testing.T) {
resource := ResourceDefinitions["azurerm_web_application_firewall_policy"]

// Verify min and max length
if resource.MinLength != 1 {
t.Errorf("Expected min length to be 1, got %d", resource.MinLength)
}

if resource.MaxLength != 80 {
t.Errorf("Expected max length to be 80, got %d", resource.MaxLength)
}

validationRegex, err := regexp.Compile(resource.ValidationRegExp)
if err != nil {
t.Fatalf("Failed to compile validation regex: %v", err)
}

// Test minimum valid length (2 chars due to start/end requirements)
minValid := "ab"
if !validationRegex.MatchString(minValid) {
t.Errorf("Minimum valid name '%s' should match validation regex", minValid)
}

// Test maximum valid length (80 chars)
// Pattern: alphanumeric start, 78 chars of allowed chars, alphanumeric or underscore end
maxValid := "a" + strings.Repeat("b", 78) + "1"
if !validationRegex.MatchString(maxValid) {
t.Errorf("Maximum valid name (80 chars) should match validation regex")
}

// Test exceeds maximum length (81 chars)
tooLong := maxValid + "x"
if validationRegex.MatchString(tooLong) {
t.Errorf("Name exceeding 80 chars should not match validation regex")
}
}

// TestWebApplicationFirewallPolicy_Scope validates the scope is global
func TestWebApplicationFirewallPolicy_Scope(t *testing.T) {
resource := ResourceDefinitions["azurerm_web_application_firewall_policy"]

if resource.Scope != "global" {
t.Errorf("Expected scope to be 'global', got '%s'", resource.Scope)
}
}
