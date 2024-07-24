package utils

import "regexp"

// remove empty strings, line breaks, and extra spaces
func CleanString(str string) string {
	// Regular expression to match multiple whitespace characters
	whitespace_regex := regexp.MustCompile(`\s+`)
	line_breakRegex := regexp.MustCompile(`(\r\n|\r|\n)`)

	// Replace all occurrences of the regex with a single space
	str = whitespace_regex.ReplaceAllString(str, " ")
	str = line_breakRegex.ReplaceAllString(str, "")
	return str

}
