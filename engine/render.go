package engine

import (
	"bytes"
	"fmt"
	"io"
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
	r.SetCtx(ctx)
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

// SetCtx sets the template context for the render context.
func (r *Render) SetCtx(ctx *TemplateCtx) *Render {
	r.Ctx = ctx
	return r
}

func (r *Render) Html(templateString string) (string, error) {
	var result strings.Builder
	var err error
	var index = 0

}

func Utf8At(doc io.Reader, target string, w io.Writer) {

}

func NextUtf8(doc io.Reader, target string, w io.Writer) {

}

func WriteUntil(doc io.Reader, target string, w io.Writer) {

}

func NextEqualOrWrite(doc io.Reader, target string, w io.Writer) {
	//
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

// RenderComponent renders a component with the given name, contents, and attributes.
func (r *Render) RenderComponent(name, contents string, attributes *map[string]string) {
	fmt.Println("COMPS: ")
	fmt.Println(name, contents)
	fmt.Println("===========")
}

// RenderTemplate renders a template with the given contents.
func (r *Render) RenderTemplate(contents string) {
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
