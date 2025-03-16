package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nantestech/note-api/internal/users"
)

type Claims struct {
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	UserID       uuid.UUID `json:"userId"`
	IsPremium    bool      `json:"isPremium"`
	PremiumUntil string    `json:"premiumUntil"`
	jwt.RegisteredClaims
}

type Config struct {
	SecretKey      string
	Issuer         string
	Audience       string
	ExpiresInHours int //hours
}

func GenerateToken(config Config, user *users.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.ExpiresInHours))

	claims := &Claims{
		Email:        user.Email,
		Name:         user.FirstName + " " + user.LastName,
		UserID:       user.ID,
		IsPremium:    user.IsPremium(),
		PremiumUntil: user.PremiumUntil.Format(time.RFC3339),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
			Audience:  []string{config.Audience},
			ID:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))

	return tokenString, err
}

func ValidateToken(config Config, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
