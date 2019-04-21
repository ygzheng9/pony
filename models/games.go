package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Game struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Criterion string    `json:"criterion" db:"criterion"`
	Options   string    `json:"options" db:"options"`
	Pairs     string    `json:"pairs" db:"pairs"`
	Weights   string    `json:"weights" db:"weights"`
	Ratios    string    `json:"ratios" db:"ratios"`
}

// String is not required by pop and may be deleted
func (g *Game) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Games is not required by pop and may be deleted
type Games []Game

// String is not required by pop and may be deleted
func (g Games) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (g *Game) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: g.Name, Name: "Name"},
		// &validators.StringIsPresent{Field: g.Criterion, Name: "Criterion"},
		// &validators.StringIsPresent{Field: g.Options, Name: "Options"},
		// &validators.StringIsPresent{Field: g.Pairs, Name: "Pairs"},
		// &validators.StringIsPresent{Field: g.Weights, Name: "Weights"},
		// &validators.StringIsPresent{Field: g.Ratios, Name: "Ratios"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (g *Game) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (g *Game) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
