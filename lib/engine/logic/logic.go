package logic

import (
	"kato-studio/katoengine/lib/utils"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func InsertData(content string, data gjson.Result) string {
	result := content
	contentLen := len(content)
	currentTag := ""

	for i := 2; i < contentLen; i++ {
		//
		if content[i] == '}' {
			result = strings.Replace(result, "{%"+currentTag+"}", data.Get(currentTag).String(), 1)
			continue
		}
		//
		if content[i-2] == '{' && content[i-1] == '%' || len(currentTag) > 0 {
			currentTag += string(content[i])
		}
	}

	return result
}

func ConvertFromStringToType(value string) interface{} {
	// Try to convert to boolean
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	// Try to convert to integer
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	// Try to convert to float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	// If all conversions fail, return the original string
	return value
}

// Handles all logic operations from templates and returns a boolean
func HandlerOperation(operation string) bool {
	opParts := strings.Split(operation, " ")
	utils.Debug("~> operation")
	utils.Debug(operation)
	utils.Print("opParts -> ")
	utils.Print(opParts)
	//
	if len(opParts) != 3 {
		if strings.ToLower(operation) == "true" {
			return true
		} else {
			return false
		}
	}

	opType := opParts[1]
	opLeft := ConvertFromStringToType(opParts[0])
	opRight := ConvertFromStringToType(opParts[2])

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
	operationSplitters := []string{"&&", "||"}
	for _, splitter := range operationSplitters {
		splitOps := strings.Split(logic, splitter)
		for _, op := range splitOps {
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
