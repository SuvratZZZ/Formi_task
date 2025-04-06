// 

package handlers

import (
	"encoding/csv"
	"formi/data"
	"formi/models"
	"formi/utils"
	"os"
	"strconv"
)

type PropertyHandler struct {
	Properties   []models.Property
	GeoAPIKey    string
	FuzzyMatcher *utils.FuzzyMatcher
}

func NewPropertyHandler(propertiesPath, locationsPath, apiKey string) (*PropertyHandler, error) {
	// Load properties
	propFile, err := os.Open(propertiesPath)
	if err != nil {
		return nil, err
	}
	defer propFile.Close()

	propReader := csv.NewReader(propFile)
	propRecords, err := propReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var properties []models.Property
	for i, record := range propRecords {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(record[1], 64)
		lon, _ := strconv.ParseFloat(record[2], 64)
		properties = append(properties, models.Property{
			Name:      record[0],
			Latitude:  lat,
			Longitude: lon,
		})
	}

	// Load locations
	locations := data.LoadIndianLocations(locationsPath)
	fuzzyMatcher := utils.NewFuzzyMatcher(locations)

	return &PropertyHandler{
		Properties:   properties,
		GeoAPIKey:    apiKey,
		FuzzyMatcher: fuzzyMatcher,
	}, nil
}

func (h *PropertyHandler) FindNearby(lat, lng float64) []models.Result {
	var results []models.Result
	for _, prop := range h.Properties {
		distance := utils.Haversine(lat, lng, prop.Latitude, prop.Longitude)
		if distance <= 50 {
			results = append(results, models.Result{
				Name:     prop.Name,
				Distance: distance,
			})
		}
	}
	return results
}