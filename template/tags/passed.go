package tags

import (
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// "passed" currently used by extend & layout tags to render passed content
var PassedTag = TemplateTag{
	Name: "passed",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		if ctx.Passed == "" {
			return pos, errs
		}

		// Render the layout template
		renderErrs := core.Render(ctx, sb, string(ctx.Passed))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}
		// reset/clear Passed context
		ctx.Passed = ""

		// Return position after the closing tag
		return pos, errs
	},
}
