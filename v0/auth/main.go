package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// AuthResponse represents the structure of the response JSON
type AuthResponse struct {
	Success     bool   `json:"success"`
	MemberToken string `json:"member_token"`
}

// Authenticate sends the username and password to the API and retrieves the member token
func Authenticate(username, password string) (string, error) {
	// Define the API URL.
	url := "https://api.boxmateapp.co.uk/member/authenticate"

	// Create a buffer to hold the multipart form data.
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the form fields.
	if err := writer.WriteField("Member_Email", username); err != nil {
		return "", fmt.Errorf("failed to write email field: %w", err)
	}
	if err := writer.WriteField("Member_Password", password); err != nil {
		return "", fmt.Errorf("failed to write password field: %w", err)
	}

	// Close the writer to finalize the body.
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Create a new POST request with the body.
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Perform the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	// Decode the response JSON.
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the authentication was successful.
	if !authResp.Success {
		return "", fmt.Errorf("authentication failed")
	}

	return authResp.MemberToken, nil
}

func main() {
	username := os.Getenv("BOXMATE_USERNAME")
	password := os.Getenv("BOXMATE_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Environment variables BOXMATE_USERNAME and BOXMATE_PASSWORD must be set")
	}

	token, err := Authenticate(username, password)
	if err != nil {
		log.Fatalf("Failed to get authentication token: %v", err)
	}

	fmt.Printf("Authentication token: %s\n", token)
}
