package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/services"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserProfile returns the current user's profile
func GetUserProfile(c echo.Context) error {
	// Get user ID from JWT middleware
	userToken := c.Get("user").(*jwt.Token)
	claims := *userToken.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	var user models.User
	var email sql.NullString
	var phoneNumber sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	query := `SELECT id, email, phone_number, full_name, role, google_id, auth_provider, 
	          phone_verified, phone_verified_at, created_at
	          FROM users WHERE id = $1`
	err := db.DB.QueryRow(query, userID).Scan(
		&user.ID, &email, &phoneNumber, &user.FullName, &user.Role,
		&googleID, &authProvider, &user.PhoneVerified, &phoneVerifiedAt, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	} else if err != nil {
		c.Logger().Errorf("Get user profile error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// Set optional fields
	if email.Valid {
		user.Email = email.String
	}
	if phoneNumber.Valid {
		user.PhoneNumber = &phoneNumber.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "phone"
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	// Get statistics
	stats, err := getUserStatistics(userID)
	if err != nil {
		c.Logger().Warnf("Failed to get user statistics: %v", err)
		// Continue without stats
	}

	response := map[string]interface{}{
		"user":  user,
		"stats": stats,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateUserProfile updates user profile information
func UpdateUserProfile(c echo.Context) error {
	// Get user ID from JWT middleware
	userToken := c.Get("user").(*jwt.Token)
	claims := *userToken.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	req := new(models.UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	// Validate full name
	if req.FullName != "" && len(req.FullName) < 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Nama lengkap minimal 3 karakter",
		})
	}

	// Validate email if provided
	if req.Email != "" {
		if !utils.IsValidEmail(req.Email) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Format email tidak valid",
			})
		}

		// Check if email is already used by another user
		var existingUserID string
		checkEmailQuery := `SELECT id FROM users WHERE email = $1 AND id != $2`
		err := db.DB.QueryRow(checkEmailQuery, req.Email, userID).Scan(&existingUserID)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Email sudah digunakan oleh user lain",
			})
		} else if err != sql.ErrNoRows {
			c.Logger().Errorf("Check email error: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Database error",
			})
		}
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.FullName != "" {
		updates = append(updates, "full_name = $"+utils.IntToString(argIndex))
		args = append(args, req.FullName)
		argIndex++
	}

	if req.Email != "" {
		updates = append(updates, "email = $"+utils.IntToString(argIndex))
		args = append(args, req.Email)
		argIndex++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tidak ada data yang diupdate",
		})
	}

	// Add user ID to args
	args = append(args, userID)
	updateQuery := "UPDATE users SET " + utils.JoinStrings(updates, ", ") + " WHERE id = $" + utils.IntToString(argIndex)

	_, err := db.DB.Exec(updateQuery, args...)
	if err != nil {
		c.Logger().Errorf("Update profile error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengupdate profil",
		})
	}

	// Get updated user
	var user models.User
	var email sql.NullString
	var phoneNumber sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	selectQuery := `SELECT id, email, phone_number, full_name, role, google_id, auth_provider, 
	                phone_verified, phone_verified_at, created_at
	                FROM users WHERE id = $1`
	err = db.DB.QueryRow(selectQuery, userID).Scan(
		&user.ID, &email, &phoneNumber, &user.FullName, &user.Role,
		&googleID, &authProvider, &user.PhoneVerified, &phoneVerifiedAt, &user.CreatedAt)

	if err != nil {
		c.Logger().Errorf("Get updated user error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data user",
		})
	}

	// Set optional fields
	if email.Valid {
		user.Email = email.String
	}
	if phoneNumber.Valid {
		user.PhoneNumber = &phoneNumber.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "phone"
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profil berhasil diupdate",
		"user":    user,
	})
}

// RequestPhoneVerification requests OTP for phone number verification
func RequestPhoneVerification(c echo.Context) error {
	// Get user ID from JWT middleware
	userToken := c.Get("user").(*jwt.Token)
	claims := *userToken.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	req := new(models.VerifyPhoneRequest)
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

	// Check if phone number is already used by another user
	var existingUserID string
	checkPhoneQuery := `SELECT id FROM users WHERE phone_number = $1 AND id != $2`
	err = db.DB.QueryRow(checkPhoneQuery, phoneNumber, userID).Scan(&existingUserID)
	if err == nil {
		return c.JSON(http.StatusConflict, models.OTPResponse{
			Success: false,
			Error:   "Nomor WhatsApp sudah digunakan oleh user lain",
		})
	} else if err != sql.ErrNoRows {
		c.Logger().Errorf("Check phone number error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error",
		})
	}

	// Check rate limiting
	canRequest, retryAfter, err := checkRateLimit(phoneNumber)
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

	// Store OTP in database
	expiresAt := utils.GetOTPExpiration()
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	query := `INSERT INTO otp_codes (phone_number, otp_code, purpose, expires_at, ip_address, user_agent, max_attempts)
	          VALUES ($1, $2, 'verify_phone', $3, $4, $5, 3)
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
		c.Logger().Warnf("Failed to update rate limit: %v", err)
	}

	// Send OTP via WhatsApp (using service from otp_auth.go)
	// Note: whatsappService is declared in otp_auth.go
	err = services.NewWhatsAppService().SendOTP(phoneNumber, otpCode)
	if err != nil {
		c.Logger().Errorf("Failed to send WhatsApp OTP: %v", err)
		// Check if it's a gateway error (chat not initiated)
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "Lid is missing in chat table") || 
		   strings.Contains(errorMsg, "chat table") ||
		   strings.Contains(errorMsg, "belum pernah") {
			// This means the number hasn't chatted with the gateway bot yet
			// Return error with helpful message
			return c.JSON(http.StatusBadRequest, models.OTPResponse{
				Success: false,
				Error:   "Nomor WhatsApp Anda belum pernah mengirim pesan ke bot kami. Silakan kirim pesan ke nomor bot terlebih dahulu, kemudian coba request OTP lagi.",
			})
		}
		// For other gateway errors, return error but keep OTP in DB for manual verification
		return c.JSON(http.StatusBadGateway, models.OTPResponse{
			Success: false,
			Error:   "Gagal mengirim OTP melalui WhatsApp. Silakan coba lagi nanti atau hubungi administrator.",
		})
	}

	return c.JSON(http.StatusOK, models.OTPResponse{
		Success:   true,
		Message:   "OTP telah dikirim ke WhatsApp Anda",
		ExpiresIn: 300, // 5 minutes
	})
}

// ConfirmPhoneVerification confirms OTP and updates user's phone number
func ConfirmPhoneVerification(c echo.Context) error {
	// Get user ID from JWT middleware
	userToken := c.Get("user").(*jwt.Token)
	claims := *userToken.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	req := new(models.ConfirmPhoneVerificationRequest)
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
	          WHERE phone_number = $1 AND otp_code = $2 AND purpose = 'verify_phone' AND is_used = false
	          ORDER BY created_at DESC
	          LIMIT 1`

	var usedAt sql.NullTime
	err = db.DB.QueryRow(query, phoneNumber, req.OTP).Scan(
		&otp.ID, &otp.PhoneNumber, &otp.OTPCode, &otp.Purpose,
		&otp.ExpiresAt, &usedAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempts)

	if err == sql.ErrNoRows {
		// OTP not found - increment attempt count
		incrementOTPAttempt(phoneNumber)
		return c.JSON(http.StatusUnauthorized, models.OTPResponse{
			Success: false,
			Error:   "OTP tidak valid",
		})
	} else if err != nil {
		c.Logger().Errorf("OTP verification query error: %v", err)
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
	now := utils.GetCurrentTime()
	updateOTPQuery := `UPDATE otp_codes SET is_used = true, used_at = $1 WHERE id = $2`
	_, err = db.DB.Exec(updateOTPQuery, now, otp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Failed to update OTP",
		})
	}

	// Check if phone number is already used by another user
	var existingUserID string
	checkPhoneQuery := `SELECT id FROM users WHERE phone_number = $1 AND id != $2`
	err = db.DB.QueryRow(checkPhoneQuery, phoneNumber, userID).Scan(&existingUserID)
	if err == nil {
		return c.JSON(http.StatusConflict, models.OTPResponse{
			Success: false,
			Error:   "Nomor WhatsApp sudah digunakan oleh user lain",
		})
	} else if err != sql.ErrNoRows {
		c.Logger().Errorf("Check phone number error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Database error",
		})
	}

	// Update user's phone number and verification status
	updateUserQuery := `UPDATE users 
	                    SET phone_number = $1, phone_verified = true, phone_verified_at = $2
	                    WHERE id = $3`
	_, err = db.DB.Exec(updateUserQuery, phoneNumber, now, userID)
	if err != nil {
		c.Logger().Errorf("Update user phone error: %v", err)
		return c.JSON(http.StatusInternalServerError, models.OTPResponse{
			Success: false,
			Error:   "Gagal mengupdate nomor WhatsApp",
		})
	}

	// Update auth_provider if needed
	var currentAuthProvider string
	getAuthProviderQuery := `SELECT auth_provider FROM users WHERE id = $1`
	err = db.DB.QueryRow(getAuthProviderQuery, userID).Scan(&currentAuthProvider)
	if err == nil {
		if currentAuthProvider == "google" {
			// User logged in via Google, now has phone too
			updateAuthProviderQuery := `UPDATE users SET auth_provider = 'phone_google' WHERE id = $1`
			db.DB.Exec(updateAuthProviderQuery, userID)
		} else if currentAuthProvider == "phone" {
			// Already phone, keep it
		} else if currentAuthProvider == "email" {
			// Legacy email login, update to phone
			updateAuthProviderQuery := `UPDATE users SET auth_provider = 'phone' WHERE id = $1`
			db.DB.Exec(updateAuthProviderQuery, userID)
		}
	}

	// Get updated user
	var user models.User
	var email sql.NullString
	var phoneNum sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	var phoneVerifiedAt sql.NullTime

	selectQuery := `SELECT id, email, phone_number, full_name, role, google_id, auth_provider, 
	                phone_verified, phone_verified_at, created_at
	                FROM users WHERE id = $1`
	err = db.DB.QueryRow(selectQuery, userID).Scan(
		&user.ID, &email, &phoneNum, &user.FullName, &user.Role,
		&googleID, &authProvider, &user.PhoneVerified, &phoneVerifiedAt, &user.CreatedAt)

	if err != nil {
		c.Logger().Errorf("Get updated user error: %v", err)
		// Return success anyway since phone was updated
		return c.JSON(http.StatusOK, models.OTPResponse{
			Success: true,
			Message: "Nomor WhatsApp berhasil diverifikasi",
		})
	}

	// Set optional fields
	if email.Valid {
		user.Email = email.String
	}
	if phoneNum.Valid {
		user.PhoneNumber = &phoneNum.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "phone"
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Nomor WhatsApp berhasil diverifikasi",
		"user":    user,
	})
}

// Helper function to get user statistics
func getUserStatistics(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count children (using parent_id column)
	var childrenCount int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM children WHERE parent_id = $1`, userID).Scan(&childrenCount)
	if err != nil {
		return nil, err
	}
	stats["children_count"] = childrenCount

	// Count measurements (measurements table uses child_id, need to join with children)
	var measurementsCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM measurements m
		JOIN children c ON m.child_id = c.id
		WHERE c.parent_id = $1
	`, userID).Scan(&measurementsCount)
	if err != nil {
		// If table doesn't exist or error, set to 0
		measurementsCount = 0
	}
	stats["measurements_count"] = measurementsCount

	// Count milestones (assuming there's a child_milestones table linked to children)
	var milestonesCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM child_milestones cm
		JOIN children c ON cm.child_id = c.id
		WHERE c.parent_id = $1
	`, userID).Scan(&milestonesCount)
	if err != nil {
		// If table doesn't exist, set to 0
		milestonesCount = 0
	}
	stats["milestones_count"] = milestonesCount

	return stats, nil
}

