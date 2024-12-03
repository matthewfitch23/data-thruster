package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/xuri/excelize/v2"
)

// Lift and RM IDs
var liftRMIDs = map[string]map[string]int{
	"Front Squat":   {"1RM": 49, "2RM": 50, "3RM": 51, "5RM": 52, "10RM": 53, "20RM": 338},
	"Back Squat":    {"1RM": 25, "2RM": 26, "3RM": 27, "5RM": 28, "8RM": 29, "10RM": 30, "20RM": 335},
	"Deadlift":      {"1RM": 1, "2RM": 2, "3RM": 3, "5RM": 4, "8RM": 5, "10RM": 6, "20RM": 334},
	"Sumo Deadlift": {"1RM": 13, "2RM": 14, "3RM": 15, "5RM": 16, "8RM": 17, "10RM": 18, "20RM": 337},
	"Bench Press":   {"1RM": 72, "2RM": 73, "3RM": 74, "5RM": 75, "8RM": 76, "10RM": 77, "20RM": 336},
	"Strict Press":  {"1RM": 130, "2RM": 131, "3RM": 132, "5RM": 133, "8RM": 331, "10RM": 134, "20RM": 333},
}

// Person structure
type Person struct {
	Name  string
	Lifts map[string]map[string]string
}

// APIResponse represents the leaderboard response
type APIResponse struct {
	Results []struct {
		UserID         int    `json:"Member_ID"`
		Name           string `json:"Member_Name"`
		ComponentName  string `json:"Component_Name"`
		ComponentScore string `json:"Component_Score"`
	} `json:"results"`
}

// FetchLeaderboard simulates an API call by reading a local JSON file
func fetchLeaderboard(exerciseID int) APIResponse {
	filename := fmt.Sprintf("./%d.json", exerciseID)
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", filename, err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(fileContent, &apiResp); err != nil {
		log.Fatalf("Failed to parse JSON from file %s: %v", filename, err)
	}

	return apiResp
}

func aggregateData() map[int]*Person {
	data := make(map[int]*Person)

	for lift, rmidMap := range liftRMIDs {
		for rm, id := range rmidMap {
			response := fetchLeaderboard(id)
			for _, result := range response.Results {
				if _, exists := data[result.UserID]; !exists {
					data[result.UserID] = &Person{
						Name:  result.Name,
						Lifts: make(map[string]map[string]string),
					}
				}
				if _, exists := data[result.UserID].Lifts[lift]; !exists {
					data[result.UserID].Lifts[lift] = make(map[string]string)
				}
				data[result.UserID].Lifts[lift][rm] = result.ComponentScore
			}
		}
	}

	// Fill missing RMs with "0"
	for _, person := range data {
		for lift := range liftRMIDs {
			if _, exists := person.Lifts[lift]; !exists {
				person.Lifts[lift] = make(map[string]string)
			}
			for rm := range liftRMIDs[lift] {
				if _, exists := person.Lifts[lift][rm]; !exists {
					person.Lifts[lift][rm] = "0"
				}
			}
		}
	}

	return data
}

func exportToExcel(data map[int]*Person) {
	f := excelize.NewFile()

	// Create a sheet for consolidated data
	sheetName := "Sheet1"

	// Write header
	header := []string{"Name", "Lift", "1RM", "2RM", "3RM", "5RM", "8RM", "10RM", "20RM"}
	for col, value := range header {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, value)
	}

	// Write data rows
	row := 2
	for _, person := range data {
		for lift, rmData := range person.Lifts {
			// Write the name
			nameCell, _ := excelize.CoordinatesToCellName(1, row)
			f.SetCellValue(sheetName, nameCell, person.Name)

			// Write the lift name
			liftCell, _ := excelize.CoordinatesToCellName(2, row)
			f.SetCellValue(sheetName, liftCell, lift)

			// Write RM data
			for col, rm := range []string{"1RM", "2RM", "3RM", "5RM", "8RM", "10RM", "20RM"} {
				scoreCell, _ := excelize.CoordinatesToCellName(col+3, row)
				f.SetCellValue(sheetName, scoreCell, rmData[rm])
			}

			row++
		}
	}

	// Save the spreadsheet
	err := f.SaveAs("LiftingData.xlsx")
	if err != nil {
		log.Fatalf("Failed to save Excel file: %v", err)
	}
	fmt.Println("Excel export complete: LiftingData.xlsx")
}

func main() {
	data := aggregateData()
	exportToExcel(data)
}
