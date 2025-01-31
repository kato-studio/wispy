package templateFuncs

import (
	"fmt"
	"html/template"
)

func GetDefaults() template.FuncMap {
	return template.FuncMap{
		"props": Props,
	}
}

func Props(v ...any) map[string]any {
	if len(v)%2 != 0 {
		panic("uneven number of key/value pairs")
	}

	m := map[string]any{}
	for i := 0; i < len(v); i += 2 {
		m[fmt.Sprint(v[i])] = v[i+1]
	}

	return m
}
