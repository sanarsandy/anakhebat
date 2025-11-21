package utils

import (
	"math"

	"github.com/jmoiron/sqlx"
)

// CalculateZScore calculates Z-score using WHO LMS method
// Formula: Z = ((value/M)^L - 1) / (L * S)
// Special case: if L = 0, use Z = ln(value/M) / S
func CalculateZScore(value, l, m, s float64) float64 {
	if l == 0 || math.Abs(l) < 0.0001 {
		// When L is close to 0, use logarithmic formula
		return math.Log(value/m) / s
	}
	return (math.Pow(value/m, l) - 1) / (l * s)
}

// GetWHOStandard fetches LMS values for a specific indicator
func GetWHOStandard(db *sqlx.DB, indicator, gender string, ageMonths int, heightCm *float64) (*WHOStandard, error) {
	var std WHOStandard

	if heightCm != nil {
		// For weight-for-height, use height
		// Round to nearest 0.5cm for lookup
		roundedHeight := math.Round(*heightCm*2) / 2
		err := db.Get(&std, `
			SELECT * FROM who_standards 
			WHERE indicator = $1 AND gender = $2 AND height_cm = $3
		`, indicator, gender, roundedHeight)
		return &std, err
	}

	// For age-based indicators
	err := db.Get(&std, `
		SELECT * FROM who_standards 
		WHERE indicator = $1 AND gender = $2 AND age_months = $3
	`, indicator, gender, ageMonths)
	return &std, err
}

// ZScoreResult contains all calculated Z-scores
type ZScoreResult struct {
	WeightForAge       float64
	HeightForAge       float64
	WeightForHeight    float64
	HeadCircumference  float64
	HasWeightForAge    bool
	HasHeightForAge    bool
	HasWeightForHeight bool
	HasHeadCirc        bool
}

// CalculateAllZScores calculates all applicable Z-scores for a measurement
func CalculateAllZScores(db *sqlx.DB, gender string, ageMonths int, weight, height, headCirc float64) (*ZScoreResult, error) {
	result := &ZScoreResult{}

	// Normalize gender to match database values
	originalGender := gender
	if gender == "L" || gender == "laki-laki" {
		gender = "male"
	} else if gender == "P" || gender == "perempuan" {
		gender = "female"
	}

	// Weight-for-age
	wfaStd, err := GetWHOStandard(db, "wfa", gender, ageMonths, nil)
	if err == nil {
		result.WeightForAge = CalculateZScore(weight, wfaStd.L, wfaStd.M, wfaStd.S)
		result.HasWeightForAge = true
	} else {
		// Log error for debugging
		println("ERROR: Failed to get WFA standard - gender:", originalGender, "->", gender, "age:", ageMonths, "error:", err.Error())
	}

	// Height-for-age
	hfaStd, err := GetWHOStandard(db, "hfa", gender, ageMonths, nil)
	if err == nil {
		result.HeightForAge = CalculateZScore(height, hfaStd.L, hfaStd.M, hfaStd.S)
		result.HasHeightForAge = true
	} else {
		println("ERROR: Failed to get HFA standard - gender:", originalGender, "->", gender, "age:", ageMonths, "error:", err.Error())
	}

	// Weight-for-height (if we have the data)
	if height > 0 {
		wfhStd, err := GetWHOStandard(db, "wfh", gender, 0, &height)
		if err == nil {
			result.WeightForHeight = CalculateZScore(weight, wfhStd.L, wfhStd.M, wfhStd.S)
			result.HasWeightForHeight = true
		} else {
			println("ERROR: Failed to get WFH standard - gender:", originalGender, "->", gender, "height:", height, "error:", err.Error())
		}
	}

	// Head circumference-for-age (if provided)
	if headCirc > 0 {
		hcfaStd, err := GetWHOStandard(db, "hcfa", gender, ageMonths, nil)
		if err == nil {
			result.HeadCircumference = CalculateZScore(headCirc, hcfaStd.L, hcfaStd.M, hcfaStd.S)
			result.HasHeadCirc = true
		} else {
			println("ERROR: Failed to get HCFA standard - gender:", originalGender, "->", gender, "age:", ageMonths, "error:", err.Error())
		}
	}

	return result, nil
}

// InterpretNutritionalStatus interprets Z-scores into human-readable status
func InterpretNutritionalStatus(wfaZ, hfaZ, wfhZ float64) (string, string, string) {
	// Weight-for-age interpretation
	var weightStatus string
	switch {
	case wfaZ < -3:
		weightStatus = "Severely Underweight / Gizi Buruk"
	case wfaZ < -2:
		weightStatus = "Underweight / Gizi Kurang"
	case wfaZ < -1:
		weightStatus = "Possible Risk of Underweight / Berisiko Gizi Kurang"
	case wfaZ <= 1:
		weightStatus = "Normal Weight / Gizi Baik"
	case wfaZ <= 2:
		weightStatus = "Possible Risk of Overweight / Berisiko Gizi Lebih"
	default:
		weightStatus = "Overweight / Gizi Lebih"
	}

	// Height-for-age interpretation (Stunting)
	var heightStatus string
	switch {
	case hfaZ < -3:
		heightStatus = "Severely Stunted / Sangat Pendek (Stunting Berat)"
	case hfaZ < -2:
		heightStatus = "Stunted / Pendek (Stunting)"
	case hfaZ < -1:
		heightStatus = "Possible Risk of Stunting / Berisiko Pendek"
	case hfaZ <= 3:
		heightStatus = "Normal Height / Tinggi Normal"
	default:
		heightStatus = "Tall / Tinggi"
	}

	// Weight-for-height interpretation (Wasting/Overweight)
	var wfhStatus string
	switch {
	case wfhZ < -3:
		wfhStatus = "Severely Wasted / Sangat Kurus (Gizi Buruk)"
	case wfhZ < -2:
		wfhStatus = "Wasted / Kurus (Gizi Kurang)"
	case wfhZ < -1:
		wfhStatus = "Possible Risk of Wasting / Berisiko Kurus"
	case wfhZ <= 1:
		wfhStatus = "Normal / Normal"
	case wfhZ <= 2:
		wfhStatus = "Possible Risk of Overweight / Berisiko Gemuk"
	case wfhZ <= 3:
		wfhStatus = "Overweight / Gemuk"
	default:
		wfhStatus = "Obese / Obesitas"
	}

	return weightStatus, heightStatus, wfhStatus
}
