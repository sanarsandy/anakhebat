package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"tukem-backend/db"
	"tukem-backend/handlers"
	customMiddleware "tukem-backend/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Database
	db.Init()

	e := EchoServer()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	e.Logger.Fatal(e.Start(":" + port))
}

func EchoServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS Configuration
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		corsOrigins = append(corsOrigins, strings.Split(corsEnv, ",")...)
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to API",
			"status":  "healthy",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "up",
		})
	})

	// Auth Routes
	auth := e.Group("/api/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.GET("/google", handlers.GetGoogleAuthURL)
	auth.GET("/google/callback", handlers.GoogleAuthCallback)

	// Protected Routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware())
	api.GET("/me", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "You are authorized!",
		})
	})

	// Add your routes here

	return e
}

