package handlers

import (
	"net/http"
	"strconv"
	"tukem-backend/db"

	"github.com/labstack/echo/v4"
)

// GetWHOStandardsForChart retrieves WHO standards data for plotting growth curves
// Query params: indicator (wfa/hfa/wfh), gender (male/female), minAge, maxAge
func GetWHOStandardsForChart(c echo.Context) error {
	// This endpoint is public (no auth required) as it only returns WHO standard data

	indicator := c.QueryParam("indicator")
	gender := c.QueryParam("gender")
	minAgeStr := c.QueryParam("minAge")
	maxAgeStr := c.QueryParam("maxAge")

	// Validate required params
	if indicator == "" || gender == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "indicator and gender are required",
		})
	}

	// Normalize gender
	if gender == "L" || gender == "laki-laki" {
		gender = "male"
	} else if gender == "P" || gender == "perempuan" {
		gender = "female"
	}

	// For weight-for-height, use height range instead of age
	if indicator == "wfh" {
		// Parse height range (default: 45-120 cm for children)
		minHeight := 45.0
		maxHeight := 120.0
		if minAgeStr != "" {
			if parsed, err := strconv.ParseFloat(minAgeStr, 64); err == nil {
				minHeight = parsed
			}
		}
		if maxAgeStr != "" {
			if parsed, err := strconv.ParseFloat(maxAgeStr, 64); err == nil {
				maxHeight = parsed
			}
		}

		// Query WHO standards for weight-for-height (uses height_cm)
		query := `
			SELECT 
				height_cm as x_value,
				sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3
			FROM who_standards
			WHERE indicator = $1 AND gender = $2 
				AND height_cm >= $3 AND height_cm <= $4
				AND height_cm IS NOT NULL
			ORDER BY height_cm ASC
		`

		type StandardPointWFH struct {
			XValue  float64 `json:"x_value" db:"x_value"`
			SD3Neg  float64 `json:"sd3neg" db:"sd3neg"`
			SD2Neg  float64 `json:"sd2neg" db:"sd2neg"`
			SD1Neg  float64 `json:"sd1neg" db:"sd1neg"`
			SD0     float64 `json:"sd0" db:"sd0"`
			SD1     float64 `json:"sd1" db:"sd1"`
			SD2     float64 `json:"sd2" db:"sd2"`
			SD3     float64 `json:"sd3" db:"sd3"`
		}

		var standards []StandardPointWFH
		err := db.DB.Select(&standards, query, indicator, gender, minHeight, maxHeight)
		if err != nil {
			c.Logger().Errorf("Failed to fetch WHO standards: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Gagal mengambil data standar WHO",
			})
		}

		// Convert to common format
		result := make([]map[string]interface{}, len(standards))
		for i, s := range standards {
			result[i] = map[string]interface{}{
				"x_value": s.XValue,
				"sd3neg":  s.SD3Neg,
				"sd2neg":  s.SD2Neg,
				"sd1neg":  s.SD1Neg,
				"sd0":     s.SD0,
				"sd1":     s.SD1,
				"sd2":     s.SD2,
				"sd3":     s.SD3,
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"indicator": indicator,
			"gender":    gender,
			"standards": result,
		})
	}

	// For age-based indicators (wfa, hfa)
	// Parse age range (default: 0-60 months)
	minAge := 0
	maxAge := 60
	if minAgeStr != "" {
		if parsed, err := strconv.Atoi(minAgeStr); err == nil {
			minAge = parsed
		}
	}
	if maxAgeStr != "" {
		if parsed, err := strconv.Atoi(maxAgeStr); err == nil {
			maxAge = parsed
		}
	}

	// Query WHO standards for age-based indicators
	query := `
		SELECT 
			age_months as x_value,
			sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3
		FROM who_standards
		WHERE indicator = $1 AND gender = $2 
			AND age_months >= $3 AND age_months <= $4
			AND age_months IS NOT NULL
		ORDER BY age_months ASC
	`

	type StandardPoint struct {
		XValue  int     `json:"x_value" db:"x_value"`
		SD3Neg  float64 `json:"sd3neg" db:"sd3neg"`
		SD2Neg  float64 `json:"sd2neg" db:"sd2neg"`
		SD1Neg  float64 `json:"sd1neg" db:"sd1neg"`
		SD0     float64 `json:"sd0" db:"sd0"`
		SD1     float64 `json:"sd1" db:"sd1"`
		SD2     float64 `json:"sd2" db:"sd2"`
		SD3     float64 `json:"sd3" db:"sd3"`
	}

	var standards []StandardPoint
	err := db.DB.Select(&standards, query, indicator, gender, minAge, maxAge)
	if err != nil {
		c.Logger().Errorf("Failed to fetch WHO standards: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data standar WHO",
		})
	}

	// Convert to common format
	result := make([]map[string]interface{}, len(standards))
	for i, s := range standards {
		result[i] = map[string]interface{}{
			"x_value": s.XValue,
			"sd3neg":  s.SD3Neg,
			"sd2neg":  s.SD2Neg,
			"sd1neg":  s.SD1Neg,
			"sd0":     s.SD0,
			"sd1":     s.SD1,
			"sd2":     s.SD2,
			"sd3":     s.SD3,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"indicator": indicator,
		"gender":    gender,
		"standards": result,
	})
}

