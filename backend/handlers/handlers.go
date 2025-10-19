package handlers

import (
	"encoding/json"
	"event-server/db"
	"event-server/models"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if event.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO events (name, description, tags, date) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.DB.QueryRow(query, event.Name, event.Description, pq.Array(event.Tags), event.Date).Scan(&event.ID)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating event: %v", err))
		http.Error(w, "Error creating event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, description, tags, date FROM events ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, pq.Array(&event.Tags), &event.Date)
		if err != nil {
			http.Error(w, "Error scanning events", http.StatusInternalServerError)
			return
		}
		events = append(events, event)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["event_id"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event models.Event
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `UPDATE events SET name=$1, description=$2, tags=$3 WHERE id=$4`
	result, err := db.DB.Exec(query, event.Name, event.Description, event.Tags, eventID)
	if err != nil {
		http.Error(w, "Error updating event", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	event.ID = eventID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func DeleteEvents(w http.ResponseWriter, r *http.Request) {
	result, err := db.DB.Exec("DELETE FROM events")
	if err != nil {
		http.Error(w, "Error deleting events", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":      "All events deleted successfully",
		"rows_deleted": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterForEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["event_id"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Check if event exists.
	var eventExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", eventID).Scan(&eventExists)
	if err != nil {
		http.Error(w, "Error checking event", http.StatusInternalServerError)
		return
	}

	if !eventExists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	var registration models.Registration
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &registration); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if registration.UserName == "" || registration.Email == "" {
		http.Error(w, "User name and email are required", http.StatusBadRequest)
		return
	}

	registration.EventID = eventID

	query := `INSERT INTO registrations (event_id, user_name, email) VALUES ($1, $2, $3) RETURNING id`
	err = db.DB.QueryRow(query, registration.EventID, registration.UserName, registration.Email).Scan(&registration.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"registrations_event_id_email_key\"" {
			http.Error(w, "User already registered for this event", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating registration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registration)
}
