package engine

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/kato-studio/wispy/internal"
)

var log = internal.GetLogger()

/*
Custom Template Functions
*/
func GetFuncs(ctx *TemplateCtx, slots map[string]string) template.FuncMap {
	return template.FuncMap(map[string]any{
		"Slot": func(slot_name string) string {
			return slots[slot_name]
		},
		"Page": func(return_value string) string {
			switch strings.ToLower(return_value) {
			case "title":
				return ctx.Page.Title
			case "head":
				return ctx.Page.InsertHead
			case "meta":
				return "<meta " + strings.Join(ctx.Page.Meta, "><meta ") + ">"
			case "css":
				return "<style>" + ctx.Page.Css + "</style>"
			case "js":
				return "<script>" + ctx.Page.Js + "</script>"
			case "lang":
				return ctx.Page.Lang
			default:
				return ""
			}
		},
	})
}

func Render(template_string string, ctx *TemplateCtx, slots map[string]string) (bytes.Buffer, error) {
	var result = bytes.Buffer{}

	base, err := template.New("template").Funcs(GetFuncs(ctx, slots)).Parse(template_string)

	if err != nil {
		return result, err
	}

	err = base.Execute(&result, ctx.Data)

	if err != nil {
		return result, err
	}

	return result, nil
}

func RenderPage(page_content bytes.Buffer, ctx *TemplateCtx) (bytes.Buffer, error) {
	_slots := map[string]string{
		"content": page_content.String(),
	}
	var _document_result = bytes.Buffer{}
	var _document = []byte{}

	// PAUSED UNTIL COMPONENTS ARE IMPLEMENTED
	// if ctx.Page.Layout != "" {
	// 	layout_file, err := os.ReadFile(ctx.Page.Layout)
	// 	if err != nil {
	// 		log.Error("Layout file not found or failed to read")
	// 		log.Warn("failed to load (" + ctx.Page.Layout + ")")
	// 	}
	// 	with_document, err := template.New("template_layout").Funcs(GetFuncs(ctx, _slots)).Parse(string(layout_file))

	// 	if err == nil { // All good?
	// 		result, err := Render(string(file), ctx, slots)
	// 		if err == nil { // All good?
	// 			return _document_result, nil
	// 		}
	// 	}
	// }

	_document_file, err := os.ReadFile("./shared/layouts/_default/_document.html")
	if err != nil {
		log.Error("Layout file not found or failed to read")
		log.Warn("A default layout should be set in \"shared/layouts/_default/_document.html\" ")
	}
	_document = _document_file

	with_document, err := template.New("template_document").Funcs(GetFuncs(ctx, _slots)).Parse(string(_document))
	if err != nil {
		return page_content, err
	}

	err = with_document.Execute(&_document_result, ctx.Data)
	if err != nil {
		return _document_result, err
	}

	return _document_result, nil
}

func RenderFile(path string, ctx *TemplateCtx) (bytes.Buffer, error) {
	var slots = map[string]string{}

	// Read File
	file, err := os.ReadFile(path)
	if err != nil {
		log.Warn("'RenderFile' FAILED: File not found " + "(" + path + ")")
		return bytes.Buffer{}, err
	}
	//
	result, err := Render(string(file), ctx, slots)
	return result, err
}
