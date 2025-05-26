package auth

import (
	"database/sql"
	"fmt"
	"log"
)

// InitUsersDB creates and migrates all necessary tables
func InitUsersDB(db *sql.DB) error {
	log.Print("Starting Database Initialization")
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Users table - core user information
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
    		uuid TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TRIGGER IF NOT EXISTS update_users_timestamp
		AFTER UPDATE ON users
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Authentication methods table
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS auth_methods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT
		);

		INSERT OR IGNORE INTO auth_methods (name, description) VALUES
			('password', 'Local password authentication'),
			('discord', 'Discord OAuth2'),
			('google', 'Google OAuth2'),
			('github', 'GitHub OAuth2');
	`)
	if err != nil {
		return fmt.Errorf("failed to create auth_methods table: %w", err)
	}

	// User authentication providers (for OAuth, etc.)
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS user_auth_providers (
			user_uuid TEXT NOT NULL,
			auth_method_id INTEGER NOT NULL,
			provider_id TEXT NOT NULL,
			provider_username TEXT,
			provider_email TEXT,
			provider_data TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_uuid, auth_method_id),
			UNIQUE (auth_method_id, provider_id),
			FOREIGN KEY (user_uuid) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (auth_method_id) REFERENCES auth_methods(id)
		);

		CREATE TRIGGER IF NOT EXISTS update_user_auth_providers_timestamp
		AFTER UPDATE ON user_auth_providers
		BEGIN
			UPDATE user_auth_providers SET updated_at = CURRENT_TIMESTAMP
			WHERE user_uuid = OLD.user_uuid AND auth_method_id = OLD.auth_method_id;
		END;
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_auth_providers table: %w", err)
	}

	// Local password credentials
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS user_passwords (
			user_uuid TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			reset_token TEXT,
			reset_token_expiry TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
		);

		CREATE TRIGGER IF NOT EXISTS update_user_passwords_timestamp
		AFTER UPDATE ON user_passwords
		BEGIN
			UPDATE user_passwords SET updated_at = CURRENT_TIMESTAMP WHERE user_uuid = OLD.user_uuid;
		END;
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_passwords table: %w", err)
	}

	// User roles/permissions tables
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT
		);

		INSERT OR IGNORE INTO roles (name, description) VALUES
			('user', 'Regular user'),
			('admin', 'System administrator'),
			('moderator', 'Content moderator');

		CREATE TABLE IF NOT EXISTS user_roles (
			user_uuid TEXT NOT NULL,
			role_id INTEGER NOT NULL,
			assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_uuid, role_id),
			FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create roles tables: %w", err)
	}

	log.Print("Database Initialization Completed Successfully")
	return tx.Commit()
}
