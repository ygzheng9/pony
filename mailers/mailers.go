package mailers

import (
	"crypto/tls"
	"log"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
)

var smtp mail.Sender
var r *render.Engine

func init() {

	// Pulling config from the env.
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	var err error

	sender, err := mail.NewSMTPSender(host, port, user, password)
	sender.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// narrow struct to interface.
	smtp = sender

	if err != nil {
		log.Fatal(err)
	}

	r = render.New(render.Options{
		HTMLLayout:   "layout.html",
		TemplatesBox: packr.New("app:mailers:templates", "../templates/mail"),
		Helpers:      render.Helpers{},
	})
}
