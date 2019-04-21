package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type GameOption struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	GameID    uuid.UUID `json:"game_id" db:"game_id"`
	Seq       int       `json:"seq" db:"seq"`
	Pairs     string    `json:"pairs" db:"pairs"`
	Weights   string    `json:"weights" db:"weights"`
	Ratios    string    `json:"ratios" db:"ratios"`
}

// String is not required by pop and may be deleted
func (g GameOption) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// GameOptions is not required by pop and may be deleted
type GameOptions []GameOption

// String is not required by pop and may be deleted
func (g GameOptions) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (g *GameOption) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
	// &validators.IntIsPresent{Field: g.Seq, Name: "Seq"},
	// &validators.StringIsPresent{Field: g.Pairs, Name: "Pairs"},
	// &validators.StringIsPresent{Field: g.Weights, Name: "Weights"},
	// &validators.StringIsPresent{Field: g.Ratios, Name: "Ratios"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (g *GameOption) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (g *GameOption) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
