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

// SeedDenverIIMilestones seeds Denver II milestones into the database
func SeedDenverIIMilestones(db *sqlx.DB) error {
	// Check if Denver II milestones already exist
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM milestones WHERE source = 'DENVER'")
	if err != nil {
		return fmt.Errorf("failed to check existing Denver II milestones: %v", err)
	}

	if count > 0 {
		log.Printf("Denver II milestones already seeded (%d entries), skipping...", count)
		return nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Construct path to Denver II seed file
	seedFilePath := filepath.Join(cwd, "data", "denver_ii_milestones.json")

	// Read seed file
	fileContent, err := ioutil.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("failed to read Denver II seed file at %s: %v", seedFilePath, err)
	}

	var milestones []models.Milestone
	if err := json.Unmarshal(fileContent, &milestones); err != nil {
		return fmt.Errorf("failed to unmarshal Denver II seed data: %v", err)
	}

	log.Printf("Seeding %d Denver II milestones...", len(milestones))

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamed(`
		INSERT INTO milestones (age_months, min_age_range, max_age_range, category, question, question_en, source, is_red_flag, pyramid_level, denver_domain)
		VALUES (:age_months, :min_age_range, :max_age_range, :category, :question, :question_en, :source, :is_red_flag, :pyramid_level, :denver_domain)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	seededCount := 0
	for _, m := range milestones {
		// Set category based on denver_domain if not set
		if m.Category == "" {
			switch m.DenverDomain {
			case nil:
				m.Category = "cognitive" // default
			default:
				switch *m.DenverDomain {
				case "PS":
					m.Category = "sensory"
				case "FM":
					m.Category = "motor"
				case "L":
					m.Category = "cognitive"
				case "GM":
					m.Category = "motor"
				default:
					m.Category = "cognitive"
				}
			}
		}

		if _, err := stmt.Exec(m); err != nil {
			log.Printf("Warning: Failed to insert Denver II milestone: %v", err)
			continue
		}
		seededCount++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Denver II milestones seeded successfully! Total: %d entries", seededCount)
	return nil
}

