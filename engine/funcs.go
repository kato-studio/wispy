package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// GetDefaultFuncs returns a map of useful template functions.
func GetDefaultFuncs(ctx *RenderCtx) template.FuncMap {
	return template.FuncMap{
		"component": func(name string, data interface{}, slot template.HTML) (template.HTML, error) {
			return Component(ctx, name, data, slot)
		},
		"include": func(args ...interface{}) (template.HTML, error) {
			// Validate arguments
			if len(args) < 1 || len(args) > 2 {
				return "", fmt.Errorf("include: expected 1 or 2 args, got %d", len(args))
			}

			// Extract component name
			name, ok := args[0].(string)
			if !ok {
				return "", fmt.Errorf("include: first arg must be a string")
			}

			// Extract optional data
			var data interface{}
			if len(args) == 2 {
				data = args[1]
			}

			// Use the closure-captured `ctx` to include the component
			return Include(ctx, name, data)
		},
		"safeHTML": SafeHTML,
		"props":    Props,
		"dict":     Props,
		"toJSON":   ToJSON,
		"ifElse":   IfElse,
	}
}

// Include inserts another template file inside the current template.
func Include(ctx *RenderCtx, name string, data interface{}) (template.HTML, error) {
	component, exists := ctx.Site.Components[name]
	if !exists {
		return "", fmt.Errorf("template for Include %s not found", name)
	}

	tmpl, err := template.New(name).Funcs(GetDefaultFuncs(ctx)).Parse(component)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// Component renders a reusable template component.
func Component(ctx *RenderCtx, name string, data interface{}, slot template.HTML) (template.HTML, error) {
	component, exists := ctx.Site.Components[name]
	if !exists {
		return "", fmt.Errorf("component %s not found", name)
	}

	tmpl, err := template.New(name).Funcs(GetDefaultFuncs(ctx)).Parse(component)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// SafeHTML marks a string as safe HTML.
func SafeHTML(htmlStr string) template.HTML {
	return template.HTML(htmlStr)
}

// Props creates a map of props from key-value pairs.
func Props(pairs ...interface{}) map[string]interface{} {
	props := make(map[string]interface{})
	for i := 0; i < len(pairs)-1; i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			continue // Skip invalid keys
		}
		props[key] = pairs[i+1]
	}
	return props
}

// ToJSON converts data to a JSON string.
func ToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// IfElse is a simple ternary-like function.
func IfElse(condition interface{}, trueVal, falseVal interface{}) interface{} {
	// Convert string to boolean if necessary
	if strCond, ok := condition.(string); ok {
		condition = strings.ToLower(strCond) == "true"
	}

	// Ensure condition is a bool
	if boolCond, ok := condition.(bool); ok {
		if boolCond {
			return trueVal
		}
		return falseVal
	}

	return falseVal
}
