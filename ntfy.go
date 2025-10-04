package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	sasl "github.com/emersion/go-sasl"
	smtp "github.com/emersion/go-smtp"
)

/*
* Upstream Discord webhook
 */

func includeAttributesDiscord(attr map[string]string) string {

	var attributes = make([]string, 0)
	for k, v := range attr {
		attributes = append(attributes, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(attributes, "\n")
}

type UpstreamDiscord struct {
	Url string
}

type DiscardMessage struct {
	Content string `json:"content"`
}

func (s *UpstreamDiscord) Send(ntfy *Notification) error {

	msg := fmt.Sprintf("%s\n\n%s", includeAttributesDiscord(ntfy.Attributes), ntfy.Msg)

	v, err := json.Marshal(DiscardMessage{Content: msg})
	if err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	res, err := http.Post(s.Url, "application/json", bytes.NewBuffer(v))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received invalid status code %d\n %v", res.StatusCode, string(resBody))
	}

	return nil
}

func (s *UpstreamDiscord) IsValid() bool {
	return s.Url != ""
}

func (s *UpstreamDiscord) Name() string {
	return "discord"
}

/*
* Upstream SMTP sender
 */

func includeAttributesSMTP(attr map[string]string) string {

	var attributes = make([]string, 0)
	for k, v := range attr {
		if k == "Subject" {
			continue
		}
		attributes = append(attributes, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(attributes, "\n")
}

type UpstreamEmail struct {
	Username   string
	Password   string
	Port       string
	Server     string
	Recipients string
}

func (s *UpstreamEmail) Send(ntfy *Notification) error {
	auth := sasl.NewPlainClient("", s.Username, s.Password)

	recipients := strings.Split(s.Recipients, ",")

	subject := ntfy.Meta["subject"]
	if subject == "" {
		subject = "New Message"
	}

	msg := strings.NewReader(fmt.Sprintf("Subject: Notification: %s\r\nFrom: %s\r\nTo: %s\r\n\r\n%s\n\n%s", subject, s.Username,
		s.Recipients, includeAttributesSMTP(ntfy.Attributes), ntfy.Msg))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.Server, s.Port), auth, s.Username, recipients, msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *UpstreamEmail) IsValid() bool {
	return s.Server != ""
}

func (s *UpstreamEmail) Name() string {
	return "email"
}
