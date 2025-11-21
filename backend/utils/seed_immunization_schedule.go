package utils

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// SeedImmunizationSchedule seeds immunization schedule data based on IDAI recommendations
func SeedImmunizationSchedule(db *sqlx.DB) error {
	// Check if data already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM immunization_schedule")
	if err != nil {
		return fmt.Errorf("failed to check existing immunization schedule: %v", err)
	}

	if count > 0 {
		log.Println("Immunization schedule already seeded, skipping...")
		return nil
	}

	// Helper function to create int pointer
	intPtr := func(i int) *int { return &i }

	// Immunization schedule data based on IDAI recommendations
	// Format: name, name_id, description, age_min_days, age_optimal_days, age_max_days, age_min_months, age_optimal_months, age_max_months, dose_number, total_doses, interval_from_previous_days, category, priority, is_required
	immunizations := []struct {
		Name                      string
		NameID                    string
		Description               string
		AgeMinDays                *int
		AgeOptimalDays            *int
		AgeMaxDays                *int
		AgeMinMonths              *int
		AgeOptimalMonths          *int
		AgeMaxMonths              *int
		DoseNumber                int
		TotalDoses                *int
		IntervalFromPreviousDays  *int
		Category                  string
		Priority                  string
		IsRequired                bool
	}{
		// Hepatitis B-0 (0-7 hari)
		{"Hepatitis B-0", "Hepatitis B-0", "Hepatitis B dosis pertama diberikan saat lahir", intPtr(0), intPtr(1), intPtr(7), intPtr(0), intPtr(0), intPtr(0), 1, intPtr(4), nil, "wajib", "high", true},

		// Polio-0 (0-7 hari)
		{"Polio-0", "Polio-0", "Polio dosis pertama diberikan saat lahir", intPtr(0), intPtr(1), intPtr(7), intPtr(0), intPtr(0), intPtr(0), 1, intPtr(5), nil, "wajib", "high", true},

		// BCG (1 bulan / dapat sejak lahir)
		{"BCG", "BCG", "BCG untuk mencegah tuberkulosis", intPtr(0), intPtr(30), intPtr(60), intPtr(0), intPtr(1), intPtr(2), 1, intPtr(1), nil, "wajib", "high", true},

		// DPT-HB-Hib (2, 3, 4 bulan)
		{"DPT-HB-Hib", "DPT-HB-Hib", "Difteri, Pertusis, Tetanus, Hepatitis B, Haemophilus influenzae type B - Dosis 1", intPtr(56), intPtr(60), intPtr(150), intPtr(2), intPtr(2), intPtr(5), 1, intPtr(4), nil, "wajib", "high", true},
		{"DPT-HB-Hib", "DPT-HB-Hib", "Difteri, Pertusis, Tetanus, Hepatitis B, Haemophilus influenzae type B - Dosis 2", intPtr(84), intPtr(90), intPtr(180), intPtr(3), intPtr(3), intPtr(6), 2, intPtr(4), intPtr(28), "wajib", "high", true},
		{"DPT-HB-Hib", "DPT-HB-Hib", "Difteri, Pertusis, Tetanus, Hepatitis B, Haemophilus influenzae type B - Dosis 3", intPtr(112), intPtr(120), intPtr(210), intPtr(4), intPtr(4), intPtr(7), 3, intPtr(4), intPtr(28), "wajib", "high", true},
		{"DPT-HB-Hib", "DPT-HB-Hib", "Difteri, Pertusis, Tetanus, Hepatitis B, Haemophilus influenzae type B - Booster", intPtr(540), intPtr(540), intPtr(720), intPtr(18), intPtr(18), intPtr(24), 4, intPtr(4), intPtr(365), "wajib", "high", true},

		// Polio (2, 3, 4 bulan)
		{"Polio", "Polio", "Polio - Dosis 1", intPtr(56), intPtr(60), intPtr(150), intPtr(2), intPtr(2), intPtr(5), 1, intPtr(5), nil, "wajib", "high", true},
		{"Polio", "Polio", "Polio - Dosis 2", intPtr(84), intPtr(90), intPtr(180), intPtr(3), intPtr(3), intPtr(6), 2, intPtr(5), intPtr(28), "wajib", "high", true},
		{"Polio", "Polio", "Polio - Dosis 3", intPtr(112), intPtr(120), intPtr(210), intPtr(4), intPtr(4), intPtr(7), 3, intPtr(5), intPtr(28), "wajib", "high", true},
		{"Polio", "Polio", "Polio - Booster", intPtr(540), intPtr(540), intPtr(720), intPtr(18), intPtr(18), intPtr(24), 4, intPtr(5), intPtr(365), "wajib", "high", true},

		// IPV (4 bulan)
		{"IPV", "IPV", "Inactivated Polio Vaccine", intPtr(112), intPtr(120), intPtr(210), intPtr(4), intPtr(4), intPtr(7), 1, intPtr(1), nil, "wajib", "medium", true},

		// Campak/MR (9 bulan)
		{"Campak/MR", "Campak/MR", "Campak dan Rubella", intPtr(270), intPtr(270), intPtr(365), intPtr(9), intPtr(9), intPtr(12), 1, intPtr(2), nil, "wajib", "high", true},

		// MR Booster (18 bulan)
		{"MR Booster", "MR Booster", "Campak dan Rubella - Booster", intPtr(540), intPtr(540), intPtr(720), intPtr(18), intPtr(18), intPtr(24), 2, intPtr(2), intPtr(270), "wajib", "high", true},

		// DPT Booster (5-7 tahun / sebelum SD)
		{"DPT Booster", "DPT Booster", "Difteri, Pertusis, Tetanus - Booster sebelum SD", intPtr(1825), intPtr(2190), intPtr(2555), intPtr(60), intPtr(72), intPtr(84), 5, intPtr(5), intPtr(365), "wajib", "medium", true},
	}

	// Insert immunizations
	query := `
		INSERT INTO immunization_schedule (
			name, name_id, description,
			age_min_days, age_optimal_days, age_max_days,
			age_min_months, age_optimal_months, age_max_months,
			dose_number, total_doses,
			interval_from_previous_days,
			category, priority, is_required,
			source, is_active
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11,
			$12,
			$13, $14, $15,
			'IDAI', true
		)
	`

	for _, imm := range immunizations {
		_, err := db.Exec(query,
			imm.Name,
			imm.NameID,
			imm.Description,
			imm.AgeMinDays,
			imm.AgeOptimalDays,
			imm.AgeMaxDays,
			imm.AgeMinMonths,
			imm.AgeOptimalMonths,
			imm.AgeMaxMonths,
			imm.DoseNumber,
			imm.TotalDoses,
			imm.IntervalFromPreviousDays,
			imm.Category,
			imm.Priority,
			imm.IsRequired,
		)
		if err != nil {
			log.Printf("Failed to insert immunization '%s' dose %d: %v", imm.Name, imm.DoseNumber, err)
			continue
		}
	}

	log.Printf("Immunization schedule seeded successfully! (%d items)", len(immunizations))
	return nil
}

