package main

import (
	"io"
	"mime"
	"net/mail"

	smtp "github.com/emersion/go-smtp"
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

	m, err := mail.ReadMessage(r)
	if err != nil {
		logger.Printf("failed to parse mail message %s\n", err)
		return nil
	}

	b, err := io.ReadAll(m.Body)
	if err != nil {
		return err
	}

	dec := new(mime.WordDecoder)
	from, _ := dec.DecodeHeader(m.Header.Get("From"))
	to, _ := dec.DecodeHeader(m.Header.Get("To"))
	subject, _ := dec.DecodeHeader(m.Header.Get("Subject"))

	attributes := &map[string]string{
		"Original From": from,
		"Original To":   to,
		"Subject":       subject,
	}

	if getEnv("SERVER_SMTP_INCLUDE_ORIGINAL_HEADER", "true") == "false" {
		attributes = &map[string]string{}
	}

	messageQueue <- &Notification{
		Msg:        string(b),
		MsgType:    "email",
		Attributes: *attributes,
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
