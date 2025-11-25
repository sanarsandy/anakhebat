package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	req := new(models.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate and normalize phone number
	phoneNumber, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Check if phone number already exists
	var existingID string
	checkQuery := `SELECT id FROM users WHERE phone_number = $1`
	err = db.DB.QueryRow(checkQuery, phoneNumber).Scan(&existingID)
	if err == nil {
		// Phone number already registered
		return c.JSON(http.StatusConflict, map[string]string{"error": "Nomor WhatsApp sudah terdaftar"})
	} else if err != sql.ErrNoRows {
		// Database error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}

	// Insert into DB - phone number only, no password
	query := `INSERT INTO users (phone_number, full_name, role, auth_provider, phone_verified) 
	          VALUES ($1, $2, 'parent', 'phone', false) 
	          RETURNING id, phone_number, full_name, role, auth_provider, created_at`
	var user models.User
	var phoneNum sql.NullString
	var authProvider sql.NullString
	
	err = db.DB.QueryRow(query, phoneNumber, req.FullName).Scan(
		&user.ID, &phoneNum, &user.FullName, &user.Role, &authProvider, &user.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
	}

	if phoneNum.Valid {
		user.PhoneNumber = &phoneNum.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "phone"
	}

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Find user - only look for email/password users (auth_provider = 'email' or 'both')
	// Google-only users should not be found here
	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider 
	          FROM users 
	          WHERE email = $1 AND (auth_provider = 'email' OR auth_provider = 'both')`
	err := db.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &googleID, &authProvider)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User belum terdaftar"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Set optional fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "email"
	}

	// Check if user has password (for 'email' or 'both' auth_provider)
	if user.PasswordHash == "" || user.PasswordHash == "NULL" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: t,
		User:  user,
	})
}
