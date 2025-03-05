package handlers

import (
	"fmt"
	"net/http"
	"prog2005_assignment1/constants"
)

func EmptyHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "This service does not provide any functionality on root path level. Please use paths <a href=\"" +
		constants.INFO_PATH + "\">" + constants.INFO_PATH + "</a> or <a href=\"" +
		constants.POPULATION_PATH + "\">" + constants.POPULATION_PATH + "</a>" //+
	//STATUS_PATH + "\">" + STATUS_PATH + "</a>."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}

}
