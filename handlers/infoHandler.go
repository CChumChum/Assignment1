package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Entry point handler for Location information
*/
func InfoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' and '"+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}

}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {

	// TODO: Check for content type

	// Instantiate decoder
	decoder := json.NewDecoder(r.Body)
	// Ensure parser fails on unknown fields (baseline way of detecting different structs than expected ones)
	// Note: This does not lead to a check whether an actually provided field is empty!
	decoder.DisallowUnknownFields()

	// Prepare empty struct to populate
	info := Info{}

	// Decode location instance --> Alternative: "err := json.NewDecoder(r.Body).Decode(&location)"
	err := decoder.Decode(&info)
	if err != nil {
		// Note: more often than not is this error due to client-side input, rather than server-side issues
		http.Error(w, "Error during decoding: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation of input (Golang does not do that itself :()

	// TODO: Write convenience function for validation

	if info.Name == "" {
		http.Error(w, "Invalid input: Field 'Name' is empty.", http.StatusBadRequest)
		return
	}

	if info.Continent == "" {
		http.Error(w, "Invalid input: Field 'Continent' not found.", http.StatusBadRequest)
		return
	}

	// Field country is not required, hence no check

	if info.Population == "" {
		http.Error(w, "Invalid input: Field 'Population' not found.", http.StatusBadRequest)
		return
	}

	if info.Languages == "" {
		http.Error(w, "Invalid input: Field 'Languages' not found.", http.StatusBadRequest)
		return
	}

	// Flat printing
	fmt.Println("Received following request:")
	fmt.Println(info)

	// Pretty printing
	output, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}

	fmt.Println("Pretty printing:")
	fmt.Println(string(output))

	// TODO: Handle content (e.g., writing to DB, process, etc.)

	// Return status code (good practice)
	http.Error(w, "OK", http.StatusOK)
}

/*
Dedicated handler for GET requests
*/
func handleGetRequest(w http.ResponseWriter, r *http.Request) {

	// Create instance of content (could be read from DB, file, etc.)
	info := Info{
		Name:       "Norway",
		Continent:  "Europe",
		Population: "5367580",
		Languages:  "Norwegian",
		Bordering:  "Sweden, Finland, Russia",
		Flag:       "https://upload.wikimedia.org/wikipedia/commons/d/d9/Flag_of_Norway.svg",
		Capital:    "Oslo",
		Cities:     "Bergen, Trondheim, Stavanger",
	}

	// Write content type header (best practice)
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	// Encode specific content --> Alternative: "err := json.NewEncoder(w).Encode(location)"
	err := encoder.Encode(info)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()

}
