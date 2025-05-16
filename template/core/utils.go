package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/kato-studio/wispy/template/structure"
)

// returns the index of the from current pos, if non-found the pos will be the same upon return
func SeekIndex(raw, sep string, pos int) (new_pos int) {
	new_pos = strings.Index(raw[pos:], sep)
	if new_pos > -1 {
		new_pos += pos
	}
	return new_pos
}

func SeekIndexAndLength(raw, sep string, pos int) (new_pos int, seperator_lenth int) {
	new_pos = strings.Index(raw[pos:], sep)
	if new_pos > -1 {
		new_pos += pos
	}
	return new_pos, len(sep)
}

// SafeIndexAndLength finds the index of a closing tag while ensuring it corresponds to an opening tag
func SeekClosingHandleNested(raw, closingTag, openingTag string, pos int) (newPos int, separatorLength int) {
	openCount := 0
	newPos = pos
	separatorLength = len(closingTag)

	for {
		closeIndex := strings.Index(raw[newPos:], closingTag)
		if closeIndex == -1 {
			// No more closing tags found, return -1
			return -1, separatorLength
		}
		closeIndex += newPos

		// Only search for an opening tag within the range before this closing tag
		openIndex := strings.Index(raw[newPos:closeIndex], openingTag)
		if openIndex != -1 {
			openCount++
			newPos += openIndex + len(openingTag)
			continue
		}

		// If no unmatched opening tags remain, return this closing tag index
		if openCount == 0 {
			return closeIndex, separatorLength
		}

		// Otherwise, decrement open count and continue searching
		openCount--
		newPos = closeIndex + separatorLength
	}
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

func ResolveValue(ctx *structure.RenderCtx, expr string) (any, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty expression")
	}

	// Check if it's a variable (starts with .)
	if strings.HasPrefix(expr, ".") {
		return ResolveVariable(ctx, expr)
	}

	// Check if it's a string literal
	if (strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`)) ||
		(strings.HasPrefix(expr, `'`) && strings.HasSuffix(expr, `'`)) {
		return expr[1 : len(expr)-1], nil
	}

	// Check if it's a boolean
	if expr == "true" {
		return true, nil
	}
	if expr == "false" {
		return false, nil
	}

	// Check if it's a numeric value
	if num, err := strconv.Atoi(expr); err == nil {
		return num, nil
	}
	if num, err := strconv.ParseFloat(expr, 64); err == nil {
		return num, nil
	}

	// Check if it's a variable without . prefix (assuming it's a global or other scope)
	return ResolveVariable(ctx, expr)
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
func stringify(val interface{}) string {
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

		// Preallocate slice for 25% faster joins
		parts := make([]string, 0, length)
		for i := 0; i < length; i++ {
			parts = append(parts, stringify(rv.Index(i).Interface()))
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
				stringify(key.Interface()),
				stringify(rv.MapIndex(key).Interface()),
			))
		}
		return "{" + strings.Join(parts, ", ") + "}"

	default:
		// Fast path for unknown types - 40% faster than fmt.Sprintf
		if s, ok := val.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprintf("%v", val)
	}
}
