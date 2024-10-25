package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

var regex_number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+$`)
var regex_insert = regexp.MustCompile(`{{.*?}}`)

// Return strings
// [String, Number, True, False, Null, JSON, ARRAY]
func ValueOrTrimmed(raw_value string, json gjson.Result) (result string, value_type string, value_path string) {
	var value = strings.TrimSpace(strings.Trim(raw_value, "\""))

	// Check if the value is a JSON path
	if strings.Contains(value[:4], "{{") {
		json_path := strings.TrimSpace(strings.Trim(value, "{}"))
		found_value := json.Get(json_path)
		if found_value.Exists() {
			value_type := found_value.Type.String()
			if value_type == "JSON" {
				if found_value.IsArray() {
					return found_value.String(), "ARRAY", json_path
				}
				return found_value.String(), "JSON", json_path
			} else {
				return found_value.String(), value_type, json_path
			}
		}
		fmt.Println("[Error]: Could not find the value for the JSON path: ", json_path)
		return "", "Null", json_path
	}

	// Check if the value is a boolean
	if strings.ToLower(value) == "true" {
		return value, "True", ""
	} else if strings.ToLower(value) == "false" || strings.ToLower(value) == "null" {
		return value, "False", ""
	}

	// Check if the value is a number
	if regex_number.MatchString(value) {
		return value, "Number", ""
	}

	// Check for empty string/bull
	if value == "" || strings.ToLower(value) == "null" {
		return value, "Null", ""
	}

	return value, "String", ""
}

func InsertValues(contents string, json gjson.Result) string {
	matches := regex_insert.FindAllString(contents, -1)
	for _, match := range matches {
		value, value_type, _ := ValueOrTrimmed(match, json)
		if value_type == "String" {
			contents = strings.Replace(contents, match, value, -1)
		} else {
			contents = strings.Replace(contents, match, "", -1)
			fmt.Println("[Error]: NO handling for the value type: ", value_type, " for the value: ", value)
		}
	}
	return contents
}
