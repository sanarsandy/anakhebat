package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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

	// Skip migrations - use init_db.sql script instead
	// Run: ./setup-database.sh or manually run backend/init_db.sql

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
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		// Support multiple origins separated by comma
		corsOrigins = append(corsOrigins, strings.Split(corsEnv, ",")...)
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
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
	// OTP Authentication Routes
	auth.POST("/request-otp", handlers.RequestOTP)
	auth.POST("/verify-otp", handlers.VerifyOTP)
	auth.POST("/resend-otp", handlers.ResendOTP)
	// Admin OTP Authentication Routes
	auth.POST("/admin/request-otp", handlers.AdminRequestOTP)
	auth.POST("/admin/verify-otp", handlers.AdminVerifyOTP)
	// Google OAuth Routes
	auth.GET("/google", handlers.GetGoogleAuthURL)
	auth.GET("/google/callback", handlers.GoogleAuthCallback)

	// Protected Routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware())
	api.GET("/me", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "You are authorized!",
		})
	})

	// User Profile Routes
	api.GET("/user/profile", handlers.GetUserProfile)
	api.PUT("/user/profile", handlers.UpdateUserProfile)
	api.POST("/user/verify-phone", handlers.RequestPhoneVerification)
	api.POST("/user/verify-phone/confirm", handlers.ConfirmPhoneVerification)

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

	// Admin Routes - Protected with admin middleware
	// Note: api group already has JWTMiddleware, so admin routes are already JWT protected
	admin := api.Group("/admin")
	admin.Use(customMiddleware.RequireAdmin())

	// Admin User Management
	admin.GET("/users", handlers.GetAdminUsers)
	admin.GET("/users/:id", handlers.GetAdminUser)
	admin.POST("/users", handlers.CreateAdminUser)
	admin.PUT("/users/:id", handlers.UpdateAdminUser)
	admin.DELETE("/users/:id", handlers.DeleteAdminUser)
	admin.POST("/users/:id/reset-password", handlers.ResetAdminUserPassword)

	// Admin Analytics
	admin.GET("/analytics/overview", handlers.GetAdminOverview)
	admin.GET("/analytics/users", handlers.GetAdminUsersAnalytics)
	admin.GET("/analytics/children", handlers.GetAdminChildrenAnalytics)
	admin.GET("/analytics/measurements", handlers.GetAdminMeasurementsAnalytics)
	admin.GET("/analytics/assessments", handlers.GetAdminAssessmentsAnalytics)
	admin.GET("/analytics/immunizations", handlers.GetAdminImmunizationsAnalytics)

	// Admin Reports
	admin.GET("/reports/users", handlers.GetUsersReport)
	admin.GET("/reports/children", handlers.GetChildrenReport)
	admin.GET("/reports/growth", handlers.GetGrowthReport)

	// Admin System Settings
	admin.GET("/settings", handlers.GetSystemSettings)
	admin.GET("/settings/:key", handlers.GetSystemSetting)
	admin.PUT("/settings/:key", handlers.UpdateSystemSetting)
	admin.PUT("/settings", handlers.UpdateSystemSettingsBatch)

	// Admin Audit Logs
	admin.GET("/audit-logs", handlers.GetAuditLogs)
	admin.GET("/audit-logs/:id", handlers.GetAuditLog)
	admin.GET("/audit-logs/export", handlers.ExportAuditLogs)

	// Admin Data Access (View All)
	admin.GET("/children", handlers.GetAdminChildren)
	admin.GET("/children/:id", handlers.GetAdminChild)
	admin.GET("/measurements", handlers.GetAdminMeasurements)
	admin.GET("/assessments", handlers.GetAdminAssessments)
	admin.GET("/immunizations", handlers.GetAdminImmunizations)

	// Admin Master Data Management
	admin.GET("/milestones", handlers.GetAdminMilestones)
	admin.GET("/milestones/:id", handlers.GetAdminMilestone)
	admin.POST("/milestones", handlers.CreateAdminMilestone)
	admin.PUT("/milestones/:id", handlers.UpdateAdminMilestone)
	admin.DELETE("/milestones/:id", handlers.DeleteAdminMilestone)

	admin.GET("/who-standards", handlers.GetAdminWHOStandards)
	admin.GET("/who-standards/:id", handlers.GetAdminWHOStandard)
	admin.POST("/who-standards", handlers.CreateAdminWHOStandard)
	admin.PUT("/who-standards/:id", handlers.UpdateAdminWHOStandard)
	admin.DELETE("/who-standards/:id", handlers.DeleteAdminWHOStandard)

	admin.GET("/stimulation-content", handlers.GetAdminStimulationContent)
	admin.GET("/stimulation-content/:id", handlers.GetAdminStimulationContentItem)
	admin.POST("/stimulation-content", handlers.CreateAdminStimulationContent)
	admin.PUT("/stimulation-content/:id", handlers.UpdateAdminStimulationContent)
	admin.DELETE("/stimulation-content/:id", handlers.DeleteAdminStimulationContent)

	admin.GET("/immunization-schedules", handlers.GetAdminImmunizationSchedules)
	admin.GET("/immunization-schedules/:id", handlers.GetAdminImmunizationSchedule)
	admin.POST("/immunization-schedules", handlers.CreateAdminImmunizationSchedule)
	admin.PUT("/immunization-schedules/:id", handlers.UpdateAdminImmunizationSchedule)
	admin.DELETE("/immunization-schedules/:id", handlers.DeleteAdminImmunizationSchedule)

	return e
}
