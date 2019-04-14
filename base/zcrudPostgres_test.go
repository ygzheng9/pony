package base

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestPrintSubCode(t *testing.T) {
	SetRoot()

	// CRUD
	err := PrintSubCode("notices")
	if err != nil {
		t.Errorf("err: %+v\n", err)
	}
}

func TestSugar(t *testing.T) {
	sugar := Sugar()
	sugar.Info("asdf")

	fire()
}
