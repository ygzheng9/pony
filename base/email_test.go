package base

import (
	"fmt"
	"log"
	"testing"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

func TestEmailLoadByID(t *testing.T) {
	SetRoot()

	u, err := uuid.FromString("2501f03e-9596-4596-b8c7-4b9f0f18d8c1")
	if err != nil {
		log.Fatalf("failed to parse UUID %v", err)
	}

	e := EmailLoadByID(u)

	err = e.markSent()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("done.")
}

func TestGenerateInvitation(t *testing.T) {
	SetRoot()
	err := GenerateInvitation()
	if err != nil {
		t.Error(err)
	}
}
func TestSendAllEmails(t *testing.T) {
	var err error
	SetRoot()

	err = SendAllEmails()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetEmailNotices(t *testing.T) {
	items, err := GetEmailNotices()
	if err != nil {
		t.Error(err)
	}

	for _, i := range items {
		fmt.Println(i)
	}

}
