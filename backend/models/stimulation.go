package models

import (
	"time"
)

// StimulationContent represents a piece of stimulation content (video/article)
type StimulationContent struct {
	ID           string    `json:"id" db:"id"`
	MilestoneID  *string   `json:"milestone_id,omitempty" db:"milestone_id"` // Optional: link to specific milestone
	Category     string    `json:"category" db:"category"`                    // sensory, motor, perception, cognitive
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description,omitempty" db:"description"`
	ContentType  string    `json:"content_type" db:"content_type"` // video, article
	URL          string    `json:"url" db:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	AgeMinMonths *int      `json:"age_min_months,omitempty" db:"age_min_months"`
	AgeMaxMonths *int      `json:"age_max_months,omitempty" db:"age_max_months"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	
	// Joined fields for API responses
	Milestone *Milestone `json:"milestone,omitempty" db:"-"`
}

// Recommendation represents a recommended stimulation content for a child
type Recommendation struct {
	Content      StimulationContent `json:"content"`
	Reason       string             `json:"reason"`       // Why this content is recommended
	Priority     string             `json:"priority"`     // high, medium, low
	RelatedMilestone *Milestone     `json:"related_milestone,omitempty"` // The milestone that triggered this recommendation
}

// RecommendationsResponse represents the API response for recommendations
type RecommendationsResponse struct {
	ChildID       string          `json:"child_id"`
	AgeMonths     int             `json:"age_months"`
	Recommendations []Recommendation `json:"recommendations"`
}

