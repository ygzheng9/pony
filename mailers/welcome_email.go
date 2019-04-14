package mailers

import (
	stdmail "net/mail"
	"os"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"

	"pony/base"
)

// SendWelcomeEmails first mail
func SendWelcomeEmails() error {
	sugar := base.Sugar()

	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "Welcome Email"

	addr := stdmail.Address{
		Name:    "想知道",
		Address: "ibmpartment@powerdekor.com.cn",
	}
	from := addr.String()
	sugar.Infow("email", "from", from)
	m.From = from

	m.To = []string{"yonggang.zheng@qq.com", "ibmpartment@powerdekor.com.cn"}

	surveyDir := envy.Get("SurveyDir", "./config/surveys")
	file := surveyDir + "/ABCD.png"
	logo, err := os.Open(file)
	if err != nil {
		sugar.Errorw("can not attached", "file", file)
		return err
	}
	m.AddEmbedded("ABCD.png", logo)

	err = m.AddBody(r.HTML("welcome_email.html"), render.Data{"msg": "第一封信"})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
