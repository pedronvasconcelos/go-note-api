package auth

type JWTConfig struct {
	SecretKey      string
	Issuer         string
	ExpiresInHours int
	Audience       string
}

type AuthResponse struct {
	Token          string `json:"token"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Provider       string `json:"provider"`
}
