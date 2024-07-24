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
	utils.Debug("Rendering content")
	utils.Print(raw_content)
	clean_content := utils.CleanString(raw_content)
	imports_string, server_funcs, content := extract.ServerLogic(clean_content)
	locations := extract.ContentScanner(content)
	operations_split := store.SmallIntStore()

	if(len(imports_string) > 0){
		utils.Debug("Imports: ")
		for _, imp := range imports_string {
			utils.Print(imp)
		}
	}

	if(len(server_funcs) > 0){
		utils.Debug("ServerFuncs: ")
		for _, server_func := range server_funcs {
			utils.Print(server_func)
		}
	}

	// split at target and then split again at end tag to get the content of the operation ensure all chunks are stored in a store
	// contentLen := len(content)
	for index, location := range locations {
		// preserve leading content
		if index == 0 {
			operations_split.Set(index, content[:location])
		}

		raw_tag_start := extract.TagName(content[location:])
		end_index := -1

		tag := extract.TagType(raw_tag_start)
		if tag == "operation" {
			name := raw_tag_start[2:]
			op_end_index := extract.FindOperationEndTag(content[location:], name)
			// uses endIndex to preserve trailing content after component
			end_index = location + op_end_index

			// contents from beginning of operation tag to end of operation closing tag
			extracted_operation := content[location:end_index]

			// render html from template operation
			parsed_operation := OperationParser(extracted_operation, name, data, components)

			operations_split.Set(location, parsed_operation)
		} else if tag == "component" {
			name := raw_tag_start[1:]

			comp_end_index := extract.ComponentEndTag(content[location:], name)
			// uses endIndex to preserve trailing content after component
			end_index = location + comp_end_index

			// contents from beginning of component tag to end of component closing tag
			extracted_component := content[location:end_index]

			results := RenderComponent(extracted_component, name, data, components)

			operations_split.Set(location, "{"+name+" Component}"+results)
		} else {
			utils.Error("failed to resolve tag type " + raw_tag_start)
			utils.Error("Exiting from engine.Render()")
			return "failed to resolve tag type " + raw_tag_start
		}

		// // preserve trailing content
		var limit int = -1
		if index+1 < len(locations) {
			limit = locations[index+1]
			operations_split.Set(location+2, content[end_index:limit])
		}

	}

	output := ""

	op_keys := operations_split.SortedKeys()
	if len(op_keys) == 0 {
		output = logic.InsertData(content, data)
	} else {
		for _, key := range op_keys {
			output += logic.InsertData(operations_split.Get(key), data)
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

	options_string, inner_Content := extract.ComponentContent(content, name)

	utils.Debug("Rendering component: " + name)
	utils.Print(options_string)
	utils.Print(inner_Content)
	utils.Print("-------")
	split := strings.Split(options_string, "{%")
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
	options_string, content := extract.OperationContent(content, tag)

	switch tag {
	case "if":
		content = IfOperation(options_string, content, context_data, components)
	case "each":
		content = EachOperation(options_string, content, context_data, components)
	}

	return content
}

/*
<%each {%data.clients} | client>

	<p>{%client}</p>

</%each>
*/
func EachOperation(options_string string, content string, data gjson.Result, components map[string][]byte) string {
	split_options := strings.Split(strings.ReplaceAll(options_string, " ", ""), "|")
	if len(split_options) != 2 {
		utils.Error("Invalid options for each operation")
		return ""
	}

	iterator_data_variable := split_options[0][2 : len(split_options[0])-1]
	var_name := split_options[1]
	data_array := data.Get(iterator_data_variable).Array()
	result := ""
	for _, item := range data_array {
		data, _ := sjson.Set(data.String(), var_name, item.String())
		result += Render(content, gjson.Parse(data), components)
	}

	return result
}

/*
<%if {%data.is_logged_in}>

	<p>User is logged in</p>

</%if>
*/
func IfOperation(options_string string, content string, data gjson.Result, components map[string][]byte) string {
	operation := logic.InsertData(options_string, data)
	value := logic.StringToBoolean(operation)
	if value {
		return Render(content, data, components)
	} else {
		return ""
	}
}
