package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nantestech/note-api/pkg/jwt"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authMiddleware struct {
	jwtConfig jwt.Config
}

func NewAuthMiddleware(jwtConfig jwt.Config) AuthMiddleware {
	return &authMiddleware{jwtConfig: jwtConfig}
}

func (m *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(m.jwtConfig, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Make user info available in request context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("isPremium", claims.IsPremium)

		// Check if user ID in token matches route parameter
		userIDParam := c.Param("userId")
		if userIDParam != "" && userIDParam != claims.UserID.String() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.Next()
	}
}
