package grifts

import (
	"fmt"
	"log"

	"github.com/markbates/grift/grift"

	"pony/base"
	"pony/mailers"
)

var _ = grift.Namespace("mails", func() {
	_ = grift.Desc("welcome", "Send welcome mail")
	_ = grift.Add("welcome", func(c *grift.Context) error {
		err := mailers.SendWelcomeEmails()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("completed. ")
		return nil
	})

	_ = grift.Desc("send", "Send all mails")
	_ = grift.Add("send", func(c *grift.Context) error {
		fmt.Println("sending all mails....")
		err := base.SendAllEmails()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("completed. ")
		return nil
	})
})
