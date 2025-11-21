package models

import (
	"time"
)

// ImmunizationSchedule represents a master immunization schedule item
type ImmunizationSchedule struct {
	ID                      string    `json:"id" db:"id"`
	Name                    string    `json:"name" db:"name"`
	NameID                  *string   `json:"name_id,omitempty" db:"name_id"`
	Description             *string   `json:"description,omitempty" db:"description"`
	
	// Timing in days
	AgeMinDays              *int      `json:"age_min_days,omitempty" db:"age_min_days"`
	AgeOptimalDays          *int      `json:"age_optimal_days,omitempty" db:"age_optimal_days"`
	AgeMaxDays              *int      `json:"age_max_days,omitempty" db:"age_max_days"`
	
	// Timing in months (for display)
	AgeMinMonths            *int      `json:"age_min_months,omitempty" db:"age_min_months"`
	AgeOptimalMonths        *int      `json:"age_optimal_months,omitempty" db:"age_optimal_months"`
	AgeMaxMonths            *int      `json:"age_max_months,omitempty" db:"age_max_months"`
	
	// Dose information
	DoseNumber              int       `json:"dose_number" db:"dose_number"`
	TotalDoses              *int      `json:"total_doses,omitempty" db:"total_doses"`
	
	// Interval information
	IntervalFromPreviousDays *int     `json:"interval_from_previous_days,omitempty" db:"interval_from_previous_days"`
	IntervalFromPreviousMonths *int   `json:"interval_from_previous_months,omitempty" db:"interval_from_previous_months"`
	
	// Category & priority
	Category                string    `json:"category" db:"category"` // wajib, tambahan
	Priority                string    `json:"priority" db:"priority"` // high, medium, low
	IsRequired              bool      `json:"is_required" db:"is_required"`
	
	// Additional info
	Notes                   *string   `json:"notes,omitempty" db:"notes"`
	Source                  string    `json:"source" db:"source"`
	IsActive                bool      `json:"is_active" db:"is_active"`
	
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

// ChildImmunization represents a recorded immunization for a child
type ChildImmunization struct {
	ID                      string    `json:"id" db:"id"`
	ChildID                 string    `json:"child_id" db:"child_id"`
	ImmunizationScheduleID  string    `json:"immunization_schedule_id" db:"immunization_schedule_id"`
	
	// Given date and age
	GivenDate               string    `json:"given_date" db:"given_date"` // YYYY-MM-DD
	GivenAtAgeDays          *int      `json:"given_at_age_days,omitempty" db:"given_at_age_days"`
	GivenAtAgeMonths        *int      `json:"given_at_age_months,omitempty" db:"given_at_age_months"`
	
	// Details
	Location                *string   `json:"location,omitempty" db:"location"`
	HealthcareFacility      *string   `json:"healthcare_facility,omitempty" db:"healthcare_facility"`
	DoctorName              *string   `json:"doctor_name,omitempty" db:"doctor_name"`
	VaccineBatchNumber      *string   `json:"vaccine_batch_number,omitempty" db:"vaccine_batch_number"`
	Notes                   *string   `json:"notes,omitempty" db:"notes"`
	
	// Status
	IsOnSchedule            *bool     `json:"is_on_schedule,omitempty" db:"is_on_schedule"`
	IsCatchUp               bool      `json:"is_catch_up" db:"is_catch_up"`
	
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Schedule                *ImmunizationSchedule `json:"schedule,omitempty" db:"-"`
}

// ImmunizationRecordRequest is the payload for recording an immunization
type ImmunizationRecordRequest struct {
	ImmunizationScheduleID string  `json:"immunization_schedule_id" validate:"required"`
	GivenDate              string  `json:"given_date" validate:"required"` // YYYY-MM-DD
	Location               *string `json:"location,omitempty"`
	HealthcareFacility     *string `json:"healthcare_facility,omitempty"`
	DoctorName             *string `json:"doctor_name,omitempty"`
	VaccineBatchNumber     *string `json:"vaccine_batch_number,omitempty"`
	Notes                  *string `json:"notes,omitempty"`
}

// ImmunizationStatus represents the status of an immunization for a child
type ImmunizationStatus struct {
	Schedule                ImmunizationSchedule `json:"schedule"`
	Status                  string               `json:"status"` // pending, completed, overdue, upcoming
	DueDate                 *string              `json:"due_date,omitempty"` // YYYY-MM-DD
	DueAgeMonths            *int                 `json:"due_age_months,omitempty"`
	DaysUntilDue            *int                 `json:"days_until_due,omitempty"`
	DaysOverdue             *int                 `json:"days_overdue,omitempty"`
	
	// If completed
	Record                  *ChildImmunization   `json:"record,omitempty"`
}

// ImmunizationScheduleResponse represents the API response for immunization schedule
type ImmunizationScheduleResponse struct {
	ChildID       string               `json:"child_id"`
	AgeMonths     int                  `json:"age_months"`
	AgeDays       int                  `json:"age_days"`
	Immunizations []ImmunizationStatus `json:"immunizations"`
	Summary       ImmunizationSummary  `json:"summary"`
}

// ImmunizationSummary contains summary statistics
type ImmunizationSummary struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Pending   int `json:"pending"`
	Overdue   int `json:"overdue"`
	Upcoming  int `json:"upcoming"`
}

