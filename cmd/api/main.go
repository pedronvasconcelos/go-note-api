package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	middleware "github.com/nantestech/note-api/internal/api/middlewares"
	"github.com/nantestech/note-api/internal/infra/postgres"
	"github.com/nantestech/note-api/internal/users"
	auth "github.com/nantestech/note-api/internal/users/auth/google"
	"github.com/nantestech/note-api/pkg/jwt"
	"gorm.io/gorm"
)

func main() {

	if os.Getenv("APP_ENV") != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	setupMiddlewares(router)
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	setupRoutes(router)
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupMiddlewares(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
}

func setupRoutes(router *gin.Engine) {

	db := setupDB()
	userRepo := users.NewUserRepository(db)
	jwtConfig := setupJWT()
	googleAuthConfig := setupGoogleAuthConfig()
	googleAuthService := auth.NewGoogleAuthService(googleAuthConfig, jwtConfig, userRepo)
	googleAuthHandler := auth.NewGoogleAuthHandler(googleAuthService, userRepo, jwtConfig)
	auth.GoogleAuthRoutes(router, googleAuthHandler)

	authMiddleware := middleware.NewAuthMiddleware(jwtConfig)
	api := router.Group("/api")
	api.Use(authMiddleware.Authenticate())
	{
		api.GET("/auth/validate-jwt", auth.ValidateJWT())
	}

}

func setupGoogleAuthConfig() auth.GoogleAuthConfig {
	googleAuthConfig := auth.GoogleAuthConfig{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
	}
	return googleAuthConfig
}

func setupDB() *gorm.DB {
	dbConfig := postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "note"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
	db, err := postgres.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func setupJWT() jwt.Config {
	jwtConfig := jwt.Config{
		SecretKey:      getEnv("JWT_SECRET", "your-secret-key"),
		Issuer:         getEnv("JWT_ISSUER", "note"),
		Audience:       getEnv("JWT_AUDIENCE", "note-web"),
		ExpiresInHours: getEnvAsInt("JWT_EXPIRES_IN", 140), // 140 hours
	}
	return jwtConfig
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
