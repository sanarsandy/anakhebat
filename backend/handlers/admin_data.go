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

// GetAdminChildren returns all children with pagination and filters
func GetAdminChildren(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	parentID := c.QueryParam("parent_id")
	gender := c.QueryParam("gender")
	search := c.QueryParam("search")

	// Build query
	query := `SELECT c.id, c.parent_id, c.name, c.dob, c.gender, c.birth_weight, 
	          c.birth_height, c.is_premature, c.gestational_age, c.created_at,
	          u.full_name as parent_name, u.email as parent_email
	          FROM children c
	          JOIN users u ON u.id = c.parent_id
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if parentID != "" {
		// Validate UUID format for parent_id filter
		if err := utils.ValidateUUID(parentID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid parent ID format"})
		}
		query += ` AND c.parent_id = $` + strconv.Itoa(argIndex)
		args = append(args, parentID)
		argIndex++
	}

	if gender != "" {
		query += ` AND c.gender = $` + strconv.Itoa(argIndex)
		args = append(args, gender)
		argIndex++
	}

	if search != "" {
		query += ` AND (c.name ILIKE $` + strconv.Itoa(argIndex) +
			` OR u.full_name ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	query += ` ORDER BY c.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminChildren query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	children := []map[string]interface{}{}
	for rows.Next() {
		var child models.Child
		var parentName sql.NullString
		var parentEmail sql.NullString

		err := rows.Scan(&child.ID, &child.ParentID, &child.Name, &child.DOB, &child.Gender,
			&child.BirthWeight, &child.BirthHeight, &child.IsPremature, &child.GestationalAge,
			&child.CreatedAt, &parentName, &parentEmail)
		if err != nil {
			continue
		}

		childMap := map[string]interface{}{
			"id":              child.ID,
			"parent_id":       child.ParentID,
			"name":            child.Name,
			"dob":             child.DOB,
			"gender":          child.Gender,
			"birth_weight":    child.BirthWeight,
			"birth_height":    child.BirthHeight,
			"is_premature":    child.IsPremature,
			"gestational_age": child.GestationalAge,
			"created_at":      child.CreatedAt,
			"parent_name":     parentName.String,
			"parent_email":    parentEmail.String,
		}
		children = append(children, childMap)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM children c WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if parentID != "" {
		countQuery += ` AND c.parent_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, parentID)
		countArgIndex++
	}

	if gender != "" {
		countQuery += ` AND c.gender = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, gender)
		countArgIndex++
	}

	if search != "" {
		countQuery += ` AND EXISTS (
			SELECT 1 FROM users u 
			WHERE u.id = c.parent_id 
			AND (c.name ILIKE $` + strconv.Itoa(countArgIndex) +
			` OR u.full_name ILIKE $` + strconv.Itoa(countArgIndex) + `)
		)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
		countArgIndex += 2
	}

	var total int
	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminChildren count error: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"children": children,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminChild returns a single child with full details
func GetAdminChild(c echo.Context) error {
	childID := c.Param("id")
	
	// Validate UUID format
	if err := utils.ValidateUUID(childID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid child ID format"})
	}

	var child models.Child
	var parentName sql.NullString
	var parentEmail sql.NullString
	var parentPhone sql.NullString

	query := `SELECT c.id, c.parent_id, c.name, c.dob, c.gender, c.birth_weight, 
	          c.birth_height, c.is_premature, c.gestational_age, c.created_at,
	          u.full_name, u.email, u.phone_number
	          FROM children c
	          JOIN users u ON u.id = c.parent_id
	          WHERE c.id = $1`
	err := db.DB.QueryRow(query, childID).Scan(
		&child.ID, &child.ParentID, &child.Name, &child.DOB, &child.Gender,
		&child.BirthWeight, &child.BirthHeight, &child.IsPremature, &child.GestationalAge,
		&child.CreatedAt, &parentName, &parentEmail, &parentPhone)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	} else if err != nil {
		c.Logger().Errorf("GetAdminChild error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Get measurements count
	var measurementsCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM measurements WHERE child_id = $1", childID).Scan(&measurementsCount)

	// Get assessments count
	var assessmentsCount int
	db.DB.QueryRow("SELECT COUNT(DISTINCT assessment_date) FROM assessments WHERE child_id = $1", childID).Scan(&assessmentsCount)

	// Get immunizations count
	var immunizationsCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM child_immunizations WHERE child_id = $1", childID).Scan(&immunizationsCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"child": child,
		"parent": map[string]interface{}{
			"id":          child.ParentID,
			"name":        parentName.String,
			"email":       parentEmail.String,
			"phone_number": parentPhone.String,
		},
		"statistics": map[string]interface{}{
			"measurements_count":  measurementsCount,
			"assessments_count":   assessmentsCount,
			"immunizations_count": immunizationsCount,
		},
	})
}

// GetAdminMeasurements returns all measurements with pagination and filters
func GetAdminMeasurements(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	childID := c.QueryParam("child_id")
	parentID := c.QueryParam("parent_id")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query
	query := `SELECT m.id, m.child_id, m.measurement_date, m.weight, m.height, 
	          m.head_circumference, m.weight_for_age_zscore, m.height_for_age_zscore, 
	          m.weight_status, m.height_status, m.created_at,
	          c.name as child_name, c.dob as child_dob, c.gender as child_gender,
	          u.full_name as parent_name, u.email as parent_email
	          FROM measurements m
	          JOIN children c ON c.id = m.child_id
	          JOIN users u ON u.id = c.parent_id
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if childID != "" {
		query += ` AND m.child_id = $` + strconv.Itoa(argIndex)
		args = append(args, childID)
		argIndex++
	}

	if parentID != "" {
		query += ` AND c.parent_id = $` + strconv.Itoa(argIndex)
		args = append(args, parentID)
		argIndex++
	}

	if dateFrom != "" {
		query += ` AND m.measurement_date >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		query += ` AND m.measurement_date <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	query += ` ORDER BY m.measurement_date DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminMeasurements query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	measurements := []map[string]interface{}{}
	for rows.Next() {
		var m models.Measurement
		var childName sql.NullString
		var childDOB sql.NullString
		var childGender sql.NullString
		var parentName sql.NullString
		var parentEmail sql.NullString
		var weightForAgeZScore sql.NullFloat64
		var heightForAgeZScore sql.NullFloat64
		var weightStatus sql.NullString
		var heightStatus sql.NullString

		err := rows.Scan(&m.ID, &m.ChildID, &m.MeasurementDate, &m.Weight, &m.Height,
			&m.HeadCircumference, &weightForAgeZScore, &heightForAgeZScore,
			&weightStatus, &heightStatus, &m.CreatedAt,
			&childName, &childDOB, &childGender, &parentName, &parentEmail)
		if err != nil {
			continue
		}

		measurementMap := map[string]interface{}{
			"id":                        m.ID,
			"child_id":                  m.ChildID,
			"measurement_date":          m.MeasurementDate,
			"weight":                    m.Weight,
			"height":                    m.Height,
			"head_circumference":        m.HeadCircumference,
			"weight_for_age_zscore":     getFloat64FromNull(weightForAgeZScore),
			"height_for_age_zscore":     getFloat64FromNull(heightForAgeZScore),
			"weight_status":             weightStatus.String,
			"height_status":             heightStatus.String,
			"created_at":                m.CreatedAt,
			"child_name":                childName.String,
			"child_dob":                 childDOB.String,
			"child_gender":              childGender.String,
			"parent_name":               parentName.String,
			"parent_email":              parentEmail.String,
		}
		measurements = append(measurements, measurementMap)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM measurements m
	               JOIN children c ON c.id = m.child_id
	               WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if childID != "" {
		countQuery += ` AND m.child_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, childID)
		countArgIndex++
	}

	if parentID != "" {
		countQuery += ` AND c.parent_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, parentID)
		countArgIndex++
	}

	if dateFrom != "" {
		countQuery += ` AND m.measurement_date >= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateFrom)
		countArgIndex++
	}

	if dateTo != "" {
		countQuery += ` AND m.measurement_date <= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateTo)
		countArgIndex++
	}

	var total int
	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminMeasurements count error: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"measurements": measurements,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminAssessments returns all assessments with pagination and filters
func GetAdminAssessments(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	childID := c.QueryParam("child_id")
	parentID := c.QueryParam("parent_id")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query
	query := `SELECT a.id, a.child_id, a.assessment_date, a.created_at,
	          c.name as child_name, c.dob as child_dob, c.gender as child_gender,
	          u.full_name as parent_name, u.email as parent_email,
	          COUNT(DISTINCT ma.milestone_id) as milestones_count
	          FROM assessments a
	          JOIN children c ON c.id = a.child_id
	          JOIN users u ON u.id = c.parent_id
	          LEFT JOIN milestone_assessments ma ON ma.assessment_id = a.id
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if childID != "" {
		query += ` AND a.child_id = $` + strconv.Itoa(argIndex)
		args = append(args, childID)
		argIndex++
	}

	if parentID != "" {
		query += ` AND c.parent_id = $` + strconv.Itoa(argIndex)
		args = append(args, parentID)
		argIndex++
	}

	if dateFrom != "" {
		query += ` AND a.assessment_date >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		query += ` AND a.assessment_date <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	query += ` GROUP BY a.id, a.child_id, a.assessment_date, a.created_at,
	          c.name, c.dob, c.gender, u.full_name, u.email
	          ORDER BY a.assessment_date DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminAssessments query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	assessments := []map[string]interface{}{}
	for rows.Next() {
		var assessmentID string
		var childID string
		var assessmentDate string
		var createdAt string
		var childName sql.NullString
		var childDOB sql.NullString
		var childGender sql.NullString
		var parentName sql.NullString
		var parentEmail sql.NullString
		var milestonesCount int

		err := rows.Scan(&assessmentID, &childID, &assessmentDate, &createdAt,
			&childName, &childDOB, &childGender, &parentName, &parentEmail, &milestonesCount)
		if err != nil {
			continue
		}

		assessmentMap := map[string]interface{}{
			"id":               assessmentID,
			"child_id":         childID,
			"assessment_date":  assessmentDate,
			"created_at":       createdAt,
			"milestones_count": milestonesCount,
			"child_name":       childName.String,
			"child_dob":        childDOB.String,
			"child_gender":     childGender.String,
			"parent_name":      parentName.String,
			"parent_email":     parentEmail.String,
		}
		assessments = append(assessments, assessmentMap)
	}

	// Get total count
	countQuery := `SELECT COUNT(DISTINCT a.id) FROM assessments a
	               JOIN children c ON c.id = a.child_id
	               WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if childID != "" {
		countQuery += ` AND a.child_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, childID)
		countArgIndex++
	}

	if parentID != "" {
		countQuery += ` AND c.parent_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, parentID)
		countArgIndex++
	}

	if dateFrom != "" {
		countQuery += ` AND a.assessment_date >= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateFrom)
		countArgIndex++
	}

	if dateTo != "" {
		countQuery += ` AND a.assessment_date <= $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, dateTo)
		countArgIndex++
	}

	var total int
	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminAssessments count error: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"assessments": assessments,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAdminImmunizations returns all immunizations with pagination and filters
func GetAdminImmunizations(c echo.Context) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	childID := c.QueryParam("child_id")
	parentID := c.QueryParam("parent_id")
	status := c.QueryParam("status")

	// Build query
	query := `SELECT ci.id, ci.child_id, ci.schedule_id, ci.immunization_date, 
	          ci.status, ci.notes, ci.created_at,
	          c.name as child_name, c.dob as child_dob, c.gender as child_gender,
	          u.full_name as parent_name, u.email as parent_email,
	          is.name_id as schedule_name, is.category as schedule_category
	          FROM child_immunizations ci
	          JOIN children c ON c.id = ci.child_id
	          JOIN users u ON u.id = c.parent_id
	          JOIN immunization_schedules is ON is.id = ci.schedule_id
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if childID != "" {
		query += ` AND ci.child_id = $` + strconv.Itoa(argIndex)
		args = append(args, childID)
		argIndex++
	}

	if parentID != "" {
		query += ` AND c.parent_id = $` + strconv.Itoa(argIndex)
		args = append(args, parentID)
		argIndex++
	}

	if status != "" {
		query += ` AND ci.status = $` + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	query += ` ORDER BY ci.immunization_date DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetAdminImmunizations query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	immunizations := []map[string]interface{}{}
	for rows.Next() {
		var immunizationID string
		var childID string
		var scheduleID string
		var immunizationDate sql.NullString
		var status sql.NullString
		var notes sql.NullString
		var createdAt string
		var childName sql.NullString
		var childDOB sql.NullString
		var childGender sql.NullString
		var parentName sql.NullString
		var parentEmail sql.NullString
		var scheduleName sql.NullString
		var scheduleCategory sql.NullString

		err := rows.Scan(&immunizationID, &childID, &scheduleID, &immunizationDate,
			&status, &notes, &createdAt, &childName, &childDOB, &childGender,
			&parentName, &parentEmail, &scheduleName, &scheduleCategory)
		if err != nil {
			continue
		}

		immunizationMap := map[string]interface{}{
			"id":                 immunizationID,
			"child_id":           childID,
			"schedule_id":        scheduleID,
			"immunization_date":  immunizationDate.String,
			"status":             status.String,
			"notes":              notes.String,
			"created_at":         createdAt,
			"child_name":         childName.String,
			"child_dob":          childDOB.String,
			"child_gender":       childGender.String,
			"parent_name":        parentName.String,
			"parent_email":       parentEmail.String,
			"schedule_name":      scheduleName.String,
			"schedule_category":  scheduleCategory.String,
		}
		immunizations = append(immunizations, immunizationMap)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM child_immunizations ci
	               JOIN children c ON c.id = ci.child_id
	               WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	if childID != "" {
		countQuery += ` AND ci.child_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, childID)
		countArgIndex++
	}

	if parentID != "" {
		countQuery += ` AND c.parent_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, parentID)
		countArgIndex++
	}

	if status != "" {
		countQuery += ` AND ci.status = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, status)
		countArgIndex++
	}

	var total int
	err = db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Logger().Errorf("GetAdminImmunizations count error: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"immunizations": immunizations,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Helper function to get float64 from sql.NullFloat64
func getFloat64FromNull(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

