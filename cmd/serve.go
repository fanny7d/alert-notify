/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Alert represents the structure of an alert from Alertmanager
type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt"`
	EndsAt      string            `json:"endsAt"`
}

// AlertmanagerPayload represents the payload sent by Alertmanager
type AlertmanagerPayload struct {
	Alerts []Alert `json:"alerts"`
}

// ChatMessage represents the structure of a message to be sent to the chat app
type ChatMessage struct {
	Text string `json:"text"`
}

func sendToChatApp(payload []byte) error {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL is not set")
	}

	log.Printf("Sending payload: %s", string(payload))

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

func formatAlertMessage(alert Alert) map[string]interface{} {
	color := "#FF0000" // Default to red for "firing"
	if alert.Status == "resolved" {
		color = "#00FF00" // Green for "resolved"
	}

	log.Printf("Alert status: %s, color: %s", alert.Status, color)

	return map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color": color,
				"text": fmt.Sprintf(
					"**Alert Name:** %s\n"+
						"**Status:** %s\n"+
						"**Description:** %s",
					alert.Labels["alertname"],
					alert.Status,
					alert.Annotations["description"],
				),
				"fields": []map[string]interface{}{
					{"short": false, "title": "Starts At", "value": alert.StartsAt},
					{"short": false, "title": "Ends At", "value": alert.EndsAt},
				},
			},
		},
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request")

	var payload AlertmanagerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	for _, alert := range payload.Alerts {
		message := formatAlertMessage(alert)

		payload, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		if err := sendToChatApp(payload); err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Println("Message sent successfully")
		}
	}

	w.WriteHeader(http.StatusOK)
}

func startWebhookServer() {
	http.HandleFunc("/webhook", alertHandler)
	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook server to receive alerts",
	Long: `This command starts an HTTP server that listens for incoming alerts
from Alertmanager and forwards them to a chat application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startWebhookServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
