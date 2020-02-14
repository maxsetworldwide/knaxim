package email

import (
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
Subject: Knaxim Reset Password
MIME-version: 1.0
Content-Type: text/plain; charset=\"UTF-8\"

A request for a password reset has been recieved for your account named %s. Follow the path to reset the password.

https://%s/profile/reset/%s

Knaxim Team
`, "\n", "\r\n")

func SendResetEmail(to []string, name, address, resetkey string) error {
	msgstr := fmt.Sprintf(resetEmail, strings.Join(to, ", "), config.V.Email.From, name, address, resetkey)
	return sendEmail(to, []byte(msgstr))
}
