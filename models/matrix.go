package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Matrix for matrix values
type Matrix struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Company    string    `json:"company" db:"company"`
	Version    string    `json:"version" db:"version"`
	Period     string    `json:"period" db:"period"`
	Matrix     string    `json:"matrix" db:"matrix"`
	Code       string    `json:"code" db:"code"`
	Value      string    `json:"value" db:"value"`
	SubmitUser string    `json:"submit_user" db:"submit_user"`
}

// String is not required by pop and may be deleted
func (m *Matrix) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Matrices is not required by pop and may be deleted
type Matrices []Matrix

// String is not required by pop and may be deleted
func (m Matrices) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Matrix) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
	// &validators.StringIsPresent{Field: m.Company, Name: "Company"},
	// &validators.StringIsPresent{Field: m.Period, Name: "Period"},
	// &validators.StringIsPresent{Field: m.Matrix, Name: "Matrix"},
	// &validators.StringIsPresent{Field: m.Code, Name: "Code"},
	// &validators.StringIsPresent{Field: m.SubmitUser, Name: "SubmitUser"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Matrix) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Matrix) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
