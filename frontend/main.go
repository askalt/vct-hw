package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Event struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Date        string   `json:"date"`
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var backendHost = getEnv("BACKEND_HOST", "localhost")
var backendPort = getEnv("BACKEND_PORT", "8080")

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/api/events", eventsHandler)
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/feed", serveFeed)
	http.HandleFunc("/create", serveCreate)

	log.Println("Frontend server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed", http.StatusFound)
}

func serveFeed(w http.ResponseWriter, r *http.Request) {
	serveStaticFile(w, r, "feed.html")
}

func serveCreate(w http.ResponseWriter, r *http.Request) {
	serveStaticFile(w, r, "create.html")
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, filename string) {
	filePath := filepath.Join("static", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "GET":
		getEvents(w, r)
	case "POST":
		createEvent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func backendUri() string {
	return fmt.Sprintf("http://%s:%s", backendHost, backendPort)
}

func eventsUri() string {
	return fmt.Sprintf("%s/events", backendUri())
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(eventsUri())
	if err != nil {
		http.Error(w, "Failed to connect to backend", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Forward to backend
	jsonData, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Failed to marshal event", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(eventsUri(), "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
