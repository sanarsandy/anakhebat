package main

import (
	"fmt"
	"os"
	"tukem-backend/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize Database
	db.Init()

	adminEmail := "admin@jurnalsikecil.com"
	adminPassword := "admin123"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	// Update password
	_, err = db.DB.Exec(`
		UPDATE users 
		SET password_hash = $1 
		WHERE email = $2
	`, string(hashedPassword), adminEmail)
	
	if err != nil {
		fmt.Printf("Error updating password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Password reset successfully for %s\n", adminEmail)
	fmt.Printf("  Password: %s\n", adminPassword)
	fmt.Printf("  Hash: %s\n", string(hashedPassword))
}


