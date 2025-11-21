package utils

import (
	"time"
)

// CalculateAgeInDays calculates the precise age in days between DOB and measurement date
func CalculateAgeInDays(dob string, measurementDate string) (int, error) {
	// Extract date portion if ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
	// This handles both "2024-11-20" and "2024-11-20T00:00:00Z" formats
	if len(measurementDate) > 10 {
		measurementDate = measurementDate[:10]
	}
	if len(dob) > 10 {
		dob = dob[:10]
	}

	dobTime, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0, err
	}

	measurementTime, err := time.Parse("2006-01-02", measurementDate)
	if err != nil {
		return 0, err
	}

	duration := measurementTime.Sub(dobTime)
	days := int(duration.Hours() / 24)

	return days, nil
}

// CalculateAgeInMonths calculates age in months (for WHO lookup)
func CalculateAgeInMonths(dob string, measurementDate string) (int, error) {
	// Extract date portion if ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
	// This handles both "2024-11-20" and "2024-11-20T00:00:00Z" formats
	if len(measurementDate) > 10 {
		measurementDate = measurementDate[:10]
	}
	if len(dob) > 10 {
		dob = dob[:10]
	}

	dobTime, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0, err
	}

	measurementTime, err := time.Parse("2006-01-02", measurementDate)
	if err != nil {
		return 0, err
	}

	years := measurementTime.Year() - dobTime.Year()
	months := int(measurementTime.Month() - dobTime.Month())
	
	totalMonths := years*12 + months

	// Adjust if the day hasn't been reached yet in the current month
	if measurementTime.Day() < dobTime.Day() {
		totalMonths--
	}

	return totalMonths, nil
}

// CalculateCorrectedAge calculates corrected age for premature babies
// Corrected age = Chronological age - (40 weeks - gestational age in weeks)
// This is used until the child reaches 24 months chronological age
// Returns: correctedAgeInDays, correctedAgeInMonths, shouldUseCorrected (bool)
func CalculateCorrectedAge(dob string, measurementDate string, isPremature bool, gestationalAgeWeeks *int) (int, int, bool, error) {
	if !isPremature || gestationalAgeWeeks == nil {
		// Not premature, return chronological age
		days, err := CalculateAgeInDays(dob, measurementDate)
		if err != nil {
			return 0, 0, false, err
		}
		months, err := CalculateAgeInMonths(dob, measurementDate)
		if err != nil {
			return 0, 0, false, err
		}
		return days, months, false, nil
	}

	// Calculate chronological age
	chronoDays, err := CalculateAgeInDays(dob, measurementDate)
	if err != nil {
		return 0, 0, false, err
	}
	chronoMonths, err := CalculateAgeInMonths(dob, measurementDate)
	if err != nil {
		return 0, 0, false, err
	}

	// Only use corrected age if chronological age < 24 months (730 days)
	if chronoDays >= 730 {
		return chronoDays, chronoMonths, false, nil
	}

	// Calculate weeks premature
	// Full term is 40 weeks, so weeks premature = 40 - gestational_age
	weeksPremature := 40 - *gestationalAgeWeeks
	if weeksPremature < 0 {
		weeksPremature = 0 // Safety check
	}

	// Convert weeks to days (1 week = 7 days)
	daysToSubtract := weeksPremature * 7

	// Calculate corrected age
	correctedDays := chronoDays - daysToSubtract
	if correctedDays < 0 {
		correctedDays = 0 // Safety check - can't have negative age
	}

	// Calculate corrected age in months
	// Approximate: divide days by 30.44 (average days per month)
	correctedMonths := int(float64(correctedDays) / 30.44)
	if correctedMonths < 0 {
		correctedMonths = 0
	}

	return correctedDays, correctedMonths, true, nil
}

// FormatAgeDisplay formats age for display (e.g., "2 years 3 months")
func FormatAgeDisplay(ageInMonths int) string {
	if ageInMonths < 0 {
		return "Invalid age"
	}

	years := ageInMonths / 12
	months := ageInMonths % 12

	if years == 0 {
		if months == 1 {
			return "1 month"
		}
		return string(rune(months)) + " months"
	}

	if months == 0 {
		if years == 1 {
			return "1 year"
		}
		return string(rune(years)) + " years"
	}

	yearStr := "year"
	if years > 1 {
		yearStr = "years"
	}

	monthStr := "month"
	if months > 1 {
		monthStr = "months"
	}

	return string(rune(years)) + " " + yearStr + " " + string(rune(months)) + " " + monthStr
}
