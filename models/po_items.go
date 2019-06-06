package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type PoItem struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	Company        string    `json:"company" db:"company"`
	PoDate         time.Time `json:"po_date" db:"po_date"`
	PoNum          string    `json:"po_num" db:"po_num"`
	VendorName     string    `json:"vendor_name" db:"vendor_name"`
	MatName        string    `json:"mat_name" db:"mat_name"`
	ItemQty        float64   `json:"item_qty" db:"item_qty"`
	ItemUnit       string    `json:"item_unit" db:"item_unit"`
	UnitPrice      float64   `json:"unit_price" db:"unit_price"`
	Cate2          string    `json:"cate2" db:"cate2"`
	Cate1          string    `json:"cate1" db:"cate1"`
	Operator       string    `json:"operator" db:"operator"`
	InboundQty     float64   `json:"inbound_qty" db:"inbound_qty"`
	OutstandingQty float64   `json:"outstanding_qty" db:"outstanding_qty"`
	InboundDate    time.Time `json:"inbound_date" db:"inbound_date"`

	LineNum string `json:"line_num" db:"line_num"`
	DnNum   string `json:"dn_num" db:"dn_num"`
	DnItem  string `json:"dn_item" db:"dn_item"`

	ItemStatus    string    `json:"item_status" db:"item_status"`
	MatCode       string    `json:"mat_code" db:"mat_code"`
	MatSpec       string    `json:"mat_spec" db:"mat_spec"`
	InboundUnit   string    `json:"inbound_unit" db:"inbound_unit"`
	ArriveBookQty float64   `json:"arrived_book_qty" db:"arrived_book_qty"`
	BookedQty     float64   `json:"booked_qty" db:"booked_qty"`
	UnbookedQty   float64   `json:"unbooked_qty" db:"unbooked_qty"`
	PlannedDate   time.Time `json:"planned_date" db:"planned_date"`
	DelayedDays   float64   `json:"delayed_days" db:"delayed_days"`
}

// String is not required by pop and may be deleted
func (p PoItem) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PoItems is not required by pop and may be deleted
type PoItems []PoItem

// String is not required by pop and may be deleted
func (p PoItems) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PoItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Company, Name: "Company"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PoItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PoItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
