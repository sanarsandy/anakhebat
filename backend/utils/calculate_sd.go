package utils

import (
	"math"

	"github.com/jmoiron/sqlx"
)

// CalculateSDValue calculates SD value from LMS parameters
// Formula: SD = M * (1 + L * S * z)^(1/L) if L != 0
//          SD = M * exp(S * z) if L = 0
func CalculateSDValue(l, m, s, z float64) float64 {
	if math.Abs(l) < 0.0001 {
		// When L is close to 0, use exponential formula
		return m * math.Exp(s*z)
	}
	// Standard LMS formula
	return m * math.Pow(1+l*s*z, 1.0/l)
}

// UpdateSDValues calculates and updates SD values for all WHO standards
func UpdateSDValues(db *sqlx.DB) error {
	// Get all WHO standards that need SD values calculated
	rows, err := db.Query(`
		SELECT id, l_value, m_value, s_value 
		FROM who_standards 
		WHERE (sd3neg = 0 OR sd3neg IS NULL)
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	type WHOStandardRow struct {
		ID     string
		LValue float64
		MValue float64
		SValue float64
	}

	var standards []WHOStandardRow
	for rows.Next() {
		var std WHOStandardRow
		if err := rows.Scan(&std.ID, &std.LValue, &std.MValue, &std.SValue); err != nil {
			continue
		}
		standards = append(standards, std)
	}

	// Update each standard with calculated SD values
	for _, std := range standards {
		sd3neg := CalculateSDValue(std.LValue, std.MValue, std.SValue, -3)
		sd2neg := CalculateSDValue(std.LValue, std.MValue, std.SValue, -2)
		sd1neg := CalculateSDValue(std.LValue, std.MValue, std.SValue, -1)
		sd0 := std.MValue // Median is always M
		sd1 := CalculateSDValue(std.LValue, std.MValue, std.SValue, 1)
		sd2 := CalculateSDValue(std.LValue, std.MValue, std.SValue, 2)
		sd3 := CalculateSDValue(std.LValue, std.MValue, std.SValue, 3)

		_, err := db.Exec(`
			UPDATE who_standards 
			SET sd3neg = $1, sd2neg = $2, sd1neg = $3, sd0 = $4, sd1 = $5, sd2 = $6, sd3 = $7
			WHERE id = $8
		`, sd3neg, sd2neg, sd1neg, sd0, sd1, sd2, sd3, std.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

