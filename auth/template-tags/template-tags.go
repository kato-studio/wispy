package tags

import (
	"fmt"
	"strings"

	template_core "github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

func delimWrap(ctx *structure.RenderCtx, value string) string {
	return strings.Join([]string{ctx.Engine.DelimStart, value, ctx.Engine.DelimEnd}, " ")
}

var UserTag = structure.TemplateTag{
	Name: "user",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error

		// Parse tag options
		options := wispy_common.SplitRespectQuotes(tag_contents)

		// Find end tag
		endTag := delimWrap(ctx, "end-user")
		endTagStart, endTagLength := template_core.SeekIndexAndLength(raw, endTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
			return pos, errs
		}

		content := raw[pos:endTagStart]
		newEndPos := endTagStart + endTagLength

		// Check authentication state requirements
		if len(options) > 0 {
			switch options[0] {
			case "logged-in":
				if ctx.UserID != "" {
					sb.WriteString(content)
					return newEndPos, nil
				}
			case "logged-out":
				if ctx.UserID == "" {
					sb.WriteString(content)
					return newEndPos, nil
				}
			default:
				errs = append(errs, fmt.Errorf("invalid props for 'user' - try setting 'logged-in' or 'logged-out'"))
			}
		} else {
			errs = append(errs, fmt.Errorf("invalid props for 'user' - try setting 'logged-in' or 'logged-out'"))
		}

		// // Handle user data fetching
		// if ctx.UserID != "" && ctx.UsersDB != nil {
		// 	if fields, exists := optionsMap["with"]; exists {
		// 		userData, err := getUserData(ctx.UsersDB, ctx.UserID, strings.Split(fields, ","))
		// 		if err != nil {
		// 			errs = append(errs, fmt.Errorf("failed to get user data: %w", err))
		// 		} else {
		// 			// Add user data to context
		// 			if ctx.Data == nil {
		// 				ctx.Data = make(map[string]interface{})
		// 			}
		// 			ctx.Data["User"] = userData
		// 		}
		// 	}
		// }

		return newEndPos, errs
	},
}
