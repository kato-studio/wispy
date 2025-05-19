package auth

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetUserByID(db *sql.DB, id string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT uuid, username, email, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func GetUserByUUID(db *sql.DB, id string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT uuid, username, email, created_at, updated_at
		FROM users WHERE uuid = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// checkUserRoles verifies if a user has any of the required roles
func CheckUserRoles(db *sql.DB, userID string, requiredRoles []string) (bool, error) {
	if len(requiredRoles) == 0 {
		return false, nil
	}

	query := `
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.name IN (` + strings.Repeat("?,", len(requiredRoles)-1) + `?)`

	args := make([]interface{}, len(requiredRoles)+1)
	args[0] = userID
	for i, role := range requiredRoles {
		args[i+1] = role
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}

	return count > 0, nil
}
