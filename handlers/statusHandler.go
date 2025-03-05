package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

var startTime = time.Now()

type StatusService struct{}

func (s *StatusService) StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("StatusHandler invoked") // Debugging log

	statusResponse := StatusResponse{
		CountriesNowApi:  s.checkAPIStatus(CountriesNowAPI + "countries/"),
		RestCountriesApi: s.checkAPIStatus(RestCountriesAPI + "all"),
		Version:          VERSION,
		Uptime:           strconv.FormatFloat(time.Since(startTime).Seconds(), 'f', 2, 64) + "s",
	}

	log.Printf("Response: %+v\n", statusResponse) // Log response before sending

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *StatusService) checkAPIStatus(apiURL string) int {
	log.Println("Checking API status for:", apiURL) // Debugging log

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println("Error contacting API:", apiURL, err) // Debugging log
		return http.StatusServiceUnavailable
	}
	defer resp.Body.Close()

	log.Println("API response for", apiURL, ":", resp.StatusCode) // Debugging log
	return resp.StatusCode
}
