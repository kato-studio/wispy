package logic

import (
	"kato-studio/katoengine/lib/utils"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

)

func InsertData(content string, data gjson.Result) string {
	result := content
	content_len := len(content)
	current_tag := ""

	for i := 2; i < content_len; i++ {
		//
		if content[i] == '}' {
			result = strings.Replace(result, "{%"+current_tag+"}", data.Get(current_tag).String(), -1)
			continue
		}
		//
		if content[i-2] == '{' && content[i-1] == '%' || len(current_tag) > 0 {
			current_tag += string(content[i])
		}
	}

	return result
}

func ConvertFromStringToType(value string) interface{} {
	// Try to convert to boolean
	if bool_val, err := strconv.ParseBool(value); err == nil {
		return bool_val
	}
	// Try to convert to integer
	if int_val, err := strconv.Atoi(value); err == nil {
		return int_val
	}
	// Try to convert to float
	if float_val, err := strconv.ParseFloat(value, 64); err == nil {
		return float_val
	}
	// If all conversions fail, return the original string
	return value
}

// Handles all logic operations from templates and returns a boolean
func HandlerOperation(operation string) bool {
	op_parts := strings.Split(operation, " ")
	utils.Debug("~> operation")
	utils.Debug(operation)
	utils.Print("opParts -> ")
	utils.Print(op_parts)
	//
	if len(op_parts) != 3 {
		if strings.ToLower(operation) == "true" {
			return true
		} else {
			return false
		}
	}

	// breaking naming convention for better readability
	opType := op_parts[1]
	opLeft := ConvertFromStringToType(op_parts[0])
	opRight := ConvertFromStringToType(op_parts[2])

	switch opType {
	case ">":
		switch left := opLeft.(type) {
		case int:
			if right, ok := opRight.(int); ok {
				return left > right
			}
		case float64:
			if right, ok := opRight.(float64); ok {
				return left > right
			}
		}
	case "<":
		switch left := opLeft.(type) {
		case int:
			if right, ok := opRight.(int); ok {
				return left < right
			}
		case float64:
			if right, ok := opRight.(float64); ok {
				return left < right
			}
		}
	case ">=":
		switch left := opLeft.(type) {
		case int:
			if right, ok := opRight.(int); ok {
				return left >= right
			}
		case float64:
			if right, ok := opRight.(float64); ok {
				return left >= right
			}
		}
	case "<=":
		switch left := opLeft.(type) {
		case int:
			if right, ok := opRight.(int); ok {
				return left <= right
			}
		case float64:
			if right, ok := opRight.(float64); ok {
				return left <= right
			}
		}
	case "==":
		return opLeft == opRight
	case "!=":
		return opLeft != opRight
	}

	return false
}

// Examples all examples are strings and should be converted to boolean
// A: len("Hello") == 5
// B: 12 + 3 == 123
// C: "true" == true
func StringToBoolean(_logic string) bool {
	// remove double all spaces & clean string
	logic := strings.Trim(strings.ReplaceAll(_logic, "  ", " "), " ")
	// handle string to boolean conversion
	// incase of string return true if string is not empty
	lowercase := strings.ToLower(logic)
	if lowercase == "true" {
		return true
	}
	if lowercase == "false" || lowercase == "" || lowercase == " " {
		return false
	}

	// otherwise return start handling logic operations
	operation_splitters := []string{"&&", "||"}
	for _, splitter := range operation_splitters {
		split_ops := strings.Split(logic, splitter)
		for _, op := range split_ops {
			result := HandlerOperation(op)
			if !result {
				return false
			}
		}
	}
	results := HandlerOperation(logic)
	utils.Debug("~> results")
	utils.Print(results)

	return results
}
