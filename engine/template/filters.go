package template

import (
	"strconv"
	"strings"
)

// UpcaseFilter converts a string value to uppercase.
func UpcaseFilter(value interface{}, args []string) interface{} {
	if s, ok := value.(string); ok {
		return strings.ToUpper(s)
	}
	return value
}

// DowncaseFilter converts a string value to lowercase.
func DowncaseFilter(value interface{}, args []string) interface{} {
	if s, ok := value.(string); ok {
		return strings.ToLower(s)
	}
	return value
}

// CapitalizeFilter converts the first character of a string to uppercase.
func CapitalizeFilter(value interface{}, args []string) interface{} {
	if s, ok := value.(string); ok && len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return value
}

// StripFilter trims spaces from a string.
func StripFilter(value interface{}, args []string) interface{} {
	if s, ok := value.(string); ok {
		return strings.TrimSpace(s)
	}
	return value
}

// TruncateFilter shortens a string to the specified length.
func TruncateFilter(value interface{}, args []string) interface{} {
	if s, ok := value.(string); ok && len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && len(s) > n {
			return s[:n]
		}
	}
	return value
}

// SliceFilter splits a string by the given delimiter (default is a comma).
func SliceFilter(value interface{}, args []string) interface{} {
	delimiter := ","
	if len(args) > 0 && args[0] != "" {
		delimiter = args[0]
	}
	if s, ok := value.(string); ok {
		parts := strings.Split(s, delimiter)
		var result []string
		for _, part := range parts {
			result = append(result, strings.TrimSpace(part))
		}
		return result
	}
	return value
}
