package handlers

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"
	"tukem-backend/db"

	"github.com/labstack/echo/v4"
)

// Report types
type UserReport struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	PhoneNumber       string    `json:"phone_number"`
	Role              string    `json:"role"`
	AuthProvider      string    `json:"auth_provider"`
	PhoneVerified     bool      `json:"phone_verified"`
	CreatedAt         time.Time `json:"created_at"`
	ChildrenCount     int       `json:"children_count"`
	MeasurementsCount int       `json:"measurements_count"`
	AssessmentsCount  int       `json:"assessments_count"`
	ImmunizationsCount int      `json:"immunizations_count"`
}

type ChildReport struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	DOB               time.Time `json:"dob"`
	Gender            string    `json:"gender"`
	BirthWeight       float64   `json:"birth_weight"`
	BirthHeight       float64   `json:"birth_height"`
	IsPremature       bool      `json:"is_premature"`
	GestationalAge    *int      `json:"gestational_age"`
	CreatedAt         time.Time `json:"created_at"`
	ParentName        string    `json:"parent_name"`
	ParentEmail       string    `json:"parent_email"`
	ParentPhone       string    `json:"parent_phone"`
	MeasurementsCount int       `json:"measurements_count"`
	AssessmentsCount  int       `json:"assessments_count"`
	ImmunizationsCount int      `json:"immunizations_count"`
}

type GrowthReport struct {
	ID                  string    `json:"id"`
	ChildID             string    `json:"child_id"`
	ChildName           string    `json:"child_name"`
	ParentName          string    `json:"parent_name"`
	MeasuredAt          time.Time `json:"measured_at"`
	AgeMonths           int       `json:"age_months"`
	Weight              float64   `json:"weight"`
	Height              float64   `json:"height"`
	HeadCircumference   *float64  `json:"head_circumference"`
	WeightForAgeZScore  *float64  `json:"weight_for_age_zscore"`
	HeightForAgeZScore  *float64  `json:"height_for_age_zscore"`
	WeightStatus        string    `json:"weight_status"`
	HeightStatus        string    `json:"height_status"`
}

// GetUsersReport generates a users report
func GetUsersReport(c echo.Context) error {
	format := c.QueryParam("format")
	if format == "" {
		format = "json"
	}

	// Get query parameters
	role := c.QueryParam("role")
	authProvider := c.QueryParam("auth_provider")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query
	query := `
		SELECT 
			u.id,
			u.email,
			u.full_name,
			u.phone_number,
			u.role,
			u.auth_provider,
			u.phone_verified,
			u.created_at,
			COUNT(DISTINCT c.id) as children_count,
			COUNT(DISTINCT m.id) as measurements_count,
			COUNT(DISTINCT a.assessment_date) as assessments_count,
			COUNT(DISTINCT ci.id) as immunizations_count
		FROM users u
		LEFT JOIN children c ON c.parent_id = u.id
		LEFT JOIN measurements m ON m.child_id = c.id
		LEFT JOIN assessments a ON a.child_id = c.id
		LEFT JOIN child_immunizations ci ON ci.child_id = c.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if role != "" {
		query += ` AND u.role = $` + strconv.Itoa(argIndex)
		args = append(args, role)
		argIndex++
	}

	if authProvider != "" {
		query += ` AND u.auth_provider = $` + strconv.Itoa(argIndex)
		args = append(args, authProvider)
		argIndex++
	}

	if dateFrom != "" {
		query += ` AND u.created_at >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		query += ` AND u.created_at <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	query += ` GROUP BY u.id, u.email, u.full_name, u.phone_number, u.role, u.auth_provider, u.phone_verified, u.created_at ORDER BY u.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetUsersReport query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var users []UserReport
	for rows.Next() {
		var u UserReport
		var email sql.NullString
		var phoneNumber sql.NullString
		var authProvider sql.NullString

		err := rows.Scan(
			&u.ID, &email, &u.FullName, &phoneNumber, &u.Role, &authProvider,
			&u.PhoneVerified, &u.CreatedAt, &u.ChildrenCount, &u.MeasurementsCount,
			&u.AssessmentsCount, &u.ImmunizationsCount,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan user row: %v", err)
			continue
		}

		if email.Valid {
			u.Email = email.String
		}
		if phoneNumber.Valid {
			u.PhoneNumber = phoneNumber.String
		}
		if authProvider.Valid {
			u.AuthProvider = authProvider.String
		}

		users = append(users, u)
	}

	if format == "csv" {
		return exportUsersCSV(c, users)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":  users,
		"total":  len(users),
		"format": format,
	})
}

// GetChildrenReport generates a children report
func GetChildrenReport(c echo.Context) error {
	format := c.QueryParam("format")
	if format == "" {
		format = "json"
	}

	// Get query parameters
	parentID := c.QueryParam("parent_id")
	gender := c.QueryParam("gender")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query
	query := `
		SELECT 
			c.id,
			c.name,
			c.dob,
			c.gender,
			c.birth_weight,
			c.birth_height,
			c.is_premature,
			c.gestational_age,
			c.created_at,
			u.full_name as parent_name,
			u.email as parent_email,
			u.phone_number as parent_phone,
			COUNT(DISTINCT m.id) as measurements_count,
			COUNT(DISTINCT a.assessment_date) as assessments_count,
			COUNT(DISTINCT ci.id) as immunizations_count
		FROM children c
		JOIN users u ON u.id = c.parent_id
		LEFT JOIN measurements m ON m.child_id = c.id
		LEFT JOIN assessments a ON a.child_id = c.id
		LEFT JOIN child_immunizations ci ON ci.child_id = c.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if parentID != "" {
		query += ` AND c.parent_id = $` + strconv.Itoa(argIndex)
		args = append(args, parentID)
		argIndex++
	}

	if gender != "" {
		query += ` AND c.gender = $` + strconv.Itoa(argIndex)
		args = append(args, gender)
		argIndex++
	}

	if dateFrom != "" {
		query += ` AND c.created_at >= $` + strconv.Itoa(argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		query += ` AND c.created_at <= $` + strconv.Itoa(argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	query += ` GROUP BY c.id, c.name, c.dob, c.gender, c.birth_weight, c.birth_height, c.is_premature, c.gestational_age, c.created_at, u.full_name, u.email, u.phone_number ORDER BY c.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetChildrenReport query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var children []ChildReport
	for rows.Next() {
		var ch ChildReport
		var parentEmail sql.NullString
		var parentPhone sql.NullString
		var gestationalAge sql.NullInt64

		err := rows.Scan(
			&ch.ID, &ch.Name, &ch.DOB, &ch.Gender, &ch.BirthWeight, &ch.BirthHeight,
			&ch.IsPremature, &gestationalAge, &ch.CreatedAt, &ch.ParentName,
			&parentEmail, &parentPhone, &ch.MeasurementsCount, &ch.AssessmentsCount,
			&ch.ImmunizationsCount,
		)
		if err != nil {
			continue
		}

		if parentEmail.Valid {
			ch.ParentEmail = parentEmail.String
		}
		if parentPhone.Valid {
			ch.ParentPhone = parentPhone.String
		}
		if gestationalAge.Valid {
			age := int(gestationalAge.Int64)
			ch.GestationalAge = &age
		}

		children = append(children, ch)
	}

	if format == "csv" {
		return exportChildrenCSV(c, children)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"children": children,
		"total":    len(children),
		"format":   format,
	})
}

// GetGrowthReport generates a growth report
func GetGrowthReport(c echo.Context) error {
	format := c.QueryParam("format")
	if format == "" {
		format = "json"
	}

	// Get query parameters
	childID := c.QueryParam("child_id")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")

	// Build query
	query := `
		SELECT 
			m.id,
			m.child_id,
			c.name as child_name,
			u.full_name as parent_name,
			m.measurement_date,
			m.age_in_months,
			m.weight,
			m.height,
			m.head_circumference,
			m.weight_for_age_zscore,
			m.height_for_age_zscore,
			m.weight_status,
			m.height_status
		FROM measurements m
		JOIN children c ON c.id = m.child_id
		JOIN users u ON u.id = c.parent_id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if childID != "" {
		query += ` AND m.child_id = $` + strconv.Itoa(argIndex)
		args = append(args, childID)
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

	query += ` ORDER BY m.measurement_date DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("GetGrowthReport query error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var measurements []GrowthReport
	for rows.Next() {
		var m GrowthReport
		var headCirc sql.NullFloat64
		var wfaZScore sql.NullFloat64
		var hfaZScore sql.NullFloat64
		var weightStatus sql.NullString
		var heightStatus sql.NullString
		var measurementDate time.Time

		err := rows.Scan(
			&m.ID, &m.ChildID, &m.ChildName, &m.ParentName, &measurementDate,
			&m.AgeMonths, &m.Weight, &m.Height, &headCirc, &wfaZScore,
			&hfaZScore, &weightStatus, &heightStatus,
		)
		if err != nil {
			c.Logger().Errorf("Failed to scan growth report row: %v", err)
			continue
		}

		// Convert measurement_date (DATE) to MeasuredAt (TIMESTAMP)
		m.MeasuredAt = time.Date(measurementDate.Year(), measurementDate.Month(), measurementDate.Day(), 0, 0, 0, 0, time.UTC)

		if headCirc.Valid {
			m.HeadCircumference = &headCirc.Float64
		}
		if wfaZScore.Valid {
			m.WeightForAgeZScore = &wfaZScore.Float64
		}
		if hfaZScore.Valid {
			m.HeightForAgeZScore = &hfaZScore.Float64
		}
		if weightStatus.Valid {
			m.WeightStatus = weightStatus.String
		} else {
			m.WeightStatus = ""
		}
		if heightStatus.Valid {
			m.HeightStatus = heightStatus.String
		} else {
			m.HeightStatus = ""
		}

		measurements = append(measurements, m)
	}

	if format == "csv" {
		return exportGrowthCSV(c, measurements)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"measurements": measurements,
		"total":        len(measurements),
		"format":       format,
	})
}

// Helper functions for CSV export
func exportUsersCSV(c echo.Context, users []UserReport) error {
	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=users_report.csv")

	writer := csv.NewWriter(c.Response())
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Email", "Full Name", "Phone Number", "Role", "Auth Provider",
		"Phone Verified", "Created At", "Children Count", "Measurements Count",
		"Assessments Count", "Immunizations Count",
	}
	if err := writer.Write(header); err != nil {
		c.Logger().Errorf("Failed to write CSV header: %v", err)
		return err
	}

	// Write data
	for _, u := range users {
		record := []string{
			u.ID, u.Email, u.FullName, u.PhoneNumber, u.Role, u.AuthProvider,
			strconv.FormatBool(u.PhoneVerified), u.CreatedAt.Format(time.RFC3339),
			strconv.Itoa(u.ChildrenCount), strconv.Itoa(u.MeasurementsCount),
			strconv.Itoa(u.AssessmentsCount), strconv.Itoa(u.ImmunizationsCount),
		}
		if err := writer.Write(record); err != nil {
			c.Logger().Errorf("Failed to write CSV record: %v", err)
			return err
		}
	}

	return nil
}

func exportChildrenCSV(c echo.Context, children []ChildReport) error {
	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=children_report.csv")

	writer := csv.NewWriter(c.Response())
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Name", "DOB", "Gender", "Birth Weight", "Birth Height",
		"Is Premature", "Gestational Age", "Created At", "Parent Name",
		"Parent Email", "Parent Phone", "Measurements Count", "Assessments Count",
		"Immunizations Count",
	}
	if err := writer.Write(header); err != nil {
		c.Logger().Errorf("Failed to write CSV header: %v", err)
		return err
	}

	// Write data
	for _, ch := range children {
		gestAge := ""
		if ch.GestationalAge != nil {
			gestAge = strconv.Itoa(*ch.GestationalAge)
		}
		record := []string{
			ch.ID, ch.Name, ch.DOB.Format(time.RFC3339), ch.Gender,
			strconv.FormatFloat(ch.BirthWeight, 'f', 2, 64),
			strconv.FormatFloat(ch.BirthHeight, 'f', 2, 64),
			strconv.FormatBool(ch.IsPremature), gestAge,
			ch.CreatedAt.Format(time.RFC3339), ch.ParentName, ch.ParentEmail,
			ch.ParentPhone, strconv.Itoa(ch.MeasurementsCount),
			strconv.Itoa(ch.AssessmentsCount), strconv.Itoa(ch.ImmunizationsCount),
		}
		if err := writer.Write(record); err != nil {
			c.Logger().Errorf("Failed to write CSV record: %v", err)
			return err
		}
	}

	return nil
}

func exportGrowthCSV(c echo.Context, measurements []GrowthReport) error {
	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=growth_report.csv")

	writer := csv.NewWriter(c.Response())
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Child ID", "Child Name", "Parent Name", "Measured At", "Age Months",
		"Weight", "Height", "Head Circumference", "Weight for Age Z-Score",
		"Height for Age Z-Score", "Weight Status", "Height Status",
	}
	if err := writer.Write(header); err != nil {
		c.Logger().Errorf("Failed to write CSV header: %v", err)
		return err
	}

	// Write data
	for _, m := range measurements {
		headCirc := ""
		if m.HeadCircumference != nil {
			headCirc = strconv.FormatFloat(*m.HeadCircumference, 'f', 2, 64)
		}
		wfaZ := ""
		if m.WeightForAgeZScore != nil {
			wfaZ = strconv.FormatFloat(*m.WeightForAgeZScore, 'f', 2, 64)
		}
		hfaZ := ""
		if m.HeightForAgeZScore != nil {
			hfaZ = strconv.FormatFloat(*m.HeightForAgeZScore, 'f', 2, 64)
		}
		record := []string{
			m.ID, m.ChildID, m.ChildName, m.ParentName,
			m.MeasuredAt.Format(time.RFC3339), strconv.Itoa(m.AgeMonths),
			strconv.FormatFloat(m.Weight, 'f', 2, 64),
			strconv.FormatFloat(m.Height, 'f', 2, 64), headCirc, wfaZ, hfaZ,
			m.WeightStatus, m.HeightStatus,
		}
		if err := writer.Write(record); err != nil {
			c.Logger().Errorf("Failed to write CSV record: %v", err)
			return err
		}
	}

	return nil
}
