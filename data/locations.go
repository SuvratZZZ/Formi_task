package data

import (
	"encoding/csv"
	"log"
	"os"
)

type IndianLocation struct {
	Name string
	Type string
}

func LoadIndianLocations(path string) []IndianLocation {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening locations file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	var locations []IndianLocation
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		locations = append(locations, IndianLocation{
			Name: record[0],
			Type: record[1],
		})
	}
	return locations
}