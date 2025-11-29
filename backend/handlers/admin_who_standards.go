package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetAdminWHOStandards returns all WHO standards with pagination and filters
func GetAdminWHOStandards(c echo.Context) error {
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

	indicator := c.QueryParam("indicator")
	gender := c.QueryParam("gender")

	// Build query
	query := `SELECT id, indicator, gender, age_months, height_cm, l_value, m_value, s_value,
	          sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3, created_at
	          FROM who_standards WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if indicator != "" {
		query += ` AND indicator = $` + strconv.Itoa(argIndex)
		args = append(args, indicator)
		argIndex++
	}

	if gender != "" {
		query += ` AND gender = $` + strconv.Itoa(argIndex)
		args = append(args, gender)
		argIndex++
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM who_standards WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if indicator != "" {
		countQuery += ` AND indicator = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, indicator)
		countArgIndex++
	}
	if gender != "" {
		countQuery += ` AND gender = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, gender)
		countArgIndex++
	}

	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminWHOStandards count error: %v", err)
		total = 0
	}

	query += ` ORDER BY indicator, gender, age_months, height_cm LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminWHOStandards query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	type WHOStandard struct {
		ID        string   `json:"id"`
		Indicator string   `json:"indicator"`
		Gender    string   `json:"gender"`
		AgeMonths *int     `json:"age_months"`
		HeightCm  *float64 `json:"height_cm"`
		LValue    float64  `json:"l_value"`
		MValue    float64  `json:"m_value"`
		SValue    float64  `json:"s_value"`
		SD3Neg    *float64 `json:"sd3neg"`
		SD2Neg    *float64 `json:"sd2neg"`
		SD1Neg    *float64 `json:"sd1neg"`
		SD0       *float64 `json:"sd0"`
		SD1       *float64 `json:"sd1"`
		SD2       *float64 `json:"sd2"`
		SD3       *float64 `json:"sd3"`
		CreatedAt string   `json:"created_at"`
	}

	var standards []WHOStandard
	for rows.Next() {
		var s WHOStandard
		var ageMonths sql.NullInt64
		var heightCm sql.NullFloat64
		var sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3 sql.NullFloat64
		var createdAt sql.NullTime

		err := rows.Scan(
			&s.ID, &s.Indicator, &s.Gender, &ageMonths, &heightCm,
			&s.LValue, &s.MValue, &s.SValue,
			&sd3neg, &sd2neg, &sd1neg, &sd0, &sd1, &sd2, &sd3, &createdAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan WHO standard row: %v", err)
			continue
		}

		if ageMonths.Valid {
			age := int(ageMonths.Int64)
			s.AgeMonths = &age
		}
		if heightCm.Valid {
			s.HeightCm = &heightCm.Float64
		}
		if sd3neg.Valid {
			s.SD3Neg = &sd3neg.Float64
		}
		if sd2neg.Valid {
			s.SD2Neg = &sd2neg.Float64
		}
		if sd1neg.Valid {
			s.SD1Neg = &sd1neg.Float64
		}
		if sd0.Valid {
			s.SD0 = &sd0.Float64
		}
		if sd1.Valid {
			s.SD1 = &sd1.Float64
		}
		if sd2.Valid {
			s.SD2 = &sd2.Float64
		}
		if sd3.Valid {
			s.SD3 = &sd3.Float64
		}
		if createdAt.Valid {
			s.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		standards = append(standards, s)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"standards": standards,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminWHOStandard returns a single WHO standard
func GetAdminWHOStandard(c echo.Context) error {
	standardID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(standardID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid standard ID format"})
	}

	type WHOStandard struct {
		ID        string   `json:"id"`
		Indicator string   `json:"indicator"`
		Gender    string   `json:"gender"`
		AgeMonths *int     `json:"age_months"`
		HeightCm  *float64 `json:"height_cm"`
		LValue    float64  `json:"l_value"`
		MValue    float64  `json:"m_value"`
		SValue    float64  `json:"s_value"`
		SD3Neg    *float64 `json:"sd3neg"`
		SD2Neg    *float64 `json:"sd2neg"`
		SD1Neg    *float64 `json:"sd1neg"`
		SD0       *float64 `json:"sd0"`
		SD1       *float64 `json:"sd1"`
		SD2       *float64 `json:"sd2"`
		SD3       *float64 `json:"sd3"`
		CreatedAt string   `json:"created_at"`
	}

	var s WHOStandard
	var ageMonths sql.NullInt64
	var heightCm sql.NullFloat64
	var sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3 sql.NullFloat64
	var createdAt sql.NullTime

	err := db.DB.QueryRow(
		`SELECT id, indicator, gender, age_months, height_cm, l_value, m_value, s_value,
		 sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3, created_at
		 FROM who_standards WHERE id = $1`,
		standardID,
	).Scan(
		&s.ID, &s.Indicator, &s.Gender, &ageMonths, &heightCm,
		&s.LValue, &s.MValue, &s.SValue,
		&sd3neg, &sd2neg, &sd1neg, &sd0, &sd1, &sd2, &sd3, &createdAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "WHO standard not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetAdminWHOStandard error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if ageMonths.Valid {
		age := int(ageMonths.Int64)
		s.AgeMonths = &age
	}
	if heightCm.Valid {
		s.HeightCm = &heightCm.Float64
	}
	if sd3neg.Valid {
		s.SD3Neg = &sd3neg.Float64
	}
	if sd2neg.Valid {
		s.SD2Neg = &sd2neg.Float64
	}
	if sd1neg.Valid {
		s.SD1Neg = &sd1neg.Float64
	}
	if sd0.Valid {
		s.SD0 = &sd0.Float64
	}
	if sd1.Valid {
		s.SD1 = &sd1.Float64
	}
	if sd2.Valid {
		s.SD2 = &sd2.Float64
	}
	if sd3.Valid {
		s.SD3 = &sd3.Float64
	}
	if createdAt.Valid {
		s.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(http.StatusOK, s)
}

// CreateAdminWHOStandard creates a new WHO standard
func CreateAdminWHOStandard(c echo.Context) error {
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		Indicator string   `json:"indicator" validate:"required"`
		Gender    string   `json:"gender" validate:"required"`
		AgeMonths *int     `json:"age_months"`
		HeightCm  *float64 `json:"height_cm"`
		LValue    float64  `json:"l_value" validate:"required"`
		MValue    float64  `json:"m_value" validate:"required"`
		SValue    float64  `json:"s_value" validate:"required"`
		SD3Neg    *float64 `json:"sd3neg"`
		SD2Neg    *float64 `json:"sd2neg"`
		SD1Neg    *float64 `json:"sd1neg"`
		SD0       *float64 `json:"sd0"`
		SD1       *float64 `json:"sd1"`
		SD2       *float64 `json:"sd2"`
		SD3       *float64 `json:"sd3"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate indicator
	validIndicators := map[string]bool{"wfa": true, "hfa": true, "wfh": true, "hcfa": true}
	if !validIndicators[req.Indicator] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid indicator"})
	}

	// Validate gender
	if req.Gender != "male" && req.Gender != "female" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid gender"})
	}

	var standardID string
	err := db.DB.QueryRow(
		`INSERT INTO who_standards (indicator, gender, age_months, height_cm, l_value, m_value, s_value,
		 sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING id`,
		req.Indicator, req.Gender, req.AgeMonths, req.HeightCm, req.LValue, req.MValue, req.SValue,
		req.SD3Neg, req.SD2Neg, req.SD1Neg, req.SD0, req.SD1, req.SD2, req.SD3,
	).Scan(&standardID)

	if err != nil {
		c.Logger().Errorf("CreateAdminWHOStandard error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create WHO standard"})
	}

	// Log audit
	standardData := map[string]interface{}{
		"indicator": req.Indicator,
		"gender":    req.Gender,
		"l_value":   req.LValue,
		"m_value":   req.MValue,
		"s_value":   req.SValue,
	}
	utils.LogAudit(adminUserID, "create", "who_standard", &standardID, nil, standardData, ipAddress, userAgent)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      standardID,
		"message": "WHO standard created successfully",
	})
}

// UpdateAdminWHOStandard updates a WHO standard
func UpdateAdminWHOStandard(c echo.Context) error {
	standardID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(standardID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid standard ID format"})
	}

	// Get existing standard for audit log
	var existing struct {
		Indicator string
		Gender    string
		LValue    float64
		MValue    float64
		SValue    float64
	}

	err := db.DB.QueryRow(
		`SELECT indicator, gender, l_value, m_value, s_value FROM who_standards WHERE id = $1`,
		standardID,
	).Scan(&existing.Indicator, &existing.Gender, &existing.LValue, &existing.MValue, &existing.SValue)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "WHO standard not found"})
	}
	if err != nil {
		c.Logger().Errorf("UpdateAdminWHOStandard get existing error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var req struct {
		Indicator *string  `json:"indicator"`
		Gender    *string  `json:"gender"`
		AgeMonths *int     `json:"age_months"`
		HeightCm  *float64 `json:"height_cm"`
		LValue    *float64 `json:"l_value"`
		MValue    *float64 `json:"m_value"`
		SValue    *float64 `json:"s_value"`
		SD3Neg    *float64 `json:"sd3neg"`
		SD2Neg    *float64 `json:"sd2neg"`
		SD1Neg    *float64 `json:"sd1neg"`
		SD0       *float64 `json:"sd0"`
		SD1       *float64 `json:"sd1"`
		SD2       *float64 `json:"sd2"`
		SD3       *float64 `json:"sd3"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Indicator != nil {
		validIndicators := map[string]bool{"wfa": true, "hfa": true, "wfh": true, "hcfa": true}
		if !validIndicators[*req.Indicator] {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid indicator"})
		}
		updateFields = append(updateFields, "indicator = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Indicator)
		argIndex++
	}
	if req.Gender != nil {
		if *req.Gender != "male" && *req.Gender != "female" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid gender"})
		}
		updateFields = append(updateFields, "gender = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Gender)
		argIndex++
	}
	if req.AgeMonths != nil {
		updateFields = append(updateFields, "age_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMonths)
		argIndex++
	}
	if req.HeightCm != nil {
		updateFields = append(updateFields, "height_cm = $"+strconv.Itoa(argIndex))
		args = append(args, *req.HeightCm)
		argIndex++
	}
	if req.LValue != nil {
		updateFields = append(updateFields, "l_value = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LValue)
		argIndex++
	}
	if req.MValue != nil {
		updateFields = append(updateFields, "m_value = $"+strconv.Itoa(argIndex))
		args = append(args, *req.MValue)
		argIndex++
	}
	if req.SValue != nil {
		updateFields = append(updateFields, "s_value = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SValue)
		argIndex++
	}
	if req.SD3Neg != nil {
		updateFields = append(updateFields, "sd3neg = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD3Neg)
		argIndex++
	}
	if req.SD2Neg != nil {
		updateFields = append(updateFields, "sd2neg = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD2Neg)
		argIndex++
	}
	if req.SD1Neg != nil {
		updateFields = append(updateFields, "sd1neg = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD1Neg)
		argIndex++
	}
	if req.SD0 != nil {
		updateFields = append(updateFields, "sd0 = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD0)
		argIndex++
	}
	if req.SD1 != nil {
		updateFields = append(updateFields, "sd1 = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD1)
		argIndex++
	}
	if req.SD2 != nil {
		updateFields = append(updateFields, "sd2 = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD2)
		argIndex++
	}
	if req.SD3 != nil {
		updateFields = append(updateFields, "sd3 = $"+strconv.Itoa(argIndex))
		args = append(args, *req.SD3)
		argIndex++
	}

	if len(updateFields) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	args = append(args, standardID)
	query := `UPDATE who_standards SET ` + updateFields[0]
	for i := 1; i < len(updateFields); i++ {
		query += `, ` + updateFields[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		c.Logger().Errorf("UpdateAdminWHOStandard error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	beforeData := map[string]interface{}{
		"indicator": existing.Indicator,
		"gender":    existing.Gender,
		"l_value":   existing.LValue,
		"m_value":   existing.MValue,
		"s_value":   existing.SValue,
	}
	afterData := map[string]interface{}{}
	if req.Indicator != nil {
		afterData["indicator"] = *req.Indicator
	}
	if req.Gender != nil {
		afterData["gender"] = *req.Gender
	}
	if req.LValue != nil {
		afterData["l_value"] = *req.LValue
	}
	if req.MValue != nil {
		afterData["m_value"] = *req.MValue
	}
	if req.SValue != nil {
		afterData["s_value"] = *req.SValue
	}
	utils.LogAudit(adminUserID, "update", "who_standard", &standardID, beforeData, afterData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "WHO standard updated successfully"})
}

// DeleteAdminWHOStandard deletes a WHO standard
func DeleteAdminWHOStandard(c echo.Context) error {
	standardID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(standardID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid standard ID format"})
	}

	// Get standard data for audit log
	var standard struct {
		Indicator string
		Gender    string
	}

	err := db.DB.QueryRow(
		`SELECT indicator, gender FROM who_standards WHERE id = $1`,
		standardID,
	).Scan(&standard.Indicator, &standard.Gender)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "WHO standard not found"})
	}
	if err != nil {
		c.Logger().Errorf("DeleteAdminWHOStandard get standard error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Delete standard
	_, err = db.DB.Exec("DELETE FROM who_standards WHERE id = $1", standardID)
	if err != nil {
		c.Logger().Errorf("DeleteAdminWHOStandard error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	standardData := map[string]interface{}{
		"indicator": standard.Indicator,
		"gender":    standard.Gender,
	}
	utils.LogAudit(adminUserID, "delete", "who_standard", &standardID, standardData, nil, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "WHO standard deleted successfully"})
}

