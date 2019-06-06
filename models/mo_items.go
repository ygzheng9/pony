package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type MoItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Company     string    `json:"company" db:"company"`
	Line        string    `json:"line" db:"line"`
	MoDate      time.Time `json:"mo_date" db:"mo_date"`
	MoNum       string    `json:"mo_num" db:"mo_num"`
	ItemNum     string    `json:"item_num" db:"item_num"`
	MatName     string    `json:"mat_name" db:"mat_name"`
	ItemQty     float64   `json:"item_qty" db:"item_qty"`
	ItemUnit    string    `json:"item_unit" db:"item_unit"`
	Cate2       string    `json:"cate2" db:"cate2"`
	Cate1       string    `json:"cate1" db:"cate1"`
	InboundDate time.Time `json:"inbound_date" db:"inbound_date"`
	Warehouse   string    `json:"warehouse" db:"warehouse"`

	WorkOrder string    `json:"work_num" db:"work_num"`
	MatCode   string    `json:"mat_code" db:"mat_code"`
	MoType    string    `json:"mo_type" db:"mo_type"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`

	Source string `json:"source" db:"source"`

	Shift        string  `json:"shift" db:"shift"`
	Step         string  `json:"step" db:"step"`
	MatSpec      string  `json:"mat_spec" db:"mat_spec"`
	MainMatQty   float64 `json:"main_mat_qty" db:"main_mat_qty"`
	InputMatQty1 float64 `json:"input_mat_qty1" db:"input_mat_qty1"`
	InputMatQty2 float64 `json:"input_mat_qty2" db:"input_mat_qty2"`
	ClaimQty1    float64 `json:"claim_qty1" db:"claim_qty1"`
	ClaimQty2    float64 `json:"claim_qty2" db:"claim_qty2"`
}

// String is not required by pop and may be deleted
func (m MoItem) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MoItems is not required by pop and may be deleted
type MoItems []MoItem

// String is not required by pop and may be deleted
func (m MoItems) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MoItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Company, Name: "Company"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MoItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MoItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
