package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var googleOAuthConfig *oauth2.Config

func init() {
	// Initialize Google OAuth config
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	
	if clientID == "" || clientSecret == "" {
		// OAuth not configured, handlers will return error
		return
	}
	
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/google/callback"
	}

	// Manual Google OAuth endpoint configuration (to avoid dependency issues)
	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func GetGoogleAuthURL(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	// Generate state token for CSRF protection
	state := generateStateToken()
	
	// Store state in session/cookie (for production, use proper session storage)
	// For now, we'll include it in the URL and verify it in callback
	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	
	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": url,
		"state":    state,
	})
}

// GoogleAuthCallback handles the OAuth callback from Google
func GoogleAuthCallback(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")
	
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization code not provided",
		})
	}

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to exchange token: " + err.Error(),
		})
	}

	// Get user info from Google
	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info: " + err.Error(),
		})
	}

	// Find or create user
	user, err := findOrCreateGoogleUser(userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create/find user: " + err.Error(),
		})
	}

	// Generate JWT
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development
	}

	t, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	
	// Encode user data as JSON for frontend
	userJSON, _ := json.Marshal(user)
	userData := string(userJSON)
	
	// Use HTML redirect with token and user data
	redirectHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Redirecting...</title>
		</head>
		<body>
			<script>
				const token = '%s';
				const userData = %s;
				const state = '%s';
				window.location.href = '%s/auth/google/callback?token=' + encodeURIComponent(token) + '&user=' + encodeURIComponent(JSON.stringify(userData)) + '&state=' + encodeURIComponent(state);
			</script>
			<p>Redirecting...</p>
		</body>
		</html>
	`, t, userData, state, frontendURL)

	return c.HTML(http.StatusOK, redirectHTML)
}

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func findOrCreateGoogleUser(googleUser *GoogleUserInfo) (*models.User, error) {
	var user models.User
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	
	// Try to find user by Google ID
	var phoneNumber sql.NullString
	var phoneVerified sql.NullBool
	var phoneVerifiedAt sql.NullTime
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, 
	          phone_number, phone_verified, phone_verified_at, created_at 
	          FROM users WHERE google_id = $1`
	err := db.DB.QueryRow(query, googleUser.ID).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&user.Role, &googleID, &authProvider, &phoneNumber, &phoneVerified, &phoneVerifiedAt, &user.CreatedAt)
	
	if err == nil {
		// User found, set optional fields
		if passwordHash.Valid {
			user.PasswordHash = passwordHash.String
		}
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if authProvider.Valid {
			user.AuthProvider = authProvider.String
		} else {
			user.AuthProvider = "google"
		}
		if phoneNumber.Valid {
			user.PhoneNumber = &phoneNumber.String
		}
		if phoneVerified.Valid {
			user.PhoneVerified = phoneVerified.Bool
		}
		if phoneVerifiedAt.Valid {
			user.PhoneVerifiedAt = &phoneVerifiedAt.Time
		}
		return &user, nil
	}
	
	if err != sql.ErrNoRows {
		// Database error
		return nil, err
	}
	
	// User not found, try to find by email
	query = `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, 
	         phone_number, phone_verified, phone_verified_at, created_at 
	         FROM users WHERE email = $1`
	err = db.DB.QueryRow(query, googleUser.Email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&user.Role, &googleID, &authProvider, &phoneNumber, &phoneVerified, &phoneVerifiedAt, &user.CreatedAt)
	
	if err == nil {
		// User exists with email
		// Check if this is a Google user (auth_provider = 'google')
		if authProvider.Valid && authProvider.String == "google" {
			// This is already a Google user, verify Google ID matches
			if googleID.Valid && googleID.String == googleUser.ID {
				// Same Google account, return user
				if passwordHash.Valid {
					user.PasswordHash = passwordHash.String
				}
				user.GoogleID = &googleID.String
				user.AuthProvider = "google"
				if phoneNumber.Valid {
					user.PhoneNumber = &phoneNumber.String
				}
				if phoneVerified.Valid {
					user.PhoneVerified = phoneVerified.Bool
				}
				if phoneVerifiedAt.Valid {
					user.PhoneVerifiedAt = &phoneVerifiedAt.Time
				}
				return &user, nil
			}
			// Different Google ID but same email - this shouldn't happen, but handle it
			return nil, fmt.Errorf("Email already registered with different Google account")
		}
		
		// User exists but is NOT a Google user (auth_provider = 'email' or 'both')
		// Don't create separate account - return error that email is already registered
		// User should use email/password login instead
		return nil, fmt.Errorf("Email sudah terdaftar dengan akun email/password. Silakan login dengan email dan password")
	}
	
	if err != sql.ErrNoRows {
		// Database error
		return nil, err
	}
	
	// Email doesn't exist, create new Google user
	fullName := googleUser.Name
	if fullName == "" {
		fullName = googleUser.GivenName + " " + googleUser.FamilyName
	}
	if fullName == "" {
		fullName = googleUser.Email
	}
	
	insertQuery := `INSERT INTO users (email, full_name, google_id, auth_provider, role) 
	                VALUES ($1, $2, $3, 'google', 'parent') 
	                RETURNING id, email, password_hash, full_name, role, google_id, auth_provider, created_at`
	var newPasswordHash sql.NullString
	var newGoogleID sql.NullString
	var newAuthProvider sql.NullString
	err = db.DB.QueryRow(insertQuery, googleUser.Email, fullName, googleUser.ID).Scan(
		&user.ID, &user.Email, &newPasswordHash, &user.FullName, 
		&user.Role, &newGoogleID, &newAuthProvider, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	// Set optional fields
	if newPasswordHash.Valid {
		user.PasswordHash = newPasswordHash.String
	}
	if newGoogleID.Valid {
		user.GoogleID = &newGoogleID.String
	}
	if newAuthProvider.Valid {
		user.AuthProvider = newAuthProvider.String
	} else {
		user.AuthProvider = "google"
	}
	
	return &user, nil
}

func generateStateToken() string {
	// Simple state token generation (for production, use crypto/rand)
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

