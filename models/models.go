// package models

// // Property represents a hotel property with geo coordinates
// type Property struct {
//     Name      string  `json:"name"`
//     Latitude  float64 `json:"latitude"`
//     Longitude float64 `json:"longitude"`
// }

// // APIResponse represents the API response structure
// type APIResponse struct {
//     Results []Result `json:"results"`
//     Error   string   `json:"error,omitempty"`
// }

// // Result represents a single search result
// type Result struct {
//     Name     string  `json:"name"`
//     Distance float64 `json:"distance_km"`
// }


package models

type Property struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type APIResponse struct {
	Results       []Result `json:"results"`
	Error         string   `json:"error,omitempty"`
	OriginalQuery string   `json:"original_query,omitempty"`
	CorrectedQuery string `json:"corrected_query,omitempty"`
	Your_latitude float64 `json:"your_lontitude,omitempty"`
	Your_longitude float64 `json:"your_longitude,omitempty"`
	QueryType     string   `json:"query_type,omitempty"`
}

type Result struct {
	Name     string  `json:"name"`
	Distance float64 `json:"distance_km"`
}