package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// API and token constants
const baseURL = "https://api.boxmateapp.co.uk"

var memberToken string

// Lift and RM IDs
var liftRMIDs = []int{49, 50, 51, 52, 53, 338, 25, 26, 27, 28, 29, 30, 335, 1, 2, 3, 4, 5, 6, 334, 13, 14, 15, 16, 17, 18, 337, 72, 73, 74, 75, 76, 77, 130, 131, 132, 133, 134, 331, 333, 336}

// Helper to create multipart form-data
func createMultipartFormData() (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("Selected_Type", "Leaderboard")
	_ = writer.WriteField("Selected_Date", "")
	_ = writer.WriteField("Selected_Session", "null")
	_ = writer.WriteField("member_token", memberToken)

	err := writer.Close()
	if err != nil {
		log.Fatalf("Failed to close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

// Fetch data from API for a given lift ID
func fetchAndSaveLiftData(liftID int) {
	memberToken = os.Getenv("BOXMATE_AUTH_TOKEN")

	url := fmt.Sprintf("%s/leaderboard/Exercise/%d", baseURL, liftID)
	body, contentType := createMultipartFormData()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("Failed to create request for ID %d: %v", liftID, err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Golang)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make API request for ID %d: %v", liftID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API returned non-OK status %d for ID %d", resp.StatusCode, liftID)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body for ID %d: %v", liftID, err)
	}

	// Save response to a JSON file
	filename := fmt.Sprintf("%d.json", liftID)
	err = os.WriteFile(filename, respBody, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON file for ID %d: %v", liftID, err)
	}

	fmt.Printf("Successfully saved data for ID %d to %s\n", liftID, filename)
}

func main() {
	for _, liftID := range liftRMIDs {
		fetchAndSaveLiftData(liftID)
	}
	fmt.Println("All data fetched and saved successfully!")
}
