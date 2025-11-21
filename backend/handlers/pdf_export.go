package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"tukem-backend/db"
	"tukem-backend/models"
	"tukem-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/labstack/echo/v4"
)

// ExportChildReport generates a PDF report for a child
func ExportChildReport(c echo.Context) error {
	childID := c.Param("id")
	c.Logger().Infof("PDF Export requested for child ID: %s", childID)

	// Get user ID from JWT to verify ownership
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	// Verify child belongs to user
	var parentID string
	err := db.DB.QueryRow("SELECT parent_id FROM children WHERE id = $1", childID).Scan(&parentID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}
	if err != nil {
		c.Logger().Errorf("Failed to verify child ownership: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify child ownership"})
	}
	if parentID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
	}

	// Get child data
	var child models.Child
	err = db.DB.QueryRow("SELECT id, name, dob, gender, birth_weight, birth_height, is_premature, gestational_age FROM children WHERE id = $1", childID).
		Scan(&child.ID, &child.Name, &child.DOB, &child.Gender, &child.BirthWeight, &child.BirthHeight, &child.IsPremature, &child.GestationalAge)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get child data"})
	}

	// Get measurements (get all, not limited)
	measurements, err := getMeasurementsForReport(childID)
	if err != nil {
		c.Logger().Errorf("Failed to get measurements: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get measurements"})
	}
	
	// If no measurements, still generate PDF with basic info
	if len(measurements) == 0 {
		c.Logger().Warnf("No measurements found for child %s", childID)
	}

	// Get assessment summary
	summary, err := getAssessmentSummaryForReport(childID)
	if err != nil {
		c.Logger().Errorf("Failed to get assessment summary: %v", err)
		// Continue even if summary fails
		summary = &models.AssessmentSummary{}
	}

	// Generate PDF
	pdf := generatePDFReport(child, measurements, summary)

	// Set response headers
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=tukem_report_%s_%s.pdf", childID[:8], time.Now().Format("20060102")))

	// Write PDF to response
	err = pdf.Output(c.Response().Writer)
	if err != nil {
		c.Logger().Errorf("Failed to write PDF: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	return nil
}

func getMeasurementsForReport(childID string) ([]models.MeasurementResponse, error) {
	query := `SELECT id, child_id, measurement_date, weight, height, head_circumference, 
		age_in_days, age_in_months, weight_for_age_zscore, height_for_age_zscore, 
		weight_for_height_zscore, head_circumference_zscore,
		nutritional_status, height_status, weight_for_height_status, created_at 
		FROM measurements WHERE child_id = $1 ORDER BY measurement_date DESC`

	rows, err := db.DB.Query(query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []models.MeasurementResponse
	for rows.Next() {
		var m models.Measurement
		var nutritionalStatus, heightStatus, wfhStatus sql.NullString
		var wfhZScore, hcZScore sql.NullFloat64

		err := rows.Scan(&m.ID, &m.ChildID, &m.MeasurementDate, &m.Weight, &m.Height, &m.HeadCircumference,
			&m.AgeInDays, &m.AgeInMonths, &m.WeightForAgeZScore, &m.HeightForAgeZScore,
			&wfhZScore, &hcZScore,
			&nutritionalStatus, &heightStatus, &wfhStatus, &m.CreatedAt)
		if err != nil {
			continue
		}

		var wfhZPtr, hcZPtr *float64
		if wfhZScore.Valid {
			wfhZPtr = &wfhZScore.Float64
		}
		if hcZScore.Valid {
			hcZPtr = &hcZScore.Float64
		}

		response := models.MeasurementResponse{
			ID:                      m.ID,
			ChildID:                 m.ChildID,
			MeasurementDate:         m.MeasurementDate,
			Weight:                  m.Weight,
			Height:                  m.Height,
			HeadCircumference:       m.HeadCircumference,
			AgeInDays:               m.AgeInDays,
			AgeInMonths:             m.AgeInMonths,
			AgeDisplay:              utils.FormatAgeDisplay(m.AgeInMonths),
			WeightForAgeZScore:      m.WeightForAgeZScore,
			HeightForAgeZScore:      m.HeightForAgeZScore,
			WeightForHeightZScore:   wfhZPtr,
			HeadCircumferenceZScore: hcZPtr,
			NutritionalStatus:       nutritionalStatus.String,
			HeightStatus:            heightStatus.String,
			WeightForHeightStatus:   wfhStatus.String,
			CreatedAt:               m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		measurements = append(measurements, response)
	}

	return measurements, nil
}

func getAssessmentSummaryForReport(childID string) (*models.AssessmentSummary, error) {
	// Fetch all assessments for this child joined with milestones
	query := `
		SELECT a.status, m.category, m.pyramid_level, m.is_red_flag, m.question
		FROM assessments a
		JOIN milestones m ON a.milestone_id = m.id
		WHERE a.child_id = $1
			AND m.source = 'KPSP'
	`

	rows, err := db.DB.Queryx(query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type AssessmentData struct {
		Status       string
		Category     string
		PyramidLevel int
		IsRedFlag    bool
		Question     string
	}

	var data []AssessmentData
	for rows.Next() {
		var d AssessmentData
		if err := rows.Scan(&d.Status, &d.Category, &d.PyramidLevel, &d.IsRedFlag, &d.Question); err != nil {
			continue
		}
		data = append(data, d)
	}

	// Calculate scores
	totalByLevel := make(map[int]int)
	completedByLevel := make(map[int]int)
	redFlags := []models.Milestone{}

	for _, d := range data {
		totalByLevel[d.PyramidLevel]++
		if d.Status == "yes" {
			completedByLevel[d.PyramidLevel]++
		}

		if d.IsRedFlag && d.Status == "no" {
			redFlags = append(redFlags, models.Milestone{
				Question: d.Question,
				Category: d.Category,
			})
		}
	}

	// Calculate percentages
	progressByCategory := make(map[string]float64)
	levelToCategory := map[int]string{
		1: "sensory",
		2: "motor",
		3: "perception",
		4: "cognitive",
	}

	for level, total := range totalByLevel {
		if total > 0 {
			cat := levelToCategory[level]
			progressByCategory[cat] = float64(completedByLevel[level]) / float64(total) * 100
		}
	}

	// Logic Warnings (Pyramid Imbalance)
	warnings := []string{}

	sensoryScore := 0.0
	if totalByLevel[1] > 0 {
		sensoryScore = float64(completedByLevel[1]) / float64(totalByLevel[1]) * 100
	}

	cognitiveScore := 0.0
	if totalByLevel[4] > 0 {
		cognitiveScore = float64(completedByLevel[4]) / float64(totalByLevel[4]) * 100
	}

	if cognitiveScore > 70 && sensoryScore < 50 {
		warnings = append(warnings, "Terdeteksi 'Lompatan Perkembangan'. Anak mahir kognitif tapi pondasi sensorik (Level 1) belum kuat. Risiko: Masalah fokus/emosi di kemudian hari.")
	}

	summary := models.AssessmentSummary{
		TotalMilestones:     len(data),
		CompletedMilestones: len(data),
		ProgressByCategory:  progressByCategory,
		RedFlagsDetected:    redFlags,
		PyramidWarnings:     warnings,
	}

	return &summary, nil
}

func generatePDFReport(child models.Child, measurements []models.MeasurementResponse, summary *models.AssessmentSummary) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTopMargin(25)
	pdf.SetLeftMargin(18)
	pdf.SetRightMargin(18)
	pdf.SetAutoPageBreak(true, 20)

	// Page 1: Header and Child Info
	pdf.AddPage()
	addHeader(pdf)
	addChildInfo(pdf, child)

	// Page 2: Growth Measurements (always show, even if empty)
	pdf.AddPage()
	addGrowthMeasurements(pdf, measurements)

	// Page 3: Developmental Assessment
	if summary.TotalMilestones > 0 {
		pdf.AddPage()
		addDevelopmentalAssessment(pdf, summary)
	}

	// Footer on last page
	addFooter(pdf)

	return pdf
}

func addHeader(pdf *gofpdf.Fpdf) {
	// Title with proper wrapping to avoid cutoff
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 102, 204) // Blue color
	pdf.SetXY(18, 25)
	pdf.MultiCell(174, 9, "Tukem - Laporan Pertumbuhan & Perkembangan Anak", "", "", false)
	
	currentY := pdf.GetY()
	pdf.SetY(currentY + 5)

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 5, fmt.Sprintf("Dibuat pada: %s", time.Now().Format("02 January 2006, 15:04 WIB")))
	pdf.Ln(10)
	
	// Draw a line separator
	pdf.Line(18, pdf.GetY(), 192, pdf.GetY())
	pdf.Ln(8)
}

func addChildInfo(pdf *gofpdf.Fpdf, child models.Child) {
	// Check if we need new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}
	
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 8, "Informasi Anak")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(245, 245, 245)
	
	// Name
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(35, 6, "Nama:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, child.Name)
	pdf.Ln(7)

	// Format DOB
	dobTime, _ := time.Parse("2006-01-02", child.DOB)
	dobFormatted := dobTime.Format("02 January 2006")
	
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(35, 6, "Tanggal Lahir:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, dobFormatted)
	pdf.Ln(7)

	genderText := "Laki-laki"
	if child.Gender == "female" {
		genderText = "Perempuan"
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(35, 6, "Jenis Kelamin:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, genderText)
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(35, 6, "Berat Lahir:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%.2f kg", child.BirthWeight))
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(35, 6, "Tinggi Lahir:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%.2f cm", child.BirthHeight))
	pdf.Ln(7)

	if child.IsPremature {
		prematureText := "Prematur"
		if child.GestationalAge != nil {
			prematureText = fmt.Sprintf("Prematur (%d minggu)", *child.GestationalAge)
		}
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(35, 6, "Status:")
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(255, 140, 0)
		pdf.Cell(0, 6, prematureText)
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(7)
	}

	pdf.Ln(8)
}

func addGrowthMeasurements(pdf *gofpdf.Fpdf, measurements []models.MeasurementResponse) {
	// Check if we need new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}
	
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Riwayat Pengukuran Pertumbuhan")
	pdf.Ln(10)

	// Table header with better formatting
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(240, 240, 240)
	
	// Header row
	pdf.Cell(20, 7, "Tanggal")
	pdf.Cell(18, 7, "Umur")
	pdf.Cell(16, 7, "Berat (kg)")
	pdf.Cell(16, 7, "Tinggi (cm)")
	pdf.Cell(16, 7, "Z BB/U")
	pdf.Cell(16, 7, "Z TB/U")
	pdf.Cell(30, 7, "Status BB/U")
	pdf.Cell(34, 7, "Status TB/U")
	pdf.Ln(7)
	
	// Draw a line under header
	pdf.Line(18, pdf.GetY()-1, 192, pdf.GetY()-1)
	pdf.Ln(2)

	// Table rows with better formatting
	pdf.SetFont("Arial", "", 7)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(220, 220, 220)
	
	// Display all measurements
	for i, m := range measurements {
		// Check if we need a new page (leave space for at least 3 more rows + summary)
		if pdf.GetY() > 240 && i < len(measurements)-1 {
			pdf.AddPage()
			// Reprint header on new page
			pdf.SetFont("Arial", "B", 7)
			pdf.SetFillColor(240, 240, 240)
			
			pdf.Cell(20, 7, "Tanggal")
			pdf.Cell(18, 7, "Umur")
			pdf.Cell(16, 7, "Berat")
			pdf.Cell(16, 7, "Tinggi")
			pdf.Cell(16, 7, "Z BB/U")
			pdf.Cell(16, 7, "Z TB/U")
			pdf.Cell(30, 7, "Status BB/U")
			pdf.Cell(34, 7, "Status TB/U")
			pdf.Ln(7)
			
			pdf.SetFont("Arial", "", 7)
			pdf.SetFillColor(255, 255, 255)
		}
		
		// Format date
		dateTime, _ := time.Parse("2006-01-02", m.MeasurementDate)
		dateFormatted := dateTime.Format("02/01/2006")
		
		// Use simple cell layout with borders
		pdf.SetFont("Arial", "", 7)
		pdf.SetFillColor(255, 255, 255)
		
		// Check if age display is too long
		ageDisplay := m.AgeDisplay
		if len(ageDisplay) > 10 {
			ageDisplay = ageDisplay[:8] + "..."
		}
		
		// Use simple Cell method for gofpdf v2
		pdf.Cell(20, 6, dateFormatted)
		pdf.Cell(18, 6, ageDisplay)
		pdf.Cell(16, 6, fmt.Sprintf("%.2f", m.Weight))
		pdf.Cell(16, 6, fmt.Sprintf("%.1f", m.Height))
		
		// Z-scores
		zScoreBB := "-"
		if m.WeightForAgeZScore != nil {
			zScoreBB = fmt.Sprintf("%.2f", *m.WeightForAgeZScore)
		}
		pdf.Cell(16, 6, zScoreBB)
		
		zScoreTB := "-"
		if m.HeightForAgeZScore != nil {
			zScoreTB = fmt.Sprintf("%.2f", *m.HeightForAgeZScore)
		}
		pdf.Cell(16, 6, zScoreTB)

		// Status (shortened to fit in cell)
		nutritionalStatus := m.NutritionalStatus
		if len(nutritionalStatus) > 18 {
			nutritionalStatus = nutritionalStatus[:16] + "..."
		}
		pdf.Cell(30, 6, nutritionalStatus)

		// Height status (shortened)
		heightStatus := m.HeightStatus
		if len(heightStatus) > 18 {
			heightStatus = heightStatus[:16] + "..."
		}
		pdf.Cell(34, 6, heightStatus)
		
		pdf.Ln(6)
	}

	// Summary statistics
	if len(measurements) > 0 {
		// Check if we need new page for summary
		if pdf.GetY() > 220 {
			pdf.AddPage()
		}
		
		pdf.Ln(10)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, "Ringkasan Statistik Pertumbuhan")
		pdf.Ln(8)
		
		// Latest measurement
		latest := measurements[0]
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "Pengukuran Terakhir:")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 9)

		pdf.SetFont("Arial", "B", 9)
		pdf.Cell(40, 5, "Tanggal:")
		pdf.SetFont("Arial", "", 9)
		pdf.Cell(50, 5, latest.MeasurementDate)
		
		pdf.SetFont("Arial", "B", 9)
		pdf.Cell(30, 5, "Umur:")
		pdf.SetFont("Arial", "", 9)
		pdf.Cell(0, 5, latest.AgeDisplay)
		pdf.Ln(6)

		pdf.SetFont("Arial", "B", 9)
		pdf.Cell(40, 5, "Berat:")
		pdf.SetFont("Arial", "", 9)
		weightText := fmt.Sprintf("%.2f kg", latest.Weight)
		if latest.WeightForAgeZScore != nil {
			weightText += fmt.Sprintf(" (Z-score: %.2f)", *latest.WeightForAgeZScore)
		}
		pdf.Cell(0, 5, weightText)
		pdf.Ln(6)

		pdf.SetFont("Arial", "B", 9)
		pdf.Cell(40, 5, "Tinggi:")
		pdf.SetFont("Arial", "", 9)
		heightText := fmt.Sprintf("%.1f cm", latest.Height)
		if latest.HeightForAgeZScore != nil {
			heightText += fmt.Sprintf(" (Z-score: %.2f)", *latest.HeightForAgeZScore)
		}
		pdf.Cell(0, 5, heightText)
		pdf.Ln(6)

		if latest.NutritionalStatus != "" {
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(40, 5, "Status Gizi:")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 5, latest.NutritionalStatus)
			pdf.Ln(6)
		}
		if latest.HeightStatus != "" {
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(40, 5, "Status Tinggi:")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 5, latest.HeightStatus)
			pdf.Ln(6)
		}
		
		// Calculate statistics
		if len(measurements) > 1 {
			pdf.Ln(8)
			pdf.SetFont("Arial", "B", 10)
			pdf.Cell(0, 6, "Statistik Seluruh Data:")
			pdf.Ln(6)
			pdf.SetFont("Arial", "", 9)
			
			// Find min/max/average
			minWeight := measurements[0].Weight
			maxWeight := measurements[0].Weight
			minHeight := measurements[0].Height
			maxHeight := measurements[0].Height
			totalWeight := 0.0
			totalHeight := 0.0
			
			for _, m := range measurements {
				if m.Weight < minWeight {
					minWeight = m.Weight
				}
				if m.Weight > maxWeight {
					maxWeight = m.Weight
				}
				if m.Height < minHeight {
					minHeight = m.Height
				}
				if m.Height > maxHeight {
					maxHeight = m.Height
				}
				totalWeight += m.Weight
				totalHeight += m.Height
			}
			
			avgWeight := totalWeight / float64(len(measurements))
			avgHeight := totalHeight / float64(len(measurements))
			
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(55, 5, "Total pengukuran:")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 5, fmt.Sprintf("%d", len(measurements)))
			pdf.Ln(6)
			
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(55, 5, "Berat Badan:")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 5, fmt.Sprintf("Min: %.2f kg | Max: %.2f kg | Rata-rata: %.2f kg", minWeight, maxWeight, avgWeight))
			pdf.Ln(6)
			
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(55, 5, "Tinggi Badan:")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 5, fmt.Sprintf("Min: %.1f cm | Max: %.1f cm | Rata-rata: %.1f cm", minHeight, maxHeight, avgHeight))
			pdf.Ln(6)
		}
	} else {
		pdf.Ln(10)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(150, 150, 150)
		pdf.Cell(0, 8, "Belum ada data pengukuran untuk anak ini.")
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(8)
	}
}

func addDevelopmentalAssessment(pdf *gofpdf.Fpdf, summary *models.AssessmentSummary) {
	// Check if we need new page
	if pdf.GetY() > 240 {
		pdf.AddPage()
	}
	
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Penilaian Perkembangan")
	pdf.Ln(10)

	// Summary stats
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Total Milestone yang Dinilai: %d", summary.TotalMilestones))
	pdf.Ln(8)

	// Progress by category
	if len(summary.ProgressByCategory) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "Progress per Kategori:")
		pdf.Ln(6)

		pdf.SetFont("Arial", "", 9)
		categoryNames := map[string]string{
			"sensory":    "Sensorik (Level 1)",
			"motor":      "Motorik (Level 2)",
			"perception": "Persepsi (Level 3)",
			"cognitive":  "Kognitif (Level 4)",
		}

		for cat, progress := range summary.ProgressByCategory {
			name := categoryNames[cat]
			if name == "" {
				name = cat
			}
			pdf.SetFont("Arial", "B", 9)
			pdf.Cell(70, 6, name+":")
			pdf.SetFont("Arial", "", 9)
			pdf.Cell(0, 6, fmt.Sprintf("%.1f%%", progress))
			pdf.Ln(7)
		}
		pdf.Ln(5)
	}

	// Red Flags
	if len(summary.RedFlagsDetected) > 0 {
		// Check if we need new page
		if pdf.GetY() > 220 {
			pdf.AddPage()
		}
		
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(255, 0, 0)
		pdf.Cell(0, 7, fmt.Sprintf("⚠ PERINGATAN: %d Red Flag Terdeteksi", len(summary.RedFlagsDetected)))
		pdf.Ln(8)
		pdf.SetTextColor(0, 0, 0)

		pdf.SetFont("Arial", "", 9)
		for i, flag := range summary.RedFlagsDetected {
			if i >= 10 { // Limit to 10 flags
				break
			}
			// Check page break
			if pdf.GetY() > 250 {
				pdf.AddPage()
			}
			
			pdf.Cell(5, 5, fmt.Sprintf("%d.", i+1))
			pdf.MultiCell(0, 5, fmt.Sprintf("%s (%s)", flag.Question, flag.Category), "", "", false)
			pdf.Ln(4)
		}
		pdf.Ln(5)
	}

	// Pyramid Warnings
	if len(summary.PyramidWarnings) > 0 {
		// Check if we need new page
		if pdf.GetY() > 220 {
			pdf.AddPage()
		}
		
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(255, 140, 0)
		pdf.Cell(0, 6, "Peringatan:")
		pdf.Ln(6)
		pdf.SetTextColor(0, 0, 0)

		pdf.SetFont("Arial", "", 9)
		for _, warning := range summary.PyramidWarnings {
			// Check page break
			if pdf.GetY() > 250 {
				pdf.AddPage()
			}
			
			pdf.MultiCell(0, 5, warning, "", "", false)
			pdf.Ln(4)
		}
	}
}

func addFooter(pdf *gofpdf.Fpdf) {
	// Add footer text at bottom
	pdf.SetY(275)
	pdf.SetFont("Arial", "I", 7)
	pdf.SetTextColor(150, 150, 150)
	footerText := "Laporan ini dibuat oleh aplikasi Tukem. Untuk konsultasi lebih lanjut, hubungi dokter spesialis anak Anda."
	pdf.MultiCell(0, 4, footerText, "", "C", false)
}

