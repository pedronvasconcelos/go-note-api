package auth

import (
	"github.com/gin-gonic/gin"
)

func GoogleAuthRoutes(router *gin.Engine, googleAuthHandler *GoogleAuthHandler) {

	auth := router.Group("/auth")
	{
		auth.GET("/google", googleAuthHandler.HandleGoogleAuth)
	}
}

// func returning string as handlerfunc
func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "JWT is valid")
	}
}
