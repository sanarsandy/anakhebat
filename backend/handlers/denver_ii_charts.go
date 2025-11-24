package handlers

import (
	"net/http"
	"strconv"
	"tukem-backend/db"
	"tukem-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetDenverIIChartData retrieves assessment data grouped by Denver II domain for charting
func GetDenverIIChartData(c echo.Context) error {
	childID := c.Param("id")

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := db.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Query assessments with Denver II domain grouping
	// Only include milestones that have been assessed for this child
	query := `
		SELECT 
			m.denver_domain,
			m.age_months,
			COUNT(DISTINCT m.id) as total_milestones,
			COUNT(CASE WHEN a.status = 'yes' THEN 1 END) as passed,
			COUNT(CASE WHEN a.status = 'no' THEN 1 END) as failed,
			COUNT(CASE WHEN a.status = 'sometimes' THEN 1 END) as sometimes
		FROM assessments a
		JOIN milestones m ON a.milestone_id = m.id
		WHERE a.child_id = $1
			AND m.denver_domain IS NOT NULL
			AND m.source = 'DENVER'
		GROUP BY m.denver_domain, m.age_months
		ORDER BY m.denver_domain, m.age_months ASC
	`

	type DomainData struct {
		Domain          string  `json:"domain" db:"denver_domain"`
		AgeMonths       int     `json:"age_months" db:"age_months"`
		TotalMilestones int     `json:"total_milestones" db:"total_milestones"`
		Passed          int     `json:"passed" db:"passed"`
		Failed          int     `json:"failed" db:"failed"`
		Sometimes       int     `json:"sometimes" db:"sometimes"`
		PassRate        float64 `json:"pass_rate"` // Calculated field
	}

	rows, err := db.DB.Queryx(query, childID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch Denver II chart data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data grafik Denver II",
		})
	}
	defer rows.Close()

	var data []DomainData
	for rows.Next() {
		var d DomainData
		if err := rows.StructScan(&d); err != nil {
			continue
		}
		// Calculate pass rate
		// "yes" = 100%, "sometimes" = 50%, "no" = 0%
		if d.TotalMilestones > 0 {
			// Pass rate = (passed * 1.0 + sometimes * 0.5) / total * 100
			passScore := float64(d.Passed)*1.0 + float64(d.Sometimes)*0.5
			d.PassRate = (passScore / float64(d.TotalMilestones)) * 100
		}
		data = append(data, d)
	}

	// Group by domain
	domainGroups := make(map[string][]DomainData)
	for _, d := range data {
		domainGroups[d.Domain] = append(domainGroups[d.Domain], d)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"child_id": childID,
		"domains":  domainGroups,
		"raw_data": data,
	})
}

// GetDenverIIMilestones retrieves all Denver II milestones (for assessment)
func GetDenverIIMilestones(c echo.Context) error {
	ageMonthsStr := c.QueryParam("age_months")
	
	ageMonths := 0
	if ageMonthsStr != "" {
		var err error
		ageMonths, err = strconv.Atoi(ageMonthsStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid age_months parameter"})
		}
	}

	query := `
		SELECT * FROM milestones 
		WHERE source = 'DENVER' 
			AND denver_domain IS NOT NULL
			AND age_months >= $1 - 3 AND age_months <= $1 + 6
		ORDER BY age_months ASC, denver_domain ASC
	`
	
	var milestones []models.Milestone
	err := db.DB.Select(&milestones, query, ageMonths)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch Denver II milestones"})
	}

	return c.JSON(http.StatusOK, milestones)
}

// GetDenverIIChartGridData retrieves all Denver II milestones with assessment status for grid chart
func GetDenverIIChartGridData(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			c.Logger().Errorf("Panic in GetDenverIIChartGridData: %v", r)
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}()
	
	childID := c.Param("id")
	c.Logger().Errorf("GetDenverIIChartGridData CALLED for childID: %s", childID) // Use Errorf to ensure it's logged

	// Get user ID from JWT to verify ownership (same pattern as GetDenverIIChartData)
	userToken := c.Get("user")
	if userToken == nil {
		c.Logger().Errorf("No user token found")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	
	user, ok := userToken.(*jwt.Token)
	if !ok {
		c.Logger().Errorf("Invalid user token type: %T", userToken)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	
	claimsMap, ok := user.Claims.(*jwt.MapClaims)
	if !ok {
		c.Logger().Errorf("Invalid claims type: %T", user.Claims)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	
	claims := *claimsMap
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.Logger().Errorf("Invalid user_id in claims: %v", claims["user_id"])
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	
	c.Logger().Infof("User ID extracted: %s", userID)

	// Verify child belongs to user
	var parentID string
	err := db.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
	if err != nil {
		c.Logger().Errorf("Child not found: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}
	if parentID != userID {
		c.Logger().Errorf("Unauthorized: parentID=%s, userID=%s", parentID, userID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}
	c.Logger().Infof("Child verified: childID=%s", childID)

	// Get child's age
	var dob string
	err = db.DB.QueryRow("SELECT dob FROM children WHERE id = $1", childID).Scan(&dob)
	if err != nil {
		c.Logger().Errorf("Failed to get DOB: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}
	c.Logger().Infof("DOB retrieved: %s", dob)

	// Fetch all Denver II milestones with assessment status
	// Get the latest assessment for each milestone
	query := `
		SELECT 
			m.id,
			m.age_months,
			m.min_age_range,
			m.max_age_range,
			m.denver_domain,
			m.question,
			m.question_en,
			m.is_red_flag,
			COALESCE(a.status, '') as assessment_status
		FROM milestones m
		LEFT JOIN (
			SELECT DISTINCT ON (milestone_id) milestone_id, status
			FROM assessments 
			WHERE child_id = $1 
			ORDER BY milestone_id, assessment_date DESC, updated_at DESC
		) a ON m.id = a.milestone_id
		WHERE m.source = 'DENVER' 
			AND m.denver_domain IS NOT NULL
		ORDER BY 
			CASE m.denver_domain
				WHEN 'PS' THEN 1
				WHEN 'FM' THEN 2
				WHEN 'L' THEN 3
				WHEN 'GM' THEN 4
			END,
			m.age_months ASC
	`

	type MilestoneWithStatus struct {
		ID                string  `json:"id" db:"id"`
		AgeMonths         int     `json:"age_months" db:"age_months"`
		MinAgeRange       *int    `json:"min_age_range" db:"min_age_range"`
		MaxAgeRange       *int    `json:"max_age_range" db:"max_age_range"`
		DenverDomain      string  `json:"denver_domain" db:"denver_domain"`
		Question          string  `json:"question" db:"question"`
		QuestionEn        string  `json:"question_en" db:"question_en"`
		IsRedFlag         bool    `json:"is_red_flag" db:"is_red_flag"`
		AssessmentStatus  string  `json:"assessment_status" db:"assessment_status"` // yes, no, sometimes, or empty
		// Calculated percentile fields for Denver II chart
		Age25Percentile   int     `json:"age_25_percentile"` // 25% of children reach this milestone
		Age50Percentile   int     `json:"age_50_percentile"` // 50% of children reach this milestone (median)
		Age75Percentile   int     `json:"age_75_percentile"` // 75% of children reach this milestone
		Age90Percentile   int     `json:"age_90_percentile"` // 90% of children reach this milestone
	}

	// Helper function to calculate percentiles from available data
	calculatePercentiles := func(ageMonths int, minAgeRange *int, maxAgeRange *int) (int, int, int, int) {
		// Use age_months as 50% percentile (median)
		age50 := ageMonths
		
		// Calculate 25% percentile
		var age25 int
		if minAgeRange != nil && *minAgeRange > 0 {
			age25 = *minAgeRange
		} else {
			// Default: 25% is 2 months before median
			age25 = ageMonths - 2
			if age25 < 0 {
				age25 = 0
			}
		}
		
		// Calculate 90% percentile
		var age90 int
		if maxAgeRange != nil && *maxAgeRange > 0 {
			age90 = *maxAgeRange
		} else {
			// Default: 90% is 2 months after median
			age90 = ageMonths + 2
		}
		
		// Calculate 75% percentile (interpolation between 50% and 90%)
		age75 := age50 + (age90-age50)*3/4
		
		return age25, age50, age75, age90
	}

	rows, err := db.DB.Queryx(query, childID)
	if err != nil {
		c.Logger().Errorf("Failed to fetch Denver II grid data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data grafik Denver II",
		})
	}
	defer rows.Close()

	var milestones []MilestoneWithStatus
	for rows.Next() {
		var m MilestoneWithStatus
		if err := rows.StructScan(&m); err != nil {
			continue
		}
		// Calculate percentiles
		m.Age25Percentile, m.Age50Percentile, m.Age75Percentile, m.Age90Percentile = 
			calculatePercentiles(m.AgeMonths, m.MinAgeRange, m.MaxAgeRange)
		milestones = append(milestones, m)
	}

	// Group by domain
	domainGroups := make(map[string][]MilestoneWithStatus)
	for _, m := range milestones {
		domainGroups[m.DenverDomain] = append(domainGroups[m.DenverDomain], m)
	}

	c.Logger().Infof("Returning data: %d milestones, %d domains", len(milestones), len(domainGroups))
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"child_id": childID,
		"dob":      dob,
		"domains":  domainGroups,
		"all_milestones": milestones,
	})
}

