package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"
)

// Haversine calculates distance between two coordinates
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth radius in km
    φ1 := lat1 * math.Pi / 180
    φ2 := lat2 * math.Pi / 180
    Δφ := (lat2 - lat1) * math.Pi / 180
    Δλ := (lon2 - lon1) * math.Pi / 180

    a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}

// GeoCodeResponse structure for OpenCage API
type GeoCodeResponse struct {
    Results []struct {
        Geometry struct {
            Lat float64 `json:"lat"`
            Lng float64 `json:"lng"`
        } `json:"geometry"`
    } `json:"results"`
}

// Geocode address using OpenCage API
func Geocode(address, apiKey string) (float64, float64, error) {
    client := &http.Client{Timeout: 2 * time.Second}
    url := "https://api.opencagedata.com/geocode/v1/json?q=" + 
           url.QueryEscape(address) + "&key=" + apiKey

    resp, err := client.Get(url)
    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response body: %v", err)
	}
	
	var result GeoCodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to decode JSON: %v", err)
	}
	
	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	// fmt.Println("response : ");
	// fmt.Println(result.Results[0].Geometry.Lat);
	// fmt.Println(result.Results[0].Geometry.Lng);
	
	return result.Results[0].Geometry.Lat, result.Results[0].Geometry.Lng, nil
}