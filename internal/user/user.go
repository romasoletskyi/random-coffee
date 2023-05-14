package user

import (
	_ "embed"
	"html/template"
	"strings"

	"github.com/romasoletskyi/random-coffee/internal/data"
	"gopkg.in/gomail.v2"
)

//go:embed confirmation-template
var confirmationEmail string

var Username, Password string

func sendMail(to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "random.coffee.manager@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer("email-smtp.eu-west-3.amazonaws.com", 587, Username, Password)
	return d.DialAndSend(m)
}

func SendConfirmationMail(form data.UserForm) error {
	t, err := template.New("make-letter").Parse(confirmationEmail)
	if err != nil {
		return err
	}

	var builder strings.Builder
	err = t.Execute(&builder, form)
	if err != nil {
		return err
	}

	return sendMail(form.Email, "Submit confirmation", builder.String())
}
