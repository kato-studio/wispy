package auth

import (
	"time"
)

// AuthConfig holds configuration for authentication
type Config_OAuth struct {
	ClientId        string
	ClientSecret    string
	RedirectURL     string
	TokenSecret     string
	TokenExpiration int
}

var dur7says = int((time.Hour * 24 * 7).Seconds())

type User struct {
	ID        int
	UUID      string
	Username  string
	Email     string
	UpdatedAt time.Time
	CreatedAt time.Time
}
