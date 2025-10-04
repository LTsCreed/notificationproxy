package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	smtp "github.com/emersion/go-smtp"
)

var messageQueue = make(chan *Notification, 5)

func httpServer() {

	http.HandleFunc("/status", HttpStatusHandler)
	http.HandleFunc("/hook", HttpWebhookHandler)

	port := getEnv("SERVER_HOOK_PORT", "8080")

	logger.Printf("HTTP server is starting: http://0.0.0.0:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		logger.Fatal(err)
	}
}

func smtpServer() {

	be := &Backend{messageQueue: messageQueue}

	s := smtp.NewServer(be)

	port := getEnv("SERVER_SMTP_PORT", "2525")

	s.Addr = fmt.Sprintf(":%s", port)
	s.Domain = "notificationproxy"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	logger.Printf("SMTP server is starting: smtp://0.0.0.0:%s", port)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func messageProcessing(upstreams *[]UpstreamServer) {

	for msg := range messageQueue {
		for _, u := range *upstreams {
			if !u.IsValid() {
				continue
			}
			if err := u.Send(msg); err != nil {
				logger.Printf("%s:: failed to send to message\n %v", u.Name(), err)
			}
		}
	}

}

func createpstreamServers() *[]UpstreamServer {

	upstreams := &[]UpstreamServer{
		&UpstreamDiscord{Url: os.Getenv("NTFY_DISCORD_URL")},
		&UpstreamEmail{
			Username:   os.Getenv("NTFY_SMTP_USERNAME"),
			Password:   os.Getenv("NTFY_SMTP_PASSWORD"),
			Recipients: os.Getenv("NTFY_SMTP_RECIPIENTS"),
			Server:     os.Getenv("NTFY_SMTP_SERVER"),
			Port:       getEnv("NTFY_SMTP_PORT", "587"),
		},
	}

	validUpstreams := []UpstreamServer{}

	for _, upstream := range *upstreams {
		if upstream.IsValid() {
			validUpstreams = append(validUpstreams, upstream)
		}
	}

	if len(validUpstreams) == 0 {
		logger.Println("Error: Please define at least one valid notification server")
	}
	return &validUpstreams

}

func main() {

	upstreams := createpstreamServers()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go httpServer()
	go smtpServer()
	go messageProcessing(upstreams)

	<-stop

}
