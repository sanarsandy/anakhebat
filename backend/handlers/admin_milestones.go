package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetAdminMilestones returns all milestones with pagination and filters
func GetAdminMilestones(c echo.Context) error {
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

	source := c.QueryParam("source")
	category := c.QueryParam("category")
	redFlag := c.QueryParam("red_flag")
	search := c.QueryParam("search")

	// Build query
	query := `SELECT id, age_months, min_age_range, max_age_range, category, question, question_en, 
	          source, is_red_flag, pyramid_level, denver_domain, created_at
	          FROM milestones WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if source != "" {
		query += ` AND source = $` + strconv.Itoa(argIndex)
		args = append(args, source)
		argIndex++
	}

	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if redFlag == "true" {
		query += ` AND is_red_flag = true`
	} else if redFlag == "false" {
		query += ` AND is_red_flag = false`
	}

	if search != "" {
		query += ` AND (question ILIKE $` + strconv.Itoa(argIndex) +
			` OR question_en ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM milestones WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if source != "" {
		countQuery += ` AND source = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, source)
		countArgIndex++
	}
	if category != "" {
		countQuery += ` AND category = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, category)
		countArgIndex++
	}
	if redFlag == "true" {
		countQuery += ` AND is_red_flag = true`
	} else if redFlag == "false" {
		countQuery += ` AND is_red_flag = false`
	}
	if search != "" {
		countQuery += ` AND (question ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR question_en ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
		countArgIndex += 2
	}

	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminMilestones count error: %v", err)
		total = 0
	}

	query += ` ORDER BY age_months ASC, pyramid_level ASC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminMilestones query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var milestones []models.Milestone
	for rows.Next() {
		var m models.Milestone
		var minAgeRange sql.NullInt64
		var maxAgeRange sql.NullInt64
		var questionEn sql.NullString
		var denverDomain sql.NullString

		err := rows.Scan(
			&m.ID, &m.AgeMonths, &minAgeRange, &maxAgeRange, &m.Category, &m.Question,
			&questionEn, &m.Source, &m.IsRedFlag, &m.PyramidLevel, &denverDomain, &m.CreatedAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan milestone row: %v", err)
			continue
		}

		if minAgeRange.Valid {
			age := int(minAgeRange.Int64)
			m.MinAgeRange = &age
		}
		if maxAgeRange.Valid {
			age := int(maxAgeRange.Int64)
			m.MaxAgeRange = &age
		}
		if questionEn.Valid {
			m.QuestionEn = questionEn.String
		}
		if denverDomain.Valid {
			m.DenverDomain = &denverDomain.String
		}

		milestones = append(milestones, m)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"milestones": milestones,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminMilestone returns a single milestone
func GetAdminMilestone(c echo.Context) error {
	milestoneID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(milestoneID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID format"})
	}

	var m models.Milestone
	var minAgeRange sql.NullInt64
	var maxAgeRange sql.NullInt64
	var questionEn sql.NullString
	var denverDomain sql.NullString

	err := db.DB.QueryRow(
		`SELECT id, age_months, min_age_range, max_age_range, category, question, question_en,
		 source, is_red_flag, pyramid_level, denver_domain, created_at
		 FROM milestones WHERE id = $1`,
		milestoneID,
	).Scan(
		&m.ID, &m.AgeMonths, &minAgeRange, &maxAgeRange, &m.Category, &m.Question,
		&questionEn, &m.Source, &m.IsRedFlag, &m.PyramidLevel, &denverDomain, &m.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Milestone not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetAdminMilestone error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if minAgeRange.Valid {
		age := int(minAgeRange.Int64)
		m.MinAgeRange = &age
	}
	if maxAgeRange.Valid {
		age := int(maxAgeRange.Int64)
		m.MaxAgeRange = &age
	}
	if questionEn.Valid {
		m.QuestionEn = questionEn.String
	}
	if denverDomain.Valid {
		m.DenverDomain = &denverDomain.String
	}

	return c.JSON(http.StatusOK, m)
}

// CreateAdminMilestone creates a new milestone
func CreateAdminMilestone(c echo.Context) error {
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		AgeMonths    int     `json:"age_months" validate:"required"`
		MinAgeRange  *int    `json:"min_age_range"`
		MaxAgeRange  *int    `json:"max_age_range"`
		Category     string  `json:"category" validate:"required"`
		Question     string  `json:"question" validate:"required"`
		QuestionEn   string  `json:"question_en"`
		Source       string  `json:"source" validate:"required"`
		IsRedFlag    bool    `json:"is_red_flag"`
		PyramidLevel int     `json:"pyramid_level" validate:"required,min=1,max=4"`
		DenverDomain *string `json:"denver_domain"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate age_months range
	if err := utils.ValidateIntRange(req.AgeMonths, 0, 60, "age_months"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validate pyramid_level range
	if err := utils.ValidateIntRange(req.PyramidLevel, 1, 4, "pyramid_level"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validate question length
	if err := utils.ValidateStringLength(req.Question, 1, 2000, "question"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validate question_en length if provided
	if req.QuestionEn != "" {
		if err := utils.ValidateStringLength(req.QuestionEn, 1, 2000, "question_en"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	// Validate min_age_range if provided
	if req.MinAgeRange != nil {
		if err := utils.ValidateIntRange(*req.MinAgeRange, 0, 60, "min_age_range"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	// Validate max_age_range if provided
	if req.MaxAgeRange != nil {
		if err := utils.ValidateIntRange(*req.MaxAgeRange, 0, 60, "max_age_range"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// Ensure max >= min if both provided
		if req.MinAgeRange != nil && *req.MaxAgeRange < *req.MinAgeRange {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "max_age_range must be greater than or equal to min_age_range"})
		}
	}

	// Validate category
	validCategories := map[string]bool{"sensory": true, "motor": true, "perception": true, "cognitive": true}
	if !validCategories[req.Category] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category"})
	}

	// Validate source
	validSources := map[string]bool{"KPSP": true, "CDC": true, "DENVER": true}
	if !validSources[req.Source] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid source"})
	}

	var milestoneID string
	err := db.DB.QueryRow(
		`INSERT INTO milestones (age_months, min_age_range, max_age_range, category, question, question_en,
		 source, is_red_flag, pyramid_level, denver_domain)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id`,
		req.AgeMonths, req.MinAgeRange, req.MaxAgeRange, req.Category, req.Question,
		req.QuestionEn, req.Source, req.IsRedFlag, req.PyramidLevel, req.DenverDomain,
	).Scan(&milestoneID)

	if err != nil {
		c.Logger().Errorf("CreateAdminMilestone error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	milestoneData := map[string]interface{}{
		"age_months":    req.AgeMonths,
		"category":      req.Category,
		"question":      req.Question,
		"source":        req.Source,
		"pyramid_level": req.PyramidLevel,
	}
	utils.LogAudit(adminUserID, "create", "milestone", &milestoneID, nil, milestoneData, ipAddress, userAgent)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id": milestoneID,
		"message": "Milestone created successfully",
	})
}

// UpdateAdminMilestone updates a milestone
func UpdateAdminMilestone(c echo.Context) error {
	milestoneID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(milestoneID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID format"})
	}

	// Get existing milestone for audit log
	var existing models.Milestone
	var minAgeRange sql.NullInt64
	var maxAgeRange sql.NullInt64
	var questionEn sql.NullString
	var denverDomain sql.NullString

	err := db.DB.QueryRow(
		`SELECT id, age_months, min_age_range, max_age_range, category, question, question_en,
		 source, is_red_flag, pyramid_level, denver_domain, created_at
		 FROM milestones WHERE id = $1`,
		milestoneID,
	).Scan(
		&existing.ID, &existing.AgeMonths, &minAgeRange, &maxAgeRange, &existing.Category,
		&existing.Question, &questionEn, &existing.Source, &existing.IsRedFlag,
		&existing.PyramidLevel, &denverDomain, &existing.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Milestone not found"})
	}
	if err != nil {
		c.Logger().Errorf("UpdateAdminMilestone get existing error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if minAgeRange.Valid {
		age := int(minAgeRange.Int64)
		existing.MinAgeRange = &age
	}
	if maxAgeRange.Valid {
		age := int(maxAgeRange.Int64)
		existing.MaxAgeRange = &age
	}
	if questionEn.Valid {
		existing.QuestionEn = questionEn.String
	}
	if denverDomain.Valid {
		existing.DenverDomain = &denverDomain.String
	}

	var req struct {
		AgeMonths    *int    `json:"age_months"`
		MinAgeRange  *int    `json:"min_age_range"`
		MaxAgeRange  *int    `json:"max_age_range"`
		Category     *string `json:"category"`
		Question     *string `json:"question"`
		QuestionEn   *string `json:"question_en"`
		Source       *string `json:"source"`
		IsRedFlag    *bool   `json:"is_red_flag"`
		PyramidLevel *int    `json:"pyramid_level"`
		DenverDomain *string `json:"denver_domain"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.AgeMonths != nil {
		if err := utils.ValidateIntRange(*req.AgeMonths, 0, 60, "age_months"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		updateFields = append(updateFields, "age_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMonths)
		argIndex++
	}
	if req.MinAgeRange != nil {
		updateFields = append(updateFields, "min_age_range = $"+strconv.Itoa(argIndex))
		args = append(args, *req.MinAgeRange)
		argIndex++
	}
	if req.MaxAgeRange != nil {
		updateFields = append(updateFields, "max_age_range = $"+strconv.Itoa(argIndex))
		args = append(args, *req.MaxAgeRange)
		argIndex++
	}
	if req.Category != nil {
		validCategories := map[string]bool{"sensory": true, "motor": true, "perception": true, "cognitive": true}
		if !validCategories[*req.Category] {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category"})
		}
		updateFields = append(updateFields, "category = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Category)
		argIndex++
	}
	if req.Question != nil {
		if err := utils.ValidateStringLength(*req.Question, 1, 2000, "question"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		updateFields = append(updateFields, "question = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Question)
		argIndex++
	}
	if req.QuestionEn != nil {
		updateFields = append(updateFields, "question_en = $"+strconv.Itoa(argIndex))
		args = append(args, *req.QuestionEn)
		argIndex++
	}
	if req.Source != nil {
		validSources := map[string]bool{"KPSP": true, "CDC": true, "DENVER": true}
		if !validSources[*req.Source] {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid source"})
		}
		updateFields = append(updateFields, "source = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Source)
		argIndex++
	}
	if req.IsRedFlag != nil {
		updateFields = append(updateFields, "is_red_flag = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsRedFlag)
		argIndex++
	}
	if req.PyramidLevel != nil {
		if err := utils.ValidateIntRange(*req.PyramidLevel, 1, 4, "pyramid_level"); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		updateFields = append(updateFields, "pyramid_level = $"+strconv.Itoa(argIndex))
		args = append(args, *req.PyramidLevel)
		argIndex++
	}
	if req.DenverDomain != nil {
		updateFields = append(updateFields, "denver_domain = $"+strconv.Itoa(argIndex))
		args = append(args, *req.DenverDomain)
		argIndex++
	}

	if len(updateFields) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	args = append(args, milestoneID)
	query := `UPDATE milestones SET ` + updateFields[0]
	for i := 1; i < len(updateFields); i++ {
		query += `, ` + updateFields[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		c.Logger().Errorf("UpdateAdminMilestone error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	beforeData := map[string]interface{}{
		"age_months":    existing.AgeMonths,
		"category":      existing.Category,
		"question":      existing.Question,
		"source":        existing.Source,
		"pyramid_level": existing.PyramidLevel,
	}
	afterData := map[string]interface{}{}
	if req.AgeMonths != nil {
		afterData["age_months"] = *req.AgeMonths
	}
	if req.Category != nil {
		afterData["category"] = *req.Category
	}
	if req.Question != nil {
		afterData["question"] = *req.Question
	}
	if req.Source != nil {
		afterData["source"] = *req.Source
	}
	if req.PyramidLevel != nil {
		afterData["pyramid_level"] = *req.PyramidLevel
	}
	utils.LogAudit(adminUserID, "update", "milestone", &milestoneID, beforeData, afterData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Milestone updated successfully"})
}

// DeleteAdminMilestone deletes a milestone
func DeleteAdminMilestone(c echo.Context) error {
	milestoneID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(milestoneID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID format"})
	}

	// Check if milestone is used in assessments
	var assessmentCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM assessments WHERE milestone_id = $1", milestoneID).Scan(&assessmentCount)
	if err != nil {
		c.Logger().Errorf("DeleteAdminMilestone check assessments error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if assessmentCount > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Cannot delete milestone: it is used in assessments",
			"assessment_count": strconv.Itoa(assessmentCount),
		})
	}

	// Get milestone data for audit log
	var milestone models.Milestone
	var minAgeRange sql.NullInt64
	var maxAgeRange sql.NullInt64
	var questionEn sql.NullString
	var denverDomain sql.NullString

	err = db.DB.QueryRow(
		`SELECT id, age_months, min_age_range, max_age_range, category, question, question_en,
		 source, is_red_flag, pyramid_level, denver_domain, created_at
		 FROM milestones WHERE id = $1`,
		milestoneID,
	).Scan(
		&milestone.ID, &milestone.AgeMonths, &minAgeRange, &maxAgeRange, &milestone.Category,
		&milestone.Question, &questionEn, &milestone.Source, &milestone.IsRedFlag,
		&milestone.PyramidLevel, &denverDomain, &milestone.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Milestone not found"})
	}

	// Delete milestone
	_, err = db.DB.Exec("DELETE FROM milestones WHERE id = $1", milestoneID)
	if err != nil {
		c.Logger().Errorf("DeleteAdminMilestone error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	milestoneData := map[string]interface{}{
		"age_months":    milestone.AgeMonths,
		"category":      milestone.Category,
		"question":      milestone.Question,
		"source":        milestone.Source,
		"pyramid_level": milestone.PyramidLevel,
	}
	utils.LogAudit(adminUserID, "delete", "milestone", &milestoneID, milestoneData, nil, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Milestone deleted successfully"})
}

