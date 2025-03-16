package auth

type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type GoogleTokenPayload struct {
	Email         string `json:"email"`
	VerifiedEmail string `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}
