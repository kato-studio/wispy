package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	// "github.com/kato-studio/wispy/wisy_common"
)

// DiscordEndpoint returns the correct oauth2.Endpoint for Discord
func DiscordEndpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:   "https://discord.com/oauth2/authorize", // Consistent with their docs
		TokenURL:  "https://discord.com/api/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
}

func GetDiscordOAuthConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		Scopes:       []string{"identify", "email"},
		Endpoint:     DiscordEndpoint(),
	}

	return config
}

func HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	// Generate random state for CSRF protection
	state, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Store state in session
	stateCookie := &http.Cookie{
		Name:  "oauth_state",
		Value: state,
		// Path:     "/",
		// MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
		// Secure:   wispy_common.IsProduction(),
		// SameSite: http.SameSiteNoneMode,
	}
	fmt.Println("Setting cookie!")
	fmt.Println(stateCookie)
	http.SetCookie(w, stateCookie)

	// Redirect to Discord OAuth
	// Generate the auth URL with all required parameters
	// Add prompt=consent to force fresh permissions
	authURL := GetDiscordOAuthConfig().AuthCodeURL(state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func HandleDiscordCallback(SessionDB, UserDB *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Verify state matches
	fmt.Println("--- Cookies ---")
	fmt.Println(r.Cookies())
	fmt.Println("------")
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "Missing state cookie", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	ctx := context.Background()
	token, err := GetDiscordOAuthConfig().Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info from Discord
	client := GetDiscordOAuthConfig().Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var discordUser DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Find or create user in database
	user, err := FindOrCreateDiscordUser(UserDB, discordUser)
	if err != nil {
		http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	if err := CreateSession(SessionDB, w, r, user); err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// // Clear the state cookie immediately after reading it
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func FindOrCreateDiscordUser(DB *sql.DB, discordUser DiscordUser) (*User, error) {
	var user User

	// Check if we already have this Discord account linked
	var authMethodID int
	err := DB.QueryRow(`
		SELECT id FROM auth_methods WHERE name = 'discord'
	`).Scan(&authMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discord auth method: %w", err)
	}

	// Look for existing provider link
	var userUUID string
	err = DB.QueryRow(`
		SELECT user_uuid FROM user_auth_providers
		WHERE auth_method_id = ? AND provider_id = ?
	`, authMethodID, discordUser.ID).Scan(&userUUID)

	if err == nil {
		// Found existing user
		return GetUserByID(DB, userUUID)
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query user auth providers: %w", err)
	}

	// Create new user
	username := discordUser.Username
	if discordUser.Discriminator != "0" {
		username = fmt.Sprintf("%s#%s", discordUser.Username, discordUser.Discriminator)
	}

	user = User{
		UUID:     uuid.New().String(),
		Username: username,
		Email:    discordUser.Email,
	}

	// Start transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user
	err = tx.QueryRow(`
		INSERT INTO users (uuid, username, email)
		VALUES (?, ?, ?)
		RETURNING uuid
	`, user.UUID, user.Username, user.Email).Scan(&user.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Link Discord account
	_, err = tx.Exec(`
		INSERT INTO user_auth_providers (
			user_uuid, auth_method_id, provider_id,
			provider_username, provider_email
		) VALUES (?, ?, ?, ?, ?)
	`, user.UUID, authMethodID, discordUser.ID,
		discordUser.Username, discordUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to link discord account: %w", err)
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
