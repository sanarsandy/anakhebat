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

// CreateMeasurement creates a new measurement for a child
func CreateMeasurement(c echo.Context) error {
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
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get child data for calculations (including premature info)
	var child models.Child
	err = db.DB.QueryRow("SELECT id, dob, gender, is_premature, gestational_age FROM children WHERE id = $1", childID).
		Scan(&child.ID, &child.DOB, &child.Gender, &child.IsPremature, &child.GestationalAge)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get child data"})
	}

	// Parse request
	req := new(models.CreateMeasurementRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Error("Bind error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format: " + err.Error()})
	}

	c.Logger().Info("Received measurement request: ", req)

	// Validate required fields
	if req.MeasurementDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "measurement_date is required"})
	}
	if req.Weight <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "weight must be greater than 0"})
	}
	if req.Height <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "height must be greater than 0"})
	}

	// Calculate age (using corrected age if premature and < 24 months)
	ageInDays, ageInMonths, useCorrected, err := utils.CalculateCorrectedAge(
		child.DOB, req.MeasurementDate, child.IsPremature, child.GestationalAge)
	if err != nil {
		c.Logger().Error("Age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format. Expected YYYY-MM-DD"})
	}

	// Also calculate chronological age for storage/display
	chronoDays, err := utils.CalculateAgeInDays(child.DOB, req.MeasurementDate)
	if err != nil {
		c.Logger().Error("Chronological age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format. Expected YYYY-MM-DD"})
	}
	chronoMonths, err := utils.CalculateAgeInMonths(child.DOB, req.MeasurementDate)
	if err != nil {
		c.Logger().Error("Chronological age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format. Expected YYYY-MM-DD"})
	}

	if useCorrected {
		c.Logger().Infof("Using corrected age for premature child - Chronological: %d days (%d months), Corrected: %d days (%d months)", 
			chronoDays, chronoMonths, ageInDays, ageInMonths)
	} else {
		c.Logger().Info("Using chronological age - Days: ", ageInDays, ", Months: ", ageInMonths)
	}

	// Calculate Z-scores using WHO standards (using corrected age if applicable)
	headCirc := 0.0
	if req.HeadCircumference != nil {
		headCirc = *req.HeadCircumference
	}
	zscores, err := utils.CalculateAllZScores(db.DB, child.Gender, ageInMonths, req.Weight, req.Height, headCirc)
	if err != nil {
		c.Logger().Warn("Failed to calculate Z-scores: ", err)
	}

	// Interpret nutritional status
	var nutritionalStatus, heightStatus, wfhStatus string
	if zscores != nil && (zscores.HasWeightForAge || zscores.HasHeightForAge) {
		// Use 0 as default if Z-score not available
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
		
		c.Logger().Infof("Z-scores calculated - WFA: %.2f (has: %v), HFA: %.2f (has: %v), WFH: %.2f (has: %v)", 
			zscores.WeightForAge, zscores.HasWeightForAge,
			zscores.HeightForAge, zscores.HasHeightForAge,
			zscores.WeightForHeight, zscores.HasWeightForHeight)
		c.Logger().Infof("Status - Weight: %s, Height: %s, WFH: %s", 
			nutritionalStatus, heightStatus, wfhStatus)
	}

	// Insert measurement (store chronological age, but use corrected age for Z-scores)
	query := `INSERT INTO measurements 
		(child_id, measurement_date, weight, height, head_circumference, age_in_days, age_in_months, 
		weight_for_age_zscore, height_for_age_zscore, weight_for_height_zscore, head_circumference_zscore,
		nutritional_status, height_status, weight_for_height_status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
		RETURNING id, created_at`
	
	// Store chronological age in database, but use corrected age for Z-score calculation (already done above)
	// So we use chronoDays and chronoMonths for storage
	storageDays := chronoDays
	storageMonths := chronoMonths

	c.Logger().Infof("Attempting to insert measurement for child %s: Date=%s, Weight=%.2f, Height=%.2f", 
		childID, req.MeasurementDate, req.Weight, req.Height)

	var measurement models.Measurement
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

	err = db.DB.QueryRow(query,
		childID, req.MeasurementDate, req.Weight, req.Height, req.HeadCircumference,
		storageDays, storageMonths, wfaZPtr, hfaZPtr, wfhZPtr, hcZPtr,
		nutritionalStatus, heightStatus, wfhStatus).
		Scan(&measurement.ID, &measurement.CreatedAt)

	if err != nil {
		c.Logger().Errorf("Database insert error: %v", err)
		c.Logger().Errorf("Query: %s", query)
		c.Logger().Errorf("Params: childID=%s, date=%s, weight=%.2f, height=%.2f", 
			childID, req.MeasurementDate, req.Weight, req.Height)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal menyimpan pengukuran. Pastikan data yang diinput valid.",
			"details": err.Error(),
		})
	}

	c.Logger().Info("Measurement created successfully with ID: ", measurement.ID)

	// Build response (use chronological age for display, but indicate if corrected age was used)
	response := models.MeasurementResponse{
		ID:                      measurement.ID,
		ChildID:                 childID,
		MeasurementDate:         req.MeasurementDate,
		Weight:                  req.Weight,
		Height:                  req.Height,
		HeadCircumference:       req.HeadCircumference,
		AgeInDays:               storageDays, // Chronological age for display
		AgeInMonths:             storageMonths,
		AgeDisplay:              utils.FormatAgeDisplay(storageMonths),
		WeightForAgeZScore:      wfaZPtr,
		HeightForAgeZScore:      hfaZPtr,
		WeightForHeightZScore:   wfhZPtr,
		HeadCircumferenceZScore: hcZPtr,
		NutritionalStatus:       nutritionalStatus,
		HeightStatus:            heightStatus,
		WeightForHeightStatus:   wfhStatus,
		CreatedAt:               measurement.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	
	// Add corrected age info if applicable (will be added to response model if needed)
	if useCorrected {
		c.Logger().Info("Z-scores calculated using corrected age for premature child")
	}

	return c.JSON(http.StatusCreated, response)
}

// GetMeasurements retrieves all measurements for a child
func GetMeasurements(c echo.Context) error {
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
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get measurements
	query := `SELECT id, child_id, measurement_date, weight, height, head_circumference, 
		age_in_days, age_in_months, weight_for_age_zscore, height_for_age_zscore, 
		weight_for_height_zscore, head_circumference_zscore,
		nutritional_status, height_status, weight_for_height_status, created_at 
		FROM measurements WHERE child_id = $1 ORDER BY measurement_date DESC`

	rows, err := db.DB.Query(query, childID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	measurements := []models.MeasurementResponse{}
	for rows.Next() {
		var m models.Measurement
		var nutritionalStatus, heightStatus, wfhStatus sql.NullString
		var wfhZScore, hcZScore sql.NullFloat64
		
		err := rows.Scan(&m.ID, &m.ChildID, &m.MeasurementDate, &m.Weight, &m.Height, &m.HeadCircumference,
			&m.AgeInDays, &m.AgeInMonths, &m.WeightForAgeZScore, &m.HeightForAgeZScore,
			&wfhZScore, &hcZScore,
			&nutritionalStatus, &heightStatus, &wfhStatus, &m.CreatedAt)
		if err != nil {
			continue
		}

		var wfhZPtr, hcZPtr *float64
		if wfhZScore.Valid {
			wfhZPtr = &wfhZScore.Float64
		}
		if hcZScore.Valid {
			hcZPtr = &hcZScore.Float64
		}

		response := models.MeasurementResponse{
			ID:                      m.ID,
			ChildID:                 m.ChildID,
			MeasurementDate:         m.MeasurementDate,
			Weight:                  m.Weight,
			Height:                  m.Height,
			HeadCircumference:       m.HeadCircumference,
			AgeInDays:               m.AgeInDays,
			AgeInMonths:             m.AgeInMonths,
			AgeDisplay:              utils.FormatAgeDisplay(m.AgeInMonths),
			WeightForAgeZScore:      m.WeightForAgeZScore,
			HeightForAgeZScore:      m.HeightForAgeZScore,
			WeightForHeightZScore:   wfhZPtr,
			HeadCircumferenceZScore: hcZPtr,
			NutritionalStatus:       nutritionalStatus.String,
			HeightStatus:            heightStatus.String,
			WeightForHeightStatus:   wfhStatus.String,
			CreatedAt:               m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		measurements = append(measurements, response)
	}

	return c.JSON(http.StatusOK, measurements)
}

// GetLatestMeasurement retrieves the most recent measurement for a child
func GetLatestMeasurement(c echo.Context) error {
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
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get latest measurement
	var m models.Measurement
	var nutritionalStatus, heightStatus, wfhStatus sql.NullString
	var wfhZScore, hcZScore sql.NullFloat64
	
	query := `SELECT id, child_id, measurement_date, weight, height, head_circumference, 
		age_in_days, age_in_months, weight_for_age_zscore, height_for_age_zscore, 
		weight_for_height_zscore, head_circumference_zscore,
		nutritional_status, height_status, weight_for_height_status, created_at 
		FROM measurements WHERE child_id = $1 ORDER BY measurement_date DESC LIMIT 1`

	err = db.DB.QueryRow(query, childID).Scan(
		&m.ID, &m.ChildID, &m.MeasurementDate, &m.Weight, &m.Height, &m.HeadCircumference,
		&m.AgeInDays, &m.AgeInMonths, &m.WeightForAgeZScore, &m.HeightForAgeZScore,
		&wfhZScore, &hcZScore,
		&nutritionalStatus, &heightStatus, &wfhStatus, &m.CreatedAt)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No measurements found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var wfhZPtr, hcZPtr *float64
	if wfhZScore.Valid {
		wfhZPtr = &wfhZScore.Float64
	}
	if hcZScore.Valid {
		hcZPtr = &hcZScore.Float64
	}

	response := models.MeasurementResponse{
		ID:                      m.ID,
		ChildID:                 m.ChildID,
		MeasurementDate:         m.MeasurementDate,
		Weight:                  m.Weight,
		Height:                  m.Height,
		HeadCircumference:       m.HeadCircumference,
		AgeInDays:               m.AgeInDays,
		AgeInMonths:             m.AgeInMonths,
		AgeDisplay:              utils.FormatAgeDisplay(m.AgeInMonths),
		WeightForAgeZScore:      m.WeightForAgeZScore,
		HeightForAgeZScore:      m.HeightForAgeZScore,
		WeightForHeightZScore:   wfhZPtr,
		HeadCircumferenceZScore: hcZPtr,
		NutritionalStatus:       nutritionalStatus.String,
		HeightStatus:            heightStatus.String,
		WeightForHeightStatus:   wfhStatus.String,
		CreatedAt:               m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteMeasurement deletes a measurement
func DeleteMeasurement(c echo.Context) error {
	childID := c.Param("id")
	measurementID := c.Param("measurementId")

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
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Delete measurement
	query := `DELETE FROM measurements WHERE id = $1 AND child_id = $2`
	result, err := db.DB.Exec(query, measurementID, childID)
	if err != nil {
		c.Logger().Errorf("Failed to delete measurement: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus pengukuran"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.Logger().Errorf("Failed to get rows affected: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus pengukuran"})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Data pengukuran tidak ditemukan"})
	}

	c.Logger().Infof("Successfully deleted measurement %s for child %s", measurementID, childID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Pengukuran berhasil dihapus"})
}

// UpdateMeasurement updates an existing measurement
func UpdateMeasurement(c echo.Context) error {
	childID := c.Param("id")
	measurementID := c.Param("measurementId")

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
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get existing measurement to verify it belongs to child
	var existingChildID string
	err = db.DB.QueryRow("SELECT child_id FROM measurements WHERE id = $1", measurementID).Scan(&existingChildID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Measurement not found"})
	}
	if existingChildID != childID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Measurement does not belong to this child"})
	}

	// Get child data for calculations (including premature info)
	var child models.Child
	err = db.DB.QueryRow("SELECT id, dob, gender, is_premature, gestational_age FROM children WHERE id = $1", childID).
		Scan(&child.ID, &child.DOB, &child.Gender, &child.IsPremature, &child.GestationalAge)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get child data"})
	}

	// Parse request
	req := new(models.CreateMeasurementRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Validate required fields
	if req.MeasurementDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "measurement_date is required"})
	}
	if req.Weight <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "weight must be greater than 0"})
	}
	if req.Height <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "height must be greater than 0"})
	}

	// Calculate age (using corrected age if premature and < 24 months)
	ageInDays, ageInMonths, useCorrected, err := utils.CalculateCorrectedAge(
		child.DOB, req.MeasurementDate, child.IsPremature, child.GestationalAge)
	if err != nil {
		c.Logger().Error("Age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format"})
	}

	// Also calculate chronological age for storage/display
	chronoDays, err := utils.CalculateAgeInDays(child.DOB, req.MeasurementDate)
	if err != nil {
		c.Logger().Error("Chronological age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format"})
	}
	chronoMonths, err := utils.CalculateAgeInMonths(child.DOB, req.MeasurementDate)
	if err != nil {
		c.Logger().Error("Chronological age calculation error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid measurement date format"})
	}

	if useCorrected {
		c.Logger().Infof("Using corrected age for premature child - Chronological: %d days (%d months), Corrected: %d days (%d months)", 
			chronoDays, chronoMonths, ageInDays, ageInMonths)
	}

	// Calculate Z-scores using WHO standards (using corrected age if applicable)
	headCirc := 0.0
	if req.HeadCircumference != nil {
		headCirc = *req.HeadCircumference
	}
	zscores, err := utils.CalculateAllZScores(db.DB, child.Gender, ageInMonths, req.Weight, req.Height, headCirc)
	if err != nil {
		c.Logger().Warn("Failed to calculate Z-scores: ", err)
	}

	// Interpret nutritional status
	var nutritionalStatus, heightStatus, wfhStatus string
	if zscores != nil && (zscores.HasWeightForAge || zscores.HasHeightForAge) {
		// Use 0 as default if Z-score not available
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

	// Update measurement (store chronological age, but use corrected age for Z-scores)
	query := `UPDATE measurements SET 
		measurement_date = $1, weight = $2, height = $3, head_circumference = $4, 
		age_in_days = $5, age_in_months = $6, 
		weight_for_age_zscore = $7, height_for_age_zscore = $8, 
		weight_for_height_zscore = $9, head_circumference_zscore = $10,
		nutritional_status = $11, height_status = $12, weight_for_height_status = $13
		WHERE id = $14 AND child_id = $15
		RETURNING created_at`

	// Store chronological age in database, but use corrected age for Z-score calculation (already done above)
	storageDays := chronoDays
	storageMonths := chronoMonths

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

	var createdAt time.Time
	err = db.DB.QueryRow(query,
		req.MeasurementDate, req.Weight, req.Height, req.HeadCircumference,
		storageDays, storageMonths, wfaZPtr, hfaZPtr, wfhZPtr, hcZPtr,
		nutritionalStatus, heightStatus, wfhStatus,
		measurementID, childID).Scan(&createdAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update measurement: " + err.Error()})
	}

	// Build response (use chronological age for display)
	response := models.MeasurementResponse{
		ID:                      measurementID,
		ChildID:                 childID,
		MeasurementDate:         req.MeasurementDate,
		Weight:                  req.Weight,
		Height:                  req.Height,
		HeadCircumference:       req.HeadCircumference,
		AgeInDays:               storageDays, // Chronological age for display
		AgeInMonths:             storageMonths,
		AgeDisplay:              utils.FormatAgeDisplay(storageMonths),
		WeightForAgeZScore:      wfaZPtr,
		HeightForAgeZScore:      hfaZPtr,
		WeightForHeightZScore:   wfhZPtr,
		HeadCircumferenceZScore: hcZPtr,
		NutritionalStatus:       nutritionalStatus,
		HeightStatus:            heightStatus,
		WeightForHeightStatus:   wfhStatus,
		CreatedAt:               createdAt.Format("2006-01-02T15:04:05Z"),
	}

	if useCorrected {
		c.Logger().Info("Z-scores calculated using corrected age for premature child")
	}

	return c.JSON(http.StatusOK, response)
}
