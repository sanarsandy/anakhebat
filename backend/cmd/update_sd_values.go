package main

import (
	"log"
	"tukem-backend/db"
	"tukem-backend/utils"
)

func main() {
	// Initialize Database
	db.Init()

	log.Println("Updating SD values for WHO standards...")
	if err := utils.UpdateSDValues(db.DB); err != nil {
		log.Fatalf("Failed to update SD values: %v", err)
	}

	log.Println("SD values updated successfully!")
}

