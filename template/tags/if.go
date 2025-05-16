package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// If take only show content if a value is true
var IfTag = TemplateTag{
	Name: "if",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := delimWrap(ctx, "endif")
		endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
			return pos, errs
		}
		content := raw[pos:endTagStart]
		newEndPos := endTagStart + endTagLength

		value, condition_errors := core.ResolveCondition(ctx, tag_contents)
		if len(condition_errors) > 0 {
			errs = append(errs, condition_errors...)
		}
		if value {
			sb.WriteString(content)
		}
		return newEndPos, errs
	},
}
