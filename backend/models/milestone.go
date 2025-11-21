package models

import (
	"time"
)

// Milestone represents a developmental milestone checklist item
type Milestone struct {
	ID           string    `json:"id" db:"id"`
	AgeMonths    int       `json:"age_months" db:"age_months"`
	MinAgeRange  *int      `json:"min_age_range,omitempty" db:"min_age_range"`
	MaxAgeRange  *int      `json:"max_age_range,omitempty" db:"max_age_range"`
	Category     string    `json:"category" db:"category"` // sensory, motor, perception, cognitive
	Question     string    `json:"question" db:"question"`
	QuestionEn   string    `json:"question_en,omitempty" db:"question_en"`
	Source       string    `json:"source" db:"source"`
	IsRedFlag    bool      `json:"is_red_flag" db:"is_red_flag"`
	PyramidLevel int       `json:"pyramid_level" db:"pyramid_level"` // 1-4
	DenverDomain *string   `json:"denver_domain,omitempty" db:"denver_domain"` // PS, FM, L, GM (Denver II)
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Assessment represents a user's answer to a milestone
type Assessment struct {
	ID             string    `json:"id" db:"id"`
	ChildID        string    `json:"child_id" db:"child_id"`
	MilestoneID    string    `json:"milestone_id" db:"milestone_id"`
	AssessmentDate string    `json:"assessment_date" db:"assessment_date"` // YYYY-MM-DD
	Status         string    `json:"status" db:"status"`                   // yes, no, sometimes
	Notes          string    `json:"notes,omitempty" db:"notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	
	// Joined fields for API responses
	Milestone *Milestone `json:"milestone,omitempty" db:"-"`
}

// AssessmentItem for batch requests
type AssessmentItem struct {
	MilestoneID string `json:"milestone_id" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=yes no sometimes"`
	Notes       string `json:"notes"`
}

// BatchAssessmentRequest is the payload for submitting multiple assessments
type BatchAssessmentRequest struct {
	AssessmentDate string           `json:"assessment_date" validate:"required"`
	Items          []AssessmentItem `json:"items" validate:"required,min=1"`
}

// AssessmentSummary contains progress data for the dashboard
type AssessmentSummary struct {
	TotalMilestones    int                    `json:"total_milestones"`
	CompletedMilestones int                   `json:"completed_milestones"`
	ProgressByCategory map[string]float64     `json:"progress_by_category"` // category -> percentage
	RedFlagsDetected   []Milestone            `json:"red_flags_detected"`
	PyramidWarnings    []string               `json:"pyramid_warnings"`
	NextMilestones     []Milestone            `json:"next_milestones"`
}
