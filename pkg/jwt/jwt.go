package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

func GenerateToken(config Config, userID uuid.UUID, name, email string, isPremium bool, premiumUntil string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.ExpiresInHours))

	claims := &Claims{
		Email:        email,
		Name:         name,
		UserID:       userID,
		IsPremium:    isPremium,
		PremiumUntil: premiumUntil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
			Audience:  []string{config.Audience},
			ID:        userID.String(),
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
