package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
)

// GetAdminStimulationContent returns all stimulation content with pagination and filters
func GetAdminStimulationContent(c echo.Context) error {
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
	contentType := c.QueryParam("content_type")
	isActive := c.QueryParam("is_active")
	search := c.QueryParam("search")

	// Build query
	query := `SELECT id, milestone_id, category, title, description, content_type, url, thumbnail_url,
	          age_min_months, age_max_months, is_active, created_at, updated_at
	          FROM stimulation_content WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if contentType != "" {
		query += ` AND content_type = $` + strconv.Itoa(argIndex)
		args = append(args, contentType)
		argIndex++
	}

	if isActive == "true" {
		query += ` AND is_active = true`
	} else if isActive == "false" {
		query += ` AND is_active = false`
	}

	if search != "" {
		query += ` AND (title ILIKE $` + strconv.Itoa(argIndex) +
			` OR description ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM stimulation_content WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if category != "" {
		countQuery += ` AND category = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, category)
		countArgIndex++
	}
	if contentType != "" {
		countQuery += ` AND content_type = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, contentType)
		countArgIndex++
	}
	if isActive == "true" {
		countQuery += ` AND is_active = true`
	} else if isActive == "false" {
		countQuery += ` AND is_active = false`
	}
	if search != "" {
		countQuery += ` AND (title ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR description ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
		countArgIndex += 2
	}

	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminStimulationContent count error: %v", err)
		total = 0
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminStimulationContent query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	type StimulationContent struct {
		ID           string  `json:"id"`
		MilestoneID  *string `json:"milestone_id"`
		Category     string  `json:"category"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		ContentType  string  `json:"content_type"`
		URL          string  `json:"url"`
		ThumbnailURL string  `json:"thumbnail_url"`
		AgeMinMonths *int    `json:"age_min_months"`
		AgeMaxMonths *int    `json:"age_max_months"`
		IsActive     bool    `json:"is_active"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	}

	var contents []StimulationContent
	for rows.Next() {
		var sc StimulationContent
		var milestoneID sql.NullString
		var description sql.NullString
		var thumbnailURL sql.NullString
		var ageMinMonths sql.NullInt64
		var ageMaxMonths sql.NullInt64
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		err := rows.Scan(
			&sc.ID, &milestoneID, &sc.Category, &sc.Title, &description, &sc.ContentType,
			&sc.URL, &thumbnailURL, &ageMinMonths, &ageMaxMonths, &sc.IsActive, &createdAt, &updatedAt,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan stimulation content row: %v", err)
			continue
		}

		if milestoneID.Valid {
			sc.MilestoneID = &milestoneID.String
		}
		if description.Valid {
			sc.Description = description.String
		}
		if thumbnailURL.Valid {
			sc.ThumbnailURL = thumbnailURL.String
		}
		if ageMinMonths.Valid {
			age := int(ageMinMonths.Int64)
			sc.AgeMinMonths = &age
		}
		if ageMaxMonths.Valid {
			age := int(ageMaxMonths.Int64)
			sc.AgeMaxMonths = &age
		}
		if createdAt.Valid {
			sc.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		if updatedAt.Valid {
			sc.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		contents = append(contents, sc)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"contents": contents,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminStimulationContentItem returns a single stimulation content
func GetAdminStimulationContentItem(c echo.Context) error {
	contentID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(contentID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid content ID format"})
	}

	type StimulationContent struct {
		ID           string  `json:"id"`
		MilestoneID  *string `json:"milestone_id"`
		Category     string  `json:"category"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		ContentType  string  `json:"content_type"`
		URL          string  `json:"url"`
		ThumbnailURL string  `json:"thumbnail_url"`
		AgeMinMonths *int    `json:"age_min_months"`
		AgeMaxMonths *int    `json:"age_max_months"`
		IsActive     bool    `json:"is_active"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	}

	var sc StimulationContent
	var milestoneID sql.NullString
	var description sql.NullString
	var thumbnailURL sql.NullString
	var ageMinMonths sql.NullInt64
	var ageMaxMonths sql.NullInt64
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	err := db.DB.QueryRow(
		`SELECT id, milestone_id, category, title, description, content_type, url, thumbnail_url,
		 age_min_months, age_max_months, is_active, created_at, updated_at
		 FROM stimulation_content WHERE id = $1`,
		contentID,
	).Scan(
		&sc.ID, &milestoneID, &sc.Category, &sc.Title, &description, &sc.ContentType,
		&sc.URL, &thumbnailURL, &ageMinMonths, &ageMaxMonths, &sc.IsActive, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stimulation content not found"})
	}
	if err != nil {
		c.Logger().Errorf("GetAdminStimulationContentItem error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if milestoneID.Valid {
		sc.MilestoneID = &milestoneID.String
	}
	if description.Valid {
		sc.Description = description.String
	}
	if thumbnailURL.Valid {
		sc.ThumbnailURL = thumbnailURL.String
	}
	if ageMinMonths.Valid {
		age := int(ageMinMonths.Int64)
		sc.AgeMinMonths = &age
	}
	if ageMaxMonths.Valid {
		age := int(ageMaxMonths.Int64)
		sc.AgeMaxMonths = &age
	}
	if createdAt.Valid {
		sc.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	if updatedAt.Valid {
		sc.UpdatedAt = updatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(http.StatusOK, sc)
}

// CreateAdminStimulationContent creates a new stimulation content
func CreateAdminStimulationContent(c echo.Context) error {
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	var req struct {
		MilestoneID  *string `json:"milestone_id"`
		Category     string  `json:"category" validate:"required"`
		Title        string  `json:"title" validate:"required"`
		Description  string  `json:"description"`
		ContentType  string  `json:"content_type" validate:"required"`
		URL          string  `json:"url" validate:"required"`
		ThumbnailURL string  `json:"thumbnail_url"`
		AgeMinMonths *int    `json:"age_min_months"`
		AgeMaxMonths *int    `json:"age_max_months"`
		IsActive     bool    `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate category
	validCategories := map[string]bool{"sensory": true, "motor": true, "perception": true, "cognitive": true}
	if !validCategories[req.Category] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category"})
	}

	// Validate content type
	if req.ContentType != "article" && req.ContentType != "video" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Content type must be 'article' or 'video'"})
	}

	var contentID string
	err := db.DB.QueryRow(
		`INSERT INTO stimulation_content (milestone_id, category, title, description, content_type, url, thumbnail_url,
		 age_min_months, age_max_months, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id`,
		req.MilestoneID, req.Category, req.Title, req.Description, req.ContentType, req.URL,
		req.ThumbnailURL, req.AgeMinMonths, req.AgeMaxMonths, req.IsActive,
	).Scan(&contentID)

	if err != nil {
		c.Logger().Errorf("CreateAdminStimulationContent error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	contentData := map[string]interface{}{
		"category":     req.Category,
		"title":        req.Title,
		"content_type": req.ContentType,
	}
	utils.LogAudit(adminUserID, "create", "stimulation_content", &contentID, nil, contentData, ipAddress, userAgent)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      contentID,
		"message": "Stimulation content created successfully",
	})
}

// UpdateAdminStimulationContent updates a stimulation content
func UpdateAdminStimulationContent(c echo.Context) error {
	contentID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(contentID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid content ID format"})
	}

	// Get existing content for audit log
	var existing struct {
		Category    string
		Title       string
		ContentType string
	}

	err := db.DB.QueryRow(
		`SELECT category, title, content_type FROM stimulation_content WHERE id = $1`,
		contentID,
	).Scan(&existing.Category, &existing.Title, &existing.ContentType)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stimulation content not found"})
	}
	if err != nil {
		c.Logger().Errorf("UpdateAdminStimulationContent get existing error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var req struct {
		MilestoneID  *string `json:"milestone_id"`
		Category     *string `json:"category"`
		Title        *string `json:"title"`
		Description  *string `json:"description"`
		ContentType  *string `json:"content_type"`
		URL          *string `json:"url"`
		ThumbnailURL *string `json:"thumbnail_url"`
		AgeMinMonths *int    `json:"age_min_months"`
		AgeMaxMonths *int    `json:"age_max_months"`
		IsActive     *bool   `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.MilestoneID != nil {
		updateFields = append(updateFields, "milestone_id = $"+strconv.Itoa(argIndex))
		args = append(args, *req.MilestoneID)
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
	if req.Title != nil {
		updateFields = append(updateFields, "title = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		updateFields = append(updateFields, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.ContentType != nil {
		if *req.ContentType != "article" && *req.ContentType != "video" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Content type must be 'article' or 'video'"})
		}
		updateFields = append(updateFields, "content_type = $"+strconv.Itoa(argIndex))
		args = append(args, *req.ContentType)
		argIndex++
	}
	if req.URL != nil {
		updateFields = append(updateFields, "url = $"+strconv.Itoa(argIndex))
		args = append(args, *req.URL)
		argIndex++
	}
	if req.ThumbnailURL != nil {
		updateFields = append(updateFields, "thumbnail_url = $"+strconv.Itoa(argIndex))
		args = append(args, *req.ThumbnailURL)
		argIndex++
	}
	if req.AgeMinMonths != nil {
		updateFields = append(updateFields, "age_min_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMinMonths)
		argIndex++
	}
	if req.AgeMaxMonths != nil {
		updateFields = append(updateFields, "age_max_months = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AgeMaxMonths)
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
	args = append(args, contentID)
	query := `UPDATE stimulation_content SET ` + updateFields[0]
	for i := 1; i < len(updateFields); i++ {
		query += `, ` + updateFields[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		c.Logger().Errorf("UpdateAdminStimulationContent error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	beforeData := map[string]interface{}{
		"category":     existing.Category,
		"title":        existing.Title,
		"content_type": existing.ContentType,
	}
	afterData := map[string]interface{}{}
	if req.Category != nil {
		afterData["category"] = *req.Category
	}
	if req.Title != nil {
		afterData["title"] = *req.Title
	}
	if req.ContentType != nil {
		afterData["content_type"] = *req.ContentType
	}
	utils.LogAudit(adminUserID, "update", "stimulation_content", &contentID, beforeData, afterData, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Stimulation content updated successfully"})
}

// DeleteAdminStimulationContent deletes a stimulation content
func DeleteAdminStimulationContent(c echo.Context) error {
	contentID := c.Param("id")
	adminUserID := c.Get("user_id").(string)
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	
	// Validate UUID format
	if err := utils.ValidateUUID(contentID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid content ID format"})
	}

	// Get content data for audit log
	var content struct {
		Category    string
		Title       string
		ContentType string
	}

	err := db.DB.QueryRow(
		`SELECT category, title, content_type FROM stimulation_content WHERE id = $1`,
		contentID,
	).Scan(&content.Category, &content.Title, &content.ContentType)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stimulation content not found"})
	}
	if err != nil {
		c.Logger().Errorf("DeleteAdminStimulationContent get content error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Delete content
	_, err = db.DB.Exec("DELETE FROM stimulation_content WHERE id = $1", contentID)
	if err != nil {
		c.Logger().Errorf("DeleteAdminStimulationContent error: %v", err)
		sanitizedErr := utils.SanitizeError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": sanitizedErr})
	}

	// Log audit
	contentData := map[string]interface{}{
		"category":     content.Category,
		"title":        content.Title,
		"content_type": content.ContentType,
	}
	utils.LogAudit(adminUserID, "delete", "stimulation_content", &contentID, contentData, nil, ipAddress, userAgent)

	return c.JSON(http.StatusOK, map[string]string{"message": "Stimulation content deleted successfully"})
}

