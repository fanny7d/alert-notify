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

func sendToChatApp(message ChatMessage) error {
	// Get the webhook URL from the environment variable
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL is not set")
	}

	// Format the message as Markdown
	payload, err := json.Marshal(map[string]string{
		"text": message.Text,
	})
	if err != nil {
		return err
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

func alertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request")

	var payload AlertmanagerPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded payload: %+v", payload)

	for _, alert := range payload.Alerts {
		// Format the message using an enhanced Markdown template
		message := ChatMessage{
			Text: fmt.Sprintf(
				"**🚨 Alert Notification 🚨**\n\n"+
					"**Alert Name:** [%s](https://your-link-here.com)\n"+
					"**Status:** %s\n"+
					"**Description:** %s\n"+
					"**Starts At:** %s\n"+
					"**Ends At:** %s\n"+
					"---",
				alert.Labels["alertname"],
				alert.Status,
				alert.Annotations["description"],
				alert.StartsAt,
				alert.EndsAt,
			),
		}

		if err := sendToChatApp(message); err != nil {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
