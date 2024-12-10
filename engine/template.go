package engine

import (
	"bytes"
	"os"
	"strings"
	"text/template"
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

func Render(template_string string, ctx *TemplateCtx) (bytes.Buffer, error) {
	var result = bytes.Buffer{}

	/*
		Custom Template Functions
	*/
	var funcs = template.FuncMap(map[string]interface{}{
		// TODO:
	})

	base, err := template.New("template").Funcs(funcs).Parse(template_string)

	if err != nil {
		return result, err
	}

	err = base.Execute(&result, ctx.Json.Value())

	if err != nil {
		return result, err
	}

	/*
		Handle Page Layout
	*/
	var layout_result = bytes.Buffer{}
	var layout = "{{Page \"Content\"}}"

	if ctx.Page.Layout != "" {
		layout_file, err := os.ReadFile(ctx.Page.Layout)
		if err != nil {
			return result, err
		}
		layout = string(layout_file)
	} else {
		// TODO: add warning/error
	}

	with_layout, err := template.New("template").Funcs(template.FuncMap{
		"Page": func(return_value string) string {
			return LayoutPageInsert(return_value, ctx, &result)
		},
	}).Parse(layout)
	if err != nil {
		return result, err
	}

	err = with_layout.Execute(&layout_result, ctx.Json.Value())
	if err != nil {
		return layout_result, err
	}

	return layout_result, nil
}

func RenderPage(path string, ctx *TemplateCtx) (bytes.Buffer, error) {
	// Read File
	file, err := os.ReadFile(path)
	if err != nil {
		return bytes.Buffer{}, err
	}
	//
	result, err := Render(string(file), ctx)
	return result, err
}
