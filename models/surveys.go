package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Survey struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	SurveyNo   string    `json:"survey_no" db:"survey_no"`
	SubmitUser string    `json:"submit_user" db:"submit_user"`
	QuestionNo string    `json:"question_no" db:"question_no"`
	Answers    string    `json:"answers" db:"answers"`
	SubmitDate string    `json:"submit_date" db:"submit_date"`
}

// String is not required by pop and may be deleted
func (s Survey) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Surveys is not required by pop and may be deleted
type Surveys []Survey

// String is not required by pop and may be deleted
func (s Surveys) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Survey) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.SurveyNo, Name: "SurveyNo"},
		// &validators.StringIsPresent{Field: s.SubmitUser, Name: "SubmitUser"},
		// &validators.StringIsPresent{Field: s.QuestionNo, Name: "QuestionNo"},
		// &validators.StringIsPresent{Field: s.Answers, Name: "Answers"},
		// &validators.StringIsPresent{Field: s.SubmitDate, Name: "SubmitDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Survey) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Survey) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
