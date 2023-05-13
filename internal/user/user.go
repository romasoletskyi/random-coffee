package user

import (
	_ "embed"
	"html/template"
	"strings"

	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

//go:embed confirmation-template
var confirmationEmail string

var Username, Password string

func sendMail(to, subject, text string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "random.coffee.manager@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer("email-smtp.eu-west-3.amazonaws.com", 587, Username, Password)
	if err := d.DialAndSend(m); err != nil {
		logrus.Error(err)
	}
}

func SendConfirmationMail(form data.UserForm) {
	t, err := template.New("make-letter").Parse(confirmationEmail)
	if err != nil {
		logrus.Error(err)
		return
	}

	var builder strings.Builder
	err = t.Execute(&builder, form)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info(builder.String())
	sendMail(form.Email, "Submit confirmation", builder.String())
}
