package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	common "github.com/kato-studio/wispy/wispy_common"
)

type SQLiteSessionsInterface struct {
	db *sql.DB
}

func NewSQLiteSessions(db *sql.DB) *SQLiteSessionsInterface {
	return &SQLiteSessionsInterface{db: db}
}

func (s *SQLiteSessionsInterface) Init() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            token TEXT PRIMARY KEY,
            user_uuid TEXT NOT NULL,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
        );

        CREATE INDEX IF NOT EXISTS idx_sessions_user_uuid ON sessions(user_uuid);
        CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
    `)
	return err
}

func (s *SQLiteSessionsInterface) Get(token string) (string, time.Time, error) {
	var userID string
	var expiresAt time.Time
	err := s.db.QueryRow(`
        SELECT user_uuid, expires_at FROM sessions
        WHERE token = ? AND expires_at > CURRENT_TIMESTAMP
    `, token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}
	return userID, expiresAt, nil
}

func (s *SQLiteSessionsInterface) Update(token string, userUUID string, duration time.Duration) error {
	expiresAt := time.Now().Add(duration)
	_, err := s.db.Exec(`
        INSERT OR REPLACE INTO sessions (token, user_uuid, expires_at)
        VALUES (?, ?, ?)
    `, token, userUUID, expiresAt)
	return err
}

func (s *SQLiteSessionsInterface) Delete(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}

// ---
// AuthManager Session Methods
// ---

// VerifySession checks if the session token is valid
func VerifySession(SessionDB *sql.DB, r *http.Request) (bool, error) {
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return false, nil // No cookie means no session
	}

	var Sessions = SQLiteSessionsInterface{db: SessionDB}

	_, expiresAt, err := Sessions.Get(sessionCookie.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Session doesn't exist
		}
		return false, fmt.Errorf("failed to verify session: %w", err)
	}

	if time.Now().After(expiresAt) {
		return false, nil // Session expired
	}

	return true, nil
}

// GetUserFromSession retrieves the authenticated user from the session cookie
func GetUserFromSession(SessionDB *sql.DB, UserDB *sql.DB, r *http.Request) (*User, error) {
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return nil, nil // No session
	}
	var Sessions = SQLiteSessionsInterface{db: SessionDB}

	userID, _, err := Sessions.Get(sessionCookie.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Session doesn't exist
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var user User
	err = UserDB.QueryRow(`
        SELECT uuid, username, email, created_at, updated_at
        FROM users WHERE id = ?
    `, userID).Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateSession creates a new authenticated session for the user
func CreateSession(SessionDB *sql.DB, w http.ResponseWriter, r *http.Request, user *User) error {
	token, err := GenerateRandomString(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	var Sessions = SQLiteSessionsInterface{db: SessionDB}

	// Store session in database
	err = Sessions.Update(token, user.UUID, time.Hour*24*7) // 1 week
	if err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "auth-session",
		Value: token,
		Path:  "/",
		// MaxAge:   int((time.Hour * 24 * 7).Seconds()), // 1 week
		HttpOnly: true,
		Secure:   common.IsProduction(),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// DestroySession removes the authenticated session
func DestroySession(SessionDB *sql.DB, w http.ResponseWriter, r *http.Request) error {
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return nil // No session to destroy
	}

	var Sessions = SQLiteSessionsInterface{db: SessionDB}

	// Remove token from storage
	err = Sessions.Delete(sessionCookie.Value)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   common.IsProduction(),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
