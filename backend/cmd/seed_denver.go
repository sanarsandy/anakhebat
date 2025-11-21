package main

import (
	"log"
	"tukem-backend/db"
	"tukem-backend/utils"
)

func main() {
	// Initialize Database
	db.Init()

	// Seed Denver II milestones
	if err := utils.SeedDenverIIMilestones(db.DB); err != nil {
		log.Fatalf("Failed to seed Denver II milestones: %v", err)
	}

	log.Println("Denver II milestones seeded successfully!")
}

