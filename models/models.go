package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

// IDT uuid
type IDT struct {
	ID uuid.UUID `json:"id" db:"id"`
}

// IDList uuid list
type IDList []IDT

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}
