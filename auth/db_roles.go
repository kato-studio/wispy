package auth

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Role management functions
func AssignRoleToUser(db *sql.DB, userID uuid.UUID, roleName string) error {
	// First get role ID
	var roleID int
	err := db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO user_roles (user_id, role_id)
		VALUES (?, ?)
	`, userID, roleID)
	return err
}

func RemoveRoleFromUser(db *sql.DB, userID uuid.UUID, roleName string) error {
	var roleID int
	err := db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	_, err = db.Exec(`
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`, userID, roleID)
	return err
}
func UserHasRole(db *sql.DB, userID uuid.UUID, roleName string) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.name = ?
		)
	`, userID, roleName).Scan(&exists)
	return exists, err
}

func GetUserRoles(db *sql.DB, userID uuid.UUID) ([]string, error) {
	rows, err := db.Query(`
		SELECT r.name FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
