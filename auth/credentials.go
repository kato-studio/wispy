package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailExists       = errors.New("email already exists")
	ErrUsernameExists    = errors.New("username already exists")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrInvalidInviteCode = errors.New("invalid invite code")
)

type Credentials struct {
	Email      string      `json:"email"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	Referrer   string      `json:"referrer"`
	InviteCode string      `json:"invite_code"`
	Data       interface{} `json:"data"`
}

func RegisterWithCredentials(UserDB *sql.DB, w http.ResponseWriter, r *http.Request, creds Credentials) (*User, error) {
	// Validate input
	if len(creds.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	// Check if email exists
	var emailCount int
	err := UserDB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE email = ?
	`, creds.Email).Scan(&emailCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailCount > 0 {
		return nil, ErrEmailExists
	}

	// Check if username exists
	var usernameCount int
	err = UserDB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE username = ?
	`, creds.Username).Scan(&usernameCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if usernameCount > 0 {
		return nil, ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := HashPassword(creds.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := UserDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create user
	var user = User{
		UUID:     ksuid.New().String(),
		Username: creds.Username,
		Email:    creds.Email,
	}
	err = tx.QueryRow(`
		INSERT INTO users (uuid, username, email)
		VALUES (?, ?, ?)
		RETURNING id, uuid, username, email, created_at, updated_at
	`, user.UUID, user.Username, user.Email).Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Store password
	fmt.Println("UUID", user.UUID, reflect.TypeOf(user.UUID))
	fmt.Println("hashedPassword", hashedPassword, reflect.TypeOf(hashedPassword))
	_, err = tx.Exec(`
		INSERT INTO user_passwords (user_uuid, password_hash)
		VALUES (?, ?)
	`, user.UUID, hashedPassword)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to store password: %w", err)
	}

	// Assign default role
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_uuid, role_id)
		SELECT ?, id FROM roles WHERE name = 'user'
	`, user.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func LoginWithCredentials(UserDB *sql.DB, w http.ResponseWriter, r *http.Request, creds Credentials) (*User, error) {
	// Find user by email or username
	var user User
	var passwordHash string

	err := UserDB.QueryRow(`
		SELECT u.uuid, u.username, u.email, u.created_at, u.updated_at, up.password_hash
		FROM users u
		JOIN user_passwords up ON u.uuid = up.user_uuid
		WHERE u.email = ? OR u.username = ?
	`, creds.Email, creds.Username).Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}
