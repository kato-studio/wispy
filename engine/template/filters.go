package template

import (
	"strconv"
	"strings"
)

// Universal template data filters function struct
type EngineFilter struct {
	Name    string
	Handler func(pipedValue any, args []string) (value any, err error)
}

// Default template data filters functions
var UpcaseFilter = EngineFilter{
	Name: "upcase",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		if s, ok := pipedValue.(string); ok {
			return strings.ToUpper(s), nil
		}
		return pipedValue, nil
	},
}

var DowncaseFilter = EngineFilter{
	Name: "downcase",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		if s, ok := pipedValue.(string); ok {
			return strings.ToLower(s), nil
		}
		return pipedValue, nil
	},
}

var CapitalizeFilter = EngineFilter{
	Name: "capitalize",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		if s, ok := pipedValue.(string); ok && len(s) > 0 {
			return strings.ToUpper(s[:1]) + s[1:], nil
		}
		return pipedValue, nil
	},
}

var StripFilter = EngineFilter{
	Name: "strip",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		if s, ok := pipedValue.(string); ok {
			return strings.TrimSpace(s), nil
		}
		return pipedValue, nil
	},
}

var TruncateFilter = EngineFilter{
	Name: "truncate",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		if s, ok := pipedValue.(string); ok && len(args) > 0 {
			if n, err := strconv.Atoi(args[0]); err == nil && len(s) > n {
				return s[:n], nil
			}
		}
		return pipedValue, nil
	},
}

var SliceFilter = EngineFilter{
	Name: "slice",
	Handler: func(pipedValue any, args []string) (value any, err error) {
		delimiter := ","
		if len(args) > 0 && args[0] != "" {
			delimiter = args[0]
		}
		if s, ok := pipedValue.(string); ok {
			parts := strings.Split(s, delimiter)
			var result []string
			for _, part := range parts {
				result = append(result, strings.TrimSpace(part))
			}
			return result, nil
		}
		return pipedValue, nil
	},
}

// Register default filters.
var DefaultTemplateFilters = []EngineFilter{
	UpcaseFilter,
	DowncaseFilter,
	CapitalizeFilter,
	StripFilter,
	TruncateFilter,
	SliceFilter,
}
