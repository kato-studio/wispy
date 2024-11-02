package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var attributes_regex = regexp.MustCompile(`([-\w:])+=({{(.*?)}}|\"(.*?)\")`)

func RenderPage(contents string, path string, ctx RenderCTX) string {
	var result = contents
	//
	path_split := strings.Split(path, "/")
	dirs_length := len(path_split) - 1
	// safety check for the path
	if dirs_length > 0 {
		for i := dirs_length; i >= 0; i-- {
			parent_layout_path := strings.Join(path_split[:i], "/") + LAYOUT_FILE
			// check if the path part is a variable
			if _, err := os.Stat(parent_layout_path); err == nil {
				layout_bytes, err := os.ReadFile(parent_layout_path)
				if err != nil {
					fmt.Println("[Error]: Could not read the layout file: ", err)
					continue
				}
				// layout found! replace the layout with the page
				this_layout := string(layout_bytes)
				result = strings.ReplaceAll(this_layout, "<_slot/>", result)
				// error handling
			} else if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				fmt.Println("[Error]: Error while checking for layout: ", err)
			}
		}
	}
	//
	return SlipEngine(result, ctx)
}

// -------------------------
// SlipEngine is the main function that will render the template
// -------------------------
func SlipEngine(template string, ctx RenderCTX) string {
	var components = make(map[string]string)
	parse_start_index := 2
	//
	imports_start_tag := "<script imports>"
	imports_end_tag := "</script>"
	imports_start_index := strings.Index(template, imports_start_tag)
	if imports_start_index > -1 {
		imports_end_index := strings.Index(template, imports_end_tag)
		// handle edge case where the end tag is not found
		if imports_end_index < 0 {
			fmt.Println("[Error]: Could not find the end tag for the imports")
			return ""
		}
		imports := template[imports_start_index+len(imports_start_tag) : imports_end_index]
		parse_start_index = imports_end_index + len(imports_end_tag)
		//
		imports_split := strings.Split(imports, " import ")
		for import_index, import_ := range imports_split {
			expected_min_length := 16
			if len(import_) < expected_min_length {
				// TODO: add debug logic
				continue
			}
			// split the import into the path and the alias
			import_split := strings.Split(import_, " from ")
			if len(import_split) != 2 {
				fmt.Printf("[Error]: The import is not in the correct format %d \n", import_index)
				continue
			}
			//
			// get the path and the alias
			raw_alias := import_split[0]
			raw_path := import_split[1]
			// removing trailing ")" and clean path
			alias := strings.Trim(raw_alias, " ")

			// TODO: error handling
			var component_folder = "./components"
			var split = strings.Split(strings.Trim(strings.TrimSpace(raw_path), "\""), "/")
			var site_name = split[0]
			var comp_path = split[1:]
			if site_name != "components" {
				// if the site name is not components then it's a site component
				component_folder = "./sites/" + site_name + "/components"
			}
			absolute_view_path, err := filepath.Abs(component_folder)
			if err != nil {
				fmt.Println("[Error]: Could not get the absolute path for the view folder")
				continue
			}
			// read the file
			read_path := absolute_view_path + "/" + strings.Join(comp_path, "/") + ".hstm"
			file, err := os.ReadFile(read_path)
			if err != nil {
				fmt.Println("[Error]: Could not import file: ", read_path)
				continue
			}
			// add the alias to the components map
			components[alias] = string(file)
		}
	}
	// find all imports with the regex pattern
	var result = ""
	// hoisted variables
	tag_name := ""
	tag_start := -1
	tag_end_end := -1
	// handle nested tags of the same type
	nested_depth := 0
	self_closing := false
	//
	for i := parse_start_index; i < len(template); i++ {
		// cache the previous characters
		prev := template[i-1]
		char := template[i]
		//
		// preserve the start
		if prev == '<' && char == '_' {
			// todo handle index not found
			test_index := strings.Index(template[i+1:], " ")
			test_name := template[i+1 : i+test_index+1]
			//
			// if not currently in a tag then start a new tag
			if tag_name == "" {
				tag_name = test_name
				tag_start = i - 1
				// remove the last character from the result
				// since it's the start of a new tag
				result = result[:len(result)-1]
				// if currently in a tag then check if it's the same tag type
				// if it is then increment the nested_depth
			} else {
				if tag_name == test_name {
					nested_depth++
				}
				continue
			}
		}
		// handle initial content if not the start of a tag
		if i == 2 {
			result += string(template[:2])
		}
		// if not currently in a tag then start a new tag
		if tag_name == "" {
			// handle initial content if not the start of a tag
			// push the character to the result
			result += string(char)
			continue
		}

		// -------
		// END OF TAG
		// -------
		// check if it's a self-closing tag
		if char == '>' {
			if prev == '/' {
				self_closing = true
				tag_end_end = i + 1
			}
		}

		// check if it's the end of the tag
		if prev == '/' && char == '_' {
			name_length := len(tag_name)
			// with the _ prefix and extra characters
			test_name := template[i+1 : i+name_length+1]
			//
			if test_name == tag_name {
				if nested_depth > 0 {
					nested_depth--
				} else {
					// +1 for the ending ">"
					tag_end_end = i + name_length + 1
					// skip the rest of the tag
					i = tag_end_end
				}
			} else {
				continue
			}
		}
		// if tag_start and tag_end are found then render the tag
		if tag_end_end > 0 || self_closing {
			// get the content of the tag
			tag_content := template[tag_start:tag_end_end]
			if tag_name == "slot" || tag_name == "slot/>" {
				continue
			}
			// check if the tag is an operation or a component
			if tag_name == "render" || tag_name == "null_placeholder" {
				result += HandleOperation(tag_name, tag_content, self_closing, ctx)
			} else {
				result += HandleComponent(tag_name, tag_content, self_closing, ctx, components)
			}
			// reset the all the tag variables
			tag_name = ""
			tag_start = -1
			tag_end_end = -1
			nested_depth = 0
			self_closing = false
			//
		}
		//
		if i == len(template)-1 && tag_name != "" {
			fmt.Println("[Info]: Could not find the end tag: ", tag_name, " at index: ", i)
			fmt.Println("[Error]: Could not find the end tag for the tag_name index: ", tag_name)
			//
			log_end_length := i + 7
			if log_end_length > len(template) {
				log_end_length = len(template)
			}
			//
			return "SAD :( \n\n\n\n" +
				"" +
				fmt.Sprintf("| tag_name: %s \n", tag_name) +
				fmt.Sprintf("| %s \n", template[i-10:log_end_length]) +
				"--------------" +
				fmt.Sprintf("| %d\n", 0) +
				fmt.Sprintf("| %d\n", 0) +
				fmt.Sprintf("| %d\n", 0)
		}
	} // end of for loop

	return InsertValues(result, ctx.Json)
}

// -------------------------
// Handles the logic for rendering operations e.g. <_if>, <_for>
// -------------------------
func HandleOperation(name string, contents string, self_closing bool, ctx RenderCTX) string {
	var result = ""
	var inner_content = ""
	var start_tag_end = -1
	for i := len(name); i < len(contents); i++ {
		char := contents[i]
		prev := contents[i-1]
		if self_closing && prev == '/' && char == '>' {
			start_tag_end = i
			break
		} else if char == '>' && (prev == '"' || prev == '}' || prev == ' ') {
			start_tag_end = i
			break
		}
	}
	// handle edge case where the end tag is not found
	if start_tag_end < 0 {
		fmt.Println("[Error]: Could not find the end tag for the operation: ", name)
		return "\n[SAD :(]\n"
	}
	if !self_closing {
		// -4 is to remove "<_{NAME}/>"
		inner_content = contents[start_tag_end+1 : len(contents)-len(name)-4]
	} else {
		// -2 is to remove "/>" from self-closing tags
		inner_content = contents[start_tag_end+1 : len(contents)-2]
	}

	attributes := attributes_regex.FindAllString(contents[:start_tag_end], -1)
	// child_attributes := []string{}
	for i, operation := range attributes {
		// regex should guarantee that the split will have 2 elements
		raw_split := strings.Split(operation, "=")
		attr_name := strings.TrimSpace(raw_split[0])
		// _ is the value of the attribute
		// current operation do not require the value
		_, value_type, value_path := ValueOrTrimmed(raw_split[1], ctx.Json)

		// If
		if attr_name == "if" {
			// get the value of the attribute
			// check if the value is a boolean
			if value_type == "True" || value_type == "String" {
				result += SlipEngine(inner_content, ctx)
			} else {
				result += ""
				break
			}
			continue
		} // ----------

		// Each + In
		if attr_name == "each" {
			continue
		}
		if attr_name == "in" && i > 0 {
			// this should be safe todo since regex should guarantee that the split will have 2 elements
			each_value_name := strings.Split(attributes[i-1], "=")
			raw_name, _, _ := ValueOrTrimmed(each_value_name[1], ctx.Json)
			item_name := strings.Trim(raw_name, "\"")

			// get the value of the attribute
			// _, value_type, value_path := ValueOrTrimmed(attr_value, ctx.Json)
			// fmt.Println("Value Path: ", value_path)
			// check if the value is a boolean
			if value_type == "ARRAY" {
				// get the array
				array := ctx.Json.Get(value_path).Array()
				for _, item := range array {
					// create a new context with the item
					new_json, err := sjson.Set(ctx.Json.Raw, item_name, item.Value())
					if err == nil {
						/*
							TODO: Don't include all context only the necessary context to component
							todo: implement hoisted context1
						*/
						new_ctx := RenderCTX{
							Json:      gjson.Parse(new_json),
							Snippet:   ctx.Snippet,
							Variables: ctx.Variables,
						}
						// render the inner content
						result += SlipEngine(inner_content, new_ctx)
					} else {
						fmt.Println("[error]: Could not pass the item to the new context: ", item)
						fmt.Println("[debug]: ", err)
					}
				}
				continue
			} else {
				fmt.Println("[error]: Could not find the array for the JSON path: ", value_path)
				result += ""
				break
			}
		} // ----------

		fmt.Println("[error]: Could not find the operation:" + attr_name)
	} // ----------

	// ----
	// result += SlipEngine(inner_content, json)
	return result
}

// -------------------------
// Handles the logic for rendering components e.g. <_Button>, <_Input>
// -------------------------
func HandleComponent(name string, contents string, self_closing bool, ctx RenderCTX, components map[string]string) string {
	var result = ""
	var inner_content = ""
	var start_tag_end = -1
	for i := len(name); i < len(contents); i++ {
		char := contents[i]
		prev := contents[i-1]
		if self_closing && prev == '/' && char == '>' {
			start_tag_end = i
			break
		} else if char == '>' && (prev == '"' || prev == '}' || prev == ' ') {
			start_tag_end = i
			break
		}
	}
	// handle edge case where the end tag is not found
	if start_tag_end < 0 {
		fmt.Println("[Error]: Could not find the end tag for the component: ", name)
		return "\n[SAD :(]\n"
	}
	if !self_closing {
		// -4 is to remove the last 4 characters "<_{NAME}/>"
		inner_content = contents[start_tag_end+1 : len(contents)-len(name)-4]
	}
	// check if the component is in the components map
	component, ok := components[name]
	// TODO: better handling for if <_slot/> has a space before closing E.g. <_slot />
	component = strings.ReplaceAll(component, "<_slot/>", inner_content)
	component = strings.ReplaceAll(component, "<_slot />", inner_content)
	if !ok {
		fmt.Println("[Error]: Could not find the component: ", name)
		return "\n[SAD :(]\n"
	}
	// render the component
	result += SlipEngine(component, ctx)
	return result
}
