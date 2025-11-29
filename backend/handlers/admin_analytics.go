package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"tukem-backend/db"

	"github.com/labstack/echo/v4"
)

// GetAdminOverview returns overview statistics for admin dashboard
func GetAdminOverview(c echo.Context) error {
	stats := make(map[string]interface{})

	// Get total users
	var totalUsers int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		c.Logger().Errorf("Failed to get total users: %v", err)
		totalUsers = 0
	}
	stats["total_users"] = totalUsers

	// Get total children
	var totalChildren int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM children").Scan(&totalChildren)
	if err != nil {
		c.Logger().Errorf("Failed to get total children: %v", err)
		totalChildren = 0
	}
	stats["total_children"] = totalChildren

	// Get total measurements
	var totalMeasurements int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM measurements").Scan(&totalMeasurements)
	if err != nil {
		c.Logger().Errorf("Failed to get total measurements: %v", err)
		totalMeasurements = 0
	}
	stats["total_measurements"] = totalMeasurements

	// Get total assessments
	var totalAssessments int
	err = db.DB.QueryRow("SELECT COUNT(DISTINCT assessment_date) FROM assessments").Scan(&totalAssessments)
	if err != nil {
		c.Logger().Errorf("Failed to get total assessments: %v", err)
		totalAssessments = 0
	}
	stats["total_assessments"] = totalAssessments

	// Get total immunizations
	var totalImmunizations int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM child_immunizations").Scan(&totalImmunizations)
	if err != nil {
		c.Logger().Errorf("Failed to get total immunizations: %v", err)
		totalImmunizations = 0
	}
	stats["total_immunizations"] = totalImmunizations

	// Get active users (users with at least 1 child)
	var activeUsers int
	err = db.DB.QueryRow("SELECT COUNT(DISTINCT parent_id) FROM children").Scan(&activeUsers)
	if err != nil {
		c.Logger().Errorf("Failed to get active users: %v", err)
		activeUsers = 0
	}
	stats["active_users"] = activeUsers

	// Get users by role
	var adminUsers int
	var parentUsers int
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminUsers)
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'parent'").Scan(&parentUsers)
	stats["admin_users"] = adminUsers
	stats["parent_users"] = parentUsers

	// Get users by auth provider
	authProviderStats := make(map[string]int)
	rows, _ := db.DB.Query("SELECT auth_provider, COUNT(*) FROM users GROUP BY auth_provider")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var provider sql.NullString
			var count int
			if rows.Scan(&provider, &count) == nil {
				if provider.Valid {
					authProviderStats[provider.String] = count
				}
			}
		}
	}
	stats["users_by_auth_provider"] = authProviderStats

	// Get children by gender
	var maleChildren int
	var femaleChildren int
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE gender = 'male'").Scan(&maleChildren)
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE gender = 'female'").Scan(&femaleChildren)
	stats["male_children"] = maleChildren
	stats["female_children"] = femaleChildren

	// Get premature vs full-term
	var prematureChildren int
	var fullTermChildren int
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE is_premature = true").Scan(&prematureChildren)
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE is_premature = false").Scan(&fullTermChildren)
	stats["premature_children"] = prematureChildren
	stats["full_term_children"] = fullTermChildren

	// Get new users in last 30 days
	var newUsersLast30Days int
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE created_at >= $1", thirtyDaysAgo).Scan(&newUsersLast30Days)
	stats["new_users_last_30_days"] = newUsersLast30Days

	// Get new children in last 30 days
	var newChildrenLast30Days int
	db.DB.QueryRow("SELECT COUNT(*) FROM children WHERE created_at >= $1", thirtyDaysAgo).Scan(&newChildrenLast30Days)
	stats["new_children_last_30_days"] = newChildrenLast30Days

	return c.JSON(http.StatusOK, stats)
}

// GetAdminUsersAnalytics returns detailed user statistics
func GetAdminUsersAnalytics(c echo.Context) error {
	stats := make(map[string]interface{})

	// Users by registration date (last 12 months)
	type MonthlyUsers struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	monthlyUsers := []MonthlyUsers{}
	rows, err := db.DB.Query(`
		SELECT TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count
		FROM users
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY TO_CHAR(created_at, 'YYYY-MM')
		ORDER BY month
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get users by month: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var mu MonthlyUsers
			if rows.Scan(&mu.Month, &mu.Count) == nil {
				monthlyUsers = append(monthlyUsers, mu)
			}
		}
	}
	stats["users_by_month"] = monthlyUsers

	// Users with most children
	type TopUsers struct {
		UserID   string `json:"user_id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Count    int    `json:"count"`
	}
	topUsers := []TopUsers{}
	rows, err = db.DB.Query(`
		SELECT u.id, u.full_name, u.email, COUNT(c.id) as count
		FROM users u
		LEFT JOIN children c ON c.parent_id = u.id
		WHERE u.role = 'parent'
		GROUP BY u.id, u.full_name, u.email
		ORDER BY count DESC
		LIMIT 10
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get top users: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var tu TopUsers
			var email sql.NullString
			if rows.Scan(&tu.UserID, &tu.FullName, &email, &tu.Count) == nil {
				if email.Valid {
					tu.Email = email.String
				}
				topUsers = append(topUsers, tu)
			}
		}
	}
	stats["top_users_by_children"] = topUsers

	// Average children per user
	var avgChildrenPerUser float64
	err = db.DB.QueryRow("SELECT COALESCE(AVG(child_count), 0) FROM (SELECT COUNT(*) as child_count FROM children GROUP BY parent_id) as counts").Scan(&avgChildrenPerUser)
	if err != nil {
		c.Logger().Errorf("Failed to get avg children per user: %v", err)
		avgChildrenPerUser = 0
	}
	stats["avg_children_per_user"] = avgChildrenPerUser

	return c.JSON(http.StatusOK, stats)
}

// GetAdminChildrenAnalytics returns detailed children statistics
func GetAdminChildrenAnalytics(c echo.Context) error {
	stats := make(map[string]interface{})

	// Children by age group
	type AgeGroup struct {
		AgeGroup string `json:"age_group"`
		Count    int    `json:"count"`
	}
	ageGroups := []AgeGroup{}
	rows, err := db.DB.Query(`
		SELECT 
			CASE 
				WHEN EXTRACT(MONTH FROM AGE(CURRENT_DATE, dob)) < 6 THEN '0-5 months'
				WHEN EXTRACT(MONTH FROM AGE(CURRENT_DATE, dob)) < 12 THEN '6-11 months'
				WHEN EXTRACT(MONTH FROM AGE(CURRENT_DATE, dob)) < 24 THEN '12-23 months'
				WHEN EXTRACT(MONTH FROM AGE(CURRENT_DATE, dob)) < 36 THEN '24-35 months'
				ELSE '36+ months'
			END as age_group,
			COUNT(*) as count
		FROM children
		GROUP BY age_group
		ORDER BY MIN(EXTRACT(MONTH FROM AGE(CURRENT_DATE, dob)))
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get children by age group: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ag AgeGroup
			if rows.Scan(&ag.AgeGroup, &ag.Count) == nil {
				ageGroups = append(ageGroups, ag)
			}
		}
	}
	stats["children_by_age_group"] = ageGroups

	// Children registration trend (last 12 months)
	type MonthlyChildren struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	monthlyChildren := []MonthlyChildren{}
	rows, err = db.DB.Query(`
		SELECT TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count
		FROM children
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY TO_CHAR(created_at, 'YYYY-MM')
		ORDER BY month
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get children by month: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var mc MonthlyChildren
			if rows.Scan(&mc.Month, &mc.Count) == nil {
				monthlyChildren = append(monthlyChildren, mc)
			}
		}
	}
	stats["children_by_month"] = monthlyChildren

	return c.JSON(http.StatusOK, stats)
}

// GetAdminMeasurementsAnalytics returns detailed measurements statistics
func GetAdminMeasurementsAnalytics(c echo.Context) error {
	stats := make(map[string]interface{})

	// Measurements by month (last 12 months)
	type MonthlyMeasurements struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	monthlyMeasurements := []MonthlyMeasurements{}
	rows, err := db.DB.Query(`
		SELECT TO_CHAR(measurement_date, 'YYYY-MM') as month, COUNT(*) as count
		FROM measurements
		WHERE measurement_date >= CURRENT_DATE - INTERVAL '12 months'
		GROUP BY TO_CHAR(measurement_date, 'YYYY-MM')
		ORDER BY month
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get measurements by month: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var mm MonthlyMeasurements
			if rows.Scan(&mm.Month, &mm.Count) == nil {
				monthlyMeasurements = append(monthlyMeasurements, mm)
			}
		}
	}
	stats["measurements_by_month"] = monthlyMeasurements

	// Average measurements per child
	var avgMeasurementsPerChild float64
	err = db.DB.QueryRow("SELECT COALESCE(AVG(measurement_count), 0) FROM (SELECT COUNT(*) as measurement_count FROM measurements GROUP BY child_id) as counts").Scan(&avgMeasurementsPerChild)
	if err != nil {
		c.Logger().Errorf("Failed to get avg measurements per child: %v", err)
		avgMeasurementsPerChild = 0
	}
	stats["avg_measurements_per_child"] = avgMeasurementsPerChild

	// Growth status distribution
	type GrowthStatus struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	growthStatuses := []GrowthStatus{}
	rows, err = db.DB.Query(`
		SELECT 
			COALESCE(weight_status, 'unknown') as status,
			COUNT(*) as count
		FROM measurements
		WHERE weight_status IS NOT NULL
		GROUP BY weight_status
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get growth status distribution: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var gs GrowthStatus
			if rows.Scan(&gs.Status, &gs.Count) == nil {
				growthStatuses = append(growthStatuses, gs)
			}
		}
	}
	stats["growth_status_distribution"] = growthStatuses

	return c.JSON(http.StatusOK, stats)
}

// GetAdminAssessmentsAnalytics returns detailed assessments statistics
func GetAdminAssessmentsAnalytics(c echo.Context) error {
	stats := make(map[string]interface{})

	// Assessments by month (last 12 months)
	type MonthlyAssessments struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	monthlyAssessments := []MonthlyAssessments{}
	rows, err := db.DB.Query(`
		SELECT TO_CHAR(assessment_date, 'YYYY-MM') as month, COUNT(DISTINCT assessment_date) as count
		FROM assessments
		WHERE assessment_date >= CURRENT_DATE - INTERVAL '12 months'
		GROUP BY TO_CHAR(assessment_date, 'YYYY-MM')
		ORDER BY month
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get assessments by month: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ma MonthlyAssessments
			if rows.Scan(&ma.Month, &ma.Count) == nil {
				monthlyAssessments = append(monthlyAssessments, ma)
			}
		}
	}
	stats["assessments_by_month"] = monthlyAssessments

	// Red flags detected
	var redFlagsCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM assessments a JOIN milestones m ON m.id = a.milestone_id WHERE m.is_red_flag = true AND a.status = 'no'").Scan(&redFlagsCount)
	if err != nil {
		c.Logger().Errorf("Failed to get red flags count: %v", err)
		redFlagsCount = 0
	}
	stats["red_flags_detected"] = redFlagsCount

	return c.JSON(http.StatusOK, stats)
}

// GetAdminImmunizationsAnalytics returns detailed immunizations statistics
func GetAdminImmunizationsAnalytics(c echo.Context) error {
	stats := make(map[string]interface{})

	// Immunizations by month (last 12 months)
	type MonthlyImmunizations struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	monthlyImmunizations := []MonthlyImmunizations{}
	rows, err := db.DB.Query(`
		SELECT TO_CHAR(given_date, 'YYYY-MM') as month, COUNT(*) as count
		FROM child_immunizations
		WHERE given_date >= CURRENT_DATE - INTERVAL '12 months'
		GROUP BY TO_CHAR(given_date, 'YYYY-MM')
		ORDER BY month
	`)
	if err != nil {
		c.Logger().Errorf("Failed to get immunizations by month: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var mi MonthlyImmunizations
			if rows.Scan(&mi.Month, &mi.Count) == nil {
				monthlyImmunizations = append(monthlyImmunizations, mi)
			}
		}
	}
	stats["immunizations_by_month"] = monthlyImmunizations

	// On-schedule vs catch-up
	var onScheduleCount int
	var catchUpCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM child_immunizations WHERE is_catch_up = false").Scan(&onScheduleCount)
	if err != nil {
		c.Logger().Errorf("Failed to get on-schedule count: %v", err)
		onScheduleCount = 0
	}
	err = db.DB.QueryRow("SELECT COUNT(*) FROM child_immunizations WHERE is_catch_up = true").Scan(&catchUpCount)
	if err != nil {
		c.Logger().Errorf("Failed to get catch-up count: %v", err)
		catchUpCount = 0
	}
	stats["on_schedule_count"] = onScheduleCount
	stats["catch_up_count"] = catchUpCount

	return c.JSON(http.StatusOK, stats)
}

