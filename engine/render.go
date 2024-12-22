package engine

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"strings"
	"text/scanner"
)

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

func ParseAttributes(s Scanner, attributes *map[string]string, _scan func() rune, key string) {
	if s.Peek() == '=' {
		_scan()
		_scan()
		(*attributes)[key] = s.TokenText()
	} else if len(key) > 1 {
		(*attributes)[key] = "true"
	}
}

type Render struct {
	AttrFuncMap      map[string]AttributeFunc
	OperationFuncMap map[string]OperationFunc
	Ctx              *TemplateCtx
}

type OperationFunc func(r *Render, values ...string) string

func (r *Render) SetOperationFunc(operation_name string, handler OperationFunc) *Render {
	r.OperationFuncMap[operation_name] = handler
	return r
}

type AttributeFunc func(value string) string

func (r *Render) SetAttributeFunc(attribute_name string, handler AttributeFunc) *Render {
	r.AttrFuncMap[attribute_name] = handler
	return r
}

func (r *Render) SetCtx(ctx *TemplateCtx) *Render {
	r.Ctx = ctx
	return r
}

func (r *Render) Html(template_string string) (string, error) {
	var result strings.Builder
	var s Scanner
	//? TODOd: add error handling and better logging for debugging
	var err error
	// var ctx = r.Ctx

	s.Init(strings.NewReader(template_string))
	var output strings.Builder
	var prev_ch rune
	var ch rune

	var _scan = func() rune {
		prev_ch = ch
		ch = s.Scan()
		return ch
	}

	for ch := _scan(); ch != scanner.EOF; ch = _scan() {
		// Used to remove tabs line breaks and chained spaces
		if ch == ' ' && s.Peek() == ' ' || ch == '\n' || ch == '\u0009' {
			continue
		}

		var name = ""
		var attributes = map[string]string{}
		var nested_lvl = 0
		var attributes_complete = false
		var contents strings.Builder
		switch true {
		// ==------------------==
		// Operation Case
		// ==------------------==
		case prev_ch == '<' && ch == 'x' && s.Peek() == '.':
			for ch = _scan(); ch != scanner.EOF; ch = _scan() {
				if ch == '-' && s.Peek() == '-' {
					_scan()
					if s.Peek() == '>' {
						_scan()
						break
					}
				}
			}

		// ==------------------==
		// Comments Case
		// Skip comments "<!--" to "-->"
		// ==------------------==
		case ch == '<' && s.Peek() == '!':
			for ch = _scan(); ch != scanner.EOF; ch = _scan() {
				if ch == '-' && s.Peek() == '-' {
					_scan()
					if s.Peek() == '>' {
						_scan()
						break
					}
				}
			}

		// ==------------------==
		// Template/Operation Case:
		// ==------------------==
		case (ch == '{' && s.Peek() == '{'):
			// Handle template tag start
			for ch = _scan(); ch != scanner.EOF; ch = _scan() {
				contents.WriteString(s.TokenText())
				if ch == '}' && s.Peek() == '}' {
					_scan()
					r.RenderTemplate(contents.String())
					break
				}
			}
		// ==------------------==
		// Components Case:
		// ==------------------==
		case ch == '<' && s.Peek() >= 'A' && s.Peek() <= 'Z':
			_scan()
			name = s.TokenText()
			attributes = map[string]string{}
			nested_lvl = 0
			attributes_complete = false
			contents.Reset()
			for _scan(); ch != scanner.EOF; ch = _scan() {
				// Used to remove tabs line breaks and chained spaces
				if ch == ' ' && s.Peek() == ' ' || ch == '\n' || ch == '\u0009' {
					continue
				}
				//
				value := s.TokenText()
				// --
				// handle nested
				if prev_ch == '<' && value == name && attributes_complete {
					nested_lvl++
				}
				// --
				// handle normal close
				if value == ">" && !attributes_complete {
					attributes_complete = true
					_scan()
					// --
					// handle self close
				} else if prev_ch == '/' {
					if value == ">" && !attributes_complete {
						attributes_complete = true
						// self-closing tag
						break
					} else if value == name {
						if nested_lvl > 0 {
							nested_lvl--
						} else {
							// skip that last character of closing tag ">"
							_scan()
							break
						}
					}
				}
				// --
				// handle attributes
				if attributes_complete {
					contents.WriteString(s.TokenText())
				} else {
					ParseAttributes(s, &attributes, _scan, value)
				}
			}
			// clean up component end tag artifact "</"
			s := strings.TrimSuffix(contents.String(), "</")
			r.RenderComponent(name, s, &attributes)
		// ==------------------==
		// Happy Flow Write String:
		// ==------------------==
		default:
			output.WriteString(s.TokenText())
		}
		//
		// set previous character variable to check against next run
		prev_ch = ch
	}

	result.WriteString(`<code>` + html.EscapeString(output.String()) + `</code>`)
	return result.String(), err
}

func (r *Render) RenderComponent(name, contents string, attributes *map[string]string) {
	fmt.Println("COMPS: ")
	fmt.Println(name, contents)
	fmt.Println("===========")
}

func (r *Render) RenderTemplate(contents string) {
	fmt.Println("Template: ")
	fmt.Println(contents)
	fmt.Println("===========")
}

func (r *Render) RenderPage(domain, page_path string) (string, error) {
	// Read File
	domain = strings.Trim(domain, "/")
	var path = "./sites" + domain + "/pages/" + PAGE_FILE
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	//
	result, err := r.Html(string(file))
	//? TODO better error handling for dev mod
	//? next js style error page
	return result, err
}

// func Render(template_string string, ctx *TemplateCtx) (bytes.Buffer, error) {
// 	var result = bytes.Buffer{}

// 	/*
// 		Custom Template Functions
// 	*/
// 	var funcs = template.FuncMap(map[string]interface{}{
// 		// TODO:
// 	})

// 	base, err := template.New("template").Funcs(funcs).Parse(template_string)

// 	if err != nil {
// 		return result, err
// 	}

// 	err = base.Execute(&result, ctx.Json.Value())

// 	if err != nil {
// 		return result, err
// 	}

// 	/*
// 		Handle Page Layout
// 	*/
// 	var layout_result = bytes.Buffer{}
// 	var layout = `{{Page "Content"}}`

// 	if ctx.Page.Layout != "" {
// 		layout_file, err := os.ReadFile(ctx.Page.Layout)
// 		if err != nil {
// 			return result, err
// 		}
// 		layout = string(layout_file)
// 	} else {
// 		// TODO: add warning/error
// 	}

// 	with_layout, err := template.New("template").Funcs(template.FuncMap{
// 		"Page": func(return_value string) string {
// 			return LayoutPageInsert(return_value, ctx, &result)
// 		},
// 	}).Parse(layout)
// 	if err != nil {
// 		return result, err
// 	}

// 	err = with_layout.Execute(&layout_result, ctx.Json.Value())
// 	if err != nil {
// 		return layout_result, err
// 	}

// 	return layout_result, nil
// }
