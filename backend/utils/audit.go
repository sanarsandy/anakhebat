package utils

import (
	"database/sql"
	"encoding/json"
	"tukem-backend/db"
)

// LogAudit logs an admin action to the audit_logs table
func LogAudit(userID, action, resourceType string, resourceID *string, beforeData, afterData interface{}, ipAddress, userAgent string) error {
	var beforeJSON, afterJSON sql.NullString

	// Convert beforeData to JSON
	if beforeData != nil {
		beforeBytes, err := json.Marshal(beforeData)
		if err == nil {
			beforeJSON = sql.NullString{String: string(beforeBytes), Valid: true}
		}
	}

	// Convert afterData to JSON
	if afterData != nil {
		afterBytes, err := json.Marshal(afterData)
		if err == nil {
			afterJSON = sql.NullString{String: string(afterBytes), Valid: true}
		}
	}

	query := `INSERT INTO audit_logs (user_id, action, resource_type, resource_id, before_data, after_data, ip_address, user_agent)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	var resourceIDPtr interface{}
	if resourceID != nil {
		resourceIDPtr = *resourceID
	}

	_, err := db.DB.Exec(query, userID, action, resourceType, resourceIDPtr,
		beforeJSON, afterJSON, ipAddress, userAgent)
	return err
}


