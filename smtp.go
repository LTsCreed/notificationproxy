package main

import (
	"io"

	smtp "github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

/*
* Downstream SMTP receiver
 */

// The Backend implements SMTP server methods.
type Backend struct {
	messageQueue chan *Notification
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{messageQueue: bkd.messageQueue}, nil
}

// A Session is returned after successful login.
type Session struct {
	messageQueue chan *Notification
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	return nil
}

func (s *Session) Data(r io.Reader) error {

	m, err := enmime.ReadEnvelope(r)
	if err != nil {
		logger.Printf("failed to parse mail message %s\n", err)
		return nil
	}

	from := m.GetHeader("From")
	to := m.GetHeader("To")
	subject := m.GetHeader("Subject")

	attributes := &map[string]string{
		"Original From": from,
		"Original To":   to,
		"Subject":       subject,
	}

	if getEnv("SERVER_SMTP_INCLUDE_ORIGINAL_HEADER", "true") == "false" {
		attributes = &map[string]string{}
	}

	contentType := "text"
	msg := m.Text

	messageQueue <- &Notification{
		Msg:         msg,
		Source:      "email",
		ContentType: contentType,
		Attributes:  *attributes,
		Meta: map[string]string{
			"subject": subject,
			"to":      to,
			"from":    from,
		},
	}

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
