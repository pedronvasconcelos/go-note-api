package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nantestech/note-api/internal/users"
	"github.com/nantestech/note-api/pkg/jwt"
)

type GoogleAuthService interface {
	GetAuth(ctx context.Context, token string) (*GoogleTokenPayload, error)
}

type googleAuthService struct {
	config     *GoogleAuthConfig
	jwtConfig  jwt.Config
	userRepo   users.Repository
	httpClient *http.Client
}

// GetAuth implements GoogleAuthService.
func (g *googleAuthService) GetAuth(ctx context.Context, token string) (*GoogleTokenPayload, error) {
	payload, err := g.validateGoogleToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewGoogleAuthService(config *GoogleAuthConfig, jwtConfig jwt.Config, userRepo users.Repository) GoogleAuthService {
	return &googleAuthService{
		config:     config,
		jwtConfig:  jwtConfig,
		userRepo:   userRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *googleAuthService) validateGoogleToken(ctx context.Context, token string) (*GoogleTokenPayload, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to verify token: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload GoogleTokenPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	if payload.Email == "" || payload.VerifiedEmail != "true" {
		return nil, errors.New("invalid token: email not verified")
	}

	return &payload, nil
}
