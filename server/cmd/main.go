package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// User represents a user in our system.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {

}

// helloHandler handles requests to the /hello endpoint.
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Received request from %s\n", r.RemoteAddr)

	airportDataMap := GetBusyAirportData()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(airportDataMap)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to load .env file")
	}
	// Create a new ServeMux to handle different endpoints.
	mux := http.NewServeMux()

	// Register our handlers for the different URL patterns.
	mux.HandleFunc("/data", getDataHandler)
	mux.HandleFunc("/subscribe", subscribeHandler)

	// Define the server address and port.
	addr := "localhost:8080"

	log.Printf("Starting server on %s\n", addr)
	// Start the server and listen for incoming requests.
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
