package main

import (
	"net/smtp"
	//"fmt"
	"html/template"
	"bytes"
)

func sendEmail(to []string, mssg []byte) error {
	if conf.Smtp.Active {
		verbose("sending email to: %v", to)
		return smtp.SendMail(
			conf.Smtp.Path,
			smtp.PlainAuth(
		 	  	conf.Smtp.Identity,
		 	  	conf.Smtp.Username,
		 	  	conf.Smtp.Password,
		 	  	conf.Smtp.Host,
	 	   ),
			conf.Smtp.From,
			to,
			mssg,
		)
	}
	return nil
}

var emailMime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

type ConfirmEmail struct {
	Uname string
	Key string
	Address string
}

var confirmSubject = "Subject: Knaxim Email Confirmation;\n"
var confirmTemplate *template.Template

func (ce *ConfirmEmail) buildMessage(from string, to string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	if _, err := buffer.WriteString("From: "); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(from); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(";\n"); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString("To: "); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(to); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(";\n"); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(confirmSubject); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(emailMime); err != nil {
		return nil, err
	}
	if err := confirmTemplate.Execute(buffer, ce); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
