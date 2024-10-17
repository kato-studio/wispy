package engine

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func CustomTemplateEngine(html string) string {
	// Regular expression to match custom elements
	re := regexp.MustCompile(`<#[A-z]((.|\n)*?)>`)

	// Process all matches
	for {
		match := re.FindStringSubmatch(html)
		if match == nil {
			break
		}

		element := match[0]
		// attrs := match[1]
		// content := match[2]

		fmt.Println("Element:", element)
		// fmt.Println("Attributes:", attrs)
		// fmt.Println("Content:", content)

		// Call the appropriate render function
		html = strings.Replace(html, match[0], "", 1)
			// If no render function is found, remove the custom element tags
			// html = strings.Replace(html, match[0], content, 1)
	}

	return html
}
 
func SlipEngine(template string, json gjson.Result) string {
	var result string = ""

  s := bufio.NewScanner(strings.NewReader(template))
  s.Split(bufio.ScanWords)
  
	var token string
	var last_token string 
	var current string = ""

	for s.Scan() {
		token = strings.Split(s.Text(), "")[0]
		text := s.Text()

		// 
		if last_token == "<" && token == "#" {
			current += fmt.Sprint(token, text)
			continue
		}else if(last_token != "#") {
			continue
		}
		if(len(current) < 1) { continue }

		// if first character after # is a capital letter
		// then it"s a component
		if(token <= "A" || token >= "Z") {
			current += ""
			fmt.Println(text)		
			// then it"s an internal operation e.i. #if, #for, etc..
		}else if(token <= "a" || token >= "z") {
			fmt.Println(text)
		}			
		// update the previous token
		last_token = token
		
	}
	

	// for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
	// 	t := s.TokenText()
	// 	//
	// 	if(t == "<") { 
	// 		arrow_tag = true 
	// 		continue
	// 	}else if(arrow_tag) { 
	// 		arrow_tag = false 
	// 		continue
	// 	}

	// 	//
	// 	if(arrow_tag && t == "#") {					
	// 		cache += t
	// 		continue
	// 	}else if(arrow_tag && t != "#") {
	// 		cache += t
	// 		continue
	// 	}

	// 	// Is it a component or operation?
	// 	if(cache == "#") {			
	// 		// if first character after # is a capital letter
	// 		// then it"s a component
	// 		if(t <= "A" || t >= "Z") {
	// 			fmt.Printf("%s: %s\n", s.Position, t)

	// 		// then it"s an internal operation e.i. #if, #for, etc..
	// 		}else {
	// 			fmt.Printf("%s: %s\n", s.Position, t)
	// 		}
// 	}

	// 	continue			
	// }

	return result
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