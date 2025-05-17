package auth

import (
	"net/http"
	"slices"
)

// Enhanced AuthMiddleware with role checking
func RoleAuthMiddleware(manager AuthManager, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := manager.GetUserFromSession(r)
			if user == nil || err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Get DB fr
			db := manager.Database
			if db == nil {
				http.Error(w, "Database connection not available", http.StatusInternalServerError)
				return
			}

			// Check if user has any of the required roles
			hasAccess := false
			userRoles, err := GetUserRoles(db, user.ID)
			if err != nil {
				http.Error(w, "Failed to check user roles", http.StatusInternalServerError)
				return
			}

			for index, requiredRole := range requiredRoles {
				if slices.Contains(userRoles, requiredRole) {
					if index == len(requiredRoles)-1 {
						hasAccess = true
					}
				} else {
					// unnecessary but feels better have it here..
					hasAccess = false
					break
				}
			}

			if !hasAccess && len(requiredRoles) > 0 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
