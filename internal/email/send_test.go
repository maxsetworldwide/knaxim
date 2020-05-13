package email

import (
	"flag"
	"os"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/config"
)

var to = flag.String("to", "", "address to send test email")

func init() {
	config.V.Email.From = "noreply@maxset.org"
	config.V.Email.Server = "sub5.mail.dreamhost.com:587"
	config.V.Email.Credential.Identity = ""
	config.V.Email.Credential.Username = "ai@maxset.org"
	config.V.Email.Credential.Password = "anM5x3B8"
	config.V.Email.Credential.Host = "sub5.mail.dreamhost.com"
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func buildTo(t *string) []string {
	return strings.Split(*t, ",")
}

func TestSend(t *testing.T) {
	if len(*to) == 0 {
		t.Fatalf("Test error: please specify an email address for which to send a test email.")
	}
	t.Log("sending to: ", to)
	t.Run("Reset", func(t *testing.T) {
		err := SendResetEmail(buildTo(to), "id:2020-02-13-23:56:51:262t", "knaxim.com", "")
		if err != nil {
			t.Fatal("unable to send email: ", err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		emailMsg := "this is a test email error message"
		config.V.ErrorEmail = *to
		err := SendErrorEmail(emailMsg)
		if err != nil {
			t.Fatal("unable to send email: ", err)
		}
	})
}
