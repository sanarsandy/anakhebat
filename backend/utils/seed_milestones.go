package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"tukem-backend/models"
)

func SeedMilestones(db *sqlx.DB) error {
	// Check if milestones already exist
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM milestones")
	if err != nil {
		return fmt.Errorf("failed to check existing milestones: %v", err)
	}

	if count > 0 {
		log.Println("Milestones already seeded, skipping...")
		return nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Construct path to seed file
	// Assuming this is run from the backend root directory
	seedFilePath := filepath.Join(cwd, "data", "milestones_seed.json")
	
	// Read seed file
	fileContent, err := ioutil.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("failed to read seed file at %s: %v", seedFilePath, err)
	}

	var milestones []models.Milestone
	if err := json.Unmarshal(fileContent, &milestones); err != nil {
		return fmt.Errorf("failed to unmarshal seed data: %v", err)
	}

	log.Printf("Seeding %d milestones...", len(milestones))

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.PrepareNamed(`
		INSERT INTO milestones (age_months, min_age_range, max_age_range, category, question, question_en, source, is_red_flag, pyramid_level, denver_domain)
		VALUES (:age_months, :min_age_range, :max_age_range, :category, :question, :question_en, :source, :is_red_flag, :pyramid_level, :denver_domain)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	for _, m := range milestones {
		if _, err := stmt.Exec(m); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert milestone: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Milestones seeded successfully!")
	return nil
}
