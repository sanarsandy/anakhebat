package models

import (
	"time"
)

type Measurement struct {
	ID                 string    `json:"id" db:"id"`
	ChildID            string    `json:"child_id" db:"child_id"`
	MeasurementDate    string    `json:"measurement_date" db:"measurement_date"` // YYYY-MM-DD
	Weight             float64   `json:"weight" db:"weight"`                      // kg
	Height             float64   `json:"height" db:"height"`                      // cm
	HeadCircumference  *float64  `json:"head_circumference,omitempty" db:"head_circumference"` // cm, optional
	AgeInDays          int       `json:"age_in_days" db:"age_in_days"`
	AgeInMonths        int       `json:"age_in_months" db:"age_in_months"`
	WeightForAgeZScore *float64  `json:"weight_for_age_zscore,omitempty" db:"weight_for_age_zscore"`
	HeightForAgeZScore *float64  `json:"height_for_age_zscore,omitempty" db:"height_for_age_zscore"`
	WeightStatus       string    `json:"weight_status,omitempty" db:"weight_status"`
	HeightStatus       string    `json:"height_status,omitempty" db:"height_status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type CreateMeasurementRequest struct {
	MeasurementDate   string   `json:"measurement_date" validate:"required"` // YYYY-MM-DD
	Weight            float64  `json:"weight" validate:"required,gt=0"`
	Height            float64  `json:"height" validate:"required,gt=0"`
	HeadCircumference *float64 `json:"head_circumference,omitempty"`
}

type MeasurementResponse struct {
	ID                      string   `json:"id"`
	ChildID                 string   `json:"child_id"`
	MeasurementDate         string   `json:"measurement_date"`
	Weight                  float64  `json:"weight"`
	Height                  float64  `json:"height"`
	HeadCircumference       *float64 `json:"head_circumference,omitempty"`
	AgeInDays               int      `json:"age_in_days"`
	AgeInMonths             int      `json:"age_in_months"`
	AgeDisplay              string   `json:"age_display"` // "2 years 3 months"
	WeightForAgeZScore      *float64 `json:"weight_for_age_zscore,omitempty"`
	HeightForAgeZScore      *float64 `json:"height_for_age_zscore,omitempty"`
	WeightForHeightZScore   *float64 `json:"weight_for_height_zscore,omitempty"`
	HeadCircumferenceZScore *float64 `json:"head_circumference_zscore,omitempty"`
	NutritionalStatus       string   `json:"nutritional_status,omitempty"`
	HeightStatus            string   `json:"height_status,omitempty"`
	WeightForHeightStatus   string   `json:"weight_for_height_status,omitempty"`
	CreatedAt               string   `json:"created_at"`
}
