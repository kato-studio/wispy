package extract

import (
	"kato-studio/katoengine/lib/utils"
	"sort"
	"strings"
)

// This as been abstracted to a function
// to make it easier to maintain and test
func ClosestString(content string, target string) int {
	return strings.Index(content, target)
}
func TagType(rawTagStart string) string {
	if rawTagStart[1] == '%' {
		return "operation"
	} else if rawTagStart[1] >= 'A' && rawTagStart[1] <= 'Z' {
		return "component"
	}
	utils.Error("Error: Could not resolve tag type content: \n" + rawTagStart)
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
		nextChar := content[j]
		if nextChar == ' ' || nextChar == '>' {
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

	contentLength := len(content)
	for i := 0; i < contentLength; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i+2 > contentLength {
			break
		}

		nextChar := content[i+1]
		// Look for template operation start tag "<%"
		// Look for component start tag "<"+CapitalLetter
		if char == '<' {
			// If next char is capital letter, then assume it's a component
			if nextChar >= 'A' && nextChar <= 'Z' || nextChar == '%' {
				results = append(results, i)
				// skip the next char
				i++
			}
		}

	}

	sort.Ints(results)

	return results

}

func ComponentContent(component string, name string) (string, string) {
	if component[len(component)-2:] == "/>" {
		return component[len(name)+2 : len(component)-2], ""
	}
	// closingTag := "</" + name + ">"
	closingTag := strings.Index(component, "</"+name+">")
	if closingTag > -1 {
		startTag := "<" + name
		// +1 to include the closing tag character ">"
		endOfStartTag := strings.Index(component, ">")

		// tag options | nested content
		return component[len(startTag):endOfStartTag], component[endOfStartTag+1 : closingTag]
	}

	utils.Error(" -> DoesComponentHaveContent")
	utils.Error("could not resolve component content: \n" + component)

	return "", ""
}

func OperationContent(operation string, name string) (string, string) {
	// closingTag := "</" + name + ">"
	closingTag := strings.Index(operation, "</%"+name+">")
	// +1 to include the closing tag character ">"
	if closingTag > -1 {
		startTag := "<%" + name
		endOfStartTag := strings.Index(operation, ">")

		// tag options | nested content
		return operation[len(startTag):endOfStartTag], operation[endOfStartTag+1 : closingTag]
	}

	utils.Warn(" -> DoesOperationHaveContent")
	utils.Error("could not resolve operation content: \n" + operation)

	return "", ""
}

func ComponentEndTag(content string, name string) int {
	contentLength := len(content)
	for i := 0; i < contentLength; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i+2 > contentLength {
			return -1
		}

		nextChar := content[i+1]
		// Look for component end tag "</"+name+">" or "/>"
		if char == '/' && nextChar == '>' {
			return i + 2
		}
		endTagLength := len(name) + 3
		if char == '<' && nextChar == name[0] && content[i:endTagLength] == "<"+name+"/>" {
			return i + endTagLength
		}
	}

	utils.Error("Error: Could not resolve component end tag content: \n" + content)
	return -1
}

func FindOperationEndTag(content string, name string) int {
	contentLength := len(content)
	nestedTagsFound := 0

	startTag := "<%" + name + ">"
	endTag := "</%" + name + ">"

	for i := 0; i < contentLength; i++ {
		char := content[i]
		// Check if we are at the end of the content
		if i == contentLength-2 {
			utils.Print("Error: i+nameLength+4")
			utils.Error("Could find operation end tag content: \n" + content)
			return -1
		}

		nextChar := content[i+1]
		thirdChar := content[i+2]
		// Look for template operation start tag "<%"
		startTagLength := len(startTag)
		endTagLength := len(endTag)
		if char == '<' && nextChar == '%' && content[i:i+startTagLength] == startTag {
			utils.Debug("nestedTagsFound?!?!? ")
			nestedTagsFound++
			// Known min length of the characters in the nested tag
			i = i + startTagLength
			continue
		}

		if char == '<' && nextChar == '/' && thirdChar == '%' && content[i:i+endTagLength] == endTag {
			if nestedTagsFound == 0 {
				return i + endTagLength
			} else {
				nestedTagsFound--
				i = i + endTagLength
			}
		}
	}

	utils.Print("Error: end of content")
	utils.Error("Could find operation end tag content: \n" + content)
	return -1

}
