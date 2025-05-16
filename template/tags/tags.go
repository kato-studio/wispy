package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// // Universal template tag function struct
//
//	type TemplateTag struct {
//		Name string
//		// render tag with given context and args, and children if a RequiresClosingTag is set
//		Render func(
//			// Executes expected logic writes results to strings.Builder and returns new index from to continue rendering
//			// - Reference to the template engine struct
//			// - Partials map,
//			// - Data map fetched via eng.GetFunc(ctx *structure.RenderCtx, key string)
//			ctx *template.RenderCtx,
//			// Finalized output is written to the string building
//			sb *strings.Builder,
//			// The inner contexts of the tag being parsed
//			// Example: "{% exampleTag ... ... ... %}"
//			// (tags using closing tag for inner content are expected to resolve the closing tag and content then return update POS int)
//			tag_contents,
//			// entire input string mainly used if tag handles closing child_content + closing_tag
//			raw string,
//			// the current position within the raw input string
//			pos int,
//		) (new_pos int, errs []error)
//	}
type TemplateTag = structure.TemplateTag

func delimWrap(ctx *structure.RenderCtx, value string) string {
	return strings.Join([]string{ctx.Engine.DelimStart, value, ctx.Engine.DelimEnd}, " ")
}

// used to skip content that should not be parsed
var CommentTag = TemplateTag{
	Name: "comment",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := delimWrap(ctx, "endcomment")
		endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
			return pos, errs
		}
		newEndPos := endTagStart + endTagLength

		return newEndPos, errs
	},
}
