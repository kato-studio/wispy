package xops

import (
	"strings"
)

type Render struct {
	Ctx struct {
		Data map[string]interface{}
	}
	Html func(string) (string, error)
}

// Example operation function
func EachOperation(r *Render, values ...string) string {
	if len(values) < 3 {
		return ""
	}
	var variableName = values[0]
	var listName = values[2]

	list, ok := r.Ctx.Data[listName].([]interface{})
	if !ok {
		return ""
	}

	var result strings.Builder
	for index, item := range list {
		r.Ctx.Data[variableName] = map[string]interface{}{
			"index": index,
			"value": item,
		}
		rendered, _ := r.Html(values[1])
		result.WriteString(rendered)
	}
	return result.String()
}

func ForOperation(r *Render, values ...string) string {
	if len(values) < 3 {
		return ""
	}
	var variableName = values[0]
	var listName = values[2]

	list, ok := r.Ctx.Data[listName].([]interface{})
	if !ok {
		return ""
	}

	var result strings.Builder
	for index, item := range list {
		r.Ctx.Data[variableName] = map[string]interface{}{
			"index": index,
			"value": item,
		}
		rendered, _ := r.Html(values[1])
		result.WriteString(rendered)
	}
	return result.String()
}
