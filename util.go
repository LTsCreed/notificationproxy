package main

import (
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var logger = log.New(os.Stdout, "notificationProxy: ", log.LstdFlags)

type Notification struct {
	Msg     string
	MsgType string
	// Additional fields to include in the message; handling is upstream-specific
	Attributes map[string]string
	// Internal metadata; must not be included in the upstream message
	Meta map[string]string
}

type UpstreamServer interface {
	Send(ntfy *Notification) error
	IsValid() bool
	Name() string
}
