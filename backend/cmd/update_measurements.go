package main

import (
	"log"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"
)

func main() {
	// Initialize database
	db.Init()
	defer db.DB.Close()

	log.Println("Starting to update existing measurements with Z-scores...")

	// Get all measurements
	var measurements []struct {
		ID            string  `db:"id"`
		ChildID       string  `db:"child_id"`
		Weight        float64 `db:"weight"`
		Height        float64 `db:"height"`
		HeadCirc      *float64 `db:"head_circumference"`
		AgeInMonths   int     `db:"age_in_months"`
	}

	query := `SELECT id, child_id, weight, height, head_circumference, age_in_months FROM measurements`
	err := db.DB.Select(&measurements, query)
	if err != nil {
		log.Fatal("Failed to fetch measurements:", err)
	}

	log.Printf("Found %d measurements to update\n", len(measurements))

	// Get child data for gender
	childGenders := make(map[string]string)
	var children []models.Child
	err = db.DB.Select(&children, "SELECT id, gender FROM children")
	if err != nil {
		log.Fatal("Failed to fetch children:", err)
	}
	for _, child := range children {
		childGenders[child.ID] = child.Gender
	}

	// Update each measurement
	updated := 0
	for _, m := range measurements {
		gender, ok := childGenders[m.ChildID]
		if !ok {
			log.Printf("Skipping measurement %s - child not found\n", m.ID)
			continue
		}

		// Calculate Z-scores
		headCirc := 0.0
		if m.HeadCirc != nil {
			headCirc = *m.HeadCirc
		}

		zscores, err := utils.CalculateAllZScores(db.DB, gender, m.AgeInMonths, m.Weight, m.Height, headCirc)
		if err != nil {
			log.Printf("Failed to calculate Z-scores for measurement %s: %v\n", m.ID, err)
			continue
		}

		// Interpret nutritional status
		var nutritionalStatus, heightStatus, wfhStatus string
		if zscores != nil && (zscores.HasWeightForAge || zscores.HasHeightForAge) {
			wfaZ := 0.0
			hfaZ := 0.0
			wfhZ := 0.0
			
			if zscores.HasWeightForAge {
				wfaZ = zscores.WeightForAge
			}
			if zscores.HasHeightForAge {
				hfaZ = zscores.HeightForAge
			}
			if zscores.HasWeightForHeight {
				wfhZ = zscores.WeightForHeight
			}
			
			nutritionalStatus, heightStatus, wfhStatus = utils.InterpretNutritionalStatus(wfaZ, hfaZ, wfhZ)
		}

		// Update database
		var wfaZPtr, hfaZPtr, wfhZPtr, hcZPtr *float64
		if zscores != nil {
			if zscores.HasWeightForAge {
				wfaZPtr = &zscores.WeightForAge
			}
			if zscores.HasHeightForAge {
				hfaZPtr = &zscores.HeightForAge
			}
			if zscores.HasWeightForHeight {
				wfhZPtr = &zscores.WeightForHeight
			}
			if zscores.HasHeadCirc {
				hcZPtr = &zscores.HeadCircumference
			}
		}

		updateQuery := `
			UPDATE measurements 
			SET weight_for_age_zscore = $1,
				height_for_age_zscore = $2,
				weight_for_height_zscore = $3,
				head_circumference_zscore = $4,
				nutritional_status = $5,
				height_status = $6,
				weight_for_height_status = $7
			WHERE id = $8
		`

		_, err = db.DB.Exec(updateQuery, wfaZPtr, hfaZPtr, wfhZPtr, hcZPtr, 
			nutritionalStatus, heightStatus, wfhStatus, m.ID)
		
		if err != nil {
			log.Printf("Failed to update measurement %s: %v\n", m.ID, err)
			continue
		}

		updated++
		log.Printf("Updated measurement %s - Status: %s\n", m.ID, nutritionalStatus)
	}

	log.Printf("Successfully updated %d out of %d measurements\n", updated, len(measurements))
}
