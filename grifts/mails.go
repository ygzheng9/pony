package grifts

import (
	"fmt"
	"log"

	"github.com/markbates/grift/grift"

	"pony/base"
	"pony/mailers"
)

var _ = grift.Namespace("mails", func() {
	grift.Desc("welcome", "Send welcome mail")
	grift.Add("welcome", func(c *grift.Context) error {
		err := mailers.SendWelcomeEmails()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("completed. ")
		return nil
	})

	grift.Desc("send", "Send all mails")
	grift.Add("send", func(c *grift.Context) error {
		fmt.Println("sending all mails....")
		err := base.SendAllEmails()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("completed. ")
		return nil
	})
})
