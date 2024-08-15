package engine

import (
	"strings"

	"github.com/valyala/fasttemplate"
)


func ComponentTemplate(template string) *fasttemplate.Template {
	return fasttemplate.New(template, "<%", "%>")
}

func VariableTemplate(template string) *fasttemplate.Template {
	return fasttemplate.New(template, "{{", "}}")
}

// split string at next given separator and return the two parts
func SplitAt(s, sep string) (string, string) {
	i := strings.Index(s, sep)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+len(sep):]
}

func SplitAtRune(s string, r rune) (string, string) {
	i := strings.IndexRune(s, r)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}