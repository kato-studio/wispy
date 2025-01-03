package engine

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// LayoutPageInsert inserts the appropriate content into the layout based on the return_value.
func LayoutPageInsert(return_value string, ctx *TemplateCtx, page_content *bytes.Buffer) string {
	switch strings.ToLower(return_value) {
	case "title":
		return ctx.Page.Title
	case "content":
		return page_content.String()
	case "head":
		return ctx.Page.Head
	case "meta":
		return ctx.Page.Meta
	case "styles":
		return ctx.Page.Styles
	case "scripts":
		return ctx.Page.Scripts
	case "lang":
		return ctx.Page.Lang
	default:
		return ""
	}
}

// ParseAttributes parses the attributes of an HTML tag.
func ParseAttributes(s *Scanner, attributes *map[string]string, _scan func() rune, key string) {
	if s.Peek() == '=' {
		_scan()
		_scan()
		(*attributes)[key] = s.TokenText()
	} else if len(key) > 1 {
		(*attributes)[key] = "true"
	}
}

// Render represents the rendering context and functions.
type Render struct {
	AttrFuncMap      map[string]AttributeFunc
	OperationFuncMap map[string]OperationFunc
	Ctx              *TemplateCtx
}

func InitEngine(ctx *TemplateCtx) Render {
	var r Render
	r.AttrFuncMap = make(map[string]AttributeFunc)
	r.OperationFuncMap = make(map[string]OperationFunc)
	r.Ctx = ctx
	return r
}

// OperationFunc defines the function signature for operations.
type OperationFunc func(r *Render, values ...string) string

// SetOperationFunc sets an operation function for the render context.
func (r *Render) SetOperationFunc(operation_name string, handler OperationFunc) *Render {
	r.OperationFuncMap[operation_name] = handler
	return r
}

// AttributeFunc defines the function signature for attribute functions.
type AttributeFunc func(value string) string

// SetAttributeFunc sets an attribute function for the render context.
func (r *Render) SetAttributeFunc(attribute_name string, handler AttributeFunc) *Render {
	r.AttrFuncMap[attribute_name] = handler
	return r
}

func (r *Render) Html(template_string string) (string, error) {
	// root variables
	var result strings.Builder
	var char rune
	var length = len(template_string)
	var err error
	var index = 0

	// flags
	var in_quotes bool
	// -------
	var tag_name = ""
	var tag_attrs = map[string]string{}
	var contents_buffer strings.Builder
	var nested = 0
	//
	// Start core loop
	for ; index <= length-1; index++ {
		if err != nil {
			fmt.Println("Error: ", err)
			err = nil
		}

		char = rune(template_string[index])
		if in_quotes {
			contents_buffer.WriteRune(char)
		} else {
			switch char {
			// ROOT CASE
			// case '"':
			// 	if in_quotes {
			// 		in_quotes = false
			// 		contents := contents_buffer.String()
			// 		fmt.Println("CONTENTS: ")
			// 		contents_buffer.Reset()
			// 		result.WriteString(contents)
			// 	} else {
			// 		in_quotes = true
			// 	}
			// 	in_quotes = true
			// ROOT CASE
			case '<':
				// Handle comments
				// ~~~~~~~~~~~~~
				if index+3 <= len(template_string) && template_string[index:index+3] == "<!-" {
					comment_end_index := strings.Index(template_string[index:], "-->")
					if comment_end_index == -1 {
						err = fmt.Errorf("Comment not closed index[", index, "]")
						break
					}
					// Update index to skip the comment content and the closing tag
					index += comment_end_index + 2 // we only need to increment by 2 because we increment by 1 at the end of the loop
					break
				}

				// Handle tags
				// ~~~~~~~~~~~~~
				if length < index+2 {
					err = fmt.Errorf("Tag not closed index[", index, "]\n -10| ", template_string[index-10:index], "\n +10| ", template_string[index:index+10])
					break
				}

				peek := template_string[index+1]
				peek_next := template_string[index+2]
				if peek == 'x' && peek_next == ':' || peek >= 'A' && peek <= 'Z' {
					break_index := strings.IndexRune(template_string[index:], ' ')
					tag_name = template_string[index+1 : index+break_index]
					// Print tag name with color
					fmt.Printf("\033[1;34mTAG NAME: [%s]\033[0m\n", tag_name)

					// Handle tag attributes
					// ~~~~~~~~~~~~~
					attr_buffer := strings.Builder{}
					for ; index <= length-1; index++ {
						char = rune(template_string[index])
						if char == '>' {
							break
						}
						if char == ' ' {
							tag_attrs[attr_buffer.String()] = ""
							attr_buffer.Reset()
						} else {
							attr_buffer.WriteRune(char)
						}
					}
					fmt.Println("-----ATTRS: ", fmt.Sprint(tag_attrs))

					// Handle tag contents
					// ~~~~~~~~~~~~~
					tested_index := index
				checkClosingTag:
					test_end_index := strings.Index(template_string[tested_index:], tag_name)
					if test_end_index == -1 {
						err = fmt.Errorf("Tag not closed index[", index, "]")
						break
					}
					pos := tested_index + test_end_index
					if template_string[pos-2] == '>' && template_string[pos-1] == '/' {
						if nested > 0 {
							nested--
						} else {
							new_pos := pos + len(tag_name) + 1
							contents_buffer.Reset()
							contents_buffer.WriteString(template_string[index:new_pos])
							r.RenderComponent(tag_name, contents_buffer.String(), &tag_attrs)
							contents_buffer.Reset()
							tag_attrs = map[string]string{}
							index = new_pos
						}
					} else {
						tested_index += len(tag_name) + 2
						nested++
						goto checkClosingTag
					}
					// ~~~~~~~~~~~~~
					// if is normal tag write contents
				} else {
					result.WriteRune(char)
				}

			default:
				if in_quotes {
					contents_buffer.WriteRune(char)
				} else {
					result.WriteRune(char)
				}
			}
		}

		/*
				switch char {
				// ROOT CASE
				case '"':
					if in_tag_attrs {
						in_quotes = !in_quotes
					}
				// ROOT CASE
				case '<':
					// peek := template_string[index+1]
					// //
					// //  Handle nested tags
					// if (is_component || is_operation) && !in_tag_name && !in_tag_attrs {
					// 	if peek == '/' {
					// 		test := template_string[index : index+len(tag_name)+2]
					// 		if test == "</"+tag_name {
					// 			if nested == 0 {
					// 				tag_name = ""
					// 				tag_attrs = []string{}
					// 				tag_contents.Reset()
					// 				fmt.Println("FINISHED \n --- tag_name: ", tag_name, fmt.Sprint(tag_attrs))
					// 			} else {
					// 				nested--
					// 			}
					// 		}
					// 	} else {
					// 		test := template_string[index : index+len(tag_name)+1]
					// 		if test == "<"+tag_name {
					// 			fmt.Println("NESTED!++", tag_name)
					// 			nested++
					// 		}
					// 	}
					// 	normalWrite()
					// }

					// test := template_string[index : index+3]
					// is_component = peek >= 'A' && peek <= 'Z'
					// if test == "<x:" {
					// 	is_operation = true
					// 	in_tag_name = true
					// } else if is_component {
					// 	fmt.Println("COMPONENT! START FOR NAME!")
					// 	in_tag_name = true
					// }

					// // ~~~~~~~~~~~~~
					// if test == "<!-" {
					// 	comment_end := strings.Index(template_string[index:], "-->")
					// 	if comment_end == -1 {
					// 		err = fmt.Errorf("Comment not closed index[", index, "]")
					// 		break
					// 	}
					// 	// Update index to skip the comment content and the closing tag
					// 	index += comment_end + 2 // we only need to increment by 2 because we increment by 1 at the end of the loop
					// }

					// if is_operation || is_component {
					// 	index++
					// 	goto redo
					// } else {
					// 	normalWrite()
					// }

				// ROOT CASE
				case ' ':
					// -----------
					//! Note to self
					//! looks like the tag name is not being recorded this is obv an issue :P
					//! might just be capturing whitespace
					// -----------
					// if is_component || is_operation {
					// 	if in_quotes {
					// 		tag_contents.WriteRune(char)
					// 	} else //
					// 	if in_tag_name {
					// 		tag_name = tag_contents.String()
					// 		tag_contents.Reset()
					// 		in_tag_name = false
					// 		fmt.Println("tag_name: ", tag_name)
					// 	} else //
					// 	if in_tag_attrs && tag_contents.Len() > 0 {
					// 		tag_attrs = append(tag_attrs, tag_contents.String())
					// 		tag_contents.Reset()
					// 	}
					// }
					// normalWrite()
				// ROOT CASE
				case '>':
					// if is_component || is_operation {
					// 	if in_tag_name {
					// 		tag_name = tag_contents.String()
					// 		tag_contents.Reset()
					// 		in_tag_name = false
					// 	}
					// 	if in_tag_attrs {
					// 		tag_attrs = append(tag_attrs, tag_contents.String())
					// 		fmt.Println("tag_attrs: ", tag_attrs)
					// 		tag_contents.Reset()
					// 		in_tag_attrs = false
					// 	}
					// } else {
					// 	normalWrite()
					// }
				case '{':
					// peek := template_string[index+1]
					// if peek == '{' {
					// 	template_end := strings.Index(template_string[index:], "}}")
					// 	if template_end == -1 {
					// 		err = fmt.Errorf("Template Tag not closed index[", index, "]")
					// 		break
					// 	}
					// 	fmt.Println("template: ", template_string[index+2:index+template_end])
					// 	index += template_end + 1 // we only need to increment by 1 because we increment by 1 at the end of the loop
					// } else {
					// 	normalWrite()
					// }
				// ROOT DEFAULT
				default:
					normalWrite()
				}
			}
		*/

	}
	return result.String(), err
}

// RenderComponent renders a component with the given name, contents, and attributes.
func (r *Render) RenderComponent(name, contents string, attributes *map[string]string) {
	/* example input that has been parsed
	<TitleThing title="I Am A Title! [Two]">
		lorem ipsum dolor sit amet consectetur adipiscing elit
	</TitleThing>
	*/
	fmt.Println("COMPS: ")
	fmt.Println(name, contents)
	fmt.Println("===========")
}

// RenderComponent renders a operation with the given name, contents, and attributes.
func (r *Render) RenderOperation(name, contents string, attributes *map[string]string) {
	/* example input that has been parsed
	<x:each link in .list>
	  boop {{link.index}}
	</x:each>
	*/
	fmt.Println("OP: ")
	fmt.Println(name, contents)
	fmt.Println("===========")
}

// RenderTemplate renders a template with the given contents.
func (r *Render) RenderTemplate(contents string) {
	/* example input that has been parsed
	{{.site.name}}
	*/
	fmt.Println("Template: ")
	fmt.Println(contents)
	fmt.Println("===========")
}

// RenderPage renders a page for the given domain and page path.
func (r *Render) RenderPage(domain, page_path string) (string, error) {
	// Read File
	domain = strings.Trim(domain, "/")
	var path = ROOT_DIR + domain + "/pages/" + PAGE_FILE
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	//
	result, err := r.Html(string(file))
	// TODO: better error handling for dev mode
	// TODO: next js style error page
	return result, err
}

// // Html renders the HTML template string.
// func (r *Render) Html(template_string string) (string, error) {
// 	var result strings.Builder
// 	var s Scanner
// // 	var err error

// 	s.Init(strings.NewReader(template_string))
// 	var output strings.Builder
// 	var txt string

// 	for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 		// Used to remove tabs line breaks and chained spaces
// 		if tok == ' ' && s.Peek() == ' ' || tok == '\n' || tok == '\u0009' {
// 			continue
// 		}
// 		var name = ""
// 		var attributes = map[string]string{}
// 		var nested_lvl = 0
// 		var attributes_complete = false
// 		var contents strings.Builder
// 		switch true {
// 		// ==------------------==
// 		// Operation Case
// 		// ==------------------==
// 		case tok == 'x' && s.Peek() == ':':
// 			var operationName strings.Builder
// 			var operationArgs []string
// 			var inArgs bool
// 			// Parse the operation name and arguments
// 			for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 				if tok == '>' {
// 					break
// 				}
// 				if tok == ' ' {
// 					inArgs = true
// 				}
// 				txt = s.TokenText()
// 				if inArgs {
// 					operationArgs = append(operationArgs, txt)
// 				} else {
// 					operationName.WriteString(txt)
// 				}
// 			}
// 			// Execute the operation function if it exists
// 			if handler, exists := r.OperationFuncMap[operationName.String()]; exists {
// 				output.WriteString(handler(r, operationArgs...))
// 			}
// 			// Skip the closing tag
// 			for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 				if tok == '<' && s.Peek() == '/' {
// 					for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 						if tok == '>' {
// 							break
// 						}
// 					}
// 					break
// 				}
// 			}
// 		// ==------------------==
// 		// Comments Case
// 		// Skip comments "<!--" to "-->"
// 		// ==------------------==
// 		case tok == '<' && s.Peek() == '!':
// 			for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 				if tok == '-' && s.Peek() == '-' {
// 					s.Scan()
// 					if s.Peek() == '>' {
// 						s.Scan()
// 						break
// 					}
// 				}
// 			}
// 		// ==------------------==
// 		// Template/Operation Case:
// 		// ==------------------==
// 		case (tok == '{' && s.Peek() == '{'):
// 			// Handle template tag start
// 			for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 				txt = s.TokenText()
// 				contents.WriteString(txt)
// 				if tok == '}' && s.Peek() == '}' {
// 					s.Scan()
// 					r.RenderTemplate(contents.String())
// 					break
// 				}
// 			}
// 		// ==------------------==
// 		// Components Case:
// 		// ==------------------==
// 		case tok == '<' && s.Peek() >= 'A' && s.Peek() <= 'Z':
// 			s.Scan()
// 			name = s.TokenText()
// 			attributes = map[string]string{}
// 			nested_lvl = 0
// 			attributes_complete = false
// 			contents.Reset()
// 			for tok := s.Scan(); tok != EOF; tok = s.Scan() {
// 				// Used to remove tabs line breaks and chained spaces
// 				if tok == ' ' && s.Peek() == ' ' || tok == '\n' || tok == '\u0009' {
// 					continue
// 				}
// 				txt = s.TokenText()
// 				s.Scan()
// 				txt2 := s.TokenText()
// 				fmt.Print("boop: ")
// 				fmt.Println(txt, txt2)
// 				// --
// 				// handle nested
// 				if tok == '<' && txt == name && attributes_complete {
// 					nested_lvl++
// 				}
// 				// --
// 				// handle normal close
// 				if tok == '>' && !attributes_complete {
// 					attributes_complete = true
// 					// --
// 					// handle self close
// 				} else if tok == '/' {
// 					if tok == '>' && !attributes_complete {
// 						attributes_complete = true
// 						// self-closing tag
// 						break
// 					} else if txt == name {
// 						if nested_lvl > 0 {
// 							nested_lvl--
// 						} else {
// 							// skip that last character of closing tag ">"
// 							s.Scan()
// 							break
// 						}
// 					}
// 				}
// 				// --
// 				// handle attributes
// 				if attributes_complete {
// 					contents.WriteString(s.TokenText())
// 				} else {
// 					ParseAttributes(&s, &attributes, s.Scan, txt)
// 				}
// 			}
// 			// clean up component end tag artifact "</"
// 			s := strings.TrimSuffix(contents.String(), "</")
// 			r.RenderComponent(name, s, &attributes)
// 		// ==------------------==
// 		// Happy Flow Write String:
// 		// ==------------------==
// 		default:
// 			output.WriteString(s.TokenText())
// 		}
// 	}

// 	result.WriteString(`<code>` + html.EscapeString(output.String()) + `</code>`)
// 	return result.String(), err
// }
