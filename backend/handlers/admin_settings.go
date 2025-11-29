package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetSystemSettings returns all system settings
func GetSystemSettings(c echo.Context) error {
	category := c.QueryParam("category")

	query := `SELECT id, key, value, type, category, description, updated_by, updated_at 
	          FROM system_settings WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	query += ` ORDER BY category, key`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetSystemSettings query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	type Setting struct {
		ID          string    `json:"id"`
		Key         string    `json:"key"`
		Value       string    `json:"value"`
		Type        string    `json:"type"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		UpdatedBy   *string   `json:"updated_by"`
		UpdatedAt   string    `json:"updated_at"`
	}

	var settings []Setting
	for rows.Next() {
		var s Setting
		var updatedBy sql.NullString
		var updatedAt sql.NullTime

		err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Type, &s.Category, &s.Description, &updatedBy, &updatedAt)
		if err != nil {
			c.Logger().Errorf("Failed to scan setting row: %v", err)
			continue
		}

		if updatedBy.Valid {
			s.UpdatedBy = &updatedBy.String
		}
		if updatedAt.Valid {
			s.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		settings = append(settings, s)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"settings": settings,
		"total":    len(settings),
	})
}

// GetSystemSetting returns a single system setting by key
func GetSystemSetting(c echo.Context) error {
	key := c.Param("key")
	
	// Validate key format (alphanumeric, underscore, hyphen)
	if key == "" || len(key) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid setting key format"})
	}

	var setting struct {
		ID          string    `json:"id"`
		Key         string    `json:"key"`
		Value       string    `json:"value"`
		Type        string    `json:"type"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		UpdatedBy   *string   `json:"updated_by"`
		UpdatedAt   string    `json:"updated_at"`
	}

	var updatedBy sql.NullString
	var updatedAt sql.NullTime

	err := db.DB.QueryRow(
		`SELECT id, key, value, type, category, description, updated_by, updated_at 
		 FROM system_settings WHERE key = $1`,
		key,
	).Scan(&setting.ID, &setting.Key, &setting.Value, &setting.Type, &setting.Category, &setting.Description, &updatedBy, &updatedAt)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Setting not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetSystemSetting error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if updatedBy.Valid {
		setting.UpdatedBy = &updatedBy.String
	}
	if updatedAt.Valid {
		setting.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(http.StatusOK, setting)
}

// UpdateSystemSetting updates a system setting
func UpdateSystemSetting(c echo.Context) error {
	key := c.Param("key")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	// Validate key format
	if key == "" || len(key) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid setting key format"})
	}

	var req struct {
		Value interface{} `json:"value"`
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("UpdateSystemSetting bind error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Convert value to string (handle number, boolean, string)
	var valueStr string
	switch v := req.Value.(type) {
	case string:
		valueStr = v
	case float64:
		// JSON numbers are parsed as float64
		valueStr = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			valueStr = "true"
		} else {
			valueStr = "false"
		}
	case nil:
		valueStr = ""
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	// Get existing setting for audit log
	var oldValue string
	err := db.DB.QueryRow("SELECT value FROM system_settings WHERE key = $1", key).Scan(&oldValue)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Setting not found"})
	}
	if err != nil {
		c.Logger().Errorf("UpdateSystemSetting get old value error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Update setting
	_, err = db.DB.Exec(
		`UPDATE system_settings 
		 SET value = $1, updated_by = $2, updated_at = CURRENT_TIMESTAMP 
		 WHERE key = $3`,
		valueStr, adminUserID, key,
	)
	if err != nil {
		c.Logger().Errorf("UpdateSystemSetting update error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update setting"})
	}

	// Log audit
	beforeData := map[string]string{"value": oldValue}
	afterData := map[string]string{"value": valueStr}
	utils.LogAudit(adminUserID, "update", "system_setting", &key, beforeData, afterData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Setting updated successfully"})
}

// UpdateSystemSettingsBatch updates multiple system settings at once
func UpdateSystemSettingsBatch(c echo.Context) error {
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		Settings map[string]string `json:"settings" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update each setting
	for key, value := range req.Settings {
		// Get old value for audit
		var oldValue string
		err := db.DB.QueryRow("SELECT value FROM system_settings WHERE key = $1", key).Scan(&oldValue)
		if err != nil {
			c.Logger().Errorf("UpdateSystemSettingsBatch get old value for %s error: %v", key, err)
			continue
		}

		// Update
		_, err = db.DB.Exec(
			`UPDATE system_settings 
			 SET value = $1, updated_by = $2, updated_at = CURRENT_TIMESTAMP 
			 WHERE key = $3`,
			value, adminUserID, key,
		)
		if err != nil {
			c.Logger().Errorf("UpdateSystemSettingsBatch update %s error: %v", key, err)
			continue
		}

		// Log audit
		beforeData := map[string]string{"value": oldValue}
		afterData := map[string]string{"value": value}
		utils.LogAudit(adminUserID, "update", "system_setting", &key, beforeData, afterData, ipAddress, userAgent)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

