package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// SlipEngine is the main function that will render the template
func SlipEngine(template string, json gjson.Result) string {
	// 
	var result = ""
	print("Template:\n", template)
	print("\n=======================================\n")
	// 
	// hoisted variables
	tag_name := ""
	tag_start := -1
	tag_end_end := -1

	// handle nested tags of the same type
	nested_depth := 0
	self_closing := false
	// 
	for i := 2; i < len(template); i++ {
		// cache the previous characters		
		prev_prev := template[i-2]
		prev := template[i-1]
		char := template[i]

		// preserve the start
		if (prev_prev == '<' && prev == '#') {
			// todo handle index not found
			test_index := strings.Index(template[i:], " ")
			test_name := template[i:i+test_index]
			// 
			// if not currently in a tag then start a new tag
			if(tag_name == "") {
				tag_name = test_name
				tag_start = i
				// 
				// if currently in a tag then check if it's the same tag type
				// if it is then increment the nested_depth
			}else{
				if(tag_name == test_name) {
					nested_depth++
				}
				continue
			}
		}

		// if currently in a tag then check if it's the end of the tag
		if (len(tag_name) > 0) {
			
			// check if it's a self-closing tag
			if(prev == '>') {
				if(prev_prev == '/') {
					self_closing = true
					tag_end_end = i
				}
			}
			// check if it's the end of the tag
			if(prev_prev == '/' && prev == '#' && char == tag_name[0]) {
				// 
				name_length := len(tag_name)
				// 
				if(template[i:i+name_length] == tag_name) {
					tag_end_end = i+name_length
				}else{
					fmt.Println("[Info]: Could not find the end tag: ", tag_name, " at index: ", i)
					fmt.Println("[Error]: Could not find the end tag for the tag_name index: ", tag_name)
					return "SAD :( \n\n\n\n" +
					""+
					fmt.Sprintf("| tag_name: %s \n",tag_name)  +
					fmt.Sprintf("| %s \n", template[i-7:i+7])  +
					"--------------"+
					fmt.Sprintf("| %d\n", 0)  +
					fmt.Sprintf("| %d\n", 0)  +
					fmt.Sprintf("| %d\n", 0)
				}
			}

			// if tag_start and tag_end are found then render the tag
			if(tag_end_end > 0 || self_closing) {
				print("tag_end_end: ", tag_end_end, "\n")
				// 
				// get the content of the tag
				tag_content := template[tag_start:tag_end_end]
				// 
				char_zero := tag_content[0]

				// check if the tag is an operation or a component
				if(char_zero >= 'A' && char_zero <= 'Z') {
					print("comp_tag_start_end: ", tag_end_end, "\n")
					result += RenderComponent(tag_name, tag_content, self_closing, json)
				}else{
					result += RenderOperation(tag_name, tag_content, self_closing, json)
				}
				// 
				print("[tag name reset] \n")

				// reset the all the tag variables
				tag_name = ""
				tag_start = -1
				// tag_start_end = -1
				tag_end_end = -1
				nested_depth = 0
				self_closing = false
				// 
			}

			if(i == 2) {result += string(template[:2])}
			//
			// if not currently in a tag then start a new tag
			if(tag_name == "") {
				// handle initial content if not the start of a tag
				// push the character to the result
				result += string(char)
			}

		} // end of tag_name check
	} // end of for loop
	return result
}

// Handles the logic for rendering operations e.g. <#if>, <#for>
func RenderOperation(name string, contents string, self_closing bool, json gjson.Result) string {
	fmt.Println("Operation: ", name)
	fmt.Println("Contents: \n", contents)
	fmt.Println("--------------------")
	return "##-OP-###"
}

// Handles the logic for rendering components e.g. <#Button>, <#Input>
func RenderComponent(name string, contents string, self_closing bool, json gjson.Result) string {
	// 
	fmt.Println("Component: ", name)
	fmt.Println("Contents: \n", contents)
	fmt.Println("--------------------")
	return "##-COMP-###"
}

// Regex example
	// Get start tag of operations & components in the template
	// - `<#[A-z]((.|\n)*?)>\n` basically <#ABC ... >\n
	// - we can tell if it"s a component if it has a capital letter after the # sign
	// - otherwise it"s an operation then we can check if the last character is "/>"
	// - to know if it"s a self-closing tag or not
	// fmt.Println("template: ", template)
	// temp_gex := regexp.MustCompile(`(?m)<#[A-z](.+?)>`) //regexp.Compile("p([a-z]+)ch")
	// for i, match := range temp_gex.FindAllString(template, -1) {
  //     fmt.Println(match, "found at index", i)
  // }