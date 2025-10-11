package main

import (
	"encoding/json"
	"fmt"
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

// Used by Docker, Kubernetes probes
func HttpStatusHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseJSON(body []byte) string {
	var bodyJSON any

	json.Unmarshal(body, &bodyJSON)

	switch v := bodyJSON.(type) {
	case map[string]any:
		var attributes = make([]string, 0)
		for k, v := range v {
			attributes = append(attributes, fmt.Sprintf("%s: %s", k, v))
		}
		return strings.Join(attributes, "\n")
	default:
		return string(body)
	}

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

	var msg string
	switch req.Header.Get("Content-Type") {
	case "application/json":
		msg = parseJSON(body)
	default:
		msg = string(body)
	}

	attributes := map[string]string{}

	for k, v := range req.URL.Query() {
		attributes[k] = strings.Join(v, ", ")
	}

	messageQueue <- &Notification{
		Msg:         msg,
		Source:      "webhook",
		Attributes:  attributes,
		ContentType: "text",
		Meta:        map[string]string{},
	}

	w.WriteHeader(http.StatusAccepted)

}
