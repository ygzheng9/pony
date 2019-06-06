package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type InvItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Company   string    `json:"company" db:"company"`
	Warehouse string    `json:"warehouse" db:"warehouse"`
	Year      string    `json:"year" db:"year"`
	Month     string    `json:"month" db:"month"`
	MatName   string    `json:"mat_name" db:"mat_name"`
	MatCode   string    `json:"mat_code" db:"mat_code"`
	MatSpec   string    `json:"mat_spec" db:"mat_spec"`
	MatStyle  string    `json:"mat_style" db:"mat_style"`
	MatType   string    `json:"mat_type" db:"mat_type"`
	TreeType  string    `json:"tree_type" db:"tree_type"`
	CustCode  string    `json:"cust_code" db:"cust_code"`
	Color     string    `json:"color" db:"color"`
	MatUnit   string    `json:"mat_unit" db:"mat_unit"`
	MatQty    float64   `json:"mat_qty" db:"mat_qty"`
	MatAmt    string    `json:"mat_amt" db:"mat_amt"`
	MatGrade  string    `json:"mat_grade" db:"mat_grade"`
	Cate1     string    `json:"cate1" db:"cate1"`
	Cate2     string    `json:"cate2" db:"cate2"`
	Surface   string    `json:"surface" db:"surface"`
	Source    string    `json:"source" db:"source"`
}

// String is not required by pop and may be deleted
func (i InvItem) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// InvItems is not required by pop and may be deleted
type InvItems []InvItem

// String is not required by pop and may be deleted
func (i InvItems) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *InvItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Company, Name: "Company"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *InvItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *InvItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
