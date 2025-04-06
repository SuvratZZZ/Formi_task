// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"formi/handlers"
// 	"formi/models"
// 	"formi/utils"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	err := godotenv.Load()
//     if err != nil {
//         log.Fatal("Error loading .env file")
//     }

//     apiKey := os.Getenv("OPENCAGE_API_KEY")
//     if apiKey == "" {
//         log.Fatal("OPENCAGE_API_KEY environment variable not set")
//     }

//     propertyHandler, err := handlers.NewPropertyHandler("data/properties.csv", apiKey)
//     if err != nil {
//         log.Fatalf("Failed to initialize properties: %v", err)
//     }

//     http.HandleFunc("/api/properties", func(w http.ResponseWriter, r *http.Request) {
//         w.Header().Set("Content-Type", "application/json")

//         query := r.URL.Query().Get("q")
//         if query == "" {
//             respondError(w, "Missing query parameter 'q'", http.StatusBadRequest)
//             return
//         }

//         lat, lng, err := utils.Geocode(query, propertyHandler.GeoAPIKey)
//         if err != nil {
//             respondJSON(w, models.APIResponse{Error: "Location not found"})
//             return
//         }

//         results := propertyHandler.FindNearby(lat, lng)
//         respondJSON(w, models.APIResponse{Results: results})
//     })

//     log.Println("Server starting on :8080")
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func respondJSON(w http.ResponseWriter, data interface{}) {
//     w.Header().Set("Content-Type", "application/json")
//     if err := json.NewEncoder(w).Encode(data); err != nil {
//         log.Printf("Error encoding JSON: %v", err)
//     }
// }

// func respondError(w http.ResponseWriter, message string, code int) {
//     w.WriteHeader(code)
//     respondJSON(w, models.APIResponse{Error: message})
// }

package main

import (
	"encoding/json"
	"formi/handlers"
	"formi/models"
	"formi/utils"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	env := os.Getenv("ENV") // Set this variable in Railway to "production"

    if env != "production" {
        // Load .env file only in non-production environments
        if err := godotenv.Load(); err != nil {
            log.Println("Error loading .env file")
        }
    }


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("OPENCAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENCAGE_API_KEY environment variable not set")
	}

	propertyHandler, err := handlers.NewPropertyHandler(
		"data/properties.csv",
		"data/indian_locations.csv",
		apiKey,
	)
	if err != nil {
		log.Fatalf("Failed to initialize properties: %v", err)
	}

	http.HandleFunc("/api/properties", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		query := r.URL.Query().Get("q")
		if query == "" {
			respondError(w, "Missing query parameter 'q'", http.StatusBadRequest)
			return
		}

		corrected, locType := propertyHandler.FuzzyMatcher.FindBestMatch(query)
		if corrected == "" {
			respondJSON(w, models.APIResponse{Error: "Location not found"})
			return
		}

		lat, lng, err := utils.Geocode(corrected, propertyHandler.GeoAPIKey)
		if err != nil {
			respondJSON(w, models.APIResponse{Error: "Location not found"})
			return
		}

		results := propertyHandler.FindNearby(lat, lng)
		respondJSON(w, models.APIResponse{
			Results:       results,
			OriginalQuery: query,
			CorrectedQuery: corrected,
			Your_latitude: lat,
			Your_longitude: lng,
			QueryType:     locType,
		})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	respondJSON(w, models.APIResponse{Error: message})
}