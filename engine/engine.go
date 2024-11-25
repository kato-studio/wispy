package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// RenderPage handles not only the rendering of the page but also the layout
// and checking for parent layouts that may also need to be rendered
func RenderPage(raw_path string, ctx RenderCTX) string {
	var result = ""
	var path_split = strings.Split(raw_path, "/")
	var site_domain = path_split[1]
	path_split = path_split[2:]
	// validate name regex (validate the site name is a domain)
	is_valid_domain := ValidateDomainName(site_domain)
	if !is_valid_domain {
		fmt.Println("[Error]: The site name is not a valid domain: ", site_domain)
		// Todo: return a 404 page or something (maybe a dev mode and prod mode response)
		return ""
	}
	var base_path = ROOT_DIR + "/" + site_domain + "/pages/"
	dirs_length := len(path_split)

	var page_path = base_path + strings.Join(path_split, "/") + PAGE_FILE
	if _, err := os.Stat(page_path); err == nil {
		page_bytes, err := os.ReadFile(page_path)
		if err != nil {
			fmt.Println("[Error]: Could not read the page file: ", err)
			return ""
		}
		result = CleanTemplate(string(page_bytes))
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("[Error]: Could not find the page file: ", page_path)
		return ""
	} else {
		fmt.Println("[Error]: Error while checking for page: ", err)
		return ""
	}

	// safety check for the path
	if dirs_length > 0 {
		// "i >= 1" because the "path_split" will always contain and empty string at the start
		// this is because the path starts with a "/" and the split function will always add an empty string
		// this is fine as we need the extra slash on the join anyways
		for i := dirs_length; i >= 1; i-- {
			if path_split[i-1] == "" {
				break
			}
			// if empty path artifact then skip
			parent_layout_path := base_path + strings.Join(path_split[:i], "/") + LAYOUT_FILE
			// check if the path part is a variable
			if _, err := os.Stat(parent_layout_path); err == nil {
				layout_bytes, err := os.ReadFile(parent_layout_path)
				if err != nil {
					fmt.Println("[Error]: Could not read the layout file: ", err)
					continue
				}
				// layout found! replace the layout with the page
				fmt.Println("[Info]: Found the layout file: ", parent_layout_path)
				this_layout := string(layout_bytes)
				result = strings.ReplaceAll(this_layout, "<_slot/>", result)
				// error handling
			} else if errors.Is(err, os.ErrNotExist) {
				fmt.Println("[Info]: Could not find the layout file: ", parent_layout_path)
				continue
			} else {
				fmt.Println("[Error]: Error while checking for layout: ", err)
			}
		}
	}
	// Insert contents to site root document & resolve head tags
	// -----
	root_document_path := base_path + DOCUMENT_FILE
	root_document, err := os.ReadFile(root_document_path)
	if err != nil {
		fmt.Println("[Error]: Could not read the root document: ", err)
	}
	//
	clean_doc := CleanTemplate(string(root_document))
	result = strings.ReplaceAll(clean_doc, "{{BODY_CONTENT}}", result)
	// Todo: Resolve head tags
	result = strings.ReplaceAll(result, "{{HEAD_CONTENT}}", "")
	return SlipEngine(result, ctx)
}

// -------------------------
// SlipEngine is the main function that will render the template
// -------------------------
// this is a simple implementation and can be optimized
func SlipEngine(template string, ctx RenderCTX) string {
	// clean the template
	template = CleanTemplate(template)

	// core imports logic
	var components = ctx.Components
	//
	imports_start_tag := "<script imports>"
	imports_end_tag := "</script>"
	//
	imports_start_index := strings.Index(template, imports_start_tag)
	imports_end_index := -1
	//
	if imports_start_index > -1 {
		imports_end_index = strings.Index(template, imports_end_tag)
		// handle edge case where the end tag is not found
		if imports_end_index < 0 {
			fmt.Println("[Error]: Could not find the end tag for the imports")
			return ""
		}
		// ensure imports onl includes contents within imports
		imports := template[imports_start_index+len(imports_start_tag) : imports_end_index]
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

			// TODO: better error handling
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
			components[alias] = CleanTemplate(string(file))
		}

		// update the context with imported components
		ctx.Components = components
		// remove the imports from the template (account for end tag length)
		// template = template[imports_start_index : imports_end_index+len(imports_end_tag)]
		template = template[:imports_start_index] + template[imports_end_index+len(imports_end_tag):]

	}

	// ----
	// End of core imports logic
	// ----
	//
	// ----
	// Core rendering logic
	// ----
	var result = ""
	// hoisted variables
	tag_name := ""
	tag_start := -1
	tag_end_end := -1
	// handle nested tags of the same type
	nested_depth := 0
	self_closing := false
	// loop starts at to prevent out of bounds error
	// when checking for the previous character
	for i := 2; i < len(template); i++ {
		// cache the previous characters
		prev := template[i-1]
		char := template[i]
		//
		// preserve the start
		if prev == '<' && char == '_' {
			// Todo: handle index not found (AKA better error handling! PS: maybe a prod & dev response modes)
			// EXAMPLE: panic: runtime error: slice bounds out of range [18:17]
			// I had used self closing end tag rather then a proper end tag
			// should should provide a better error message
			//
			test_index := strings.Index(template[i+1:], " ")
			test_name := template[i+1 : i+test_index+1]
			//
			// if not currently in a tag then start a new tag
			if tag_name == "" {
				tag_name = test_name
				tag_start = i - 1
				// remove the last character from the result since it's the start of a new tag
				result = result[:len(result)-1]
				// if currently in a tag then check if it's the same tag type if it is then increment the nested_depth
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
	//
	return InsertValues(result, ctx.Json)
}

// -------------------------
// Handles the logic for rendering operations e.g. <_if>, <_for>
// -------------------------
// this is a simple implementation and can be optimized
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
	//
	attributes := regex_attributes.FindAllString(contents[:start_tag_end], -1)
	// child_attributes := []string{}
	for i, operation := range attributes {
		// regex should guarantee that the split will have 2 elements
		raw_split := strings.Split(operation, "=")
		attr_name := strings.TrimSpace(raw_split[0])
		// _ is the value of the attribute
		// current operation do not require the value
		_, value_type, value_path := ValueOrTrimmed(raw_split[1], ctx.Json)
		//
		// HANDLE: If
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
		//
		// HANDLE: Each + In
		if attr_name == "each" {
			continue
		}
		if attr_name == "in" && i > 0 {
			// this should be safe todo since regex should guarantee that the split will have 2 elements
			each_value_name := strings.Split(attributes[i-1], "=")
			raw_name, _, _ := ValueOrTrimmed(each_value_name[1], ctx.Json)
			item_name := strings.Trim(raw_name, "\"")
			//
			// get the value of the attribute
			// _, value_type, value_path := ValueOrTrimmed(attr_value, ctx.Json)
			// fmt.Println("Value Path: ", value_path)
			// check if the value is a boolean
			if value_type == "ARRAY" {
				// get the array
				array := ctx.Json.Get(value_path).Array()
				for _, item := range array {
					// create a new context with the item
					new_json, err := sjson.Set("{}", item_name, item.Value())
					if err == nil {
						// TODO: Don't include all context only the necessary context to component
						// TODO: implement hoisted context
						new_ctx := RenderCTX{
							Json:       gjson.Parse(new_json),
							Components: ctx.Components,
							Head:       ctx.Head,
							Variables:  ctx.Variables,
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
	// HANDLE: (Other operations)
	// ----
	// tbd (to be done)
	// ----
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
	if !ok {
		fmt.Println("[Error]: Could not find the component: ", name)
		return "\n[SAD :(]\n"
	}

	// ----
	// Basic handling for passing child content
	// TODO: better handling for if <_slot/> has a space before closing E.g. <_slot />
	component = strings.ReplaceAll(component, "<_slot/>", inner_content)
	component = strings.ReplaceAll(component, "<_slot />", inner_content)

	// ----
	// Handle passing props/data to the component
	var attributes = regex_attributes.FindAllString(contents[:start_tag_end], -1)
	new_json, err := sjson.Set("{}", "attributes", "")

	for _, item := range attributes {
		each_value_name := strings.Split(item, "=")
		name := each_value_name[0] // shouldn't be needed... strings.TrimSpace()
		parsed_value, _, _ := ValueOrTrimmed(each_value_name[1], ctx.Json)
		//
		new_json, err = sjson.Set(new_json, name, parsed_value)
		if err != nil {
			fmt.Println("[Error]: Could not pass the item to the new context: ", name)
		}
	}
	if err != nil {
		fmt.Println("[Error]: Could not set all props for the component: ", name)
	}
	// create a new context with passed props
	new_ctx := RenderCTX{
		Json:       gjson.Parse(new_json),
		Components: ctx.Components,
		Head:       ctx.Head,
		Variables:  ctx.Variables,
	}

	// ----
	// Render with the new context
	result += SlipEngine(component, new_ctx)
	return result
}
