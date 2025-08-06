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

// getUsersHandler handles requests to the /users endpoint.
// func getUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	// Only allow GET requests.
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Hardcoded list of users. In a real application, this would come from a database.
// 	users := []User{
// 		{ID: 1, Name: "Alice", City: "New York"},
// 		{ID: 2, Name: "Bob", City: "London"},
// 		{ID: 3, Name: "Charlie", City: "Paris"},
// 	}

// 	// Set the Content-Type header to indicate that we're sending JSON.
// 	w.Header().Set("Content-Type", "application/json")

// 	// Encode the users slice to JSON and write it to the response writer.
// 	if err != nil {
// 		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
// 		return
// 	}
// }

// helloHandler handles requests to the /hello endpoint.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	mux.HandleFunc("/hello", helloHandler)

	// Define the server address and port.
	addr := "localhost:8080"

	log.Printf("Starting server on %s\n", addr)
	// Start the server and listen for incoming requests.
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
