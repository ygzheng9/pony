package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type SoItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Company   string    `json:"company" db:"company"`
	OrderNum  string    `json:"order_num" db:"order_num"`
	CustNum   string    `json:"cust_num" db:"cust_num"`
	Category  string    `json:"category" db:"category"`
	Serial    string    `json:"serial" db:"serial"`
	MatName   string    `json:"mat_name" db:"mat_name"`
	MatModel  string    `json:"mat_model" db:"mat_model"`
	ItemQty   float64   `json:"item_qty" db:"item_qty"`
	Period    string    `json:"period" db:"period"`
	SalesType string    `json:"sales_type" db:"sales_type"`
	MoveType  string    `json:"move_type" db:"move_type"`
	Warehouse string    `json:"warehouse" db:"warehouse"`
	WhDate    time.Time `json:"wh_date" db:"wh_date"`
	DocDate   time.Time `json:"doc_date" db:"doc_date"`
	BookParty string    `json:"book_party" db:"book_party"`
	Remark    string    `json:"remark" db:"remark"`
	Source    string    `json:"source" db:"source"`
}

// String is not required by pop and may be deleted
func (s SoItem) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// SoItems is not required by pop and may be deleted
type SoItems []SoItem

// String is not required by pop and may be deleted
func (s SoItems) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *SoItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Company, Name: "Company"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *SoItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *SoItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
