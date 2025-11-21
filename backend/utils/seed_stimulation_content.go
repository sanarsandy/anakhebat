package utils

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// SeedStimulationContent seeds sample stimulation content for intervention recommendations
func SeedStimulationContent(db *sqlx.DB) error {
	// Check if content already exists
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM stimulation_content")
	if err != nil {
		return fmt.Errorf("failed to check existing content: %v", err)
	}

	if count > 0 {
		log.Println("Stimulation content already seeded, skipping...")
		return nil
	}

	// Sample stimulation content data
	contents := []struct {
		Category     string
		Title        string
		Description  string
		ContentType  string
		URL          string
		AgeMinMonths *int
		AgeMaxMonths *int
	}{
		// Sensory stimulation content
		{
			Category:     "sensory",
			Title:        "Stimulasi Sensorik untuk Bayi 0-6 Bulan",
			Description:  "Panduan lengkap stimulasi sensorik untuk bayi usia 0-6 bulan, termasuk stimulasi penglihatan, pendengaran, dan sentuhan",
			ContentType:  "article",
			URL:          "https://example.com/stimulasi-sensorik-0-6-bulan",
			AgeMinMonths: intPtr(0),
			AgeMaxMonths: intPtr(6),
		},
		{
			Category:     "sensory",
			Title:        "Video: Stimulasi Visual untuk Bayi Baru Lahir",
			Description:  "Video tutorial cara melakukan stimulasi visual untuk meningkatkan kemampuan penglihatan bayi",
			ContentType:  "video",
			URL:          "https://www.youtube.com/watch?v=example1",
			AgeMinMonths: intPtr(0),
			AgeMaxMonths: intPtr(3),
		},
		// Motor stimulation content
		{
			Category:     "motor",
			Title:        "Latihan Motorik Kasar: Menengkurapkan Bayi",
			Description:  "Panduan langkah demi langkah melatih bayi untuk tengkurap, termasuk tips keamanan dan variasi latihan",
			ContentType:  "article",
			URL:          "https://example.com/latihan-tengkurap",
			AgeMinMonths: intPtr(2),
			AgeMaxMonths: intPtr(6),
		},
		{
			Category:     "motor",
			Title:        "Video: Stimulasi Motorik Halus 6-12 Bulan",
			Description:  "Video tutorial stimulasi motorik halus seperti memegang benda, memindahkan benda dari satu tangan ke tangan lain",
			ContentType:  "video",
			URL:          "https://www.youtube.com/watch?v=example2",
			AgeMinMonths: intPtr(6),
			AgeMaxMonths: intPtr(12),
		},
		// Perception stimulation content
		{
			Category:     "perception",
			Title:        "Mengembangkan Persepsi Ruang pada Bayi",
			Description:  "Cara mengembangkan kemampuan persepsi ruang bayi melalui permainan dan aktivitas sehari-hari",
			ContentType:  "article",
			URL:          "https://example.com/persepsi-ruang",
			AgeMinMonths: intPtr(9),
			AgeMaxMonths: intPtr(18),
		},
		{
			Category:     "perception",
			Title:        "Permainan Mengenal Warna dan Bentuk",
			Description:  "Ide permainan sederhana untuk membantu bayi mengenal warna dan bentuk dasar",
			ContentType:  "article",
			URL:          "https://example.com/mengenal-warna-bentuk",
			AgeMinMonths: intPtr(12),
			AgeMaxMonths: intPtr(24),
		},
		// Cognitive stimulation content
		{
			Category:     "cognitive",
			Title:        "Stimulasi Kognitif: Belajar Sebab-Akibat",
			Description:  "Cara mengajarkan konsep sebab-akibat pada balita melalui permainan interaktif",
			ContentType:  "article",
			URL:          "https://example.com/sebab-akibat",
			AgeMinMonths: intPtr(18),
			AgeMaxMonths: intPtr(36),
		},
		{
			Category:     "cognitive",
			Title:        "Video: Permainan Mengasah Memori Balita",
			Description:  "Video tutorial permainan sederhana untuk mengasah kemampuan memori dan konsentrasi balita",
			ContentType:  "video",
			URL:          "https://www.youtube.com/watch?v=example3",
			AgeMinMonths: intPtr(24),
			AgeMaxMonths: intPtr(48),
		},
		// General age-appropriate content
		{
			Category:     "sensory",
			Title:        "Stimulasi Multi-Sensor untuk Tumbuh Kembang Optimal",
			Description:  "Panduan stimulasi yang melibatkan berbagai indera untuk mendukung tumbuh kembang bayi secara menyeluruh",
			ContentType:  "article",
			URL:          "https://example.com/multi-sensor",
			AgeMinMonths: intPtr(0),
			AgeMaxMonths: intPtr(12),
		},
		{
			Category:     "motor",
			Title:        "Tahapan Perkembangan Motorik: Panduan Orang Tua",
			Description:  "Panduan lengkap tahapan perkembangan motorik kasar dan halus dari bayi hingga balita",
			ContentType:  "article",
			URL:          "https://example.com/tahapan-motorik",
			AgeMinMonths: intPtr(0),
			AgeMaxMonths: intPtr(36),
		},
	}

	// Insert content
	query := `
		INSERT INTO stimulation_content (category, title, description, content_type, url, age_min_months, age_max_months, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
	`

	for _, content := range contents {
		_, err := db.Exec(query,
			content.Category,
			content.Title,
			content.Description,
			content.ContentType,
			content.URL,
			content.AgeMinMonths,
			content.AgeMaxMonths,
		)
		if err != nil {
			log.Printf("Failed to insert stimulation content '%s': %v", content.Title, err)
			// Continue with other content
			continue
		}
	}

	log.Printf("Stimulation content seeded successfully! (%d items)", len(contents))
	return nil
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

