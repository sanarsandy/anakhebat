package main

import (
	"log"
	"tukem-backend/db"
	"tukem-backend/utils"
)

func main() {
	// Initialize Database
	db.Init()

	// Clear existing WHO standards
	log.Println("Clearing existing WHO standards...")
	_, err := db.DB.Exec("DELETE FROM who_standards")
	if err != nil {
		log.Fatalf("Failed to clear WHO standards: %v", err)
	}
	log.Println("WHO standards cleared")

	// Re-seed WHO standards
	if err := utils.SeedWHOStandards(db.DB); err != nil {
		log.Fatalf("Failed to seed WHO standards: %v", err)
	}

	log.Println("WHO standards re-seeded successfully!")
}

