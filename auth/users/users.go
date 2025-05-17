package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserStorage defines the interface for user persistence
type UserStorage interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User, password ...string) error
	GetOrCreateByDiscord(ctx context.Context, discordUser *DiscordUser) (*User, error)
}

// User represents an authenticated user
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string // hashed
	//
	// base (most tables should contain these values)
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      map[string]interface{} // Additional provider-specific data
}
