package controller

import (
	"encoding/json"
	"fmt"
	"kai-app/api/models"
	"kai-app/api/service"
	"kai-app/arch/database"
	"net/http"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type ScanRequest struct {
		Repo  string   `json:"repo"`
		Files []string `json:"files"`
	}

	var req ScanRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Initialize Downloader
	downloader, err := service.NewDownloader(req.Repo, "main")
	if err != nil {
		http.Error(w, "Failed to initialize downloader", http.StatusInternalServerError)
		return
	}

	downloader.DownloadFilesConcurrently(req.Files)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Files processed and saved successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// QueryHandler handles the vulnerability query requests
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// type QueryRequest struct {
	// 	Filters struct {
	// 		Severity string `json:"severity"`
	// 	} `json:"filters"`
	// }

	type QueryRequest struct {
		Filters map[string]interface{} `json:"filters"` // Accept any filter dynamically
	}

	var req QueryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Initialize DB connection
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	query := db.Model(&models.Vulnerability{})

	for key, value := range req.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Query the database based on filters
	var vulnerabilities []models.Vulnerability
	if err := query.Find(&vulnerabilities).Error; err != nil {
		http.Error(w, "Failed to fetch vulnerabilities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vulnerabilities)
}
