package auth

import (
	"database/sql"
	"fmt"
)

// InitDB creates and migrates all necessary tables
func InitDB(db *sql.DB) error {
	// createTables sets up the required database tables with proper normalization
	// Use a transaction for atomic table creation
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Users table - core user information
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TRIGGER IF NOT EXISTS update_users_timestamp
		AFTER UPDATE ON users
		FOR EACH ROW
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
			user_id TEXT NOT NULL,
			auth_method_id INTEGER NOT NULL,
			provider_id TEXT NOT NULL,  // Unique ID from provider (e.g., Discord ID)
			provider_username TEXT,
			provider_email TEXT,
			provider_data TEXT,  // JSON blob of additional provider data
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, auth_method_id),
			UNIQUE (auth_method_id, provider_id),  // One provider ID per auth method
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (auth_method_id) REFERENCES auth_methods(id)
		);

		CREATE TRIGGER IF NOT EXISTS update_user_auth_providers_timestamp
		AFTER UPDATE ON user_auth_providers
		FOR EACH ROW
		BEGIN
			UPDATE user_auth_providers SET updated_at = CURRENT_TIMESTAMP
			WHERE user_id = OLD.user_id AND auth_method_id = OLD.auth_method_id;
		END;
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_auth_providers table: %w", err)
	}

	// Local password credentials (separate from OAuth providers)
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS user_passwords (
			user_id TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			reset_token TEXT,
			reset_token_expiry TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TRIGGER IF NOT EXISTS update_user_passwords_timestamp
		AFTER UPDATE ON user_passwords
		FOR EACH ROW
		BEGIN
			UPDATE user_passwords SET updated_at = CURRENT_TIMESTAMP WHERE user_id = OLD.user_id;
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
			user_id TEXT NOT NULL,
			role_id INTEGER NOT NULL,
			assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create roles tables: %w", err)
	}

	// Create indexes for performance
	_, err = tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_user_auth_providers_provider ON user_auth_providers(auth_method_id, provider_id);
		CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return tx.Commit()
}

// Example of how this would work with different providers:
// // Adding a new authentication method becomes trivial:
// _, err = db.Exec(`INSERT OR IGNORE INTO auth_methods (name, description) VALUES (?, ?)`,
// 	"apple", "Apple Sign In")

// // Querying a user by provider:
// var userID string
// err = db.QueryRow(`
// 	SELECT user_id FROM user_auth_providers
// 	WHERE auth_method_id = (SELECT id FROM auth_methods WHERE name = ?)
// 	AND provider_id = ?`,
// 	"discord", "123456789").Scan(&userID)
