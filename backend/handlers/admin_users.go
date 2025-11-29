package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// GetAdminUsers returns all users with pagination and filters
func GetAdminUsers(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	role := c.QueryParam("role")
	search := c.QueryParam("search")
	authProvider := c.QueryParam("auth_provider")

	// Build query
	query := `SELECT id, email, phone_number, full_name, role, auth_provider, 
	          phone_verified, created_at 
	          FROM users WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if role != "" {
		query += ` AND role = $` + strconv.Itoa(argIndex)
		args = append(args, role)
		argIndex++
	}

	if authProvider != "" {
		query += ` AND auth_provider = $` + strconv.Itoa(argIndex)
		args = append(args, authProvider)
		argIndex++
	}

	if search != "" {
		query += ` AND (full_name ILIKE $` + strconv.Itoa(argIndex) +
			` OR email ILIKE $` + strconv.Itoa(argIndex) +
			` OR phone_number ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminUsers query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var user models.User
		var email sql.NullString
		var phoneNumber sql.NullString
		var authProvider sql.NullString

		err := rows.Scan(&user.ID, &email, &phoneNumber, &user.FullName, &user.Role,
			&authProvider, &user.PhoneVerified, &user.CreatedAt)
		if err != nil {
			continue
		}

		if email.Valid {
			user.Email = email.String
		}
		if phoneNumber.Valid {
			user.PhoneNumber = &phoneNumber.String
		}
		if authProvider.Valid {
			user.AuthProvider = authProvider.String
		}

		userMap := map[string]interface{}{
			"id":             user.ID,
			"email":          user.Email,
			"phone_number":   user.PhoneNumber,
			"full_name":      user.FullName,
			"role":           user.Role,
			"auth_provider":  user.AuthProvider,
			"phone_verified": user.PhoneVerified,
			"created_at":     user.CreatedAt,
		}
		users = append(users, userMap)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if role != "" {
		countQuery += ` AND role = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, role)
		countArgIndex++
	}

	if authProvider != "" {
		countQuery += ` AND auth_provider = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, authProvider)
		countArgIndex++
	}

	if search != "" {
		countQuery += ` AND (full_name ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR email ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR phone_number ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern, searchPattern)
		countArgIndex += 3
	}

	var total int
	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminUsers count error: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminUser returns a single user with details and statistics
func GetAdminUser(c echo.Context) error {
	userID := c.Param("id")

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
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	} else if err != nil {
		c.Logger().Errorf("GetAdminUser error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

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
	}
	if phoneVerifiedAt.Valid {
		user.PhoneVerifiedAt = &phoneVerifiedAt.Time
	}

	// Get statistics
	stats, err := getUserStatisticsForAdmin(userID)
	if err != nil {
		c.Logger().Warnf("Failed to get user statistics: %v", err)
		stats = map[string]interface{}{}
	}

	// Get children count
	var childrenCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE parent_id = $1", userID).Scan(&childrenCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":          user,
		"statistics":    stats,
		"children_count": childrenCount,
	})
}

// CreateAdminUser creates a new user (admin only)
func CreateAdminUser(c echo.Context) error {
	// Get admin user ID for audit log
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FullName    string `json:"full_name" validate:"required"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role" validate:"required,oneof=parent admin"`
		AuthProvider string `json:"auth_provider" validate:"oneof=email phone google"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate
	if req.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "full_name is required"})
	}
	
	// Validate full_name length
	if err := utils.ValidateStringLength(req.FullName, 1, 255, "full_name"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	if req.Role != "parent" && req.Role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role"})
	}
	
	// Validate email format and length
	if req.Email != "" {
		if !strings.Contains(req.Email, "@") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email format"})
		}
		if err := utils.ValidateStringLength(req.Email, 3, 255, "email"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	// Check if email or phone already exists
	if req.Email != "" {
		var existingID string
		err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
	}

	if req.PhoneNumber != "" {
		var existingID string
		err := db.DB.QueryRow("SELECT id FROM users WHERE phone_number = $1", req.PhoneNumber).Scan(&existingID)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Phone number already exists"})
		}
	}

	// Hash password if provided
	var passwordHash sql.NullString
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}
		passwordHash = sql.NullString{String: string(hashedPassword), Valid: true}
	}

	// Set default auth_provider
	if req.AuthProvider == "" {
		if req.Email != "" && req.Password != "" {
			req.AuthProvider = "email"
		} else if req.PhoneNumber != "" {
			req.AuthProvider = "phone"
		}
	}

	// Insert user
	var user models.User
	var email sql.NullString
	var phoneNumber sql.NullString
	var authProvider sql.NullString

	insertQuery := `INSERT INTO users (email, password_hash, full_name, phone_number, role, auth_provider)
	                VALUES ($1, $2, $3, $4, $5, $6)
	                RETURNING id, email, phone_number, full_name, role, auth_provider, created_at`

	err := db.DB.QueryRow(insertQuery,
		req.Email, passwordHash, req.FullName, req.PhoneNumber, req.Role, req.AuthProvider).
		Scan(&user.ID, &email, &phoneNumber, &user.FullName, &user.Role, &authProvider, &user.CreatedAt)

	if err != nil {
		c.Logger().Errorf("CreateAdminUser error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	if email.Valid {
		user.Email = email.String
	}
	if phoneNumber.Valid {
		user.PhoneNumber = &phoneNumber.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	}

	// Log audit
	utils.LogAudit(adminUserID, "create", "user", &user.ID, nil, user, ipAddress, userAgent)

	return c.JSON(http.StatusCreated, user)
}

// UpdateAdminUser updates a user (admin only)
func UpdateAdminUser(c echo.Context) error {
	userID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Get existing user for audit log
	var existingUser models.User
	var email sql.NullString
	var phoneNumber sql.NullString
	var authProvider sql.NullString

	selectQuery := `SELECT id, email, phone_number, full_name, role, auth_provider, created_at
	                FROM users WHERE id = $1`
	err := db.DB.QueryRow(selectQuery, userID).Scan(
		&existingUser.ID, &email, &phoneNumber, &existingUser.FullName,
		&existingUser.Role, &authProvider, &existingUser.CreatedAt)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if email.Valid {
		existingUser.Email = email.String
	}
	if phoneNumber.Valid {
		existingUser.PhoneNumber = &phoneNumber.String
	}
	if authProvider.Valid {
		existingUser.AuthProvider = authProvider.String
	}

	var req struct {
		Email       *string `json:"email"`
		FullName    *string `json:"full_name"`
		PhoneNumber *string `json:"phone_number"`
		Role        *string `json:"role"`
		AuthProvider *string `json:"auth_provider"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Email != nil {
		// Check if email already exists (excluding current user)
		var existingID string
		err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1 AND id != $2", *req.Email, userID).Scan(&existingID)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
		updates = append(updates, "email = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Email)
		argIndex++
	}

	if req.FullName != nil {
		updates = append(updates, "full_name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.FullName)
		argIndex++
	}

	if req.PhoneNumber != nil {
		// Check if phone already exists (excluding current user)
		var existingID string
		err := db.DB.QueryRow("SELECT id FROM users WHERE phone_number = $1 AND id != $2", *req.PhoneNumber, userID).Scan(&existingID)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Phone number already exists"})
		}
		updates = append(updates, "phone_number = $"+strconv.Itoa(argIndex))
		args = append(args, *req.PhoneNumber)
		argIndex++
	}

	if req.Role != nil {
		if *req.Role != "parent" && *req.Role != "admin" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role"})
		}
		// Prevent admin from changing their own role from admin
		if userID == adminUserID && *req.Role != "admin" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot change your own role from admin"})
		}
		
		// Check if this is the last admin and trying to change role
		if *req.Role != "admin" {
			var adminCount int
			err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin' AND id != $1", userID).Scan(&adminCount)
			if err != nil {
				c.Logger().Errorf("Failed to check admin count: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate operation"})
			}
			if adminCount == 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot change role: at least one admin must remain"})
			}
		}
		
		updates = append(updates, "role = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Role)
		argIndex++
	}

	if req.AuthProvider != nil {
		updates = append(updates, "auth_provider = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AuthProvider)
		argIndex++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	// Add userID to args
	args = append(args, userID)
	updateQuery := `UPDATE users SET ` + strings.Join(updates, ", ") + ` WHERE id = $` + strconv.Itoa(argIndex)

	_, err = db.DB.Exec(updateQuery, args...)
	if err != nil {
		c.Logger().Errorf("UpdateAdminUser error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Get updated user
	var updatedUser models.User
	selectQuery = `SELECT id, email, phone_number, full_name, role, auth_provider, created_at
	               FROM users WHERE id = $1`
	err = db.DB.QueryRow(selectQuery, userID).Scan(
		&updatedUser.ID, &email, &phoneNumber, &updatedUser.FullName,
		&updatedUser.Role, &authProvider, &updatedUser.CreatedAt)
	if err != nil {
		c.Logger().Errorf("Get updated user error: %v", err)
	}

	if email.Valid {
		updatedUser.Email = email.String
	}
	if phoneNumber.Valid {
		updatedUser.PhoneNumber = &phoneNumber.String
	}
	if authProvider.Valid {
		updatedUser.AuthProvider = authProvider.String
	}

	// Log audit
	utils.LogAudit(adminUserID, "update", "user", &userID, existingUser, updatedUser, ipAddress, userAgent)

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteAdminUser deletes a user (admin only)
func DeleteAdminUser(c echo.Context) error {
	userID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Validate UUID format
	if err := utils.ValidateUUID(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	// Prevent admin from deleting themselves
	if userID == adminUserID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete your own account"})
	}

	// Get existing user for audit log
	var existingUser models.User
	var email sql.NullString
	var phoneNumber sql.NullString
	var authProvider sql.NullString

	selectQuery := `SELECT id, email, phone_number, full_name, role, auth_provider, created_at
	                FROM users WHERE id = $1`
	err := db.DB.QueryRow(selectQuery, userID).Scan(
		&existingUser.ID, &email, &phoneNumber, &existingUser.FullName,
		&existingUser.Role, &authProvider, &existingUser.CreatedAt)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if email.Valid {
		existingUser.Email = email.String
	}
	if phoneNumber.Valid {
		existingUser.PhoneNumber = &phoneNumber.String
	}
	if authProvider.Valid {
		existingUser.AuthProvider = authProvider.String
	}

	// Check if this is the last admin
	if existingUser.Role == "admin" {
		var adminCount int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount)
		if err != nil {
			c.Logger().Errorf("Failed to check admin count: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate operation"})
		}
		if adminCount <= 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete the last admin user"})
		}
	}

	// Delete user (cascade will delete children)
	_, err = db.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		c.Logger().Errorf("DeleteAdminUser error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	utils.LogAudit(adminUserID, "delete", "user", &userID, existingUser, nil, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// ResetAdminUserPassword resets a user's password (admin only)
func ResetAdminUserPassword(c echo.Context) error {
	userID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	
	// Validate UUID format
	if err := utils.ValidateUUID(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if len(req.NewPassword) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 6 characters"})
	}

	// Check if user exists
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Update password
	_, err = db.DB.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(hashedPassword), userID)
	if err != nil {
		c.Logger().Errorf("ResetAdminUserPassword error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reset password"})
	}

	// Log audit
	auditData := map[string]string{"action": "password_reset"}
	utils.LogAudit(adminUserID, "update", "user", &userID, nil, auditData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// getUserStatisticsForAdmin is a helper function to get user statistics for admin
func getUserStatisticsForAdmin(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get children count
	var childrenCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE parent_id = $1", userID).Scan(&childrenCount)
	if err != nil {
		return nil, err
	}
	stats["children_count"] = childrenCount

	// Get measurements count
	var measurementsCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM measurements m
		JOIN children c ON c.id = m.child_id
		WHERE c.parent_id = $1
	`, userID).Scan(&measurementsCount)
	if err != nil {
		return nil, err
	}
	stats["measurements_count"] = measurementsCount

	// Get assessments count
	var assessmentsCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(DISTINCT assessment_date) FROM assessments a
		JOIN children c ON c.id = a.child_id
		WHERE c.parent_id = $1
	`, userID).Scan(&assessmentsCount)
	if err != nil {
		return nil, err
	}
	stats["assessments_count"] = assessmentsCount

	// Get immunizations count
	var immunizationsCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM child_immunizations ci
		JOIN children c ON c.id = ci.child_id
		WHERE c.parent_id = $1
	`, userID).Scan(&immunizationsCount)
	if err != nil {
		return nil, err
	}
	stats["immunizations_count"] = immunizationsCount

	return stats, nil
}

