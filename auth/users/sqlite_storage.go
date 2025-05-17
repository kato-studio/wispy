package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteUserStorage implements UserStorage for SQLite
type SQLiteUserStorage struct {
	db *sql.DB
}

// GetByUsername retrieves a user by username
func (s *SQLiteUserStorage) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *SQLiteUserStorage) GetByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Create creates a new user with optional password
func (s *SQLiteUserStorage) Create(ctx context.Context, user *User, password ...string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Insert base user
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, username, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`,
		user.ID.String(),
		user.Username,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Add password if provided
	if len(password) > 0 && password[0] != "" {
		hashedPassword, err := HashPassword(password[0])
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_passwords (user_id, password_hash)
			VALUES (?, ?)
		`,
			user.ID.String(),
			hashedPassword,
		)
		if err != nil {
			return fmt.Errorf("failed to set password: %w", err)
		}
	}

	return tx.Commit()
}

// GetOrCreateByProvider gets or creates a user from OAuth provider data
func (s *SQLiteUserStorage) GetOrCreateByProvider(ctx context.Context, providerName string, providerUser *User) (*User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if provider user exists
	if providerUser.ID.String() != "" {
		var userID string
		err = tx.QueryRowContext(ctx, `
			SELECT user_id FROM user_auth_providers
			WHERE auth_method_id = (SELECT id FROM auth_methods WHERE name = ?)
			AND provider_id = ?
		`, providerName, providerUser.ID).Scan(&userID)

		if err == nil {
			// Provider user exists, return the associated user
			return s.GetByID(ctx, userID)
		}
	} else {
		var userID string
		err = tx.QueryRowContext(ctx, `
			SELECT user_id FROM user_auth_providers
			WHERE auth_method_id = (SELECT id FROM auth_methods WHERE name = ?)
			AND provider_email = ?
		`, providerName, providerUser.Email).Scan(&userID)

		if err == nil {
			// Provider user exists, return the associated user
			return s.GetByID(ctx, userID)
		}
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query provider user: %w", err)
	}

	// Create new user
	user := &User{
		ID:        uuid.New(),
		Username:  providerUser.Username,
		Email:     providerUser.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user
	if err := s.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert provider data to JSON
	providerData, err := json.Marshal(providerUser.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal provider data: %w", err)
	}

	// Link provider account
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_auth_providers
		(user_id, auth_method_id, provider_id, provider_username, provider_email, provider_data)
		VALUES (?,
			(SELECT id FROM auth_methods WHERE name = ?),
			?, ?, ?, ?
		)
	`,
		user.ID.String(),
		providerName,
		providerUser.ID,
		providerUser.Username,
		providerUser.Email,
		providerData,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create provider user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

// GetOrCreateByDiscord gets or creates a user from Discord OAuth data
func (s *SQLiteUserStorage) GetOrCreateByDiscord(ctx context.Context, discordUser *DiscordUser) (*User, error) {
	providerUser := &User{
		Username: discordUser.Username,
		Email:    discordUser.Email,
		Data: map[string]interface{}{
			"id":     discordUser.ID,
			"avatar": discordUser.Avatar,
		},
	}
	return s.GetOrCreateByProvider(ctx, "discord", providerUser)
}

// VerifyPassword checks if the password matches for a user
func (s *SQLiteUserStorage) VerifyPassword(ctx context.Context, username, password string) (bool, error) {
	var hashedPassword string
	err := s.db.QueryRowContext(ctx, `
		SELECT up.password_hash FROM user_passwords up
		JOIN users u ON up.user_id = u.id
		WHERE u.username = ?
	`, username).Scan(&hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify password: %w", err)
	}

	return CheckPassword(password, hashedPassword), nil
}

// SetPassword updates or sets a user's password
func (s *SQLiteUserStorage) SetPassword(ctx context.Context, userID string, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO user_passwords (user_id, password_hash, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id) DO UPDATE SET
			password_hash = excluded.password_hash,
			updated_at = CURRENT_TIMESTAMP
	`, userID, hashedPassword)

	return err
}

func (s *SQLiteUserStorage) DB() *sql.DB {
	return s.db
}

// NewSQLiteUserStorage creates a new SQLite user storage
func NewSQLiteUserStorage(dbPath string) (*SQLiteUserStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign key support
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &SQLiteUserStorage{db: db}, nil
}
