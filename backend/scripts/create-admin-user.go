package main

import (
	"fmt"
	"os"
	"tukem-backend/db"
	"tukem-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize Database
	db.Init()

	// Get admin credentials from environment or use defaults
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@jurnalsikecil.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		adminName = "Admin Jurnal Si Kecil"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	// Check if user already exists
	var existingID string
	err = db.DB.QueryRow("SELECT id FROM users WHERE email = $1", adminEmail).Scan(&existingID)
	if err == nil {
		// User exists, update to admin
		_, err = db.DB.Exec(`
			UPDATE users 
			SET role = 'admin', password_hash = $1, full_name = $2
			WHERE email = $3
		`, string(hashedPassword), adminName, adminEmail)
		if err != nil {
			fmt.Printf("Error updating user to admin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ User %s updated to admin role\n", adminEmail)
	} else {
		// User doesn't exist, create new admin user
		_, err = db.DB.Exec(`
			INSERT INTO users (email, password_hash, full_name, role, auth_provider)
			VALUES ($1, $2, $3, 'admin', 'email')
		`, adminEmail, string(hashedPassword), adminName)
		if err != nil {
			fmt.Printf("Error creating admin user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Admin user created: %s\n", adminEmail)
	}

	fmt.Printf("\nAdmin credentials:\n")
	fmt.Printf("  Email: %s\n", adminEmail)
	fmt.Printf("  Password: %s\n", adminPassword)
	fmt.Printf("\n⚠️  Please change the password after first login!\n")
}


