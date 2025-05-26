package auth

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kato-studio/wispy/template"
	"github.com/kato-studio/wispy/wispy_common"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// Define some example colors for logging
const (
	colorCyan  = "\033[36m"
	colorGrey  = "\033[90m"
	colorReset = "\033[0m"
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

// Sets up basic sql-lite database
func SetupUserAndSessionDB() (UsersDB *sql.DB, SessionsDB *sql.DB) {
	// SETUP DATABASES FOR AUTH
	var err error
	// Ensure database directory exists
	if err = os.MkdirAll("./db", 0755); err != nil {
		log.Fatalf("failed to create db directory: %s", err)
	}

	// -- DB's --
	wispy_common.CreateFileIfNotExists("./db/auth_users_and_roles.db", nil)
	UsersDB, err = sql.Open("sqlite3", "./db/auth_users_and_roles.db")
	// Initialize databases with connection pooling
	if err != nil {
		log.Fatalf("failed to open users database: %s", err)
	}
	//
	err = InitUsersDB(UsersDB)
	if err != nil {
		log.Fatalf("failed to initialize users database: %s", err)
	}
	//
	wispy_common.CreateFileIfNotExists("./db/auth_sessions.db", nil)
	SessionsDB, err = sql.Open("sqlite3", "./db/auth_sessions.db")
	if err != nil {
		log.Fatalf("failed to open sessions database: %s", err)
		UsersDB.Close() // Clean up already opened DB
	}
	authSessions := NewSQLiteSessions(SessionsDB)
	err = authSessions.Init()
	if err != nil {
		log.Fatalf("failed to initialize sessions database: %s", err)
	}

	// Verify database connections
	if err = UsersDB.Ping(); err != nil {
		cleanupDBs(UsersDB, SessionsDB)
		log.Fatalf("users database connection failed: %s", err)
	}

	if err = SessionsDB.Ping(); err != nil {
		cleanupDBs(UsersDB, SessionsDB)
		log.Fatalf("sessions database connection failed: %s", err)
	}

	return UsersDB, SessionsDB
}

// cleanupDBs closes database connections
func cleanupDBs(databases ...*sql.DB) {
	for _, db := range databases {
		if db != nil {
			db.Close()
		}
	}
}

func SiteAuthRouteHandler(engine *structure.TemplateEngine, SessionsDB *sql.DB, UserDB *sql.DB, w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	domain := r.Host

	// Look up the site structure for the domain
	site, exists := engine.SiteMap[domain]
	if !exists {
		http.Error(w, fmt.Sprintf("domain %s not found", domain), http.StatusNotFound)
		return
	}

	scopedDirectory := filepath.Join(engine.SITES_DIR, site.Domain)
	// Handle public content
	path := filepath.Clean(r.URL.Path)
	// if file extension check if there is a valid file in public directory to serve
	if filepath.Ext(path) != "" {
		// Serve public content if available
		target := filepath.Join(scopedDirectory, "public", path)
		_, err := os.Stat(target)
		if err != nil {
			if os.IsNotExist(err) {
				// File doesn't exist, continue
				fmt.Println("404:", target)
			}
		} else {
			// File exists and is not a directory - serve it
			http.ServeFile(w, r, target)
			return
		}
	}
	//
	data := map[string]any{}
	ctx := engine.InitCtx(scopedDirectory, &site, data)

	// -------- Auth code here --------
	validSession, userID, _, getSessionErr := VerifyAndGetSession(SessionsDB, r)
	if getSessionErr != nil {
		slog.Error("VerifyAndGetSession failed!" + getSessionErr.Error())
		return
	}

	// Render the route
	// Set up the rendering context using NewRenderCtx (which initializes Internal automatically).
	ctx.UsersDB = UserDB
	if validSession {
		ctx.UserID = userID
	}
	// -------- ------------- --------

	//
	page, err := template.RenderRoute(engine, ctx, r.URL.Path, data, w, r)
	if err != nil {
		slog.Error("Rendering Route using \"RenderRoute()\"" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Measure rendering and styling times
	renderTime := time.Now()
	var results bytes.Buffer
	results.WriteString(page)

	// Log performance metrics
	colorize := func(dur time.Duration) string {
		return fmt.Sprintf("%s%v%s", colorCyan, dur, colorGrey)
	}

	fmt.Printf("%s[Render: %s | Total: %s]%s\n",
		colorGrey,
		colorize(renderTime.Sub(startTime)),
		colorize(time.Since(startTime)),
		colorReset,
	)

	// Write the final HTML to the response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(results.Bytes())
}
