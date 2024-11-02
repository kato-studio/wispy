package engine

import (
	"fmt"
	"os"
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

// -----===================-----
func fileClosure(dir_path string, contents string, ctx RenderCTX) error {
	renderedContents := RenderPage(contents, dir_path, ctx)
	//
	if renderedContents == "" {
		fmt.Println("[warn]: Could not render the page: ", dir_path)
		return nil
	}
	output_dir := strings.Replace(dir_path, "./sites", "./static_sites", 1)
	output_dir = strings.Replace(output_dir, "+page.hstm", "", 1)
	//
	output_path := strings.Replace(dir_path, "./sites", "./static_sites", 1)
	output_path = strings.Replace(output_path, "+page.hstm", "index.html", 1)
	//
	dir_err := os.MkdirAll(output_dir, 0755)
	err := os.WriteFile(output_path, []byte(renderedContents), 0644)
	if err != nil {
		fmt.Println("[Error]: Could not write the file: ", output_path)
		if dir_err != nil {
			fmt.Println("[Error]: Could not create the directory: ", output_dir)
		} else {
			fmt.Println("[Error]: Could not write the file: ", output_path)
		}
		fmt.Println("----------")
		return err
	}
	return nil
}

func dirClosure(dir_path string, ctx RenderCTX) error {
	folder_items, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println("[Error]: Could not read the directory: ", dir_path)
		return err
	}
	for _, item := range folder_items {
		this_path := dir_path + "/" + item.Name()
		if item.IsDir() {
			dirClosure(this_path, ctx)
		} else {
			if item.Name() != "+page.hstm" {
				continue
			}
			contentBytes, err := os.ReadFile(this_path)
			if err != nil {
				fmt.Println("[Error]: Could not read the file: ", this_path)
				return err
			}
			fileClosure(this_path, string(contentBytes), ctx)
		}
	}
	return nil
}
func RenderAllSites(sitesDir string, ctx RenderCTX) error {
	files, err := os.ReadDir(sitesDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		dirClosure(sitesDir+"/"+file.Name()+"/pages", ctx)
	}

	return nil
}
