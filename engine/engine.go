package engine

import (
	"fmt"
	"strings"
	"unicode"

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
		prev := template[i-1]
		char := template[i]

		// preserve the start
		if (prev == '<' && char == '#') {
			// todo handle index not found
			test_index := strings.Index(template[i+1:], " ")
			test_name := template[i+1:i+test_index+1]
			// 
			// if not currently in a tag then start a new tag
			if(tag_name == "") {
				tag_name = test_name
				tag_start = i-1
				// remove the last character from the result 
				// since it's the start of a new tag
				result = result[:len(result)-1]
			// if currently in a tag then check if it's the same tag type
			// if it is then increment the nested_depth
			}else{
				if(tag_name == test_name) {
					nested_depth++
				}
				continue
			}
		}

		// handle initial content if not the start of a tag
		if(i == 2) {result += string(template[:2])}
		// if not currently in a tag then start a new tag
		if(tag_name == "") {
			// handle initial content if not the start of a tag
			// push the character to the result
			result += string(char)
			continue
		}

		// -------
		// CORE LOGIC
		// -------
		// check if it's a self-closing tag
		if(char == '>') {
			if(prev == '/') {
				self_closing = true
				tag_end_end = i+1
			}
		}
		
		// check if it's the end of the tag
		if(prev == '/' && char == '#') {
			// 
			name_length := len(tag_name)
			// something doesn't check out in this block
			// TODO logs (test_name) is actually the end tag 
			// with the # prefix and extra characters
			test_name := template[i+1:i+name_length+1]
			// 
			fmt.Printf("[end] tag_name: %s, test_name: %s \n", tag_name, test_name)
			if(test_name == tag_name) {
				if(nested_depth > 0) {
					nested_depth--
				}else{
					// +1 for the ending ">"
					tag_end_end = i+name_length+1
					// skip the rest of the tag
					i = tag_end_end
				}
			}else{
				continue
			}
		}
		// if tag_start and tag_end are found then render the tag
		if(tag_end_end > 0 || self_closing) {
			// 
			// get the content of the tag
			tag_content := template[tag_start:tag_end_end]

			// check if the tag is an operation or a component
			// check if the first character of the tag is a capital letter
			if(unicode.IsUpper(rune(tag_name[0]))) {
				result += RenderComponent(tag_name, tag_content, self_closing, json)
			}else{
				result += RenderOperation(tag_name, tag_content, self_closing, json)
			}
			// 
			fmt.Printf("[reset] after %s \n", tag_name)

			// reset the all the tag variables
			tag_name = ""
			tag_start = -1
			tag_end_end = -1
			nested_depth = 0
			self_closing = false
			// 
		}
		// 
		if(i == len(template)-1 && tag_name != "") {
			fmt.Println("[Info]: Could not find the end tag: ", tag_name, " at index: ", i)
			fmt.Println("[Error]: Could not find the end tag for the tag_name index: ", tag_name)
			
			log_end_length := i+7;
			if(log_end_length > len(template)) {
				log_end_length = len(template)
			}

			return "SAD :( \n\n\n\n" +
				""+
				fmt.Sprintf("| tag_name: %s \n",tag_name)  +
				fmt.Sprintf("| %s \n", template[i-10:log_end_length])  +
				"--------------"+
				fmt.Sprintf("| %d\n", 0)  +
				fmt.Sprintf("| %d\n", 0)  +
				fmt.Sprintf("| %d\n", 0)
		}
	} // end of for loop
	return result
}

// Handles the logic for rendering operations e.g. <#if>, <#for>
func RenderOperation(name string, contents string, self_closing bool, json gjson.Result) string {
	//
	fmt.Println("Operation: ", name)
	fmt.Println("Contents: \n", contents)
	fmt.Println("--------------------")
	return "###-OP-###"
}

// Handles the logic for rendering components e.g. <#Button>, <#Input>
func RenderComponent(name string, contents string, self_closing bool, json gjson.Result) string {
	//
	fmt.Println("Component: ", name)
	fmt.Println("Contents: \n", contents)
	fmt.Println("--------------------")
	return "###-COMP-###"
}