package engine

// import (
// 	"fmt"
// 	"regexp"
// 	"slices"
// 	"strings"
//

// 	"github.com/tidwall/gjson"

// )

// By default, the engine will remove all comments from the template
// todo: add specific comments "option" tag to preserve comments for on render
// func RemoveComments(contents string) string {
// 	// Remove comments
// 	contents = strings.ReplaceAll(contents, "<!--", "<!---")
// 	contents = strings.ReplaceAll(contents, "-->", "--->")
// 	contents = regexp.MustCompile(`<!---.*?--->`).ReplaceAllString(contents, "")
// 	return contents
// }

// // Clean the template by removing comments, newlines, and extra spaces
// func CleanTemplate(contents string) string {
// 	// Remove comments
// 	contents = RemoveComments(contents)
// 	// Remove newlines
// 	contents = strings.ReplaceAll(contents, "\n", "")
// 	contents = strings.ReplaceAll(contents, "\r", "")
// 	// Remove extra spaces
// 	contents = regexp.MustCompile(`\s+`).ReplaceAllString(contents, " ")
// 	return contents
// }

// Return strings
// [String, Number, True, False, Null, JSON, ARRAY]
// func ValueOrTrimmed(raw_value string, json gjson.Result) (result string, value_type string, value_path string) {
// 	var value = strings.TrimSpace(strings.Trim(raw_value, "\""))

// 	// Check if the value is a JSON path
// 	if strings.Contains(value, "{{") {
// 		json_path := strings.TrimSpace(strings.Trim(value, "{}"))
// 		found_value := json.Get(json_path)
// 		if found_value.Exists() {
// 			value_type := found_value.Type.String()
// 			if value_type == "JSON" {
// 				if found_value.IsArray() {
// 					return found_value.String(), "ARRAY", json_path
// 				}
// 				return found_value.String(), "JSON", json_path
// 			} else {
// 				return found_value.String(), value_type, json_path
// 			}
// 		}
// 		// Debug
// 		// fmt.Println("[Error]: Could not find the value for the JSON path: ", json_path)
// 		return "", "Null", json_path
// 	}

// 	// Check if the value is a boolean
// 	if strings.ToLower(value) == "true" {
// 		return value, "True", ""
// 	} else if strings.ToLower(value) == "false" || strings.ToLower(value) == "null" {
// 		return value, "False", ""
// 	}

// 	// Check if the value is a number
// 	if regex_number.MatchString(value) {
// 		return value, "Number", ""
// 	}

// 	// Check for empty string/bull
// 	if value == "" || strings.ToLower(value) == "null" {
// 		return value, "Null", ""
// 	}

// 	return value, "String", ""
// }

// Insert values into the template from json data
// Similar to Go's text/template package (benchmarks would be interesting)
// this is a simple implementation and can be optimized (so is the entire engine :P )
// func InsertValues(contents string, json gjson.Result) string {
// 	matches := regex_insert.FindAllString(contents, -1)
// 	for _, match := range matches {
// 		value, value_type, _ := ValueOrTrimmed(match, json)
// 		if value_type == "String" {
// 			contents = strings.Replace(contents, match, value, -1)
// 		} else {
// 			contents = strings.Replace(contents, match, "", -1)
// 			// Debug
// 			// fmt.Println("[Error]: NO handling for the value type: ", value_type, " for the value: ", value)
// 		}
// 	}
// 	return contents
// }

// validate domain name (for client sites directory)
// func ValidateDomainName(domain string) bool {
// 	if strings.Contains(domain, ".") {
// 		split_domain := strings.Split(domain, ".")
// 		domain_name := split_domain[len(split_domain)-1]
// 		//
// 		if slices.Contains(DOMAINS[:], domain_name) {
// 			return true
// 		} else {
// 			fmt.Println("[Error]: Invalid domain name: ", domain_name)
// 			return false
// 		}
// 	}
// 	fmt.Println("[Error]: Invalid domain: ", domain)
// 	return false
// }
