package helpers

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"
)

var (
	tmpl *template.Template
)

func init() {
	var err error
	if tmpl, err = template.ParseGlob("../templates/*.tmpl"); err != nil {
		log.Fatal("Unable to load template files")
	}
}

// EmailSender sends emails
type EmailSender interface {
	EmailSend(toAddresses []string, subject string, body string) error
}

// EmailServerConfigurer configures the email sending service
type EmailServerConfigurer struct {
	EmailHost     string
	EmailPort     string
	EmailFrom     string
	EmailLogin    string
	EmailPassword string
}

// EmailMessageDataProvider provides data to render an email
type EmailMessageDataProvider struct {
	From    string
	To      string
	Subject string
	Body    string
}

// NewEmailSender provides an instance of an EmailSender given an EmailServerConfigurer
func NewEmailSender(conf EmailServerConfigurer) EmailSender {
	return &emailSender{conf, smtp.SendMail}
}

type emailSender struct {
	conf EmailServerConfigurer
	send func(string, smtp.Auth, string, []string, []byte) error
}

// EmailSend sends an email given message parameters
func (e *emailSender) EmailSend(toAddresses []string, subject string, body string) error {
	buf := new(bytes.Buffer)
	emdp := EmailMessageDataProvider{
		From:    e.conf.EmailFrom,
		To:      strings.Join(toAddresses, ", "),
		Subject: subject,
		Body:    body,
	}

	if err := tmpl.ExecuteTemplate(buf, "email.raw.tmpl", emdp); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", e.conf.EmailLogin, e.conf.EmailPassword, e.conf.EmailHost)
	hostAndPort := fmt.Sprintf("%s:%s", e.conf.EmailHost, e.conf.EmailPort)

	if err := e.send(hostAndPort, auth, e.conf.EmailFrom, toAddresses, buf.Bytes()); err != nil {
		return err
	}

	return nil
}
