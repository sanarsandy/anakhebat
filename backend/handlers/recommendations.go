package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetRecommendations fetches personalized recommendations for a child
// based on their milestone assessments and age
func GetRecommendations(c echo.Context) error {
	childID := c.Param("id")

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := db.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}
	if err != nil {
		c.Logger().Errorf("Failed to verify child ownership: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify child ownership"})
	}
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get child data
	var child models.Child
	err = db.DB.QueryRow("SELECT id, dob, is_premature, gestational_age FROM children WHERE id = $1", childID).
		Scan(&child.ID, &child.DOB, &child.IsPremature, &child.GestationalAge)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get child data"})
	}

	// Calculate current age (using corrected age if applicable)
	// Use current date for age calculation
	today := time.Now().Format("2006-01-02")
	_, ageInMonths, useCorrected, err := utils.CalculateCorrectedAge(
		child.DOB, today, child.IsPremature, child.GestationalAge)
	if err != nil {
		c.Logger().Errorf("Failed to calculate age: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate child age"})
	}

	c.Logger().Infof("Getting recommendations for child %s, age: %d months (corrected: %v)", childID, ageInMonths, useCorrected)

	// Get incomplete milestones for this child
	incompleteMilestones, err := getIncompleteMilestones(childID, ageInMonths)
	if err != nil {
		c.Logger().Errorf("Failed to get incomplete milestones: %v", err)
		// Continue anyway, we can still provide general recommendations
		incompleteMilestones = []models.Milestone{}
	}

	// Get recommendations based on incomplete milestones and category needs
	recommendations, err := getRecommendationsForChild(childID, ageInMonths, incompleteMilestones)
	if err != nil {
		c.Logger().Errorf("Failed to get recommendations: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get recommendations"})
	}

	return c.JSON(http.StatusOK, models.RecommendationsResponse{
		ChildID:        childID,
		AgeMonths:      ageInMonths,
		Recommendations: recommendations,
	})
}

// getIncompleteMilestones fetches milestones that are not yet completed (status = "no" or not assessed)
func getIncompleteMilestones(childID string, ageMonths int) ([]models.Milestone, error) {
	query := `
		SELECT m.* FROM milestones m
		LEFT JOIN assessments a ON a.milestone_id = m.id AND a.child_id = $1
		WHERE (m.age_months <= $2 + 3) -- Include milestones up to 3 months ahead
		  AND (a.status IS NULL OR a.status = 'no')
		ORDER BY m.age_months ASC, m.pyramid_level ASC
		LIMIT 20
	`

	var milestones []models.Milestone
	err := db.DB.Select(&milestones, query, childID, ageMonths)
	if err != nil {
		return nil, err
	}

	return milestones, nil
}

// getRecommendationsForChild fetches stimulation content recommendations
func getRecommendationsForChild(childID string, ageMonths int, incompleteMilestones []models.Milestone) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	// Strategy 1: Get recommendations for incomplete milestones
	milestoneIDs := make(map[string]bool)
	categories := make(map[string]bool)
	for _, milestone := range incompleteMilestones {
		milestoneIDs[milestone.ID] = true
		categories[milestone.Category] = true
	}

	// Get specific recommendations for milestones
	if len(milestoneIDs) > 0 {
		ids := make([]string, 0, len(milestoneIDs))
		for id := range milestoneIDs {
			ids = append(ids, id)
		}

		query := `
			SELECT * FROM stimulation_content
			WHERE milestone_id = ANY($1)
			  AND is_active = true
			ORDER BY created_at DESC
		`

		var contents []models.StimulationContent
		err := db.DB.Select(&contents, query, ids)
		if err == nil {
			for _, content := range contents {
				// Find related milestone
				var relatedMilestone *models.Milestone
				for _, m := range incompleteMilestones {
					if m.ID == *content.MilestoneID {
						relatedMilestone = &m
						break
					}
				}

				recommendations = append(recommendations, models.Recommendation{
					Content:          content,
					Reason:           "Direkomendasikan berdasarkan milestone yang belum tercapai",
					Priority:         "high",
					RelatedMilestone: relatedMilestone,
				})
			}
		}
	}

	// Strategy 2: Get category-based recommendations for categories that need stimulation
	if len(categories) > 0 {
		cats := make([]string, 0, len(categories))
		for cat := range categories {
			cats = append(cats, cat)
		}

		query := `
			SELECT * FROM stimulation_content
			WHERE category = ANY($1)
			  AND (age_min_months IS NULL OR age_min_months <= $2)
			  AND (age_max_months IS NULL OR age_max_months >= $2)
			  AND milestone_id IS NULL -- General category recommendations, not specific to milestone
			  AND is_active = true
			ORDER BY created_at DESC
			LIMIT 10
		`

		var contents []models.StimulationContent
		err := db.DB.Select(&contents, query, cats, ageMonths)
		if err == nil {
			for _, content := range contents {
				// Avoid duplicates
				isDuplicate := false
				for _, rec := range recommendations {
					if rec.Content.ID == content.ID {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
				categoryNames := map[string]string{
					"sensory":    "Sensorik",
					"motor":      "Motorik",
					"perception": "Persepsi",
					"cognitive":  "Kognitif",
				}
				categoryName := content.Category
				if name, ok := categoryNames[content.Category]; ok {
					categoryName = name
				}

					recommendations = append(recommendations, models.Recommendation{
						Content:  content,
						Reason:   "Direkomendasikan untuk stimulasi kategori " + categoryName,
						Priority: "medium",
					})
				}
			}
		}
	}

	// Strategy 3: Get age-appropriate general recommendations if we don't have many yet
	if len(recommendations) < 5 {
		query := `
			SELECT * FROM stimulation_content
			WHERE (age_min_months IS NULL OR age_min_months <= $1)
			  AND (age_max_months IS NULL OR age_max_months >= $1)
			  AND milestone_id IS NULL
			  AND is_active = true
			ORDER BY created_at DESC
			LIMIT 5
		`

		var contents []models.StimulationContent
		err := db.DB.Select(&contents, query, ageMonths)
		if err == nil {
			for _, content := range contents {
				// Avoid duplicates
				isDuplicate := false
				for _, rec := range recommendations {
					if rec.Content.ID == content.ID {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					recommendations = append(recommendations, models.Recommendation{
						Content:  content,
						Reason:   "Direkomendasikan untuk usia " + strconv.Itoa(ageMonths) + " bulan",
						Priority: "low",
					})
				}
			}
		}
	}

	// Limit total recommendations
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations, nil
}

