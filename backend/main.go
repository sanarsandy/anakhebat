package main

import (
	"log"
	"net/http"
	"os"
	"tukem-backend/db"
	"tukem-backend/handlers"
	customMiddleware "tukem-backend/middleware"
	"tukem-backend/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Database
	db.Init()

	// Apply Migrations in order
	migrations := []string{
		"001_init_schema.sql",
		"002_children_table.sql",
		"003_measurements_table.sql",
		"004_milestones_tables.sql",
		"005_who_standards.sql",
		"006_add_denver_domain.sql",
		"007_stimulation_content.sql",
		"008_immunization_tables.sql",
	}

	for _, migration := range migrations {
		if err := utils.ApplyMigration(db.DB, migration); err != nil {
			// Log error but don't fail if it's just "table already exists"
			// For this simple implementation, we assume the SQL has IF NOT EXISTS
			log.Printf("Warning: Migration %s might have failed (or already exists): %v", migration, err)
		}
	}

	// Seed Data
	if err := utils.SeedMilestones(db.DB); err != nil {
		log.Printf("Warning: Milestones seeding failed: %v", err)
	}

	if err := utils.SeedWHOStandards(db.DB); err != nil {
		log.Printf("Warning: WHO standards seeding failed: %v", err)
	}
	if err := utils.SeedDenverIIMilestones(db.DB); err != nil {
		log.Printf("Warning: Denver II milestones seeding failed: %v", err)
	}
	
	if err := utils.SeedStimulationContent(db.DB); err != nil {
		log.Printf("Warning: Stimulation content seeding failed: %v", err)
	}
	
	if err := utils.SeedImmunizationSchedule(db.DB); err != nil {
		log.Printf("Warning: Immunization schedule seeding failed: %v", err)
	}

	e := EchoServer()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	e.Logger.Fatal(e.Start(":" + port))
}

func EchoServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS Configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders: []string{"Content-Disposition"}, // Expose Content-Disposition for file downloads
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Tukem API",
			"status":  "healthy",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "up",
		})
	})

	// Auth Routes
	auth := e.Group("/api/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)

	// Protected Routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware())
	api.GET("/me", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "You are authorized!",
		})
	})

	// Milestone Routes
	milestoneHandler := handlers.NewMilestoneHandler(db.DB) // Assuming db.DB is the sqlx.DB instance
	api.GET("/milestones", milestoneHandler.GetMilestones)
	api.GET("/milestones/denver-ii", handlers.GetDenverIIMilestones)
	
	// Recommendations Route - MUST be BEFORE /children/:id to avoid route conflict
	api.GET("/children/:id/recommendations", handlers.GetRecommendations)
	
	// Immunization Routes - MUST be BEFORE /children/:id to avoid route conflict
	api.GET("/children/:id/immunizations", handlers.GetImmunizationSchedule)
	api.POST("/children/:id/immunizations", handlers.RecordImmunization)
	
	// Children Routes (general routes first)
	api.POST("/children", handlers.CreateChild)
	api.GET("/children", handlers.GetChildren)
	
	// Denver II Routes - MUST be registered BEFORE /children/:id to avoid route conflict
	// Echo matches routes in order, so more specific routes must come first
	api.GET("/children/:id/denver-ii/chart-data", handlers.GetDenverIIChartData)
	api.GET("/children/:id/denver-ii/grid-data", handlers.GetDenverIIChartGridData)
	
	// PDF Export Route - MUST be BEFORE /children/:id and other /children/:id/* routes to avoid route conflict
	api.GET("/children/:id/export-pdf", handlers.ExportChildReport)
	
	// Measurement Routes (must come before /children/:id to avoid conflict)
	api.POST("/children/:id/measurements", handlers.CreateMeasurement)
	api.GET("/children/:id/measurements", handlers.GetMeasurements)
	api.GET("/children/:id/measurements/latest", handlers.GetLatestMeasurement)
	api.PUT("/children/:id/measurements/:measurementId", handlers.UpdateMeasurement)
	api.DELETE("/children/:id/measurements/:measurementId", handlers.DeleteMeasurement)
	
	// Assessment Routes (must come before /children/:id to avoid conflict)
	api.GET("/children/:id/assessments", milestoneHandler.GetChildAssessments)
	api.PUT("/children/:id/assessments/batch", milestoneHandler.BatchUpsertAssessments)
	api.GET("/children/:id/assessments/summary", milestoneHandler.GetAssessmentSummary)
	
	// Children detail routes (must come after ALL specific /children/:id/* routes)
	api.GET("/children/:id", handlers.GetChild)
	api.PUT("/children/:id", handlers.UpdateChild)
	api.DELETE("/children/:id", handlers.DeleteChild)

	// Growth Charts Routes
	api.GET("/who-standards", handlers.GetWHOStandardsForChart)

	return e
}
