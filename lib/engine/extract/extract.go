package extract

import (
	"fmt"
	"kato-studio/katoengine/lib/utils"
	"sort"
	"strings"

)

// This as been abstracted to a function
// to make it easier to maintain and test
func ClosestString(content string, target string) int {
	return strings.Index(content, target)
}

// ---
// import Header from '@/Header'
// import Footer from '@/Footer'
// 
// %git_users = Fetch('GET:https://api.github.com/users','Content-Type: application/json','{"foo":"bar"}')
// ---
func ServerLogic(content string) (imports_strings []string, server_funcs []string, remaining_content string)  {
	imports_strings = []string{}
	server_funcs = []string{}
	remaining_content = ""
	if(len(content) > 6 && content[:3] == "---"){
		result := content[3:len(content)-6]
		import_parts := strings.Split(result, "import ")
		for _, import_string := range import_parts {
			if strings.Contains(import_string, " from ") {
				import_split := strings.Split(import_string, " from ")
				import_string := strings.Join(import_split,":")
				end_index := ClosestString(import_string, " ")
				end_of_logic_section := ClosestString(import_string, "---")
				if (end_index > 0) {
					server_func_strings := ""
					if(end_of_logic_section > -1){
						server_func_strings = import_string[end_index+1:end_of_logic_section]
						for _, server_func_string := range strings.Split(server_func_strings, ")") {
							// 2 is an arbitrary number but it works
							if(len(server_func_string) > 2){
								server_funcs = append(server_funcs, server_func_string+")")
							}
						}
					}
					// 
					imports_strings = append(imports_strings, import_string[:end_index])
					// 
					remaining_content = import_string[end_of_logic_section+3:]
					// 
				}else{
					imports_strings = append(imports_strings, import_string)
				}
			}
		}
	}else{
		remaining_content = content
	}
	return imports_strings, server_funcs, remaining_content
}

func TagType(raw_tag_start string) string {
	if raw_tag_start[1] == '%' {
		return "operation"
	} else if raw_tag_start[1] >= 'A' && raw_tag_start[1] <= 'Z' {
		return "component"
	}
	utils.Error("Error: Could not resolve tag type content: \n" + raw_tag_start)
	utils.Error("returning unknown from extract.TagType()")
	return "unknown"
}

func TagName(content string) string {
	tag := ""

	for j := 0; j < 65; j++ {
		// Check if we are at the end of the content
		if j == 65 {
			utils.Error("Error: Could not resolve component name content: \n" + content[:64])
			break
		}

		// Check if we are at the end of the component name
		next_char := content[j]
		if next_char == ' ' || next_char == '>' {
			tag = content[:j]
			break
		}

	}

	return tag
}

func ComponentName(content string) string {
	return TagName(content[1:])
}

func OperationName(content string) string {
	return TagName(content[2:])
}

// This function scans the content for components and template operations
// It returns the indexes of the start of the component and template operation
func ContentScanner(content string) []int {
	results := []int{}
	content_length := len(content)
	for i := 0; i < content_length; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i+2 > content_length {
			break
		}

		next_char := content[i+1]
		// Look for template operation start tag "<%"
		// Look for component start tag "<"+CapitalLetter
		if char == '<' {
			// If next char is capital letter, then assume it's a component
			if next_char >= 'A' && next_char <= 'Z' || next_char == '%' {
				results = append(results, i)
				// skip the next char
				i++
			}
		}
	}
	sort.Ints(results)

	return results

}

func ComponentContent(component string, name string) (options_string string, slot_content string) {
	if component[len(component)-2:] == "/>" {
		return component[len(name)+2 : len(component)-2], ""
	}
	// closing_tag := "</" + name + ">"
	closing_tag := ClosestString(component, "</"+name+">")
	if closing_tag > -1 {
		start_tag := "<" + name
		// +1 to include the closing tag character ">"
		end_of_start_tag := ClosestString(component, ">")

		// tag options | nested content
		return component[len(start_tag):end_of_start_tag], component[end_of_start_tag+1 : closing_tag]
	}

	utils.Error(" -> DoesComponentHaveContent")
	utils.Error("could not resolve component content: \n" + component)

	return "", ""
}

func OperationContent(operation string, name string) (tag_options string, nested_content string) {
	// closingTag := "</" + name + ">"
	closing_tag := ClosestString(operation, "</%"+name+">")
	// +1 to include the closing tag character ">"
	if closing_tag > -1 {
		start_tag := "<%" + name
		end_of_start_tag := ClosestString(operation, ">")

		// tag options | nested content
		return operation[len(start_tag):end_of_start_tag], operation[end_of_start_tag+1 : closing_tag]
	}

	utils.Warn("-> DoesOperationHaveContent")
	fmt.Print(operation)
	utils.Error("could not resolve operation content: \n" + operation)

	return "", ""
}

func ComponentEndTag(content string, name string) int {
	utils.Print("ComponentEndTag")
	utils.Print(content)
	utils.Print("--------------")

	start_index := len(name) + 2
	content_length := len(content)
	tag_start_flag := false
	for i := start_index; i < content_length; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i+2 > content_length {
			utils.Error("Error: Could not resolve component end tag content: \n" + content)
			return content_length
		}

		next_char := content[i+1]

		// Look for component end tag "</"+name+">" or "/>"
		if char == '/' && next_char == '>' && !tag_start_flag {
			return i + 2
		}

		// check is second char a letter to avoid false positives
		char_is_letter := next_char >= 'A' && next_char <= 'Z' || next_char >= 'a' && next_char <= 'z'
		if char == '<' && char_is_letter && !tag_start_flag {
			utils.Print("tag_start_flag = true")
			utils.Print("because: "+content[i:i+4])

			tag_start_flag = true
		}
		
		end_tag_length := len(name) + 3
		if char == '<' && next_char == '/' && content[i:i+end_tag_length] == "</"+name+">" {
		utils.Print("test:"+content[i:i+end_tag_length])
			return i + end_tag_length
		}
	}
	
	utils.Error("Error: Could not resolve component end tag content: \n" + content)
	return +1
}

func ImportsMap(imports_string []string) map[string]string {
	imports_map := map[string]string{}	
	for _, imp := range imports_string {
		split := strings.Split(imp, ":")
		// remove any wrapping quotes
		path := strings.Trim(split[1], "'")
		path = strings.Trim(path, "\"")
		imports_map[split[0]] = path
	}

	return imports_map
}

func FindOperationEndTag(content string, name string) int {
	content_length := len(content)
	nested_tags_found := 0

	start_tag := "<%" + name + ">"
	end_tag := "</%" + name + ">"

	for i := 0; i < content_length; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i == content_length-2 {
			utils.Print("Error: i+nameLength+4")
			utils.Error("Could find operation end tag content: \n" + content)
			return -1
		}

		next_char := content[i+1]
		third_char := content[i+2]
		// Look for template operation start tag "<%"
		start_tag_length := len(start_tag)
		end_tag_length := len(end_tag)
		if char == '<' && next_char == '%' && content[i:i+start_tag_length] == start_tag {
			utils.Debug("nestedTagsFound?!?!? ")
			nested_tags_found++
			// Known min length of the characters in the nested tag
			i = i + start_tag_length
			continue
		}

		if char == '<' && next_char == '/' && third_char == '%' && content[i:i+end_tag_length] == end_tag {
			if nested_tags_found == 0 {
				return i + end_tag_length
			} else {
				nested_tags_found--
				i = i + end_tag_length
			}
		}
	}

	utils.Print("Error: end of content")
	utils.Error("Could find operation end tag content: \n" + content)
	return -1

}
