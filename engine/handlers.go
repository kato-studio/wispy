package engine

import (
	"io"

	"github.com/tidwall/gjson"
)

func TemplateVariables(template string, json gjson.Result) string {
	varTmp := VariableTemplate(template)
	result := varTmp.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		var firstRun = tag[0] 
		// preserve scoped variables in template
		// scoped variables are used 
		if(firstRun == '$') {
			return w.Write([]byte("{{"+tag+"}}"))
		}

		if(json.Get(tag).Exists()) {
			return w.Write([]byte(json.Get(tag).String()))
		}
		return w.Write([]byte(""))
	})
	return result
}

func TemplateFunctions(tag string, rest string) string{
	var firstRun = tag[0] 
	if(firstRun > 'A' && firstRun < 'Z') {
		// Component
		return "COMPONENT {"+tag+"}"
	}else if(tag == "if") {
		// If
		return "IF HERE"
	}else if(tag == "each") {
		// Each
		return "EACH HERE"
	}else{
		return "UNKNOWN"
	}
}