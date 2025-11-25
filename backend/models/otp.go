package models

import (
	"time"
)

type OTPCode struct {
	ID           string    `json:"id" db:"id"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	OTPCode      string    `json:"otp_code" db:"otp_code"`
	Purpose      string    `json:"purpose" db:"purpose"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	UsedAt       *time.Time `json:"used_at,omitempty" db:"used_at"`
	IsUsed       bool      `json:"is_used" db:"is_used"`
	AttemptCount int       `json:"attempt_count" db:"attempt_count"`
	MaxAttempts  int       `json:"max_attempts" db:"max_attempts"`
	IPAddress    *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string   `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type OTPRateLimit struct {
	ID           string    `json:"id" db:"id"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	RequestCount int       `json:"request_count" db:"request_count"`
	WindowStart  time.Time `json:"window_start" db:"window_start"`
	WindowEnd    time.Time `json:"window_end" db:"window_end"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

