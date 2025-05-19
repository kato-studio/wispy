package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// -- GenerateRandomString creates a random string for session tokens
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Password Utilities
// -- HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// -- CheckPassword compares a password with its hashed version
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// OAuth2 Utilities
// // -- HandleDiscordLogin initiates the Discord OAuth2 flow
// func (a *AuthService) HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
// 	err := a.oauthServer.HandleAuthorizeRequest(w, r)
// 	if err != nil {
// 		http.Error(w, "Failed to start OAuth flow", http.StatusInternalServerError)
// 		return
// 	}
// }
