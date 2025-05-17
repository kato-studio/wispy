package users

// DiscordUser represents user data from Discord OAuth
type DiscordUser struct {
	ID       string
	Username string
	Email    string
	Avatar   string
}
