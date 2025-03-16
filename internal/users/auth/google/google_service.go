package auth

import (
	"context"
	"time"

	"net/http"

	"github.com/nantestech/note-api/internal/users"
	"github.com/nantestech/note-api/internal/users/"
	"github.com/nantestech/note-api/pkg/jwt"
)

type GoogleAuthService interface {
	GetAuth(ctx context.Context, token string) (*string, error)
}

type googleAuthService struct {
	config     *GoogleAuthConfig
	jwtConfig  jwt.Config
	userRepo   users.Repository
	httpClient *http.Client
}

func NewGoogleAuthService(config *GoogleAuthConfig, jwtConfig jwt.Config, userRepo users.Repository) GoogleAuthService {
	return &googleAuthService{
		config:     config,
		jwtConfig:  jwtConfig,
		userRepo:   userRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
