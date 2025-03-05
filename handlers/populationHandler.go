package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPopulationRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Only 'GET' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleGetPopulationRequest(w http.ResponseWriter, r *http.Request) {
	// Adjusted URL Path Parsing: Split the URL into parts.
	parts := strings.Split(r.URL.Path, "/")

	// Ensure there are enough parts in the URL and that we have a 2-character ISO code.
	if len(parts) < 5 || len(parts[4]) != 2 {
		log.Printf("Invalid or missing country code in URL: %v", parts)
		http.Error(w, "Invalid or missing country code in URL", http.StatusBadRequest)
		return
	}

	// Extract the ISO code (e.g., 'no' for Norway).
	isoCode := strings.ToUpper(parts[4])
	countryName, fetchErr := fetchCountryName(isoCode) // Use fetchErr for error handling
	if fetchErr != nil {
		log.Printf("No country mapping found for ISO code: %s", isoCode)
		http.Error(w, "Invalid ISO code. Use a valid 2-letter ISO code.", http.StatusBadRequest)
		return
	}
	log.Printf("Fetching population data for country: %s", countryName)

	// Handle optional query parameter for limiting the years.
	limitQuery := r.URL.Query().Get("limit")
	var startYear, endYear int
	var err error
	if limitQuery != "" {
		limits := strings.Split(limitQuery, "-")
		if len(limits) != 2 {
			log.Printf("Invalid limit format: %s", limitQuery)
			http.Error(w, "Invalid year limit format", http.StatusBadRequest)
			return
		}
		startYear, err = strconv.Atoi(limits[0])
		if err != nil || startYear < 0 {
			log.Printf("Invalid start year: %s", limits[0])
			http.Error(w, "Invalid start year", http.StatusBadRequest)
			return
		}
		endYear, err = strconv.Atoi(limits[1])
		if err != nil || endYear < startYear {
			log.Printf("Invalid end year: %s", limits[1])
			http.Error(w, "Invalid end year", http.StatusBadRequest)
			return
		}
	}

	// Fetch population data from external API
	populationURL := CountriesNowAPI + "countries/population"
	requestBody, _ := json.Marshal(map[string]interface{}{"country": countryName})

	req, err := http.NewRequest("POST", populationURL, strings.NewReader(string(requestBody)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Server error while creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request to the external API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request to CountriesNowAPI: %v", err)
		http.Error(w, "Server error while making request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("CountriesNow API response status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		http.Error(w, "Failed to fetch population data", http.StatusNotFound)
		return
	}

	// Decode the population data response
	var populationData GetPopulationData
	if err := json.NewDecoder(resp.Body).Decode(&populationData); err != nil {
		log.Printf("Error decoding response: %v", err)
		http.Error(w, "Error decoding population data", http.StatusInternalServerError)
		return
	}

	// Filter and calculate the population data
	filteredValues := []map[string]int{}
	var total int
	var count int

	for _, entry := range populationData.Data.PopulationInfo {
		year, exists := entry["year"]
		if !exists {
			continue
		}
		// Filter based on the specified year range (if provided)
		if limitQuery != "" && (year < startYear || year > endYear) {
			continue
		}
		filteredValues = append(filteredValues, entry)
		total += entry["value"]
		count++
	}

	// Calculate the mean population if there are valid entries
	mean := 0
	if count > 0 {
		mean = total / count
	}

	// Prepare the response data
	response := PopulationInfoResponse{
		Mean:   mean,
		Values: filteredValues,
	}

	// Set the response headers and send the data
	w.Header().Add("Content-Type", JsonHeader)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func fetchCountryName(isoCode string) (string, error) {
	restCountriesURL := RestCountriesAPI + "alpha/" + isoCode

	resp, err := http.Get(restCountriesURL)
	if err != nil {
		return "", fmt.Errorf("error making request to RestCountries API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("RestCountries API returned status: %d", resp.StatusCode)
	}

	// Log the response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	log.Printf("Response from RestCountries API: %s", bodyBytes)

	// The API returns an array of objects, so we must parse it as an array
	var countries []struct {
		Name struct {
			Common   string `json:"common"`
			Official string `json:"official"`
		} `json:"name"`
	}

	if err := json.Unmarshal(bodyBytes, &countries); err != nil {
		return "", fmt.Errorf("error decoding RestCountries response: %v", err)
	}

	// Ensure we received at least one country
	if len(countries) == 0 {
		return "", fmt.Errorf("no country data found for ISO code: %s", isoCode)
	}

	return countries[0].Name.Common, nil
}
