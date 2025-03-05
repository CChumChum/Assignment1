package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"prog2005_assignment1/constants"
	"prog2005_assignment1/structs"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetInfoRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleGetInfoRequest(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || len(parts[4]) != 2 {
		log.Printf("Invalid or missing country code in URL: %v", parts)
		http.Error(w, "Invalid or missing country code in URL", http.StatusBadRequest)
		return
	}

	isoCode := strings.ToUpper(parts[4]) // Convert to uppercase for consistency

	// Validate the ISO code (ensure it's 2 uppercase letters)
	if !isValidIsoCode(isoCode) {
		log.Printf("Invalid ISO code: %s", isoCode)
		http.Error(w, "Invalid ISO code. Please provide a valid 2-letter ISO code.", http.StatusBadRequest)
		return
	}

	// Log the ISO code for debugging
	log.Printf("ISO Code: %s is valid", isoCode)

	// Fetch general country info
	client := &http.Client{}
	defer client.CloseIdleConnections()

	restURL := constants.RestCountriesAPI + "alpha/" + isoCode

	restCountriesRequest, err := http.NewRequest(http.MethodGet, restURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, constants.RESPONSE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	restCountriesRequest.Header.Add("Content-Type", "application/json")

	restResponse, err := client.Do(restCountriesRequest)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		http.Error(w, constants.DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	defer restResponse.Body.Close()

	// Log status code
	log.Printf("RestCountries API response status: %s", restResponse.Status)

	// Check if the status code is 200 (OK), indicating the country exists
	if restResponse.StatusCode != http.StatusOK {
		log.Printf("Country with ISO code %s not found, received status %s", isoCode, restResponse.Status)
		http.Error(w, "Country not found for the provided ISO code.", http.StatusNotFound)
		return
	}

	var countries []structs.RestCountriesResponse // Expect an array

	decoder := json.NewDecoder(restResponse.Body)
	if decodeErr := decoder.Decode(&countries); decodeErr != nil {
		log.Printf("Error decoding response: %v", decodeErr)
		http.Error(w, constants.DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	// Ensure we got at least one country in the response
	if len(countries) == 0 {
		log.Printf("Empty response from API")
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	// Convert array to a single object
	countryInfo := countries[0]

	// Log the data to see what's being returned
	log.Printf("Country Name: %s", countryInfo.Name.CountryName)
	log.Printf("Flag: %s", countryInfo.Flag.Png)
	log.Printf("Population: %d", countryInfo.Population)
	log.Printf("Capital: %v", countryInfo.Capital)
	log.Printf("Languages: %v", countryInfo.Languages)
	log.Printf("Continents: %v", countryInfo.Continents)

	// Validate the data
	if invalidRESTRequest(countryInfo) {
		log.Printf("Invalid values in restCountriesData")
		http.Error(w, constants.GENERIC_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	// Fetch cities
	countryPayload := map[string]string{"country": countryInfo.Name.CountryName}
	countriesNowBody, err := json.Marshal(countryPayload)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		http.Error(w, constants.RESPONSE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	countriesNowRequest, err := http.NewRequest(http.MethodPost, constants.CountriesNowAPI, bytes.NewBuffer(countriesNowBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, constants.DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	countriesNowRequest.Header.Add("Content-Type", "application/json")

	countriesNowResponse, err := client.Do(countriesNowRequest)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer countriesNowResponse.Body.Close()

	// Log status code
	log.Printf("CountriesNow API response status: %s", countriesNowResponse.Status)

	decoder = json.NewDecoder(countriesNowResponse.Body)
	var countriesNowData structs.Cities
	if decodeErr := decoder.Decode(&countriesNowData); decodeErr != nil {
		log.Printf("Error decoding cities JSON: %v", decodeErr.Error())
		http.Error(w, constants.DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	// Sort cities alphabetically before applying the limit
	sort.Strings(countriesNowData.Cities)

	// Apply limit if provided
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Invalid limit parameter: %v", err)
			http.Error(w, constants.INTEGER_FAULT, http.StatusBadRequest)
			return
		}

		if limit > len(countriesNowData.Cities) { // Prevent slicing beyond bounds
			limit = len(countriesNowData.Cities)
		}
		countriesNowData.Cities = countriesNowData.Cities[:limit]
	}

	// Prepare response
	infoResponse := structs.InfoResponse{
		Name:       countryInfo.Name.CountryName,
		Continents: countryInfo.Continents,
		Population: countryInfo.Population,
		Languages:  countryInfo.Languages,
		Bordering:  countryInfo.Bordering,
		Flag:       countryInfo.Flag.Png,
		Capital:    countryInfo.Capital[0],
		Cities:     countriesNowData.Cities,
	}

	// Send response
	w.Header().Add("Content-Type", constants.JsonHeader)
	encoder := json.NewEncoder(w)
	if encodeErr := encoder.Encode(infoResponse); encodeErr != nil {
		log.Printf("Error encoding response: %v", encodeErr)
		http.Error(w, constants.ENCODE_SERVER_ERROR, http.StatusInternalServerError)
	}
}

// Function to validate the ISO code (2 uppercase letters)
func isValidIsoCode(isoCode string) bool {
	if len(isoCode) != 2 {
		return false
	}
	for _, char := range isoCode {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}

// Function to validate REST response
func invalidRESTRequest(data structs.RestCountriesResponse) bool {
	return data.Name.CountryName == "" ||
		data.Flag.Png == "" ||
		data.Population <= 0 ||
		data.Capital[0] == "" ||
		len(data.Languages) == 0 ||
		len(data.Continents) == 0
}
