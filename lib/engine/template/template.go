package template

import (
	"fmt"
	"html/template"
	"kato-studio/katoengine/lib/engine/extract"
	"kato-studio/katoengine/lib/engine/logic"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

)

// --------------
//	COR RENDER
// --------------
func Render(raw_content string, data gjson.Result, components map[string][]byte) string {
	clean_content := utils.CleanString(raw_content)
	imports_string, server_funcs, content := extract.ServerLogic(clean_content)
	locations := extract.ContentScanner(content)
	operations_split := store.SmallIntStore()

	if(len(imports_string) > 0){
		utils.Debug("Imports: ")
		for _, imp := range imports_string {
			utils.Debug(imp)
		}
	}

	if(len(server_funcs) > 0){
		utils.Debug("ServerFuncs: ")
		for _, server_func := range server_funcs {
			utils.Debug(server_func)
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
			// strip {% from tag name
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
			// strip < from tag name
			name := raw_tag_start[1:]

			comp_end_index := extract.ComponentEndTag(content[location:], name)
			// uses endIndex to preserve trailing content after component
			end_index = location + comp_end_index

			// contents from beginning of component tag to end of component closing tag
			extracted_component := content[location:end_index]

			// imports_map := map[string]string{}
			// for _, imp := range imports_string {
			// 	split := strings.Split(imp, ":")
			// 	// remove any wrapping quotes
			// 	path := strings.Trim(split[1], "'")
			// 	path = strings.Trim(path, "\"")
			// 	imports_map[split[0]] = path
			// }
			imports_map := extract.ImportsMap(imports_string)

			results := RenderComponent(extracted_component, name, data, imports_map, components)

			operations_split.Set(location, results)
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

// <Header data-hello="boop" onclick="alert('hello')" {% text:"Header Text!" subtext:"Subtext!" %} />
func RenderComponent(content string, name string, data gjson.Result, imports map[string]string, components map[string][]byte) string {
	result := ""
	options_string, inner_Content := extract.ComponentContent(content, name)
	split := strings.Split(options_string, "{%")
	attributes := split[0]
	data_props := ""
	component_path := imports[name]+".kato"
	raw_component := components[component_path]
	result = string(raw_component)
	
	if(raw_component == nil){
		utils.Error("Component not found: name: "+name+ " path: " + component_path)
		return ""
	}

	// handle replacing child content via slot
	result = strings.Replace(result, "<slot/>", inner_Content, -1)
	result = strings.Replace(result, "@root", attributes, 1)
	
	// if there are data properties in the component
	if len(split) > 1 {
		data_props = strings.Trim(strings.Replace(split[1], "%}", "", -1), " ")
		json_string := strings.Replace("{" + data_props + "}",",}","}",-1)

		isValidJson := gjson.Valid(json_string)
		if(isValidJson){
			for key, value := range gjson.Parse(json_string).Map() {
				dataString := data.String()
				dataString, _ = sjson.Set(dataString, key, value)
				data = gjson.Parse(dataString)
			}
		}else{
			utils.Error("Invalid JSON string: " + json_string)
		}
	}

	render_component := Render(result, data, components)

	return render_component
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

	utils.Debug("EachOperation")
	utils.Debug(options_string)
	utils.Debug(content)
	utils.Debug("----------")

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


func SlipScanner(content string) []int {
	// find all components
	locations := []int{}

	content_length := len(content)-1
	for i := 0; i < content_length; i++ {
		// Check if we are at the end of the content

		char := content[i]
		next_char := content[i+1]

		if char == '<' {
			// If next char is capital letter, then assume it's a component
			if next_char >= 'A' && next_char <= 'Z' {
				locations = append(locations, i)
				// skip the next char
				i = i + 5
			}
		}
	}

	return locations
}

func SlipEngine(name string, page_bytes []byte, json_data gjson.Result) string {
		out := new(strings.Builder)

		page_raw := string(page_bytes)
		var err error
		var page_template *template.Template

		// imports_string, server_funcs,
		_, _, page_html := extract.ServerLogic(page_raw)
		// imports_map := extract.ImportsMap(imports_string)

		component_indexs := SlipScanner(page_html)
		utils.Print("component_indexs: "+fmt.Sprint(component_indexs))
		for _, location := range component_indexs {
			component_name := extract.ComponentName(page_html[location:])
			clean_name := strings.Replace(strings.Replace(component_name," ", "",-1),"\n", "",-1)
			component_end_location := extract.ComponentEndTag(page_html[location+2:], clean_name)

			component_content := page_html[location:location+component_end_location+2]
			page_template, err = template.New(clean_name).Parse(component_content)
			utils.Fatal(err)
		}

		err = page_template.ExecuteTemplate(out, "page", json_data.Value())
		utils.Fatal(err)

		return out.String()
}

func LoadTemplateComponents(page_html string, paths []string) string {
	const components_dir = "./templates/components/"
	
	for _, path := range paths {
		component_bytes, _ := os.ReadFile(components_dir + path)
		component_html := string(component_bytes)
		page_html = "{{define \"@" + path + "\"}}" + component_html + "{{end}}" + page_html
	}

	return page_html
}

