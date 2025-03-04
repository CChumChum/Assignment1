package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handleGetInfoRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' and '"+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}

}

func handleGetInfoRequest(w http.ResponseWriter, r *http.Request) {

	isoCode := r.URL.Query().Get("two_letter_country_code")

	if len(isoCode) != 2 {
		log.Printf("Iso code is invalid: %v", isoCode)
		http.Error(w, "Iso code is invalid", http.StatusBadRequest)
		return
	}

	restURL := RestCountriesAPI + "alpha/" + isoCode

	client := &http.Client{}
	defer client.CloseIdleConnections()

	restCountriesRequest, err := http.NewRequest(http.MethodGet, restURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, RESPONSE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	restCountriesRequest.Header.Add("Content-Type", "application/json")

	restResponse, err := client.Do(restCountriesRequest)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		http.Error(w, DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(restResponse.Body)

	var countryInfo []RestCountriesResponse

	if decodeErr := decoder.Decode(&countryInfo); decodeErr != nil {
		log.Printf("Error decoding response: %v", decodeErr)
		http.Error(w, DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	if invalidRESTRequest(countryInfo[0]) {
		log.Printf("Invalid values in restCountriesData")
		http.Error(w, GENERIC_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	countryPayload := map[string]string{
		"country": countryInfo[0].Name.CountryName,
	}

	countriesNowBody, err := json.Marshal(countryPayload)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		http.Error(w, RESPONSE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	countriesNowRequest, err := http.NewRequest(http.MethodPost, CountriesNowAPI, bytes.NewBuffer(countriesNowBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	countriesNowRequest.Header.Add("Content-Type", "application/json")

	countriesNowResponse, err := client.Do(countriesNowRequest)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	decoder = json.NewDecoder(countriesNowResponse.Body)

	var countriesNowData Cities
	if decodeErr := decoder.Decode(&countriesNowData); decodeErr != nil {
		log.Printf("Error during decoding of countriesNowData json body: %v", decodeErr.Error())
		http.Error(w, DECODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	// Sort and limit cities
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Invalid limit parameter: %v", err)
			http.Error(w, INTEGER_FAULT, http.StatusBadRequest)
			return
		}

		if limit > len(countriesNowData.Cities) { // Prevent slicing beyond bounds
			limit = len(countriesNowData.Cities)
		}

		sort.Strings(countriesNowData.Cities)
		countriesNowData.Cities = countriesNowData.Cities[:limit]
	}

	// Prepare response
	infoResponse := InfoResponse{
		Name:       countryInfo[0].Name,
		Continents: countryInfo[0].Continents,
		Population: countryInfo[0].Population,
		Languages:  countryInfo[0].Languages,
		Bordering:  countryInfo[0].Bordering,
		Flag:       countryInfo[0].Flag.Png,
		Capital:    countryInfo[0].Capital[0],
		Cities:     cities.Cities,
	}

	w.Header().Add("Content-Type", JsonHeader)

	encoder := json.NewEncoder(w)
	if encodeErr := encoder.Encode(infoResponse); encodeErr != nil {
		log.Printf("Error encoding response: %v", encodeErr)
		http.Error(w, ENCODE_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

}

func invalidRESTRequest(data RestCountriesResponse) bool {
	return data.Name.CountryName == "" ||
		data.Flag.Png == "" ||
		data.Population <= 0 ||
		data.Capital[0] == "" ||
		len(data.Languages) == 0 ||
		len(data.Continents) == 0
}
