package user

import (
	"os"
	"testing"

	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/stretchr/testify/require"
)

func TestInvitation(t *testing.T) {
	form := data.PairForm{
		ActionLink: "https://www.google.com/",
		Left:       data.PairInfo{Name: "Roman", Email: "abc@gmail.com", Contact: "@ABC", Bio: "Hello"},
		Right:      data.PairInfo{Name: "Yaroslav", Email: "xyz@gmail.com", Contact: "", Bio: "Testing bio!"},
	}

	mail, err := PrepareMail(invitationEmail, form)
	require.NoError(t, err)

	file, err := os.Create("invitation.html")
	require.NoError(t, err)

	_, err = file.WriteString(mail)
	require.NoError(t, err)
}
