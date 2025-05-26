package auth

import (
	"database/sql"
	"fmt"
	"slices"
)

func GetUserByID(db *sql.DB, id string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT uuid, username, email, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.UUID,
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

func GetUserByUUID(db *sql.DB, uuid string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT uuid, username, email, created_at, updated_at
		FROM users WHERE uuid = ?
	`, uuid).Scan(
		&user.UUID,
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

	roles, err := GetUserRoles(db, userID)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}

	// check if user is missing required role
	for _, reqRole := range requiredRoles {
		if slices.Contains(roles, reqRole) == false {
			return false, fmt.Errorf("user: [%s] is missing role (%s)  \n", userID, reqRole)
		}
	}

	// query := `
	// 	SELECT COUNT(*)
	// 	FROM user_roles ur
	// 	JOIN roles r ON ur.role_id = r.id
	// 	WHERE ur.user_id = ? AND r.name IN (` + strings.Repeat("?,", len(requiredRoles)-1) + `?)`

	// fmt.Println(query)

	// args := make([]any, len(requiredRoles)+1)
	// args[0] = userID
	// for i, role := range requiredRoles {
	// 	args[i+1] = role
	// }

	// var count int
	// err := db.QueryRow(query, args...).Scan(&count)
	// if err != nil {
	// 	return false, fmt.Errorf("database query failed: %w", err)
	// }
	return len(roles) > 0, nil
}
