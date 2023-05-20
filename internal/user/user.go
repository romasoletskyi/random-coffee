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

//go:embed invitation-template
var invitationEmail string

var Username, Password string

func SendMail(to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "random.coffee.manager@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer("smtp-relay.sendinblue.com", 587, Username, Password)
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

	return SendMail(form.Email, "Submit confirmation", builder.String())
}

func SendInvitationMail(form data.PairForm) error {
	t, err := template.New("make-letter").Parse(invitationEmail)
	if err != nil {
		return err
	}

	var builder strings.Builder
	err = t.Execute(&builder, form)
	if err != nil {
		return err
	}

	return SendMail(form.Left.Email, "Pair invitation", builder.String())
}
