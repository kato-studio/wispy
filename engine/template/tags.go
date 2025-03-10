package template

import (
	"fmt"
	"strings"
)

// Universal template tag function struct
type TemplateTag struct {
	Name string
	// render tag with given context and args, and children if a RequiresClosingTag is set
	Render func(
		// The context of this render execution including...
		// - Reference to the template engine struct
		// - Partials map,
		// - Data map fetched via eng.GetFunc(ctx *RenderCtx, key string)
		ctx *RenderCtx,
		// Finalized output is written to the string building
		sb *strings.Builder,
		// The inner contexts of the tag being parsed
		// Example: "{% exampleTag ... ... ... %}"
		// (tags using closing tag for inner content are expected to resolve the closing tag and content then return update POS int)
		tag_contents,
		// entire input string mainly used if tag handles closing child_content + closing_tag
		raw string,
		// the current position within the raw input string
		pos int,
	) (new_pos int, errs []error)
}

// Default template tag functions
var PartialTag = TemplateTag{
	Name: "partial",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {

		return pos, errs
	},
}

var IfTag = TemplateTag{
	Name: "if",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := IndexAt(raw, "{% endif %}", pos)
		if endTag == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", "{% endtag %}"))
			return pos, errs
		}

		sb.WriteString("[[" + tag_contents + "]]")
		return pos, errs
	},
}

var ForTag = TemplateTag{
	Name: "for",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		sb.WriteString("[[" + tag_contents + "]]")
		return pos, errs
	},
}

var RootTag = TemplateTag{
	Name: "root",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {

		return pos, errs
	},
}

var LayoutTag = TemplateTag{
	Name: "layout",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// check if layout has
		if _, ok := ctx.InternalFlags["layout"]; ok {
			errs = append(errs, fmt.Errorf("Wanring: layout alread set for this page %s", ":("))
		} else {

			// startDelim, endDelim := FindDelim(ctx,raw,pos)
			Render(ctx, sb, raw)
		}

		return pos, errs
	},
}

// var PageTag = TemplateTag{
// 	Name: "page",
// 	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {

// 		return pos, errs
// 	},
// }

var DefaultTemplateTags = []TemplateTag{
	IfTag,
	ForTag,
	PartialTag,
}

var DefaultEngineTags = []TemplateTag{
	IfTag,
	ForTag,
	PartialTag,
	// PageTag,
	RootTag,
}
