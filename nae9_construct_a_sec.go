package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Tracker represents a single tracker entity
type Tracker struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	Data      string `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// SecureMiddleware handles authentication and authorization
func SecureMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token and authenticate user
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Authorize user based on token
		if !Authorize(token) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Authorize returns true if the token is valid, false otherwise
func Authorize(token string) bool {
	// Implement token verification logic here
	return true // Replace with actual logic
}

func main() {
	// Initialize database connection
	db, err := gorm.Open("sqlite3", "./trackers.db")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	// Create table for trackers
	db.AutoMigrate(&Tracker{})

	// Initialize router
	r := mux.NewRouter()

	// Protected routes
	r.HandleFunc("/trackers", GetTrackers).Methods("GET")
	r.HandleFunc("/trackers/{id}", GetTracker).Methods("GET")
	r.HandleFunc("/trackers", CreateTracker).Methods("POST")
	r.HandleFunc("/trackers/{id}", UpdateTracker).Methods("PUT")
	r.HandleFunc("/trackers/{id}", DeleteTracker).Methods("DELETE")

	// Apply security middleware to protected routes
	r.Use(SecureMiddleware)

	// Start server
	http.ListenAndServe(":8080", r)
}

// GetTrackers returns all trackers
func GetTrackers(w http.ResponseWriter, r *http.Request) {
	var trackers []Tracker
	db.Find(&trackers)
	json.NewEncoder(w).Encode(trackers)
}

// GetTracker returns a single tracker
func GetTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tracker Tracker
	db.First(&tracker, params["id"])
	json.NewEncoder(w).Encode(tracker)
}

// CreateTracker creates a new tracker
func CreateTracker(w http.ResponseWriter, r *http.Request) {
	var tracker Tracker
	json.NewDecoder(r.Body).Decode(&tracker)
	db.Create(&tracker)
	w.WriteHeader(http.StatusCreated)
}

// UpdateTracker updates an existing tracker
func UpdateTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tracker Tracker
	db.First(&tracker, params["id"])
	json.NewDecoder(r.Body).Decode(&tracker)
	db.Save(&tracker)
}

// DeleteTracker deletes a tracker
func DeleteTracker(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tracker Tracker
	db.First(&tracker, params["id"])
	db.Delete(&tracker)
}