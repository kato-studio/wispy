package extract

import (
	"strings"
)

// This as been abstracted to a function
// to make it easier to maintain and test
func ClosestString(content string, target string) int {
	return strings.Index(content, target)
}

func Name(content string) string {
	tag := ""

	for j := 0; j < 65; j++ {
		// Check if we are at the end of the content
		if j == 65 {
			panic("Error: could not resolve component name content: \n" + content[:64])
		}

		// Check if we are at the end of the component name
		nextChar := content[j]
		if nextChar == ' ' || nextChar == '>' || nextChar == '\n' {
			tag = content[:j]
			break
		}
	}

	return tag
}
func ComponentName(content string) string {
	return Name(content[1:])
}
func OperationName(content string) string {
	return Name(content[2:])
}

type ContentScannerResult struct {
	Components  []int
	TemplateOps []int
}

// This function scans the content for components and template operations
// It returns the indexes of the start of the component and template operation
func ContentScanner(content string) ContentScannerResult {
	result := ContentScannerResult{
		Components:  []int{},
		TemplateOps: []int{},
	}

	contentLength := len(content)
	for i := 0; i < contentLength; i++ {
		char := content[i]

		// Check if we are at the end of the content
		if i+2 > contentLength {
			break
		}

		// Look for template operation start tag "<%"
		if char == '<' && content[i+1] == '%' {
			result.TemplateOps = append(result.TemplateOps, i)
		}

		// Look for component start tag "<"+CapitalLetter
		if char == '<' {
			nextChar := content[i+1]
			// If next char is capital letter, then assume it's a component
			if nextChar >= 'A' && nextChar <= 'Z' {
				result.Components = append(result.Components, i)
			}
		}
	}

	return result
}

func FindOperationEndTag(content string, tag string) (int, int) {
	startString := "<%" + tag
	endString := "</%" + tag + ">"

	tagLength := len(startString)
	endTagLength := len(endString + startString)
	firstGuess := strings.Index(content[tagLength:], endString)
	nextTagOfSameType := strings.Index(content[tagLength:], startString)

	if firstGuess == -1 {
		panic("Error: could not resolve component name content: \n" + content)
	}

	// If there is no other tag of the same type, then return the first guess
	// or if the next tag of the same type is after the first guess
	if nextTagOfSameType == -1 || nextTagOfSameType > firstGuess {
		return firstGuess + endTagLength, 0
	}

	return FindOperationEndTag(content[nextTagOfSameType+5:], tag)

	// EHHHH --------
	// // If there is another tag of the same type, then find the end tag of that tag
	// // and continue until there are no more tags of the same type
	// result, nestedTagCount := FindOperationEndTag(content[nextTagOfSameType+5:], tag)

	// // nestedCount represents the number of nested operations of the same type
	// return result + endTagLength, nestedTagCount + 1
}
