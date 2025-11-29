package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateUUID validates if a string is a valid UUID format
func ValidateUUID(uuid string) error {
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(strings.ToLower(uuid)) {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}

// ValidateStringLength validates string length
func ValidateStringLength(s string, min, max int, fieldName string) error {
	length := len(strings.TrimSpace(s))
	if length < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	if length > max {
		return fmt.Errorf("%s must be at most %d characters", fieldName, max)
	}
	return nil
}

// ValidateIntRange validates integer range
func ValidateIntRange(value int, min, max int, fieldName string) error {
	if value < min {
		return fmt.Errorf("%s must be at least %d", fieldName, min)
	}
	if value > max {
		return fmt.Errorf("%s must be at most %d", fieldName, max)
	}
	return nil
}

// SanitizeError sanitizes error messages to avoid exposing sensitive information
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}
	
	errMsg := err.Error()
	
	// Remove database-specific error details
	if strings.Contains(errMsg, "pq:") {
		// PostgreSQL errors - return generic message
		if strings.Contains(errMsg, "duplicate key") {
			return "Record already exists"
		}
		if strings.Contains(errMsg, "foreign key") {
			return "Cannot perform operation: related records exist"
		}
		if strings.Contains(errMsg, "violates") {
			return "Invalid data provided"
		}
		return "Database operation failed"
	}
	
	// Return original error if it's a validation error (safe to expose)
	return errMsg
}


