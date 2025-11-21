package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetImmunizationSchedule fetches immunization schedule and status for a child
func GetImmunizationSchedule(c echo.Context) error {
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

	// Calculate current age
	today := time.Now().Format("2006-01-02")
	ageInDays, ageInMonths, _, err := utils.CalculateCorrectedAge(
		child.DOB, today, child.IsPremature, child.GestationalAge)
	if err != nil {
		c.Logger().Errorf("Failed to calculate age: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to calculate child age"})
	}

	c.Logger().Infof("Getting immunization schedule for child %s, age: %d days (%d months)", childID, ageInDays, ageInMonths)

	// Get all active immunization schedules
	schedules, err := getActiveImmunizationSchedules()
	if err != nil {
		c.Logger().Errorf("Failed to get immunization schedules: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get immunization schedules"})
	}

	// Get completed immunizations for this child
	completedImmunizations, err := getCompletedImmunizations(childID)
	if err != nil {
		c.Logger().Errorf("Failed to get completed immunizations: %v", err)
		completedImmunizations = map[string]*models.ChildImmunization{} // Empty map
	}

	// Calculate status for each immunization
	statuses := []models.ImmunizationStatus{}
	summary := models.ImmunizationSummary{
		Total: len(schedules),
	}

	for _, schedule := range schedules {
		status := calculateImmunizationStatus(schedule, child.DOB, ageInDays, ageInMonths, completedImmunizations[schedule.ID])
		
		statuses = append(statuses, status)

		// Update summary
		switch status.Status {
		case "completed":
			summary.Completed++
		case "pending":
			summary.Pending++
		case "overdue":
			summary.Overdue++
		case "upcoming":
			summary.Upcoming++
		}
	}

	return c.JSON(http.StatusOK, models.ImmunizationScheduleResponse{
		ChildID:       childID,
		AgeMonths:     ageInMonths,
		AgeDays:       ageInDays,
		Immunizations: statuses,
		Summary:       summary,
	})
}

// RecordImmunization records an immunization given to a child
func RecordImmunization(c echo.Context) error {
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
	err = db.DB.QueryRow("SELECT id, dob FROM children WHERE id = $1", childID).
		Scan(&child.ID, &child.DOB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get child data"})
	}

	// Parse request
	req := new(models.ImmunizationRecordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Validate required fields
	if req.ImmunizationScheduleID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "immunization_schedule_id is required"})
	}
	if req.GivenDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "given_date is required"})
	}

	// Get immunization schedule
	var schedule models.ImmunizationSchedule
	err = db.DB.Get(&schedule, "SELECT * FROM immunization_schedule WHERE id = $1", req.ImmunizationScheduleID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Immunization schedule not found"})
	}

	// Calculate age at given date
	givenAgeDays, givenAgeMonths, _, err := utils.CalculateCorrectedAge(
		child.DOB, req.GivenDate, child.IsPremature, child.GestationalAge)
	if err != nil {
		c.Logger().Errorf("Failed to calculate age at given date: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid given_date format"})
	}

	// Check if on schedule
	isOnSchedule := false
	if schedule.AgeOptimalDays != nil {
		// Consider on schedule if within ±7 days of optimal age
		diff := givenAgeDays - *schedule.AgeOptimalDays
		if diff >= -7 && diff <= 7 {
			isOnSchedule = true
		}
	}

	// Check if catch-up
	isCatchUp := false
	if schedule.AgeOptimalDays != nil && givenAgeDays > *schedule.AgeOptimalDays+7 {
		isCatchUp = true
	}

	// Insert immunization record
	query := `
		INSERT INTO child_immunizations (
			child_id, immunization_schedule_id,
			given_date, given_at_age_days, given_at_age_months,
			location, healthcare_facility, doctor_name, vaccine_batch_number, notes,
			is_on_schedule, is_catch_up
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (child_id, immunization_schedule_id) 
		DO UPDATE SET
			given_date = EXCLUDED.given_date,
			given_at_age_days = EXCLUDED.given_at_age_days,
			given_at_age_months = EXCLUDED.given_at_age_months,
			location = EXCLUDED.location,
			healthcare_facility = EXCLUDED.healthcare_facility,
			doctor_name = EXCLUDED.doctor_name,
			vaccine_batch_number = EXCLUDED.vaccine_batch_number,
			notes = EXCLUDED.notes,
			is_on_schedule = EXCLUDED.is_on_schedule,
			is_catch_up = EXCLUDED.is_catch_up,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	var recordID string
	err = db.DB.QueryRow(query,
		childID, req.ImmunizationScheduleID,
		req.GivenDate, givenAgeDays, givenAgeMonths,
		req.Location, req.HealthcareFacility, req.DoctorName, req.VaccineBatchNumber, req.Notes,
		isOnSchedule, isCatchUp,
	).Scan(&recordID)

	if err != nil {
		c.Logger().Errorf("Failed to record immunization: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to record immunization"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":      recordID,
		"message": "Imunisasi berhasil dicatat",
		"status":  "completed",
	})
}

// Helper functions

func getActiveImmunizationSchedules() ([]models.ImmunizationSchedule, error) {
	query := `SELECT * FROM immunization_schedule WHERE is_active = true ORDER BY age_optimal_days ASC, dose_number ASC`
	var schedules []models.ImmunizationSchedule
	err := db.DB.Select(&schedules, query)
	return schedules, err
}

func getCompletedImmunizations(childID string) (map[string]*models.ChildImmunization, error) {
	query := `
		SELECT * FROM child_immunizations
		WHERE child_id = $1
	`
	
	var records []models.ChildImmunization
	err := db.DB.Select(&records, query, childID)
	if err != nil {
		return nil, err
	}

	// Map by schedule ID for quick lookup
	completed := make(map[string]*models.ChildImmunization)
	for i := range records {
		completed[records[i].ImmunizationScheduleID] = &records[i]
	}

	return completed, nil
}

func calculateImmunizationStatus(
	schedule models.ImmunizationSchedule,
	dob string,
	currentAgeDays int,
	currentAgeMonths int,
	record *models.ChildImmunization,
) models.ImmunizationStatus {
	status := models.ImmunizationStatus{
		Schedule: schedule,
	}

	// If already completed
	if record != nil {
		status.Status = "completed"
		status.Record = record
		return status
	}

	// Calculate due date and status
	if schedule.AgeOptimalDays == nil {
		status.Status = "pending"
		return status
	}

	// Parse DOB and calculate due date
	dobTime, err := time.Parse("2006-01-02", dob)
	if err != nil {
		status.Status = "pending"
		return status
	}

	dueDate := dobTime.AddDate(0, 0, *schedule.AgeOptimalDays)
	dueDateStr := dueDate.Format("2006-01-02")
	status.DueDate = &dueDateStr

	if schedule.AgeOptimalMonths != nil {
		status.DueAgeMonths = schedule.AgeOptimalMonths
	}

	// Calculate days until due or days overdue
	today := time.Now()
	daysDiff := int(today.Sub(dueDate).Hours() / 24)

	if daysDiff < 0 {
		// Future - upcoming
		status.Status = "upcoming"
		daysUntil := -daysDiff
		status.DaysUntilDue = &daysUntil
		if daysUntil <= 7 {
			status.Status = "upcoming" // Within 7 days
		}
	} else if daysDiff <= 7 {
		// Within 7 days of due date
		status.Status = "upcoming"
		daysUntil := daysDiff
		status.DaysUntilDue = &daysUntil
	} else if schedule.AgeMaxDays != nil && currentAgeDays <= *schedule.AgeMaxDays {
		// Overdue but still within max age
		status.Status = "overdue"
		status.DaysOverdue = &daysDiff
	} else {
		// Past max age - still pending but late
		status.Status = "overdue"
		if schedule.AgeMaxDays != nil {
			overdue := currentAgeDays - *schedule.AgeMaxDays
			status.DaysOverdue = &overdue
		}
	}

	// Default to pending if status not set
	if status.Status == "" {
		status.Status = "pending"
	}

	return status
}

