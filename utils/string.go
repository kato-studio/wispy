package utils

import (
	"regexp"
	"strings"
)

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

// split string at next given separator and return the two parts
func SplitAt(s, sep string) (string, string) {
	i := strings.Index(s, sep)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+len(sep):]
}

func SplitAtFirst(s string, seps []string) (string, string) {
	for _, sep := range seps {
		i := strings.Index(s, sep)
		if i != -1 {
			return s[:i], s[i+len(sep):]
		}
	}
	return s, ""
}

func SplitAtRune(s string, r rune) (string, string) {
	i := strings.IndexRune(s, r)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func FindComponentTagEnd(s string, start int) int {
	inQuotes := false
	tagName := ""
	isClosingTag := false
	isSelfClosing := false

	// Find the tag name
	for i := start + 1; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '>' {
			tagName = s[start+1 : i]
			break
		}
	}

	for i := start; i < len(s); i++ {
		switch s[i] {
		case '"':
			inQuotes = !inQuotes
		case '/':
			if !inQuotes && i+1 < len(s) && s[i+1] == '>' {
				isSelfClosing = true
			}
		case '<':
			if !inQuotes && i+1 < len(s) && s[i+1] == '/' {
				isClosingTag = true
			}
		case '>':
			if !inQuotes {
				if isSelfClosing {
					return i
				}
				if isClosingTag {
					closingTagName := s[i-len(tagName)-2 : i-1]
					if closingTagName == tagName {
						return i
					}
				}
			}
		}
	}
	return -1 // Not found
}