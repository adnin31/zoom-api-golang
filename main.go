package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

// Home route handler
func homeHandler(write http.ResponseWriter, r *http.Request) {
	write.Header().Set("Content-Type", "application/json")

	json.NewEncoder(write).Encode(Response{Message: "Test API!"})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET")

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
