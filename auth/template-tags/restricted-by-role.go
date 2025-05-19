package tags

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kato-studio/wispy/auth"
	"github.com/kato-studio/wispy/wispy_common"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

var RestrictedByRole = structure.TemplateTag{
	Name: "restricted-by-role",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error
		UserID := ctx.UserID

		// Parse tag options
		options := wispy_common.SplitRespectQuotes(tag_contents)
		optionsMap := wispy_common.ParseKeyValuePairs(options)

		if UserID == "" {
			errs = append(errs, fmt.Errorf("[warning] 'UserID' is not set user is not logged-in!"))
			if url, exists := optionsMap["redirect"]; exists {
				fmt.Println("redirected to -> ", url)
				http.Redirect(*ctx.ResponseWriter, ctx.Request, url, http.StatusSeeOther)
			} else {
				fmt.Println(":( no redirect url found!")
			}
			return len(raw), errs
		}
		if err := ctx.UsersDB.Ping(); err != nil {
			errs = append(errs, fmt.Errorf("ctx.UsersDB was nil \"restricted-by-role\" failed no content rendered"))
			return len(raw), errs
		}

		// Check for required roles parameter
		rolesString, exists := optionsMap["roles"]
		if !exists {
			errs = append(errs, fmt.Errorf("'roles' parameter is required for role-access tag"))
			return len(raw), errs
		}

		// Check user roles against required roles using UsersDB
		requiredRoles := strings.Split(rolesString, " ")
		hasAccess, err := auth.CheckUserRoles(ctx.UsersDB, UserID, requiredRoles)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to check user roles: %w", err))
			return len(raw), errs
		}

		// If user has access, render the content
		if hasAccess {
			return 0, nil
		}

		writer := *ctx.ResponseWriter
		writer.WriteHeader(401)
		writer.Write([]byte("You do not have the required role to access this route."))

		// If no access, don't render anything
		return len(raw), errs
	},
}
