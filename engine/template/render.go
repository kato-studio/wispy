package template

import (
	"fmt"
	"strings"
)

// Render processes the raw template string and writes the rendered output to sb (String Builder)
func Render(ctx *RenderCtx, sb *strings.Builder, raw string) (errs []error) {
	// Starting/Opening Delimiter
	var ds = ctx.Engine.DelimStart
	var de = ctx.Engine.DelimEnd
	pos := 0
	length := len(raw)

	for pos < length {
		// Find the next occurrence of a variable or tag start delimiter.
		next := IndexAt(raw, ds, pos)
		// If no more delimiters found, append the remaining text and break.
		if next >= length || next == -1 {
			sb.WriteString(raw[pos:])
			break
		}
		// Append literal text between the current position and the next delimiter.
		sb.WriteString(raw[pos:next])
		pos = next
		//
		//* Core tag logic -----------
		// find bracket closing delim
		endDelim := IndexAt(raw, de, pos)
		if endDelim == -1 {
			errs = append(errs, fmt.Errorf("missing closing delimit not found %d", pos))
			break
		} else {
			endDelim = endDelim + len(de)
		}
		// Extract the contents of the variable.
		tag_contents := CleanTemplateTag(ctx, raw[pos:endDelim])
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
