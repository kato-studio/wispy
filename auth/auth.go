package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/kato-studio/wispy/auth/users"
	"github.com/kato-studio/wispy/utilities"
)

type AuthManager struct {
	Manager     *manage.Manager
	Database    *sql.DB
	UserStorage users.UserStorage
	Services    map[string](*AuthService)
}

func (am *AuthManager) Init(userStorage users.UserStorage, Database *sql.DB) {
	//
	am.UserStorage = userStorage
	am.Database = Database
	//
	am.Manager.MustTokenStorage(store.NewFileTokenStore("./db/auth-session-tokens.buntdb"))
	//
	am.Manager = manage.NewManager()
	// default implementation
	am.Manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	am.Manager.MapAccessGenerate(generates.NewAccessGenerate())
}

// NewAuthService creates a new authentication service
func (am *AuthManager) NewAuthService(serviceName string, config AuthConfig) {
	clientStore := store.NewClientStore()
	clientStore.Set(config.ClientId, &models.Client{
		ID:     config.ClientId,
		Secret: config.ClientSecret,
		Domain: config.RedirectURL,
	})
	am.Manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(am.Manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	am.Services[serviceName] = &AuthService{
		oauthServer: srv,
		config:      config,
	}
}

// Verify Token from session cookie
func (am *AuthManager) VarifySession(r *http.Request) (bool, error) {
	// Get session cookie
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return false, err
	}

	// Verify we have a valid session token
	tokenInfo, err := am.Manager.LoadAccessToken(r.Context(), sessionCookie.Value)
	if err != nil {
		return false, err
	}

	tokenInfo.GetAccess()

	return true, err
}

// GetUserFromSession retrieves the authenticated user from the session cookie
func (am *AuthManager) GetUserFromSession(r *http.Request) (*users.User, error) {
	// Get session cookie
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return nil, err
	}

	// Verify we have a valid session token
	tokenInfo, err := am.Manager.LoadAccessToken(r.Context(), sessionCookie.Value)
	if err != nil {
		return nil, err
	}

	// Get user ID from token
	userID := tokenInfo.GetUserID()
	if userID == "" {
		return nil, errors.New("no user ID in token")
	}

	// Get user from storage
	user, err := am.UserStorage.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// CreateSession creates a new authenticated session for the user
func (am *AuthManager) CreateSession(w http.ResponseWriter, r *http.Request, user *users.User) error {
	// Generate new access token
	token, err := am.Manager.GenerateAccessToken(
		r.Context(),
		oauth2.PasswordCredentials,
		&oauth2.TokenGenerateRequest{
			UserID: user.ID.String(),
		},
	)
	if err != nil {
		return err
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-session",
		Value:    token.GetAccess(),
		Path:     "/",
		MaxAge:   int((time.Hour * 24 * 7).Seconds()), // 1 week
		HttpOnly: true,
		Secure:   true, // Enable in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// DestroySession removes the authenticated session
func (am *AuthManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	sessionCookie, err := r.Cookie("auth-session")
	if err != nil {
		return nil // No session to destroy
	}

	// Remove token from storage
	err = am.Manager.RemoveAccessToken(r.Context(), sessionCookie.Value)
	if err != nil {
		return err
	}

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   utilities.IsProduction(),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
