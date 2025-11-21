package models

import (
	"time"
)

type Child struct {
	ID              string    `json:"id" db:"id"`
	ParentID        string    `json:"parent_id" db:"parent_id"`
	Name            string    `json:"name" db:"name"`
	DOB             string    `json:"dob" db:"dob"` // Date of Birth (YYYY-MM-DD)
	Gender          string    `json:"gender" db:"gender"` // 'male' or 'female'
	BirthWeight     float64   `json:"birth_weight" db:"birth_weight"` // kg
	BirthHeight     float64   `json:"birth_height" db:"birth_height"` // cm
	IsPremature     bool      `json:"is_premature" db:"is_premature"`
	GestationalAge  *int      `json:"gestational_age,omitempty" db:"gestational_age"` // weeks (if premature)
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type CreateChildRequest struct {
	Name           string  `json:"name" validate:"required"`
	DOB            string  `json:"dob" validate:"required"` // YYYY-MM-DD
	Gender         string  `json:"gender" validate:"required,oneof=male female"`
	BirthWeight    float64 `json:"birth_weight" validate:"required,gt=0"`
	BirthHeight    float64 `json:"birth_height" validate:"required,gt=0"`
	IsPremature    bool    `json:"is_premature"`
	GestationalAge *int    `json:"gestational_age,omitempty"`
}

type UpdateChildRequest struct {
	Name           string  `json:"name"`
	DOB            string  `json:"dob"`
	Gender         string  `json:"gender"`
	BirthWeight    float64 `json:"birth_weight"`
	BirthHeight    float64 `json:"birth_height"`
	IsPremature    bool    `json:"is_premature"`
	GestationalAge *int    `json:"gestational_age,omitempty"`
}
