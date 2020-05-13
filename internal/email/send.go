package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"git.maxset.io/web/knaxim/internal/config"
)

func sendEmail(to []string, msg []byte) error {
	return smtp.SendMail(config.V.Email.Server, smtp.PlainAuth(config.V.Email.Credential.Identity, config.V.Email.Credential.Username, config.V.Email.Credential.Password, config.V.Email.Credential.Host), config.V.Email.From, to, msg)
}

var resetEmail = strings.ReplaceAll(`To: %s
From: %s
Subject: CloudEdison Reset Password
MIME-version: 1.0
Content-Type: text/plain; charset=\"UTF-8\"

A request for a password reset has been recieved for your account named %s. Follow the path to reset the password.

https://%s/reset/%s

CloudEdison Team
`, "\n", "\r\n")

// SendResetEmail sends an email with a reset link
// to: list of emails
// name: username of the account with the password reset
// address: web address of the site to visit to reset the password
// resetkey: the key needed to reset the password
func SendResetEmail(to []string, name, address, resetkey string) error {
	msgstr := fmt.Sprintf(resetEmail, strings.Join(to, ", "), config.V.Email.From, name, address, resetkey)
	return sendEmail(to, []byte(msgstr))
}

var errorEmail = strings.ReplaceAll(`To: %s
From: %s
Subject: Automated Error Report
MIME-version: 1.0
Content-Type: text/plain; charset=\"UTF-8\"

This is an automated email from the CloudEdison server regarding the occurrence of a server error.

%s
`, "\n", "\r\n")

// SendErrorEmail sends an email containing the given message to the address
// specified as ErrorEmail in the config.
// If no email address has been set, an error will be returned.
func SendErrorEmail(msg string) error {
	if len(config.V.ErrorEmail) == 0 {
		return errors.New("no error email address has been set")
	}
	msgstr := fmt.Sprintf(errorEmail, config.V.ErrorEmail, config.V.Email.From, msg)
	return sendEmail(strings.Split(config.V.ErrorEmail, ","), []byte(msgstr))
}
