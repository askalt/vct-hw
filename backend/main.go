package main

import (
	"event-server/db"
	"event-server/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize a database.
	db.InitDB()
	defer db.DB.Close()

	// Create a router.
	router := mux.NewRouter()

	// Routes.
	router.HandleFunc("/events", handlers.CreateEvent).Methods("POST")
	router.HandleFunc("/events", handlers.GetEvents).Methods("GET")
	router.HandleFunc("/events/{event_id}", handlers.UpdateEvent).Methods("PUT")
	router.HandleFunc("/events", handlers.DeleteEvents).Methods("DELETE")
	router.HandleFunc("/register/{event_id}", handlers.RegisterForEvent).Methods("PUT")

	// Start a server.
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
