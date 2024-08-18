package engine

import (
	"fmt"
	"io"
	"kato-studio/katoengine/utils"

	"github.com/tidwall/gjson"

)

func SlipEngine(template string, json gjson.Result) string {
	// Handle variables
	output := TemplateVariables(template, json)
	// Handle components
	comTmp := ComponentTemplate(output)
	output = comTmp.ExecuteFuncString(func(w io.Writer, extracted string) (int, error) {
		utils.Print("extracted")
		utils.Print(extracted)
		
		tag, rest := SplitAtRune(extracted, ' ')
		result := TemplateFunctions(tag, rest)

		return w.Write([]byte(result))
	})

	utils.Print("-----------")
	fmt.Printf("%s", output)

	return output
}

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