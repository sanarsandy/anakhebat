package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AdminRequestOTP handles OTP request for admin login
func AdminRequestOTP(c echo.Context) error {
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

	// Check if user is registered AND is admin
	var userID string
	var role string
	checkUserQuery := `SELECT id, role FROM users WHERE phone_number = $1`
	err = db.DB.QueryRow(checkUserQuery, phoneNumber).Scan(&userID, &role)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "Nomor WhatsApp tidak terdaftar sebagai admin",
		})
	} else if err != nil {
		c.Logger().Errorf("Admin user check error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error",
		})
	}

	// Check if user is admin
	if role != "admin" && role != "super_admin" {
		return c.JSON(http.StatusForbidden, models.OTPResponse{
			Success: false,
			Error:   "Akses ditolak. Hanya admin yang dapat login.",
		})
	}

	// Check rate limiting (max 3 requests per 15 minutes)
	canRequest, retryAfter, err := checkAdminRateLimit(phoneNumber)
	if err != nil {
		c.Logger().Errorf("Rate limit check error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error",
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

	// Store OTP in database with purpose 'admin_login'
	expiresAt := utils.GetOTPExpiration()
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	query := `INSERT INTO otp_codes (phone_number, otp_code, purpose, expires_at, ip_address, user_agent, max_attempts)
	          VALUES ($1, $2, 'admin_login', $3, $4, $5, 3)
	          RETURNING id`
	var otpID string
	err = db.DB.QueryRow(query, phoneNumber, otpCode, expiresAt, ipAddress, userAgent).Scan(&otpID)
	if err != nil {
		c.Logger().Errorf("Failed to store OTP: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to store OTP",
		})
	}

	// Update rate limit
	_ = updateAdminRateLimit(phoneNumber)

	// Send OTP via WhatsApp
	err = whatsappService.SendOTP(phoneNumber, otpCode)
	if err != nil {
		c.Logger().Errorf("Failed to send WhatsApp OTP for admin: %v", err)
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "Lid is missing in chat table") ||
			strings.Contains(errorMsg, "chat table") ||
			strings.Contains(errorMsg, "belum pernah") {
			return c.JSON(http.StatusBadRequest, models.OTPResponse{
				Success: false,
				Error:   "Nomor WhatsApp belum pernah mengirim pesan ke bot. Silakan kirim pesan ke nomor bot terlebih dahulu.",
			})
		}
		return c.JSON(http.StatusBadGateway, models.OTPResponse{
			Success: false,
			Error:   "Gagal mengirim OTP melalui WhatsApp. Silakan coba lagi nanti.",
		})
	}

	return c.JSON(http.StatusOK, models.OTPResponse{
		Success:   true,
		Message:   "OTP telah dikirim ke WhatsApp Anda",
		ExpiresIn: 300, // 5 minutes
	})
}

// AdminVerifyOTP handles OTP verification for admin login
func AdminVerifyOTP(c echo.Context) error {
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

	// Find OTP record with purpose 'admin_login'
	var otp models.OTPCode
	query := `SELECT id, phone_number, otp_code, purpose, expires_at, used_at, is_used, attempt_count, max_attempts
	          FROM otp_codes
	          WHERE phone_number = $1 AND otp_code = $2 AND purpose = 'admin_login' AND is_used = false
	          ORDER BY created_at DESC
	          LIMIT 1`

	var usedAt sql.NullTime
	err = db.DB.QueryRow(query, phoneNumber, req.OTP).Scan(
		&otp.ID, &otp.PhoneNumber, &otp.OTPCode, &otp.Purpose,
		&otp.ExpiresAt, &usedAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempts)

	if err == sql.ErrNoRows {
		incrementAdminOTPAttempt(phoneNumber)
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "OTP tidak valid",
		})
	} else if err != nil {
		c.Logger().Errorf("Admin OTP verification query error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error",
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

	// Find admin user
	user, err := findAdminByPhone(phoneNumber)
	if err != nil {
		c.Logger().Errorf("Failed to find admin user: %v", err)
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "Admin tidak ditemukan",
		})
	}

	// Generate JWT
	token, err := generateAdminJWT(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
	}

	// Log admin login for audit
	logAdminLogin(user.ID, c.RealIP(), c.Request().UserAgent())

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Helper functions for admin OTP

func findAdminByPhone(phoneNumber string) (*models.User, error) {
	var user models.User
	var email sql.NullString
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider, 
	          phone_verified, phone_verified_at, created_at
	          FROM users WHERE phone_number = $1 AND (role = 'admin' OR role = 'super_admin')`
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

func checkAdminRateLimit(phoneNumber string) (bool, int, error) {
	windowStart := time.Now().Add(-15 * time.Minute)
	windowEnd := time.Now()

	var count int
	query := `SELECT COALESCE(SUM(request_count), 0) 
	          FROM otp_rate_limits 
	          WHERE phone_number = $1 AND window_start >= $2 AND window_end <= $3`
	err := db.DB.QueryRow(query, phoneNumber, windowStart, windowEnd).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, 0, err
	}

	if count >= 3 {
		retryAfter := int(time.Until(windowEnd.Add(15 * time.Minute)).Seconds())
		return false, retryAfter, nil
	}

	return true, 0, nil
}

func updateAdminRateLimit(phoneNumber string) error {
	windowStart := time.Now().Truncate(15 * time.Minute)
	windowEnd := windowStart.Add(15 * time.Minute)

	query := `INSERT INTO otp_rate_limits (phone_number, request_count, window_start, window_end, updated_at)
	          VALUES ($1, 1, $2, $3, $4)
	          ON CONFLICT (phone_number, window_start) 
	          DO UPDATE SET request_count = otp_rate_limits.request_count + 1, updated_at = $4`
	_, err := db.DB.Exec(query, phoneNumber, windowStart, windowEnd, time.Now())
	return err
}

func incrementAdminOTPAttempt(phoneNumber string) {
	query := `UPDATE otp_codes 
	          SET attempt_count = attempt_count + 1 
	          WHERE phone_number = $1 AND purpose = 'admin_login' AND is_used = false 
	          ORDER BY created_at DESC 
	          LIMIT 1`
	db.DB.Exec(query, phoneNumber)
}

func generateAdminJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["is_admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 24 hours for admin

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	return token.SignedString([]byte(jwtSecret))
}

func logAdminLogin(userID, ipAddress, userAgent string) {
	query := `INSERT INTO audit_logs (user_id, action, resource_type, ip_address, user_agent, details)
	          VALUES ($1, 'admin_login', 'auth', $2, $3, '{"method": "otp"}')`
	db.DB.Exec(query, userID, ipAddress, userAgent)
}

