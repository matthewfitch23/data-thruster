package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	// Define the API URL
	url := "https://api.boxmateapp.co.uk/member/authenticate"

	// Define the multipart form-data body
	body := fmt.Sprintf(`------WebKitFormBoundaryoA6q2wZUjNogmwHC
Content-Disposition: form-data; name="Member_Email"

%s
------WebKitFormBoundaryoA6q2wZUjNogmwHC
Content-Disposition: form-data; name="Member_Password"

%s
------WebKitFormBoundaryoA6q2wZUjNogmwHC--`, username, password)

	// Create a new POST request with the body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Host", "api.boxmateapp.co.uk")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Origin", "capacitor://localhost")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryoA6q2wZUjNogmwHC")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	// Decode the response JSON
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the authentication was successful
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
