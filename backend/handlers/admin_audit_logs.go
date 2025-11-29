package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"tukem-backend/db"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetAuditLogs returns audit logs with pagination and filters
func GetAuditLogs(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	userID := c.QueryParam("user_id")
	action := c.QueryParam("action")
	resourceType := c.QueryParam("resource_type")
	resourceID := c.QueryParam("resource_id")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")
	
	// Validate UUID format for user_id filter if provided
	if userID != "" {
		if err := utils.ValidateUUID(userID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
		}
	}

	// Build query
	query := `
		SELECT 
			al.id,
			al.user_id,
			u.full_name as user_name,
			u.email as user_email,
			al.action,
			al.resource_type,
			al.resource_id,
			al.before_data,
			al.after_data,
			al.ip_address,
			al.user_agent,
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if userID != "" {
		query += ` AND al.user_id = $` + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}

	if action != "" {
		query += ` AND al.action = $` + strconv.Itoa(argIndex)
		args = append(args, action)
		argIndex++
	}

	if resourceType != "" {
		query += ` AND al.resource_type = $` + strconv.Itoa(argIndex)
		args = append(args, resourceType)
		argIndex++
	}

	if resourceID != "" {
		query += ` AND al.resource_id = $` + strconv.Itoa(argIndex)
		args = append(args, resourceID)
		argIndex++
	}

	if dateFrom != "" {
		query += ` AND al.created_at >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		query += ` AND al.created_at <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM audit_logs al WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if userID != "" {
		countQuery += ` AND al.user_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, userID)
		countArgIndex++
	}
	if action != "" {
		countQuery += ` AND al.action = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, action)
		countArgIndex++
	}
	if resourceType != "" {
		countQuery += ` AND al.resource_type = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, resourceType)
		countArgIndex++
	}
	if resourceID != "" {
		countQuery += ` AND al.resource_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, resourceID)
		countArgIndex++
	}
	if dateFrom != "" {
		countQuery += ` AND al.created_at >= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateFrom)
		countArgIndex++
	}
	if dateTo != "" {
		countQuery += ` AND al.created_at <= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateTo)
		countArgIndex++
	}

	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAuditLogs count error: %v", err)
		total = 0
	}

	query += ` ORDER BY al.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAuditLogs query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	type AuditLog struct {
		ID          string                 `json:"id"`
		UserID      *string                `json:"user_id"`
		UserName    string                 `json:"user_name"`
		UserEmail   string                 `json:"user_email"`
		Action      string                 `json:"action"`
		ResourceType string                `json:"resource_type"`
		ResourceID  *string                `json:"resource_id"`
		BeforeData  map[string]interface{} `json:"before_data"`
		AfterData   map[string]interface{} `json:"after_data"`
		IPAddress   string                 `json:"ip_address"`
		UserAgent   string                 `json:"user_agent"`
		CreatedAt   string                 `json:"created_at"`
	}

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		var userID sql.NullString
		var userName sql.NullString
		var userEmail sql.NullString
		var resourceID sql.NullString
		var beforeData sql.NullString
		var afterData sql.NullString
		var ipAddress sql.NullString
		var userAgent sql.NullString
		var createdAt time.Time

		err := rows.Scan(
			&log.ID, &userID, &userName, &userEmail, &log.Action, &log.ResourceType,
			&resourceID, &beforeData, &afterData, &ipAddress, &userAgent, &createdAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan audit log row: %v", err)
			continue
		}

		if userID.Valid {
			log.UserID = &userID.String
		}
		if userName.Valid {
			log.UserName = userName.String
		}
		if userEmail.Valid {
			log.UserEmail = userEmail.String
		}
		if resourceID.Valid {
			log.ResourceID = &resourceID.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}
		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}
		log.CreatedAt = createdAt.Format("2006-01-02T15:04:05Z07:00")

		// Parse JSON data
		if beforeData.Valid && beforeData.String != "" {
			json.Unmarshal([]byte(beforeData.String), &log.BeforeData)
		}
		if afterData.Valid && afterData.String != "" {
			json.Unmarshal([]byte(afterData.String), &log.AfterData)
		}

		logs = append(logs, log)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAuditLog returns a single audit log entry
func GetAuditLog(c echo.Context) error {
	logID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(logID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid log ID format"})
	}

	type AuditLog struct {
		ID          string                 `json:"id"`
		UserID      *string                `json:"user_id"`
		UserName    string                 `json:"user_name"`
		UserEmail   string                 `json:"user_email"`
		Action      string                 `json:"action"`
		ResourceType string                `json:"resource_type"`
		ResourceID  *string                `json:"resource_id"`
		BeforeData  map[string]interface{} `json:"before_data"`
		AfterData   map[string]interface{} `json:"after_data"`
		IPAddress   string                 `json:"ip_address"`
		UserAgent   string                 `json:"user_agent"`
		CreatedAt   string                 `json:"created_at"`
	}

	var log AuditLog
	var userID sql.NullString
	var userName sql.NullString
	var userEmail sql.NullString
	var resourceID sql.NullString
	var beforeData sql.NullString
	var afterData sql.NullString
	var ipAddress sql.NullString
	var userAgent sql.NullString
	var createdAt time.Time

	err := db.DB.QueryRow(
		`SELECT 
			al.id, al.user_id, u.full_name, u.email, al.action, al.resource_type,
			al.resource_id, al.before_data, al.after_data, al.ip_address, al.user_agent, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		WHERE al.id = $1`,
		logID,
	).Scan(
		&log.ID, &userID, &userName, &userEmail, &log.Action, &log.ResourceType,
		&resourceID, &beforeData, &afterData, &ipAddress, &userAgent, &createdAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Audit log not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetAuditLog error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if userID.Valid {
		log.UserID = &userID.String
	}
	if userName.Valid {
		log.UserName = userName.String
	}
	if userEmail.Valid {
		log.UserEmail = userEmail.String
	}
	if resourceID.Valid {
		log.ResourceID = &resourceID.String
	}
	if ipAddress.Valid {
		log.IPAddress = ipAddress.String
	}
	if userAgent.Valid {
		log.UserAgent = userAgent.String
	}
	log.CreatedAt = createdAt.Format("2006-01-02T15:04:05Z07:00")

	// Parse JSON data
	if beforeData.Valid && beforeData.String != "" {
		json.Unmarshal([]byte(beforeData.String), &log.BeforeData)
	}
	if afterData.Valid && afterData.String != "" {
		json.Unmarshal([]byte(afterData.String), &log.AfterData)
	}

	return c.JSON(http.StatusOK, log)
}

// ExportAuditLogs exports audit logs to CSV
func ExportAuditLogs(c echo.Context) error {
	// Get query parameters (same filters as GetAuditLogs)
	userID := c.QueryParam("user_id")
	action := c.QueryParam("action")
	resourceType := c.QueryParam("resource_type")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query (same as GetAuditLogs but without pagination)
	query := `
		SELECT 
			al.id,
			u.full_name as user_name,
			u.email as user_email,
			al.action,
			al.resource_type,
			al.resource_id,
			al.ip_address,
			al.user_agent,
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if userID != "" {
		query += ` AND al.user_id = $` + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}
	if action != "" {
		query += ` AND al.action = $` + strconv.Itoa(argIndex)
		args = append(args, action)
		argIndex++
	}
	if resourceType != "" {
		query += ` AND al.resource_type = $` + strconv.Itoa(argIndex)
		args = append(args, resourceType)
		argIndex++
	}
	if dateFrom != "" {
		query += ` AND al.created_at >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}
	if dateTo != "" {
		query += ` AND al.created_at <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	query += ` ORDER BY al.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("ExportAuditLogs query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=audit_logs.csv")

	writer := csv.NewWriter(c.Response())
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "User Name", "User Email", "Action", "Resource Type", "Resource ID",
		"IP Address", "User Agent", "Created At",
	}
	if err := writer.Write(header); err != nil {
		c.Logger().Errorf("Failed to write CSV header: %v", err)
		return err
	}

	// Write data
	for rows.Next() {
		var id, userName, userEmail, action, resourceType, resourceID, ipAddress, userAgent string
		var createdAt time.Time
		var userNameNull, userEmailNull, resourceIDNull, ipAddressNull, userAgentNull sql.NullString

		err := rows.Scan(
			&id, &userNameNull, &userEmailNull, &action, &resourceType,
			&resourceIDNull, &ipAddressNull, &userAgentNull, &createdAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan audit log row: %v", err)
			continue
		}

		if userNameNull.Valid {
			userName = userNameNull.String
		}
		if userEmailNull.Valid {
			userEmail = userEmailNull.String
		}
		if resourceIDNull.Valid {
			resourceID = resourceIDNull.String
		}
		if ipAddressNull.Valid {
			ipAddress = ipAddressNull.String
		}
		if userAgentNull.Valid {
			userAgent = userAgentNull.String
		}

		record := []string{
			id, userName, userEmail, action, resourceType, resourceID,
			ipAddress, userAgent, createdAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			c.Logger().Errorf("Failed to write CSV record: %v", err)
			return err
		}
	}

	return nil
}

