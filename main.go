package main

import (
	"log"
	"net/http"
	"os"
	"prog2005_assignment1/handlers"
)

func main() {

	// Extract PORT variable from the environment variables
	port := os.Getenv("PORT")

	// Override port with default port if not provided (e.g. local deployment)
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Instantiate the router
	router := http.NewServeMux()

	// Instantiate the status Service
	statusService := &handlers.StatusService{}

	// Assign path for diagnostics handlers with different patterns
	router.HandleFunc(handlers.DEFAULT_PATH, handlers.EmptyHandler)
	router.HandleFunc(handlers.INFO_PATH, handlers.InfoHandler)
	router.HandleFunc(handlers.POPULATION_PATH, handlers.PopulationHandler)
	router.HandleFunc(handlers.STATUS_PATH, statusService.StatusHandler)

	// Start HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
