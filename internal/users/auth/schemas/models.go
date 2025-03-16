package auth

type JWTConfig struct {
	SecretKey      string
	Issuer         string
	ExpiresInHours int
	Audience       string
}

type AuthResponse struct {
	Token          string `json:"token"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Provider       string `json:"provider"`
	UserId         string `json:"userId"`
}
