// here is my golang definitions file that is referenced to generate a trie of class names and their css properties, let's make small changes to use css variables that exist in the theme where it makes sense to do so ```
package atomicstyle

import (
	"fmt"
	"strings"
	"unicode"
)

type ClassRecipe struct {
	Attribute string
}

func ProcessRecipe(input string) (string, bool) {
	lastChar := input[len(input)-1]

	// Check if number currently all recipes end in a number or rounded bracket
	if unicode.IsDigit(rune(lastChar)) {
		j := strings.LastIndexByte(input, byte('-'))
		start := input[:j+1]
		value := input[j+1:]
		prefix := ""
		exists := false
		var rc ClassRecipe
		if prefix, exists = strings.CutPrefix(input, "text-"); exists {
			rc = ClassRecipe{Attribute: "color: color-mix(in oklab, var(--color-%s) %s%%, transparent)"}
		} else if prefix, exists = strings.CutPrefix(input, "bg-"); exists {
			rc = ClassRecipe{Attribute: "background-color: color-mix(in oklab, var(--color-%s) %s%%, transparent)"}
		}
		if exists {
			if input[len(input)-2] == '/' {
				opacity := input[len(input)-1:]
				return fmt.Sprintf(rc.Attribute, prefix[:len(prefix)-4], opacity), true
			} else if input[len(input)-3] == '/' {
				opacity := input[len(input)-2:]
				return fmt.Sprintf(rc.Attribute, prefix[:len(prefix)-3], opacity), true
			} else {
				return "", false
			}
		}

		// Standard recipes
		if rc, exists = StdRecipes[start]; exists {
			return fmt.Sprintf(rc.Attribute, value), true
		} else {
			return "", false
		}

		// Handle rounded bracket recipes
	} else if rune(lastChar) == ')' {
		j := strings.LastIndexByte(input, byte('-'))
		prefix := input[:j+1]
		value := input[j+1:]
		exists := false
		var rc ClassRecipe

		// Standard recipes
		if rc, exists = BracketRecipes[prefix]; exists {
			return fmt.Sprintf(rc.Attribute, value), true
		} else {
			return "", false
		}
	}
	return "", false
}

var StdRecipes = map[string]ClassRecipe{
	// Width and Height
	"w-":     {Attribute: "width: calc(var(--spacing) * %s);"},
	"h-":     {Attribute: "height: calc(var(--spacing) * %s);"},
	"max-w-": {Attribute: "max-width: calc(var(--spacing) * %s);"},
	"max-h-": {Attribute: "max-height: calc(var(--spacing) * %s);"},
	"min-w-": {Attribute: "min-width: calc(var(--spacing) * %s);"},
	"min-h-": {Attribute: "min-height: calc(var(--spacing) * %s);"},

	// Padding
	"p-":  {Attribute: "padding: calc(var(--spacing) * %s);"},
	"px-": {Attribute: "padding-inline: calc(var(--spacing) * %s)"},
	"py-": {Attribute: "padding-block: calc(var(--spacing) * %s);"},
	"pt-": {Attribute: "padding-top: calc(var(--spacing) * %s);"},
	"pb-": {Attribute: "padding-bottom: calc(var(--spacing) * %s);"},
	"pl-": {Attribute: "padding-left: calc(var(--spacing) * %s);"},
	"pr-": {Attribute: "padding-right: calc(var(--spacing) * %s);"},

	// Margin (with negative value support)
	"m-":  {Attribute: "margin: calc(var(--spacing) * %s);"},
	"mx-": {Attribute: "margin-inline: calc(var(--spacing) * %s);"},
	"my-": {Attribute: "margin-block: calc(var(--spacing) * %s);"},
	"mt-": {Attribute: "margin-top: calc(var(--spacing) * %s);"},
	"mb-": {Attribute: "margin-bottom: calc(var(--spacing) * %s);"},
	"ml-": {Attribute: "margin-left: calc(var(--spacing) * %s);"},
	"mr-": {Attribute: "margin-right: calc(var(--spacing) * %s);"},
	//
	"-m-":  {Attribute: "margin: calc(var(--spacing) * %s * -1);"},
	"-mt-": {Attribute: "margin-top: calc(var(--spacing) * %s * -1);"},
	"-mb-": {Attribute: "margin-bottom: calc(var(--spacing) * %s * -1);"},
	"-ml-": {Attribute: "margin-left: calc(var(--spacing) * %s * -1);"},
	"-mr-": {Attribute: "margin-right: calc(var(--spacing) * %s * -1);"},

	// Grid
	"grid-":  {Attribute: "display: grid;"},
	"cols-":  {Attribute: "grid-template-columns: repeat(%s, minmax(0, 1fr));"},
	"rows-":  {Attribute: "grid-template-rows: repeat(%s, minmax(0, 1fr));"},
	"gap-":   {Attribute: "gap: calc(var(--spacing) * %s);"},
	"gap-x-": {Attribute: "column-gap: calc(var(--spacing) * %s);"},
	"gap-y-": {Attribute: "row-gap: calc(var(--spacing) * %s);"},

	// Absolute
	"top-":     {Attribute: "top: calc(var(--spacing) * %s);"},
	"bottom-":  {Attribute: "bottom: calc(var(--spacing) * %s);"},
	"left-":    {Attribute: "left: calc(var(--spacing) * %s);"},
	"right-":   {Attribute: "right: calc(var(--spacing) * %s);"},
	"-top-":    {Attribute: "top: calc(var(--spacing) * %s * -1);"},
	"-bottom-": {Attribute: "bottom: calc(var(--spacing) * %s * -1);"},
	"-left-":   {Attribute: "left: calc(var(--spacing) * %s * -1);"},
	"-right-":  {Attribute: "right: calc(var(--spacing) * %s * -1);"},

	// TODO:  Add support for percentages
	// TODO:  Add css variable using rounded brackets I.E. "text-(--custom-color)" = "color: var(--custom-color)"
}

var BracketRecipes = map[string]ClassRecipe{

	// Other
	"bg-url": {Attribute: "background: url(%s)"},
}
