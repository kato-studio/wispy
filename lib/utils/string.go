package utils

import "regexp"

// remove empty strings, line breaks, and extra spaces
func CleanString(str string) string {
	// Regular expression to match multiple whitespace characters
	whitespaceRegex := regexp.MustCompile(`\s+`)
	lineBreakRegex := regexp.MustCompile(`(\r\n|\r|\n)`)

	// Replace all occurrences of the regex with a single space
	str = whitespaceRegex.ReplaceAllString(str, " ")
	str = lineBreakRegex.ReplaceAllString(str, "")
	return str

}
