package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"tukem-backend/models"
	"tukem-backend/utils"
)

type MilestoneHandler struct {
	DB *sqlx.DB
}

func NewMilestoneHandler(db *sqlx.DB) *MilestoneHandler {
	return &MilestoneHandler{DB: db}
}

// GetMilestones fetches milestones based on age window logic
func (h *MilestoneHandler) GetMilestones(c echo.Context) error {
	ageMonthsStr := c.QueryParam("age_months")
	childID := c.QueryParam("child_id") // Optional child_id to get corrected age
	
	// Default to 0 if not provided
	ageMonths := 0
	if ageMonthsStr != "" {
		var err error
		ageMonths, err = strconv.Atoi(ageMonthsStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid age_months parameter"})
		}
	}
	
	// If child_id is provided, calculate corrected age if applicable
	if childID != "" {
		// Get user ID from JWT to verify ownership
		user := c.Get("user").(*jwt.Token)
		claims := *user.Claims.(*jwt.MapClaims)
		userID := claims["user_id"].(string)
		
		// Verify child belongs to user and get child data
		var child models.Child
		err := h.DB.QueryRow("SELECT id, dob, is_premature, gestational_age FROM children WHERE id = $1 AND parent_id = $2", 
			childID, userID).Scan(&child.ID, &child.DOB, &child.IsPremature, &child.GestationalAge)
		if err == nil {
			// Calculate corrected age using current date as measurement date
			today := time.Now().Format("2006-01-02")
			_, correctedMonths, useCorrected, err := utils.CalculateCorrectedAge(
				child.DOB, today, child.IsPremature, child.GestationalAge)
			if err == nil && useCorrected {
				ageMonths = correctedMonths
				c.Logger().Infof("Using corrected age %d months for premature child %s", ageMonths, childID)
			}
		}
	}

	// Logic: Fetch milestones for:
	// 1. Current age target (e.g., if child is 7 months, target is 9 months or 6 months depending on proximity)
	// For simplicity in this MVP, we'll fetch milestones where age_months is close to the child's age
	// plus any previous milestones that haven't been completed (if child_id is provided, but here we just fetch general milestones)
	
	// Refined Logic:
	// Fetch milestones where age_months is within a range.
	// If age is 7 months, we might want to see 6mo and 9mo milestones.
	
	// Let's fetch milestones for the closest standard age groups
	// Standard ages: 3, 6, 9, 12, 18, 24, 36, 48, 60
	
	query := `
		SELECT * FROM milestones 
		WHERE age_months >= $1 - 3 AND age_months <= $1 + 6
		ORDER BY age_months ASC, pyramid_level ASC
	`
	
	var milestones []models.Milestone
	err := h.DB.Select(&milestones, query, ageMonths)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch milestones"})
	}

	return c.JSON(http.StatusOK, milestones)
}

// GetChildAssessments fetches all assessments for a child
func (h *MilestoneHandler) GetChildAssessments(c echo.Context) error {
	childID := c.Param("id")

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := h.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
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
	
	query := `
		SELECT 
			a.id,
			a.child_id,
			a.milestone_id,
			a.assessment_date,
			a.status,
			a.notes,
			a.created_at,
			a.updated_at,
			m.id as "milestone.id",
			m.age_months as "milestone.age_months",
			m.category as "milestone.category",
			m.question as "milestone.question",
			m.pyramid_level as "milestone.pyramid_level",
			m.source as "milestone.source",
			m.denver_domain as "milestone.denver_domain"
		FROM assessments a
		JOIN milestones m ON a.milestone_id = m.id
		WHERE a.child_id = $1
		ORDER BY a.assessment_date DESC, a.created_at DESC
	`
	
	rows, err := h.DB.Queryx(query, childID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch assessments: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch assessments"})
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		var m models.Milestone
		
		err := rows.Scan(
			&a.ID,
			&a.ChildID,
			&a.MilestoneID,
			&a.AssessmentDate,
			&a.Status,
			&a.Notes,
			&a.CreatedAt,
			&a.UpdatedAt,
			&m.ID,
			&m.AgeMonths,
			&m.Category,
			&m.Question,
			&m.PyramidLevel,
			&m.Source,
			&m.DenverDomain,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan assessment: %v", err)
			continue
		}
		
		a.Milestone = &m
		assessments = append(assessments, a)
	}

	return c.JSON(http.StatusOK, assessments)
}

// BatchUpsertAssessments handles bulk update of assessments
func (h *MilestoneHandler) BatchUpsertAssessments(c echo.Context) error {
	childID := c.Param("id")

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := h.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
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
	
	var req models.BatchAssessmentRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Manual validation
	if req.AssessmentDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tanggal penilaian harus diisi"})
	}

	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Minimal 1 item penilaian harus diisi"})
	}

	// Validate each item
	validStatuses := map[string]bool{"yes": true, "no": true, "sometimes": true}
	for i, item := range req.Items {
		if item.MilestoneID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "milestone_id tidak boleh kosong pada item ke-" + strconv.Itoa(i+1),
			})
		}
		if !validStatuses[item.Status] {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Status harus 'yes', 'no', atau 'sometimes' pada item ke-" + strconv.Itoa(i+1),
			})
		}
	}

	c.Logger().Infof("Saving %d assessments for child %s on date %s", len(req.Items), childID, req.AssessmentDate)

	tx, err := h.DB.Beginx()
	if err != nil {
		c.Logger().Errorf("Failed to begin transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memulai transaksi database"})
	}
	defer tx.Rollback()

	// Prepare upsert statement
	// PostgreSQL specific ON CONFLICT syntax
	stmt, err := tx.PrepareNamed(`
		INSERT INTO assessments (child_id, milestone_id, assessment_date, status, notes, updated_at)
		VALUES (:child_id, :milestone_id, :assessment_date, :status, :notes, NOW())
		ON CONFLICT (child_id, milestone_id) 
		DO UPDATE SET 
			status = :status,
			notes = :notes,
			assessment_date = :assessment_date,
			updated_at = NOW()
	`)
	if err != nil {
		c.Logger().Errorf("Failed to prepare statement: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menyiapkan query database"})
	}

	for i, item := range req.Items {
		assessment := map[string]interface{}{
			"child_id":        childID,
			"milestone_id":    item.MilestoneID,
			"assessment_date": req.AssessmentDate,
			"status":          item.Status,
			"notes":           item.Notes,
		}
		
		if _, err := stmt.Exec(assessment); err != nil {
			c.Logger().Errorf("Failed to upsert assessment for milestone %s (item %d): %v", item.MilestoneID, i+1, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Gagal menyimpan penilaian untuk milestone: " + item.MilestoneID + ". Pastikan milestone_id valid.",
				"details": err.Error(),
			})
		}
	}

	if err := tx.Commit(); err != nil {
		c.Logger().Errorf("Failed to commit transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan ke database"})
	}

	c.Logger().Infof("Successfully saved %d assessments for child %s", len(req.Items), childID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Penilaian berhasil disimpan",
		"count":   strconv.Itoa(len(req.Items)),
	})
}

// GetAssessmentSummary calculates pyramid health and returns summary
func (h *MilestoneHandler) GetAssessmentSummary(c echo.Context) error {
	childID := c.Param("id")

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := h.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
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

	// 1. Fetch all assessments for this child joined with milestones
	// Only include KPSP milestones for pyramid calculation (Denver II uses different domain system)
	query := `
		SELECT a.status, m.category, m.pyramid_level, m.is_red_flag, m.question
		FROM assessments a
		JOIN milestones m ON a.milestone_id = m.id
		WHERE a.child_id = $1
			AND m.source = 'KPSP'
	`

	rows, err := h.DB.Queryx(query, childID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch assessment data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch assessment data"})
	}
	defer rows.Close()

	// Data structures for calculation
	type AssessmentData struct {
		Status       string
		Category     string
		PyramidLevel int
		IsRedFlag    bool
		Question     string
	}

	var data []AssessmentData
	for rows.Next() {
		var d AssessmentData
		if err := rows.Scan(&d.Status, &d.Category, &d.PyramidLevel, &d.IsRedFlag, &d.Question); err != nil {
			continue
		}
		data = append(data, d)
	}

	// 2. Calculate scores
	totalByLevel := make(map[int]int)
	completedByLevel := make(map[int]int)
	redFlags := []models.Milestone{}
	
	for _, d := range data {
		totalByLevel[d.PyramidLevel]++
		if d.Status == "yes" {
			completedByLevel[d.PyramidLevel]++
		}
		
		if d.IsRedFlag && d.Status == "no" {
			redFlags = append(redFlags, models.Milestone{
				Question: d.Question,
				Category: d.Category,
			})
		}
	}

	// 3. Calculate percentages
	progressByCategory := make(map[string]float64)
	// Map levels to categories for response
	levelToCategory := map[int]string{
		1: "sensory",
		2: "motor",
		3: "perception",
		4: "cognitive",
	}

	for level, total := range totalByLevel {
		if total > 0 {
			cat := levelToCategory[level]
			progressByCategory[cat] = float64(completedByLevel[level]) / float64(total) * 100
		}
	}

	// 4. Logic Warnings (Pyramid Imbalance)
	warnings := []string{}
	
	// Check: High Cognitive (Level 4) but Low Sensory (Level 1)
	// Thresholds: Cognitive > 70% and Sensory < 50%
	sensoryScore := 0.0
	if totalByLevel[1] > 0 {
		sensoryScore = float64(completedByLevel[1]) / float64(totalByLevel[1]) * 100
	}
	
	cognitiveScore := 0.0
	if totalByLevel[4] > 0 {
		cognitiveScore = float64(completedByLevel[4]) / float64(totalByLevel[4]) * 100
	}

	if cognitiveScore > 70 && sensoryScore < 50 {
		warnings = append(warnings, "Terdeteksi 'Lompatan Perkembangan'. Anak mahir kognitif tapi pondasi sensorik (Level 1) belum kuat. Risiko: Masalah fokus/emosi di kemudian hari.")
	}

	summary := models.AssessmentSummary{
		TotalMilestones:     len(data),
		CompletedMilestones: len(data), // This logic might need adjustment, currently just count of assessed items
		ProgressByCategory:  progressByCategory,
		RedFlagsDetected:    redFlags,
		PyramidWarnings:     warnings,
	}

	return c.JSON(http.StatusOK, summary)
}
