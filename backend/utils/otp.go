package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateOTP generates a cryptographically secure 6-digit OTP
func GenerateOTP() (string, error) {
	bytes := make([]byte, 3) // 3 bytes = 6 hex digits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	otp := binary.BigEndian.Uint32(append([]byte{0}, bytes...))
	return fmt.Sprintf("%06d", otp%1000000), nil
}

// ValidatePhoneNumber validates and normalizes phone number
// Returns normalized phone number in international format (+6281234567890)
func ValidatePhoneNumber(phone string) (string, error) {
	// Remove spaces, dashes, parentheses, etc.
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Handle Indonesian format (08xx -> +628xx)
	if strings.HasPrefix(cleaned, "08") && len(cleaned) >= 10 {
		cleaned = "+62" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "8") && len(cleaned) >= 9 {
		cleaned = "+62" + cleaned
	} else if strings.HasPrefix(cleaned, "62") && !strings.HasPrefix(cleaned, "+62") {
		cleaned = "+" + cleaned
	}

	// Must start with + (international format)
	if !strings.HasPrefix(cleaned, "+") {
		return "", fmt.Errorf("nomor telepon harus dimulai dengan + (format internasional)")
	}

	// Validate length (10-15 digits after +)
	if len(cleaned) < 11 || len(cleaned) > 16 {
		return "", fmt.Errorf("panjang nomor telepon tidak valid (harus 10-15 digit setelah kode negara)")
	}

	// Validate digits only after +
	if matched, _ := regexp.MatchString(`^\+[0-9]{10,15}$`, cleaned); !matched {
		return "", fmt.Errorf("nomor telepon mengandung karakter tidak valid")
	}

	return cleaned, nil
}

// GetOTPExpiration returns the expiration time for OTP (5 minutes from now)
func GetOTPExpiration() time.Time {
	return time.Now().Add(5 * time.Minute)
}

// IsOTPExpired checks if OTP has expired
func IsOTPExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IntToString converts integer to string
func IntToString(n int) string {
	return fmt.Sprintf("%d", n)
}

// JoinStrings joins string slice with separator
func JoinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// GetCurrentTime returns current time
func GetCurrentTime() time.Time {
	return time.Now()
}

