package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

var AssignTag = TemplateTag{
	Name: "assign",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		parts := strings.SplitN(tag_contents, "=", 2)
		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("invalid assign syntax: expected '{VAR} = {VALUE}', got %q", tag_contents))
			return pos, errs
		}

		variable := strings.TrimSpace(parts[0])
		valueExpr := strings.TrimSpace(parts[1])

		var value any
		var err error
		// Resolve the value to be assigned from string expression
		value, err = core.ResolveValue(ctx, valueExpr)

		if err != nil {
			errs = append(errs, fmt.Errorf("could not resolve value for assignment: %v", err))
			return pos, errs
		}

		if variable[0] == '.' {
			path := strings.Split(strings.TrimPrefix(variable, "."), ".")
			err := insertNestedValue(ctx.Data, path, value)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to set nested value: %v", err))
			}
		} else {
			ctx.Data[variable] = value
		}

		return pos, errs
	},
}
