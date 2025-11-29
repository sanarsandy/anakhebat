package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetAdminImmunizationSchedules returns all immunization schedules with pagination and filters
func GetAdminImmunizationSchedules(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	category := c.QueryParam("category")
	priority := c.QueryParam("priority")
	isRequired := c.QueryParam("is_required")
	search := c.QueryParam("search")

	// Build query
	query := `SELECT id, name, name_id, description, age_min_days, age_optimal_days, age_max_days,
	          age_min_months, age_optimal_months, age_max_months, dose_number, total_doses,
	          interval_from_previous_days, interval_from_previous_months, category, priority,
	          is_required, notes, source, is_active, created_at, updated_at
	          FROM immunization_schedule WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if priority != "" {
		query += ` AND priority = $` + strconv.Itoa(argIndex)
		args = append(args, priority)
		argIndex++
	}

	if isRequired == "true" {
		query += ` AND is_required = true`
	} else if isRequired == "false" {
		query += ` AND is_required = false`
	}

	if search != "" {
		query += ` AND (name ILIKE $` + strconv.Itoa(argIndex) +
			` OR name_id ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM immunization_schedule WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if category != "" {
		countQuery += ` AND category = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, category)
		countArgIndex++
	}
	if priority != "" {
		countQuery += ` AND priority = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, priority)
		countArgIndex++
	}
	if isRequired == "true" {
		countQuery += ` AND is_required = true`
	} else if isRequired == "false" {
		countQuery += ` AND is_required = false`
	}
	if search != "" {
		countQuery += ` AND (name ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR name_id ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
		countArgIndex += 2
	}

	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminImmunizationSchedules count error: %v", err)
		total = 0
	}

	query += ` ORDER BY age_optimal_days ASC, dose_number ASC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminImmunizationSchedules query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	type ImmunizationSchedule struct {
		ID                        string  `json:"id"`
		Name                      string  `json:"name"`
		NameID                    string  `json:"name_id"`
		Description               string  `json:"description"`
		AgeMinDays                *int    `json:"age_min_days"`
		AgeOptimalDays            *int    `json:"age_optimal_days"`
		AgeMaxDays                *int    `json:"age_max_days"`
		AgeMinMonths              *int    `json:"age_min_months"`
		AgeOptimalMonths          *int    `json:"age_optimal_months"`
		AgeMaxMonths              *int    `json:"age_max_months"`
		DoseNumber                int     `json:"dose_number"`
		TotalDoses                *int    `json:"total_doses"`
		IntervalFromPreviousDays  *int    `json:"interval_from_previous_days"`
		IntervalFromPreviousMonths *int   `json:"interval_from_previous_months"`
		Category                  string  `json:"category"`
		Priority                  string  `json:"priority"`
		IsRequired                bool    `json:"is_required"`
		Notes                     string  `json:"notes"`
		Source                    string  `json:"source"`
		IsActive                  bool    `json:"is_active"`
		CreatedAt                 string  `json:"created_at"`
		UpdatedAt                 string  `json:"updated_at"`
	}

	var schedules []ImmunizationSchedule
	for rows.Next() {
		var s ImmunizationSchedule
		var nameID sql.NullString
		var description sql.NullString
		var ageMinDays, ageOptimalDays, ageMaxDays sql.NullInt64
		var ageMinMonths, ageOptimalMonths, ageMaxMonths sql.NullInt64
		var totalDoses sql.NullInt64
		var intervalDays, intervalMonths sql.NullInt64
		var notes sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		err := rows.Scan(
			&s.ID, &s.Name, &nameID, &description, &ageMinDays, &ageOptimalDays, &ageMaxDays,
			&ageMinMonths, &ageOptimalMonths, &ageMaxMonths, &s.DoseNumber, &totalDoses,
			&intervalDays, &intervalMonths, &s.Category, &s.Priority, &s.IsRequired,
			&notes, &s.Source, &s.IsActive, &createdAt, &updatedAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan immunization schedule row: %v", err)
			continue
		}

		if nameID.Valid {
			s.NameID = nameID.String
		}
		if description.Valid {
			s.Description = description.String
		}
		if ageMinDays.Valid {
			age := int(ageMinDays.Int64)
			s.AgeMinDays = &age
		}
		if ageOptimalDays.Valid {
			age := int(ageOptimalDays.Int64)
			s.AgeOptimalDays = &age
		}
		if ageMaxDays.Valid {
			age := int(ageMaxDays.Int64)
			s.AgeMaxDays = &age
		}
		if ageMinMonths.Valid {
			age := int(ageMinMonths.Int64)
			s.AgeMinMonths = &age
		}
		if ageOptimalMonths.Valid {
			age := int(ageOptimalMonths.Int64)
			s.AgeOptimalMonths = &age
		}
		if ageMaxMonths.Valid {
			age := int(ageMaxMonths.Int64)
			s.AgeMaxMonths = &age
		}
		if totalDoses.Valid {
			doses := int(totalDoses.Int64)
			s.TotalDoses = &doses
		}
		if intervalDays.Valid {
			interval := int(intervalDays.Int64)
			s.IntervalFromPreviousDays = &interval
		}
		if intervalMonths.Valid {
			interval := int(intervalMonths.Int64)
			s.IntervalFromPreviousMonths = &interval
		}
		if notes.Valid {
			s.Notes = notes.String
		}
		if createdAt.Valid {
			s.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		if updatedAt.Valid {
			s.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		schedules = append(schedules, s)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"schedules": schedules,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminImmunizationSchedule returns a single immunization schedule
func GetAdminImmunizationSchedule(c echo.Context) error {
	scheduleID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(scheduleID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid schedule ID format"})
	}

	type ImmunizationSchedule struct {
		ID                        string  `json:"id"`
		Name                      string  `json:"name"`
		NameID                    string  `json:"name_id"`
		Description               string  `json:"description"`
		AgeMinDays                *int    `json:"age_min_days"`
		AgeOptimalDays            *int    `json:"age_optimal_days"`
		AgeMaxDays                *int    `json:"age_max_days"`
		AgeMinMonths              *int    `json:"age_min_months"`
		AgeOptimalMonths          *int    `json:"age_optimal_months"`
		AgeMaxMonths              *int    `json:"age_max_months"`
		DoseNumber                int     `json:"dose_number"`
		TotalDoses                *int    `json:"total_doses"`
		IntervalFromPreviousDays  *int    `json:"interval_from_previous_days"`
		IntervalFromPreviousMonths *int   `json:"interval_from_previous_months"`
		Category                  string  `json:"category"`
		Priority                  string  `json:"priority"`
		IsRequired                bool    `json:"is_required"`
		Notes                     string  `json:"notes"`
		Source                    string  `json:"source"`
		IsActive                  bool    `json:"is_active"`
		CreatedAt                 string  `json:"created_at"`
		UpdatedAt                 string  `json:"updated_at"`
	}

	var s ImmunizationSchedule
	var nameID sql.NullString
	var description sql.NullString
	var ageMinDays, ageOptimalDays, ageMaxDays sql.NullInt64
	var ageMinMonths, ageOptimalMonths, ageMaxMonths sql.NullInt64
	var totalDoses sql.NullInt64
	var intervalDays, intervalMonths sql.NullInt64
	var notes sql.NullString
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	err := db.DB.QueryRow(
		`SELECT id, name, name_id, description, age_min_days, age_optimal_days, age_max_days,
		 age_min_months, age_optimal_months, age_max_months, dose_number, total_doses,
		 interval_from_previous_days, interval_from_previous_months, category, priority,
		 is_required, notes, source, is_active, created_at, updated_at
		 FROM immunization_schedule WHERE id = $1`,
		scheduleID,
	).Scan(
		&s.ID, &s.Name, &nameID, &description, &ageMinDays, &ageOptimalDays, &ageMaxDays,
		&ageMinMonths, &ageOptimalMonths, &ageMaxMonths, &s.DoseNumber, &totalDoses,
		&intervalDays, &intervalMonths, &s.Category, &s.Priority, &s.IsRequired,
		&notes, &s.Source, &s.IsActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Immunization schedule not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetAdminImmunizationSchedule error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if nameID.Valid {
		s.NameID = nameID.String
	}
	if description.Valid {
		s.Description = description.String
	}
	if ageMinDays.Valid {
		age := int(ageMinDays.Int64)
		s.AgeMinDays = &age
	}
	if ageOptimalDays.Valid {
		age := int(ageOptimalDays.Int64)
		s.AgeOptimalDays = &age
	}
	if ageMaxDays.Valid {
		age := int(ageMaxDays.Int64)
		s.AgeMaxDays = &age
	}
	if ageMinMonths.Valid {
		age := int(ageMinMonths.Int64)
		s.AgeMinMonths = &age
	}
	if ageOptimalMonths.Valid {
		age := int(ageOptimalMonths.Int64)
		s.AgeOptimalMonths = &age
	}
	if ageMaxMonths.Valid {
		age := int(ageMaxMonths.Int64)
		s.AgeMaxMonths = &age
	}
	if totalDoses.Valid {
		doses := int(totalDoses.Int64)
		s.TotalDoses = &doses
	}
	if intervalDays.Valid {
		interval := int(intervalDays.Int64)
		s.IntervalFromPreviousDays = &interval
	}
	if intervalMonths.Valid {
		interval := int(intervalMonths.Int64)
		s.IntervalFromPreviousMonths = &interval
	}
	if notes.Valid {
		s.Notes = notes.String
	}
	if createdAt.Valid {
		s.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if updatedAt.Valid {
		s.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(http.StatusOK, s)
}

// CreateAdminImmunizationSchedule creates a new immunization schedule
func CreateAdminImmunizationSchedule(c echo.Context) error {
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		Name                      string  `json:"name" validate:"required"`
		NameID                    string  `json:"name_id"`
		Description               string  `json:"description"`
		AgeMinDays                *int    `json:"age_min_days"`
		AgeOptimalDays            *int    `json:"age_optimal_days"`
		AgeMaxDays                *int    `json:"age_max_days"`
		AgeMinMonths              *int    `json:"age_min_months"`
		AgeOptimalMonths          *int    `json:"age_optimal_months"`
		AgeMaxMonths              *int    `json:"age_max_months"`
		DoseNumber                int     `json:"dose_number" validate:"required"`
		TotalDoses                *int    `json:"total_doses"`
		IntervalFromPreviousDays  *int    `json:"interval_from_previous_days"`
		IntervalFromPreviousMonths *int   `json:"interval_from_previous_months"`
		Category                  string  `json:"category"`
		Priority                  string  `json:"priority"`
		IsRequired                bool    `json:"is_required"`
		Notes                     string  `json:"notes"`
		Source                    string  `json:"source"`
		IsActive                  bool    `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Set defaults
	if req.Category == "" {
		req.Category = "wajib"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if req.Source == "" {
		req.Source = "IDAI"
	}

	var scheduleID string
	err := db.DB.QueryRow(
		`INSERT INTO immunization_schedule (name, name_id, description, age_min_days, age_optimal_days, age_max_days,
		 age_min_months, age_optimal_months, age_max_months, dose_number, total_doses,
		 interval_from_previous_days, interval_from_previous_months, category, priority,
		 is_required, notes, source, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		 RETURNING id`,
		req.Name, req.NameID, req.Description, req.AgeMinDays, req.AgeOptimalDays, req.AgeMaxDays,
		req.AgeMinMonths, req.AgeOptimalMonths, req.AgeMaxMonths, req.DoseNumber, req.TotalDoses,
		req.IntervalFromPreviousDays, req.IntervalFromPreviousMonths, req.Category, req.Priority,
		req.IsRequired, req.Notes, req.Source, req.IsActive,
	).Scan(&scheduleID)

	if err != nil {
		c.Logger().Errorf("CreateAdminImmunizationSchedule error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	scheduleData := map[string]interface{}{
		"name":        req.Name,
		"category":    req.Category,
		"priority":    req.Priority,
		"dose_number": req.DoseNumber,
	}
	utils.LogAudit(adminUserID, "create", "immunization_schedule", &scheduleID, nil, scheduleData, ipAddress, userAgent)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      scheduleID,
		"message": "Immunization schedule created successfully",
	})
}

// UpdateAdminImmunizationSchedule updates an immunization schedule
func UpdateAdminImmunizationSchedule(c echo.Context) error {
	scheduleID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(scheduleID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid schedule ID format"})
	}

	// Get existing schedule for audit log
	var existing struct {
		Name     string
		Category string
		Priority string
	}

	err := db.DB.QueryRow(
		`SELECT name, category, priority FROM immunization_schedule WHERE id = $1`,
		scheduleID,
	).Scan(&existing.Name, &existing.Category, &existing.Priority)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Immunization schedule not found"})
	}
	if err != nil {
		c.Logger().Errorf("UpdateAdminImmunizationSchedule get existing error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var req struct {
		Name                      *string `json:"name"`
		NameID                    *string `json:"name_id"`
		Description               *string `json:"description"`
		AgeMinDays                *int    `json:"age_min_days"`
		AgeOptimalDays            *int    `json:"age_optimal_days"`
		AgeMaxDays                *int    `json:"age_max_days"`
		AgeMinMonths              *int    `json:"age_min_months"`
		AgeOptimalMonths          *int    `json:"age_optimal_months"`
		AgeMaxMonths              *int    `json:"age_max_months"`
		DoseNumber                *int    `json:"dose_number"`
		TotalDoses                *int    `json:"total_doses"`
		IntervalFromPreviousDays  *int    `json:"interval_from_previous_days"`
		IntervalFromPreviousMonths *int   `json:"interval_from_previous_months"`
		Category                  *string `json:"category"`
		Priority                  *string `json:"priority"`
		IsRequired                *bool   `json:"is_required"`
		Notes                     *string `json:"notes"`
		Source                    *string `json:"source"`
		IsActive                  *bool   `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		updateFields = append(updateFields, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.NameID != nil {
		updateFields = append(updateFields, "name_id = $"+strconv.Itoa(argIndex))
		args = append(args, *req.NameID)
		argIndex++
	}
	if req.Description != nil {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.AgeMinDays != nil {
		updateFields = append(updateFields, "age_min_days = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMinDays)
		argIndex++
	}
	if req.AgeOptimalDays != nil {
		updateFields = append(updateFields, "age_optimal_days = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeOptimalDays)
		argIndex++
	}
	if req.AgeMaxDays != nil {
		updateFields = append(updateFields, "age_max_days = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMaxDays)
		argIndex++
	}
	if req.AgeMinMonths != nil {
		updateFields = append(updateFields, "age_min_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMinMonths)
		argIndex++
	}
	if req.AgeOptimalMonths != nil {
		updateFields = append(updateFields, "age_optimal_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeOptimalMonths)
		argIndex++
	}
	if req.AgeMaxMonths != nil {
		updateFields = append(updateFields, "age_max_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMaxMonths)
		argIndex++
	}
	if req.DoseNumber != nil {
		updateFields = append(updateFields, "dose_number = $"+strconv.Itoa(argIndex))
		args = append(args, *req.DoseNumber)
		argIndex++
	}
	if req.TotalDoses != nil {
		updateFields = append(updateFields, "total_doses = $"+strconv.Itoa(argIndex))
		args = append(args, *req.TotalDoses)
		argIndex++
	}
	if req.IntervalFromPreviousDays != nil {
		updateFields = append(updateFields, "interval_from_previous_days = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IntervalFromPreviousDays)
		argIndex++
	}
	if req.IntervalFromPreviousMonths != nil {
		updateFields = append(updateFields, "interval_from_previous_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IntervalFromPreviousMonths)
		argIndex++
	}
	if req.Category != nil {
		updateFields = append(updateFields, "category = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Category)
		argIndex++
	}
	if req.Priority != nil {
		updateFields = append(updateFields, "priority = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Priority)
		argIndex++
	}
	if req.IsRequired != nil {
		updateFields = append(updateFields, "is_required = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsRequired)
		argIndex++
	}
	if req.Notes != nil {
		updateFields = append(updateFields, "notes = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}
	if req.Source != nil {
		updateFields = append(updateFields, "source = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Source)
		argIndex++
	}
	if req.IsActive != nil {
		updateFields = append(updateFields, "is_active = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(updateFields) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, scheduleID)
	query := `UPDATE immunization_schedule SET ` + updateFields[0]
	for i := 1; i < len(updateFields); i++ {
		query += `, ` + updateFields[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		c.Logger().Errorf("UpdateAdminImmunizationSchedule error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	beforeData := map[string]interface{}{
		"name":     existing.Name,
		"category": existing.Category,
		"priority": existing.Priority,
	}
	afterData := map[string]interface{}{}
	if req.Name != nil {
		afterData["name"] = *req.Name
	}
	if req.Category != nil {
		afterData["category"] = *req.Category
	}
	if req.Priority != nil {
		afterData["priority"] = *req.Priority
	}
	utils.LogAudit(adminUserID, "update", "immunization_schedule", &scheduleID, beforeData, afterData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Immunization schedule updated successfully"})
}

// DeleteAdminImmunizationSchedule deletes an immunization schedule
func DeleteAdminImmunizationSchedule(c echo.Context) error {
	scheduleID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(scheduleID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid schedule ID format"})
	}

	// Check if schedule is used in child_immunizations
	var immunizationCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM child_immunizations WHERE immunization_schedule_id = $1", scheduleID).Scan(&immunizationCount)
	if err != nil {
		c.Logger().Errorf("DeleteAdminImmunizationSchedule check immunizations error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if immunizationCount > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":            "Cannot delete immunization schedule: it is used in child immunizations",
			"immunization_count": strconv.Itoa(immunizationCount),
		})
	}

	// Get schedule data for audit log
	var schedule struct {
		Name     string
		Category string
		Priority string
	}

	err = db.DB.QueryRow(
		`SELECT name, category, priority FROM immunization_schedule WHERE id = $1`,
		scheduleID,
	).Scan(&schedule.Name, &schedule.Category, &schedule.Priority)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Immunization schedule not found"})
	}
	if err != nil {
		c.Logger().Errorf("DeleteAdminImmunizationSchedule get schedule error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Delete schedule
	_, err = db.DB.Exec("DELETE FROM immunization_schedule WHERE id = $1", scheduleID)
	if err != nil {
		c.Logger().Errorf("DeleteAdminImmunizationSchedule error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	scheduleData := map[string]interface{}{
		"name":     schedule.Name,
		"category": schedule.Category,
		"priority": schedule.Priority,
	}
	utils.LogAudit(adminUserID, "delete", "immunization_schedule", &scheduleID, scheduleData, nil, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Immunization schedule deleted successfully"})
}

