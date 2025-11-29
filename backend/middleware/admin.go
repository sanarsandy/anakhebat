package middleware

import (
	"net/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// RequireAdmin checks if user has admin role
// This middleware should be used after JWTMiddleware
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user from JWT token (set by JWTMiddleware)
			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			token, ok := user.(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// Extract claims - echo-jwt stores claims as *jwt.MapClaims
			var claims jwt.MapClaims
			if ptrClaims, ok := token.Claims.(*jwt.MapClaims); ok {
				claims = *ptrClaims
			} else if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
				claims = mapClaims
			} else {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// Check role
			role, ok := claims["role"].(string)
			if !ok || role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Admin access required",
				})
			}

			// Store user info in context for handlers
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("user_email", email)
			}
			c.Set("user_role", role)

			return next(c)
		}
	}
}

