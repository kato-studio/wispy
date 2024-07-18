package template

import (
	"kato-studio/katoengine/lib/utils"
)

func OperationParser(content string, tag string, context_data interface{}) string {

	switch tag {
	case "if":
		content = IfOperation(content, context_data)
	case "each":
		content = EachOperation(content, context_data)
	}

	return content
}

/*
<%each {%data.clients} | client>

	<p>{%client}</p>

</%each>
*/
func EachOperation(content string, data interface{}) string {
	utils.Debug("Each Operation")
	return "{Each Operation}"
}

/*
<%if {%data.is_logged_in}>

	<p>User is logged in</p>

</%if>
*/
func IfOperation(content string, data interface{}) string {
	utils.Debug("If Operation")
	return "{If Operation}"
}
