package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// VerifyAndGetSession checks for a valid session and returns the session info
func VerifyAndGetSession(sessionDB *sql.DB, r *http.Request) (valid bool, userID string, expiresAt time.Time, err error) {
	// Get session cookie
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, "", time.Time{}, nil // No session cookie
		}
		return false, "", time.Time{}, fmt.Errorf("failed to get session cookie: %w", err)
	}

	// Initialize sessions interface
	sessions := SQLiteSessionsInterface{db: sessionDB}

	// Get session from database
	userID, expiresAt, err = sessions.Get(sessionCookie.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", time.Time{}, nil // Session doesn't exist
		}
		return false, "", time.Time{}, fmt.Errorf("failed to verify session: %w", err)
	}

	// Check expiration
	if time.Now().After(expiresAt) {
		return false, "", time.Time{}, nil // Session expired
	}

	return true, userID, expiresAt, nil
}

// Middleware version for http handlers
func SessionMiddleware(sessionDB *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valid, userID, _, err := VerifyAndGetSession(sessionDB, r)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add userID to request context
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Convenience wrapper that returns just the validation status
func IsSessionValid(sessionDB *sql.DB, r *http.Request) (bool, error) {
	valid, _, _, err := VerifyAndGetSession(sessionDB, r)
	return valid, err
}
