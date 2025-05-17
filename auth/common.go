package auth

import (
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4/server"
)

// AuthConfig holds configuration for authentication
type AuthConfig struct {
	ClientId        string
	ClientSecret    string
	RedirectURL     string
	TokenSecret     string
	TokenExpiration time.Duration
}

// AuthService provides authentication utilities
type AuthService struct {
	oauthServer *server.Server
	config      AuthConfig
}

var dur7says, _ = time.ParseDuration("7d")
var DiscordConfig = AuthConfig{
	ClientId:        os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret:    os.Getenv("DISCORD_CLIENT_SECRET"),
	TokenSecret:     os.Getenv("DISCORD_REDIRECT_URL"),
	TokenExpiration: dur7says,
}
