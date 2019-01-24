package helpers

import (
	"net/smtp"
	"testing"
)

func TestEmailSender(t *testing.T) {
	f, r := mockSend(nil)
	sender := &emailSender{send: f}
	body := "Hello!"
	subject := "Hello!"
	err := sender.EmailSend([]string{"me@example.com"}, subject, body)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expected := "From: \nTo: me@example.com\nSubject: Hello!\n\nHello!"

	if string(r.msg) != expected {
		t.Errorf("wrong message body. expected: %s got: %s", expected, r.msg)
	}
}

func mockSend(errToReturn error) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*r = emailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, r
}

type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}
