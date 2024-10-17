package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"

)

func FindEndTag(content string, tagName string) string {
	// Regular expression to match custom elements
	re := regexp.MustCompile(`</#[A-z]((.|\n)*?)>`)

	// Process all matches
	for {
		match := re.FindStringSubmatch(content)
		if match == nil { break }

		element := match[0]

		fmt.Println("Element:", element)
		// fmt.Println("Attributes:", attrs)
		// fmt.Println("Content:", content)

		// Call the appropriate render function
		content = strings.Replace(content, match[0], "", 1)
		// If no render function is found, remove the custom element tags
		// html = strings.Replace(html, match[0], content, 1)
	}

	return content 
}
 
func SlipEngine(template string, json gjson.Result) string {
		// Regular expression to match custom elements
	re := regexp.MustCompile(`<#[A-z]((.|\n)*?)>`)

	// Process all matches
	for {
		match := re.FindStringSubmatch(template)
		if match == nil { break }

		element := match[0]

		fmt.Println("Element:", element)
		// fmt.Println("Attributes:", attrs)
		// fmt.Println("Content:", content)

		// Call the appropriate render function
		template = strings.Replace(template, match[0], "", 1)
		// If no render function is found, remove the custom element tags
		// html = strings.Replace(html, match[0], content, 1)
		template = FindEndTag(template, element)
	}

	return template
}

func RenderComponent(component string, attributes string, children string, json gjson.Result) string {
	// 
	return SlipEngine(component, json)
}

// Regex example
	// Get start tag of operations & components in the template
	// - `<#[A-z]((.|\n)*?)>\n` basically <#ABC ... >\n
	// - we can tell if it"s a component if it has a capital letter after the # sign
	// - otherwise it"s an operation then we can check if the last character is "/>"
	// - to know if it"s a self-closing tag or not
	// fmt.Println("template: ", template)
	// temp_gex := regexp.MustCompile(`(?m)<#[A-z](.+?)>`) //regexp.Compile("p([a-z]+)ch")
	// for i, match := range temp_gex.FindAllString(template, -1) {
  //     fmt.Println(match, "found at index", i)
  // }