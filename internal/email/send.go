package email

import (
	"fmt"
	"strings"

	"git.maxset.io/web/knaxim/internal/config"
)

func sendEmail(to []string, msg []byte) error {
	return nil
}

var resetEmail = strings.ReplaceAll(`To: %s;
From: %s;
Subject: Knaxim Reset Password;

A request for a password reset has been recieved for your account named %s. Follow the path to reset the password.

https://%s/profile/reset/%s

Knaxim Team
`, "\n", "\r\n")

func SendResetEmail(to []string, name string, resetkey string) error {
	msgstr := fmt.Sprintf(resetEmail, strings.Join(to, ", "), config.V.Email.From, name, resetkey)
	return sendEmail(to, []byte(msgstr))
}
