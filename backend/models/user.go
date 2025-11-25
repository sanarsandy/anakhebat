package models

import (
	"time"
)

type User struct {
	ID              string     `json:"id" db:"id"`
	Email           string     `json:"email" db:"email"`
	PhoneNumber     *string    `json:"phone_number,omitempty" db:"phone_number"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FullName        string     `json:"full_name" db:"full_name"`
	Role            string     `json:"role" db:"role"`
	GoogleID        *string    `json:"-" db:"google_id"`
	AuthProvider    string     `json:"auth_provider" db:"auth_provider"`
	PhoneVerified   bool       `json:"phone_verified" db:"phone_verified"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty" db:"phone_verified_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	IsNewUser bool   `json:"is_new_user,omitempty"`
}

// OTP Request Models
type RequestOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	OTP         string `json:"otp" validate:"required,len=6"`
}

type OTPResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
	RetryAfter int   `json:"retry_after,omitempty"`
}

// Profile Request Models
type UpdateProfileRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

type VerifyPhoneRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type ConfirmPhoneVerificationRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	OTP         string `json:"otp" validate:"required,len=6"`
}
