package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"bytes"
	"encoding/json"
)

func TestSendToChatApp(t *testing.T) {
	// Create a test server to mock the webhook endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Set the environment variable to the test server URL
	os.Setenv("WEBHOOK_URL", ts.URL)

	// Create a sample chat message
	message := ChatMessage{
		Content: "Test message",
	}

	// Call the function to test
	err := sendToChatApp(message)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestAlertHandler(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(alertHandler))
	defer ts.Close()

	// Set the environment variable to a mock URL
	os.Setenv("WEBHOOK_URL", "http://example.com")

	// Create a sample alert payload
	alertPayload := AlertmanagerPayload{
		Alerts: []Alert{
			{
				Status: "firing",
				Labels: map[string]string{"alertname": "TestAlert"},
				Annotations: map[string]string{"description": "This is a test alert"},
			},
		},
	}

	payloadBytes, err := json.Marshal(alertPayload)
	if err != nil {
		t.Fatalf("Failed to marshal alert payload: %v", err)
	}

	// Send a POST request to the test server
	resp, err := http.Post(ts.URL+"/webhook", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
} 