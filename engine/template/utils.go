package template

import (
	"strings"
)

// concatenate strings
func ConcatStrings(s ...string) string {
	return strings.Join(s, "")
}

// returns the index of the from current pos, if non-found the pos will be the same upon return
func IndexAt(s, sep string, pos int) int {
	idx := strings.Index(s[pos:], sep)
	if idx > -1 {
		idx += pos
	}
	return idx
}

// SplitRespectQuotes splits a string by spaces while respecting quoted substrings and removing empty values.
func SplitRespectQuotes(s string) []string {
	var result []string
	var current strings.Builder
	inQuote := false

	for _, char := range s {
		switch char {
		case '"':
			inQuote = !inQuote // Toggle quote mode
		case ' ':
			if !inQuote {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				continue
			}
		}
		current.WriteRune(char)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
