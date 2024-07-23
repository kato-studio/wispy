package template

import (
	"kato-studio/katoengine/lib/engine/extract"
	"kato-studio/katoengine/lib/engine/logic"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// --------------
//
//	COR RENDER
//
// --------------
func Render(raw_content string, data gjson.Result, components map[string][]byte) string {
	content := utils.CleanString(raw_content)

	locations := extract.ContentScanner(content)

	operationsSplit := store.SmallIntStore()

	// split at target and then split again at end tag to get the content of the operation ensure all chunks are stored in a store
	// contentLen := len(content)
	for index, location := range locations {
		// preserve leading content
		if index == 0 {
			operationsSplit.Set(index, content[:location])
		}

		rawTagStart := extract.TagName(content[location:])
		endIndex := -1

		tag := extract.TagType(rawTagStart)
		if tag == "operation" {
			name := rawTagStart[2:]
			opEndIndex := extract.FindOperationEndTag(content[location:], name)
			// uses endIndex to preserve trailing content after component
			endIndex = location + opEndIndex

			// contents from beginning of operation tag to end of operation closing tag
			extractedOperation := content[location:endIndex]

			// render html from template operation
			parsedOperation := OperationParser(extractedOperation, name, data, components)

			operationsSplit.Set(location, parsedOperation)
		} else if tag == "component" {
			name := rawTagStart[1:]

			compEndIndex := extract.ComponentEndTag(content[location:], name)
			// uses endIndex to preserve trailing content after component
			endIndex = location + compEndIndex

			// contents from beginning of component tag to end of component closing tag
			extractedComponent := content[location:endIndex]

			results := RenderComponent(extractedComponent, name, data, components)

			operationsSplit.Set(location, "{"+name+" Component}"+results)
		} else {
			utils.Error("failed to resolve tag type " + rawTagStart)
			utils.Error("Exiting from engine.Render()")
			return "failed to resolve tag type " + rawTagStart
		}

		// // preserve trailing content
		var limit int = -1
		if index+1 < len(locations) {
			limit = locations[index+1]
			operationsSplit.Set(location+2, content[endIndex:limit])
		}

	}

	output := ""

	opKeys := operationsSplit.SortedKeys()
	if len(opKeys) == 0 {
		output = logic.InsertData(content, data)
	} else {
		for _, key := range opKeys {
			output += logic.InsertData(operationsSplit.Get(key), data)
		}
	}

	return output
}

// --------------
//
//	Components
//
// --------------

// <Header data-hello="boop" onclick="alert('hello')" {% text="Header Text!" subtext="Subtext!" %} />
func RenderComponent(content string, name string, data gjson.Result, components map[string][]byte) string {
	result := ""

	optionsString, innerContent := extract.ComponentContent(content, name)

	utils.Debug("Rendering component: " + name)
	utils.Print(optionsString)
	utils.Print(innerContent)
	utils.Print("-------")
	// strings.Split(optionsString, "{%")
	split := strings.Split(optionsString, "{%")
	attributes := split[0]
	data_props := ""

	// if there are data properties in the component
	if len(split) > 1 {
		data_props = strings.Trim(strings.Replace(split[1], "%}", "", -1), " ")
	}

	utils.Print(attributes)
	utils.Print(data_props)
	utils.Print("-------")

	return result
}

// --------------
//
//	Operations
//
// --------------
func OperationParser(content string, tag string, context_data gjson.Result, components map[string][]byte) string {
	optionsString, content := extract.OperationContent(content, tag)

	switch tag {
	case "if":
		content = IfOperation(optionsString, content, context_data, components)
	case "each":
		content = EachOperation(optionsString, content, context_data, components)
	}

	return content
}

/*
<%each {%data.clients} | client>

	<p>{%client}</p>

</%each>
*/
func EachOperation(optionsString string, content string, data gjson.Result, components map[string][]byte) string {
	splitOptions := strings.Split(strings.ReplaceAll(optionsString, " ", ""), "|")
	if len(splitOptions) != 2 {
		utils.Error("Invalid options for each operation")
		return ""
	}

	iteratorDataVariable := splitOptions[0][2 : len(splitOptions[0])-1]
	varName := splitOptions[1]
	dataArray := data.Get(iteratorDataVariable).Array()
	result := ""
	for _, item := range dataArray {
		data, _ := sjson.Set(data.String(), varName, item.String())
		result += Render(content, gjson.Parse(data), components)
	}

	return result
}

/*
<%if {%data.is_logged_in}>

	<p>User is logged in</p>

</%if>
*/
func IfOperation(optionsString string, content string, data gjson.Result, components map[string][]byte) string {
	operation := logic.InsertData(optionsString, data)
	value := logic.StringToBoolean(operation)
	if value {
		return Render(content, data, components)
	} else {
		return ""
	}
}
