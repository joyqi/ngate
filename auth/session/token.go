package session

// Token represents a standard token structure that can be stored in a session
type Token struct {
	// AccessToken represents the access token
	AccessToken string `json:"access_token"`

	// RefreshToken represents the refresh token
	RefreshToken string `json:"refresh_token"`

	// ExpiresAt is the time at which the token expires
	ExpiresAt int64 `json:"expires_at"`

	// Groups is a list of groups
	Groups []string `json:"groups,flow"`
}
