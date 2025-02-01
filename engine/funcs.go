package engine

import (
	"bytes"
	"encoding/json"
	"html/template"
	"path/filepath"
)

func GetDefaultFuncs(e *EngineCtx) template.FuncMap {
	return template.FuncMap{
		"component": func(e *EngineCtx, name string, data interface{}) (template.HTML, error) {
			return Component(e, name, data)
		},
		"include": func(e *EngineCtx, name string, data interface{}) (template.HTML, error) {
			return Include(e, name, data)
		},
		"safeHTML": func(htmlStr string) template.HTML {
			return SafeHTML(htmlStr)
		},
		"props": func(pairs ...interface{}) map[string]interface{} {
			return Props(pairs)
		},
		"toJSON": func(data interface{}) (string, error) {
			return ToJSON(data)
		},
		"ifElse": func(condition bool, trueVal, falseVal interface{}) interface{} {
			return IfElse(condition, trueVal, falseVal)
		},
	}
}

// Component renders a reusable template component.
func Component(e *EngineCtx, name string, data interface{}) (template.HTML, error) {
	tmplPath := filepath.Join("templates", name+".html")
	tmpl, err := template.New(name).Funcs(GetDefaultFuncs(e)).ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// Include inserts another template file inside the current template.
func Include(e *EngineCtx, name string, data interface{}) (template.HTML, error) {
	return Component(e, name, data)
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
func IfElse(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
