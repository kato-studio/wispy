package engine

import (
	"bytes"
	"fmt"
	"unicode"
)

// smartSplit splits a []byte at spaces, respecting quoted sections & round brackets sections as well keeping the quotes in the output
func smartSplit(input []byte) []string {
	var result []string
	var current []byte
	inQuotes := false
	inBrackets := false

	for i := 0; i < len(input); i++ {
		char := input[i]
		switch char {
		case ' ':
			if inQuotes || inBrackets {
				// Inside quotes or brackets, keep spaces as part of the token
				current = append(current, char)
			} else if len(current) > 0 {
				// Outside quotes/brackets, space signals end of a token
				result = append(result, string(current))
				current = nil
			}
		case '"':
			// Toggle the inQuotes state
			inQuotes = !inQuotes
			current = append(current, char) // Include the quote in the token
		case '(':
			inBrackets = true
			current = append(current, char, ')')
		case ')':
			if inBrackets {
				inBrackets = false
			}
		default:
			current = append(current, char)
		}
	}

	// Append the last token if it's not empty
	if len(current) > 0 {
		result = append(result, string(current))
	}
	fmt.Println("results from smartSplit()", result)
	return result
}

func parseAttributes(input []byte) map[string]string {
	attributes := make(map[string]string)
	var key, value []byte
	var inKey, inValue, inQuotes bool
	var quoteChar byte

	for _, b := range input {
		switch {
		case b == '=' && inKey: // End of the key
			inKey = false
			inValue = true

		case (b == '"' || b == '\'') && inValue: // Start or end of a quoted value
			if !inQuotes {
				inQuotes = true
				quoteChar = b
			} else if b == quoteChar { // End of quoted value
				inQuotes = false
				inValue = false
				attributes[string(bytes.TrimSpace(key))] = string(value)
				key, value = nil, nil // Reset for the next attribute
			} else {
				value = append(value, b)
			}

		case inKey: // Accumulate key bytes
			key = append(key, b)

		case inValue && inQuotes: // Accumulate value bytes inside quotes
			value = append(value, b)

		case b != ' ' && !inKey && !inValue: // Start of a new key
			inKey = true
			key = append(key, b)
		}
	}

	return attributes
}

func isComponent(name string) bool {
	return len(name) > 0 && unicode.IsUpper(rune(name[0]))
}

// Detects the type of a tag based on its characteristics.
func detectTagType(tagName string) string {
	if isComponent(tagName) {
		return "component"
	}
	var minLen2 = len(tagName) > 2
	if minLen2 && tagName[:2] == "x:" {
		return "operation"
	}
	return "element"
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
}
