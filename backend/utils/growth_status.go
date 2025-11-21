package utils

// InterpretWeightStatus interprets weight-for-age Z-score
func InterpretWeightStatus(zscore float64) string {
	if zscore < -3 {
		return "Severely Underweight"
	} else if zscore < -2 {
		return "Underweight"
	} else if zscore <= 2 {
		return "Normal Weight"
	} else if zscore <= 3 {
		return "Possible Risk of Overweight"
	} else {
		return "Overweight"
	}
}

// InterpretHeightStatus interprets height-for-age Z-score
func InterpretHeightStatus(zscore float64) string {
	if zscore < -3 {
		return "Severely Stunted"
	} else if zscore < -2 {
		return "Stunted"
	} else if zscore <= 2 {
		return "Normal Height"
	} else {
		return "Tall"
	}
}

// GetStatusColor returns color code for status display
func GetStatusColor(zscore float64) string {
	if zscore < -2 || zscore > 2 {
		return "danger" // Red
	} else if zscore < -1 || zscore > 1 {
		return "warning" // Amber
	}
	return "success" // Green
}
