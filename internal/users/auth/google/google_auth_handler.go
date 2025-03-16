package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nantestech/note-api/internal/users"
	auth "github.com/nantestech/note-api/internal/users/auth/schemas"
	"github.com/nantestech/note-api/pkg/jwt"
)

type GoogleAuthHandler struct {
	googleAuthService GoogleAuthService
	userRepo          users.Repository
	jwtConfig         jwt.Config
}

func newGoogleAuthHandler(googleAuthService GoogleAuthService, userRepo users.Repository, jwtConfig jwt.Config) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		googleAuthService: googleAuthService,
		userRepo:          userRepo,
		jwtConfig:         jwtConfig,
	}
}

func getRequestString(c *gin.Context) string {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return ""
	}

	return code
}

func (h *GoogleAuthHandler) createAuthResponse(token string, profilePic string, user *users.User) auth.AuthResponse {
	return auth.AuthResponse{
		Token:          token,
		UserId:         user.ID.String(),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		ProfilePicture: profilePic,
		Provider:       "google",
	}
}

func (h *GoogleAuthHandler) HandleGoogleAuth(c *gin.Context) {
	code := getRequestString(c)

	payload, err := h.googleAuthService.GetAuth(c.Request.Context(), code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByEmail(c.Request.Context(), payload.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		user = users.NewUser(payload.GivenName, payload.FamilyName, payload.Email)
		err = h.userRepo.Add(c.Request.Context(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	token, err := jwt.GenerateToken(h.jwtConfig, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.createAuthResponse(token, payload.Picture, user)
	c.JSON(http.StatusOK, response)
}
