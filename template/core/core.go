package core

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

// // reperesnets the settings/presets of the current template engine instances
//
//	type TemplateEngine struct {
//		// starting deliminator - default "{%"
//		DelimStart string
//		// starting deliminator  - default "%}"
//		DelimEnd string
//		// used to trim the start and end delim as well as leading and trailing whitespace
//		CutSet string // generated by NewTemplateEngine() using default deliminator's
//		// map to check template tags against when rendering
//		// TagMap map[string]TemplateTag
//		// map to check template filters against when rendering
//		FilterMap map[string]core.TemplateFilter
//	}
//
// Render processes the raw template string and writes the rendered output to sb (String Builder)
func Render(ctx *structure.RenderCtx, sb *strings.Builder, raw string) (errs []error) {
	// Starting/Opening Delimiter
	// var ds = ctx.Engine.DelimStart
	var de = ctx.Engine.DelimEnd
	pos := 0
	length := len(raw)

	for pos < length {
		startDelim, endDelim := FindDelim(ctx, raw, pos)
		// If no more delimiters found, append the remaining text and break.
		if startDelim >= len(raw) || startDelim == -1 {
			sb.WriteString(raw[pos:])
			break
		}
		// Append literal text between the current position and the next delimiter.
		sb.WriteString(raw[pos:startDelim])
		pos = startDelim
		//
		//* Core tag logic -----------
		if endDelim == -1 {
			errs = append(errs, fmt.Errorf("missing closing delimit not found %d", pos))
			break
		} else {
			endDelim = endDelim + len(de)
		}
		// Extract the contents of the variable.
		tag_contents := strings.Trim(raw[pos:endDelim], ctx.Engine.CutSet)
		if strings.HasPrefix(tag_contents, ".") {
			if err := ResolveFiltersIfAny(ctx, sb, tag_contents); err != nil {
				errs = append(errs, err)
			}
			pos = endDelim
		} else {
			var tagErrs []error
			pos = endDelim
			pos, tagErrs = ResolveTag(ctx, sb, pos, tag_contents, raw)
			if len(tagErrs) > 0 {
				errs = append(errs, tagErrs...)
			}
		}
	}
	return errs
}
