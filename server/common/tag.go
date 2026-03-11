package common

import (
	"regexp"
	"strings"
)

// GenerateTagAbbreviation generates a machine-friendly abbreviation from a tag name
// by replacing spaces with dashes, removing non-alphanumerics (except dashes),
// and converting to lowercase.
func GenerateTagAbbreviation(name string) string {
	// Convert to lowercase
	abbr := strings.ToLower(name)

	// Replace spaces with dashes
	abbr = strings.ReplaceAll(abbr, " ", "-")

	// Remove all non-alphanumeric characters except dashes
	reg := regexp.MustCompile("[^a-z0-9-]+")
	abbr = reg.ReplaceAllString(abbr, "")

	// Remove consecutive dashes
	reg = regexp.MustCompile("-+")
	abbr = reg.ReplaceAllString(abbr, "-")

	// Trim leading/trailing dashes
	abbr = strings.Trim(abbr, "-")

	return abbr
}
