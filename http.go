package main

import (
	"io"
	"net/http"
	"strings"
)

/*
* Downstream SMTP receiver
 */

type HookMessage struct {
	IsHostIncluded bool
	Host           string
	Message        string
}

// Used by Dokcer, Kubernetes probes
func HttpStatusHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HttpWebhookHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	attributes := map[string]string{}

	for k, v := range req.URL.Query() {
		attributes[k] = strings.Join(v, ", ")
	}

	messageQueue <- &Notification{
		Msg:        string(body),
		Source:     "webhook",
		Attributes: attributes,
		Meta:       map[string]string{},
	}

	w.WriteHeader(http.StatusAccepted)

}
