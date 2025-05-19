package wispy_common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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

// ParseKeyValuePairs handles the key=value pairs parsing logic
func ParseKeyValuePairs(pairs []string) map[string]string {
	options := make(map[string]string)

	for _, pair := range pairs {
		if strings.Contains(pair, "=") {
			parts := strings.SplitN(pair, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			options[key] = value
		} else {
			options[pair] = "true"
		}
	}

	return options
}

func ParseDataPath(parts []string, value any) any {
	var current_val any
	var exists bool
	for _, part := range parts {
		if part == "" {
			continue // Skip empty parts (e.g., leading dot).
		}
		// Check if the current value is a map.
		if m, ok := value.(map[string]any); ok {
			if current_val, exists = m[part]; exists {
				value = current_val
			} else {
				// Variable not found in Props, fall back to Data.
				return nil
			}
		} else {
			// Invalid path in Props, fall back to Data.
			return nil
		}
	}
	return value
}

// this stringify function curerntly handles non-string values when resolving variables from ctx.Data & ctx.Props
func Stringify(val interface{}) string {
	if val == nil {
		return ""
	}

	// Handle primitives quickly without reflection
	switch v := val.(type) {
	case string:
		return v
	case bool:
		if v {
			return "True"
		}
		return "False"
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}

	// Handle composite types with reflection
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		if length == 0 {
			return "[]"
		}

		// Preallocate slice for faster joins
		parts := make([]string, 0, length)
		for i := 0; i < length; i++ {
			parts = append(parts, Stringify(rv.Index(i).Interface()))
		}
		return "[" + strings.Join(parts, ", ") + "]"

	case reflect.Map:
		keys := rv.MapKeys()
		if len(keys) == 0 {
			return "{}"
		}

		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, fmt.Sprintf(
				"%s: %s",
				Stringify(key.Interface()),
				Stringify(rv.MapIndex(key).Interface()),
			))
		}
		return "{" + strings.Join(parts, ", ") + "}"

	default:
		// Fast path for unknown types - faster than fmt.Sprintf
		if s, ok := val.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprintf("%v", val)
	}
}
