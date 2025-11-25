package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/services"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var whatsappService = services.NewWhatsAppService()

// RequestOTP handles OTP request
func RequestOTP(c echo.Context) error {
	req := new(models.RequestOTPRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// Validate and normalize phone number
	phoneNumber, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Check if user is registered
	var userID string
	checkUserQuery := `SELECT id FROM users WHERE phone_number = $1`
	err = db.DB.QueryRow(checkUserQuery, phoneNumber).Scan(&userID)
	if err == sql.ErrNoRows {
		// User not registered
		return c.JSON(http.StatusNotFound, models.OTPResponse{
			Success: false,
			Error:   "User belum terdaftar. Silakan daftar terlebih dahulu.",
		})
	} else if err != nil {
		c.Logger().Errorf("User check error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
	}

	// Check rate limiting (max 3 requests per 15 minutes)
	canRequest, retryAfter, err := checkRateLimit(phoneNumber)
	if err != nil {
		c.Logger().Errorf("Rate limit check error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
	}

	if !canRequest {
		return c.JSON(http.StatusTooManyRequests, models.OTPResponse{
			Success:    false,
			Error:      "Terlalu banyak permintaan. Silakan coba lagi nanti",
			RetryAfter: retryAfter,
		})
	}

	// Generate OTP
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to generate OTP",
		})
	}

	// Store OTP in database
	expiresAt := utils.GetOTPExpiration()
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	query := `INSERT INTO otp_codes (phone_number, otp_code, purpose, expires_at, ip_address, user_agent, max_attempts)
	          VALUES ($1, $2, 'login', $3, $4, $5, 3)
	          RETURNING id`
	var otpID string
	err = db.DB.QueryRow(query, phoneNumber, otpCode, expiresAt, ipAddress, userAgent).Scan(&otpID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to store OTP",
		})
	}

	// Update rate limit
	err = updateRateLimit(phoneNumber)
	if err != nil {
		// Log error but don't fail the request
		// OTP is already sent
	}

	// Send OTP via WhatsApp
	err = whatsappService.SendOTP(phoneNumber, otpCode)
	if err != nil {
		// Log error - this is important for debugging
		c.Logger().Errorf("Failed to send WhatsApp OTP: %v", err)
		// Still return success since OTP is stored in database
		// User can still verify OTP manually if needed
		// In production, you might want to handle this differently
	}

	return c.JSON(http.StatusOK, models.OTPResponse{
		Success:   true,
		Message:   "OTP telah dikirim ke WhatsApp Anda",
		ExpiresIn: 300, // 5 minutes
	})
}

// VerifyOTP handles OTP verification and login/registration
func VerifyOTP(c echo.Context) error {
	req := new(models.VerifyOTPRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// Validate and normalize phone number
	phoneNumber, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Find OTP record
	var otp models.OTPCode
	query := `SELECT id, phone_number, otp_code, purpose, expires_at, used_at, is_used, attempt_count, max_attempts
	          FROM otp_codes
	          WHERE phone_number = $1 AND otp_code = $2 AND is_used = false
	          ORDER BY created_at DESC
	          LIMIT 1`
	
	var usedAt sql.NullTime
	err = db.DB.QueryRow(query, phoneNumber, req.OTP).Scan(
		&otp.ID, &otp.PhoneNumber, &otp.OTPCode, &otp.Purpose,
		&otp.ExpiresAt, &usedAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempts)
	
	if err == sql.ErrNoRows {
		// OTP not found - increment attempt count on latest OTP for this phone
		incrementOTPAttempt(phoneNumber)
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "OTP tidak valid",
		})
	} else if err != nil {
		c.Logger().Errorf("OTP verification query error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
	}

	// Check if OTP is expired
	if utils.IsOTPExpired(otp.ExpiresAt) {
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "OTP sudah kadaluarsa",
		})
	}

	// Check if max attempts exceeded
	if otp.AttemptCount >= otp.MaxAttempts {
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "Terlalu banyak percobaan. Silakan request OTP baru",
		})
	}

	// Mark OTP as used
	now := time.Now()
	updateQuery := `UPDATE otp_codes SET is_used = true, used_at = $1 WHERE id = $2`
	_, err = db.DB.Exec(updateQuery, now, otp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to update OTP",
		})
	}

	// Find user (must be registered)
	user, err := findPhoneUser(phoneNumber)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.OTPResponse{
			Success: false,
			Error:   "User belum terdaftar. Silakan daftar terlebih dahulu.",
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to find user: " + err.Error(),
		})
	}

	// Generate JWT
	token, err := generateJWT(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// ResendOTP handles OTP resend request
func ResendOTP(c echo.Context) error {
	req := new(models.RequestOTPRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// Validate and normalize phone number
	phoneNumber, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.OTPResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Check if user is registered
	var userID string
	checkUserQuery := `SELECT id FROM users WHERE phone_number = $1`
	err = db.DB.QueryRow(checkUserQuery, phoneNumber).Scan(&userID)
	if err == sql.ErrNoRows {
		// User not registered
		return c.JSON(http.StatusNotFound, models.OTPResponse{
			Success: false,
			Error:   "User belum terdaftar. Silakan daftar terlebih dahulu.",
		})
	} else if err != nil {
		c.Logger().Errorf("User check error (resend): %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
	}

	// Check rate limiting
	canRequest, retryAfter, err := checkRateLimit(phoneNumber)
	if err != nil {
		c.Logger().Errorf("Rate limit check error (resend): %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error: " + err.Error(),
		})
	}

	if !canRequest {
		return c.JSON(http.StatusTooManyRequests, models.OTPResponse{
			Success:    false,
			Error:      "Terlalu banyak permintaan. Silakan coba lagi nanti",
			RetryAfter: retryAfter,
		})
	}

	// Invalidate old OTPs for this phone number
	_, err = db.DB.Exec(`UPDATE otp_codes SET is_used = true WHERE phone_number = $1 AND is_used = false`, phoneNumber)
	if err != nil {
		// Log but continue
	}

	// Generate new OTP
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to generate OTP",
		})
	}

	// Store new OTP
	expiresAt := utils.GetOTPExpiration()
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	query := `INSERT INTO otp_codes (phone_number, otp_code, purpose, expires_at, ip_address, user_agent, max_attempts)
	          VALUES ($1, $2, 'login', $3, $4, $5, 3)
	          RETURNING id`
	var otpID string
	err = db.DB.QueryRow(query, phoneNumber, otpCode, expiresAt, ipAddress, userAgent).Scan(&otpID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to store OTP",
		})
	}

	// Update rate limit
	err = updateRateLimit(phoneNumber)
	if err != nil {
		// Log error but don't fail the request
	}

	// Send OTP via WhatsApp
	err = whatsappService.SendOTP(phoneNumber, otpCode)
	if err != nil {
		// Log error - this is important for debugging
		c.Logger().Errorf("Failed to send WhatsApp OTP (resend): %v", err)
		// Still return success since OTP is stored in database
	}

	return c.JSON(http.StatusOK, models.OTPResponse{
		Success:   true,
		Message:   "OTP baru telah dikirim ke WhatsApp Anda",
		ExpiresIn: 300,
	})
}

// Helper functions

// findPhoneUser finds a user by phone number (does not create new user)
func findPhoneUser(phoneNumber string) (*models.User, error) {
	var user models.User
	var email sql.NullString
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, 
	          phone_verified, phone_verified_at, created_at
	          FROM users WHERE phone_number = $1`
	err := db.DB.QueryRow(query, phoneNumber).Scan(
		&user.ID, &email, &passwordHash, &user.FullName, &user.Role,
		&googleID, &authProvider, &user.PhoneVerified, &phoneVerifiedAt, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	if email.Valid {
		user.Email = email.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "phone"
	}
	user.PhoneNumber = &phoneNumber

	// Update phone verification if not already verified
	if !user.PhoneVerified {
		now := time.Now()
		updateQuery := `UPDATE users SET phone_verified = true, phone_verified_at = $1 WHERE id = $2`
		db.DB.Exec(updateQuery, now, user.ID)
		user.PhoneVerified = true
		user.PhoneVerifiedAt = &now
	}

	return &user, nil
}

func checkRateLimit(phoneNumber string) (bool, int, error) {
	// Check requests in last 15 minutes
	windowStart := time.Now().Add(-15 * time.Minute)
	windowEnd := time.Now()

	var count int
	query := `SELECT COALESCE(SUM(request_count), 0) 
	          FROM otp_rate_limits 
	          WHERE phone_number = $1 AND window_start >= $2 AND window_end <= $3`
	err := db.DB.QueryRow(query, phoneNumber, windowStart, windowEnd).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		// Log the actual error for debugging
		return false, 0, fmt.Errorf("rate limit query failed: %w", err)
	}

	// Max 3 requests per 15 minutes
	if count >= 3 {
		// Calculate retry after (seconds until next window)
		retryAfter := int(time.Until(windowEnd.Add(15 * time.Minute)).Seconds())
		return false, retryAfter, nil
	}

	return true, 0, nil
}

func updateRateLimit(phoneNumber string) error {
	// Create or update rate limit record for current 15-minute window
	windowStart := time.Now().Truncate(15 * time.Minute)
	windowEnd := windowStart.Add(15 * time.Minute)

	query := `INSERT INTO otp_rate_limits (phone_number, request_count, window_start, window_end, updated_at)
	          VALUES ($1, 1, $2, $3, $4)
	          ON CONFLICT (phone_number, window_start) 
	          DO UPDATE SET request_count = otp_rate_limits.request_count + 1, updated_at = $4`
	_, err := db.DB.Exec(query, phoneNumber, windowStart, windowEnd, time.Now())
	return err
}

func incrementOTPAttempt(phoneNumber string) {
	query := `UPDATE otp_codes 
	          SET attempt_count = attempt_count + 1 
	          WHERE phone_number = $1 AND is_used = false 
	          ORDER BY created_at DESC 
	          LIMIT 1`
	db.DB.Exec(query, phoneNumber)
}

func findOrCreatePhoneUser(phoneNumber string) (*models.User, bool, error) {
	var user models.User
	var email sql.NullString
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	// Try to find user by phone number
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, 
	          phone_verified, phone_verified_at, created_at
	          FROM users WHERE phone_number = $1`
	err := db.DB.QueryRow(query, phoneNumber).Scan(
		&user.ID, &email, &passwordHash, &user.FullName, &user.Role,
		&googleID, &authProvider, &user.PhoneVerified, &phoneVerifiedAt, &user.CreatedAt)

	if err == nil {
		// User found
		if email.Valid {
			user.Email = email.String
		}
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if authProvider.Valid {
			user.AuthProvider = authProvider.String
		} else {
			user.AuthProvider = "phone"
		}
		user.PhoneNumber = &phoneNumber

		// Update phone verification if not already verified
		if !user.PhoneVerified {
			now := time.Now()
			updateQuery := `UPDATE users SET phone_verified = true, phone_verified_at = $1 WHERE id = $2`
			db.DB.Exec(updateQuery, now, user.ID)
			user.PhoneVerified = true
			user.PhoneVerifiedAt = &now
		}

		// Update auth_provider if needed
		if user.AuthProvider == "phone" || user.AuthProvider == "google" {
			// Keep as is
		} else if user.GoogleID != nil && user.AuthProvider != "phone_google" {
			updateQuery := `UPDATE users SET auth_provider = 'phone_google' WHERE id = $1`
			db.DB.Exec(updateQuery, user.ID)
			user.AuthProvider = "phone_google"
		}

		return &user, false, nil
	}

	if err != sql.ErrNoRows {
		return nil, false, err
	}

	// User not found, create new user
	now := time.Now()
	insertQuery := `INSERT INTO users (phone_number, full_name, role, auth_provider, phone_verified, phone_verified_at)
	                VALUES ($1, $2, 'parent', 'phone', true, $3)
	                RETURNING id, email, password_hash, full_name, role, google_id, auth_provider, 
	                          phone_verified, phone_verified_at, created_at`
	
	var newEmail sql.NullString
	var newPasswordHash sql.NullString
	var newGoogleID sql.NullString
	var newAuthProvider sql.NullString
	var newPhoneVerifiedAt sql.NullTime

	err = db.DB.QueryRow(insertQuery, phoneNumber, "User", now).Scan(
		&user.ID, &newEmail, &newPasswordHash, &user.FullName, &user.Role,
		&newGoogleID, &newAuthProvider, &user.PhoneVerified, &newPhoneVerifiedAt, &user.CreatedAt)

	if err != nil {
		return nil, false, err
	}

	user.PhoneNumber = &phoneNumber
	user.PhoneVerified = true
	user.PhoneVerifiedAt = &now
	user.AuthProvider = "phone"

	if newEmail.Valid {
		user.Email = newEmail.String
	}
	if newGoogleID.Valid {
		user.GoogleID = &newGoogleID.String
	}

	return &user, true, nil
}

func generateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	return token.SignedString([]byte(jwtSecret))
}

