package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
)

// WHOStandard represents a WHO growth standard entry
type WHOStandard struct {
	ID        string   `json:"-" db:"id"`  // Database ID (not in JSON)
	Indicator string   `json:"indicator" db:"indicator"`
	Gender    string   `json:"gender" db:"gender"`
	AgeMonths *int     `json:"age_months,omitempty" db:"age_months"`
	HeightCm  *float64 `json:"height_cm,omitempty" db:"height_cm"`
	L         float64  `json:"l" db:"l_value"`
	M         float64  `json:"m" db:"m_value"`
	S         float64  `json:"s" db:"s_value"`
	SD3Neg    float64  `json:"sd3neg" db:"sd3neg"`
	SD2Neg    float64  `json:"sd2neg" db:"sd2neg"`
	SD1Neg    float64  `json:"sd1neg" db:"sd1neg"`
	SD0       float64  `json:"sd0" db:"sd0"`
	SD1       float64  `json:"sd1" db:"sd1"`
	SD2       float64  `json:"sd2" db:"sd2"`
	SD3       float64  `json:"sd3" db:"sd3"`
	CreatedAt string   `json:"-" db:"created_at"`  // Database timestamp (not in JSON)
}

// RawWHOData represents the raw JSON format from WHO files
type RawWHOData struct {
	Month   *int     `json:"month"`
	AgeMonths *int   `json:"age_months"`
	L       float64  `json:"L"`
	M       float64  `json:"M"`
	S       float64  `json:"S"`
	LValue  float64  `json:"l"`
	MValue  float64  `json:"m"`
	SValue  float64  `json:"s"`
}

// SeedWHOStandards loads WHO growth standards data into the database
func SeedWHOStandards(db *sqlx.DB) error {
	// Check if both WFA and HFA are already seeded
	var wfaCount, hfaCount int
	err := db.Get(&wfaCount, "SELECT COUNT(*) FROM who_standards WHERE indicator = 'wfa'")
	if err != nil {
		return err
	}
	err = db.Get(&hfaCount, "SELECT COUNT(*) FROM who_standards WHERE indicator = 'hfa'")
	if err != nil {
		return err
	}

	// Only skip if we have substantial data for both indicators (more than 50 records each)
	if wfaCount > 50 && hfaCount > 50 {
		log.Printf("WHO standards already seeded (WFA: %d, HFA: %d entries), skipping...", wfaCount, hfaCount)
		return nil
	}
	
	// If HFA is missing, we need to seed it
	if hfaCount == 0 {
		log.Printf("HFA data missing, will seed HFA data...")
	}

	log.Println("Seeding WHO standards...")

	// List of data files to seed - using full data files
	fileConfigs := []struct {
		file     string
		indicator string
		gender   string
	}{
		{"data/who/wfa_boys_0_60.json", "wfa", "male"},
		{"data/who/wfa_girls_0_60.json", "wfa", "female"},
		{"data/who/hfa_boys_0_60.json", "hfa", "male"},
		{"data/who/hfa_girls_0_60.json", "hfa", "female"},
		// Fallback to sample files if full files don't exist
		{"data/who/wfa_sample.json", "wfa", "male"},
		{"data/who/hfa_sample.json", "hfa", "male"},
	}

	totalSeeded := 0
	seededFiles := make(map[string]bool)

	for _, config := range fileConfigs {
		// Skip if we already seeded this file type
		key := config.indicator + "_" + config.gender
		if seededFiles[key] {
			continue
		}

		data, err := ioutil.ReadFile(config.file)
		if err != nil {
			log.Printf("Warning: Could not read %s: %v", config.file, err)
			continue
		}

		// Try to parse as array of WHOStandard first (sample format)
		var standards []WHOStandard
		if err := json.Unmarshal(data, &standards); err == nil {
			// Sample format - check if it has indicator field
			if len(standards) > 0 && standards[0].Indicator != "" {
				// This is the sample format with indicator already set
				for _, std := range standards {
					_, err := db.NamedExec(`
						INSERT INTO who_standards 
						(indicator, gender, age_months, height_cm, l_value, m_value, s_value,
						 sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3)
						VALUES 
						(:indicator, :gender, :age_months, :height_cm, :l_value, :m_value, :s_value,
						 :sd3neg, :sd2neg, :sd1neg, :sd0, :sd1, :sd2, :sd3)
						ON CONFLICT (indicator, gender, age_months, height_cm) DO NOTHING
					`, std)
					if err != nil {
						log.Printf("Warning: Failed to insert WHO standard: %v", err)
						continue
					}
					totalSeeded++
				}
				log.Printf("Seeded %s successfully (%d entries)", config.file, len(standards))
				seededFiles[key] = true
				continue
			}
		}

		// Try raw format (wfa_boys_0_60.json format with month, L, M, S)
		var rawData []RawWHOData
		if err := json.Unmarshal(data, &rawData); err == nil && len(rawData) > 0 {
			for _, raw := range rawData {
				// Determine age_months
				var ageMonths *int
				if raw.Month != nil {
					ageMonths = raw.Month
				} else if raw.AgeMonths != nil {
					ageMonths = raw.AgeMonths
				}

				// Determine L, M, S values
				var l, m, s float64
				if raw.L != 0 || raw.M != 0 || raw.S != 0 {
					l = raw.L
					m = raw.M
					s = raw.S
				} else {
					l = raw.LValue
					m = raw.MValue
					s = raw.SValue
				}

				if ageMonths == nil || m == 0 {
					continue
				}

				// Calculate SD values from LMS parameters
				sd3neg := CalculateSDValue(l, m, s, -3)
				sd2neg := CalculateSDValue(l, m, s, -2)
				sd1neg := CalculateSDValue(l, m, s, -1)
				sd0 := m // Median is always M
				sd1 := CalculateSDValue(l, m, s, 1)
				sd2 := CalculateSDValue(l, m, s, 2)
				sd3 := CalculateSDValue(l, m, s, 3)

				std := WHOStandard{
					Indicator: config.indicator,
					Gender:    config.gender,
					AgeMonths: ageMonths,
					L:         l,
					M:         m,
					S:         s,
					SD3Neg:    sd3neg,
					SD2Neg:    sd2neg,
					SD1Neg:    sd1neg,
					SD0:       sd0,
					SD1:       sd1,
					SD2:       sd2,
					SD3:       sd3,
				}

				_, err := db.NamedExec(`
					INSERT INTO who_standards 
					(indicator, gender, age_months, height_cm, l_value, m_value, s_value,
					 sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3)
					VALUES 
					(:indicator, :gender, :age_months, :height_cm, :l_value, :m_value, :s_value,
					 :sd3neg, :sd2neg, :sd1neg, :sd0, :sd1, :sd2, :sd3)
					ON CONFLICT (indicator, gender, age_months, height_cm) DO NOTHING
				`, std)
				if err != nil {
					log.Printf("Warning: Failed to insert WHO standard: %v", err)
					continue
				}
				totalSeeded++
			}
			log.Printf("Seeded %s successfully (%d entries)", config.file, len(rawData))
			seededFiles[key] = true
		} else {
			log.Printf("Warning: Could not parse %s in any known format: %v", config.file, err)
		}
	}

	log.Printf("WHO standards seeded successfully! Total: %d entries", totalSeeded)
	return nil
}
