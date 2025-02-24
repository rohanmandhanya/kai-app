package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Scan represents a vulnerability scan record
type ScanResult struct {
	ID              uint            `gorm:"primaryKey"`
	ResourceName    string          `json:"resource_name"`
	ResourceType    string          `json:"resource_type"`
	ScanID          string          `json:"scan_id" gorm:"unique"`
	ScanStatus      string          `json:"scan_status"`
	Timestamp       string          `json:"timestamp"`
	ScanMetadata    ScanMetadata    `gorm:"embedded"`
	Summary         ScanSummary     `gorm:"embedded"`
	Vulnerabilities []Vulnerability `gorm:"foreignKey:ScanResultID"`
}

// ScanMetadata stores scan configuration details -- Currently not working
type ScanMetadata struct {
	ExcludedPaths   string `json:"excluded_paths"` // Store as CSV
	PoliciesVersion string `json:"policies_version"`
	ScannerVersion  string `json:"scanner_version"`
	ScanningRules   string `json:"scanning_rules"` // Store as CSV
}

// ScanSummary stores summary of vulnerabilities
type ScanSummary struct {
	Compliant            bool `json:"compliant"`
	FixableCount         int  `json:"fixable_count"`
	TotalVulnerabilities int  `json:"total_vulnerabilities"`
	CriticalCount        int  `json:"severity_counts.CRITICAL"`
	HighCount            int  `json:"severity_counts.HIGH"`
	MediumCount          int  `json:"severity_counts.MEDIUM"`
	LowCount             int  `json:"severity_counts.LOW"`
}

// Vulnerability represents a detected vulnerability
type Vulnerability struct {
	ID             uint     `gorm:"primaryKey"`
	ScanResultID   uint     `json:"scan_result_id"` // Foreign key
	CVEID          string   `json:"id"`
	PackageName    string   `json:"package_name"`
	CurrentVersion string   `json:"current_version"`
	FixedVersion   string   `json:"fixed_version"`
	Severity       string   `json:"severity"`
	CVSS           float64  `json:"cvss"`
	Description    string   `json:"description"`
	Status         string   `json:"status"`
	PublishedDate  string   `json:"published_date"`
	Link           string   `json:"link"`
	RiskFactors    []string `json:"risk_factors" gorm:"serializer:json"`
}

func toCSV(slice []string) string {
	return fmt.Sprintf("%s", slice)
}

func InsertMultipleScans(db *gorm.DB, jsonData []byte, sourceFile string) error {

	var FileResponse struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	// Parse JSON to extract Base64 encoded content
	if err := json.Unmarshal(jsonData, &FileResponse); err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	// Decode the Base64-encoded scan results
	decodedJSON, err := base64.StdEncoding.DecodeString(FileResponse.Content)
	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	// var scanResults []dto.ScanResultsWrapper

	var scanResultsWrappers []struct {
		ScanResults ScanResult `json:"scanResults"`
	}

	// Parse the decoded JSON string into interface
	if err := json.Unmarshal(decodedJSON, &scanResultsWrappers); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// // pretty print the interface
	// prettyJSON, err := json.MarshalIndent(scanResultsWrappers, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error printing JSON: %v", err)
	// }
	// log.Println(string(prettyJSON))

	// Start a database transaction
	tx := db.Begin()

	for _, wrapper := range scanResultsWrappers {
		scanResult := wrapper.ScanResults

		// Convert slices to CSV
		// scanResult.ScanMetadata.ExcludedPaths = toCSV(scanResult.ScanMetadata.ExcludedPaths)
		// scanResult.ScanMetadata.ScanningRules = toCSV(scanResult.ScanMetadata.ScanningRules)

		// Save the scan result
		db.Create(&scanResult)

	}

	// Commit transaction
	return tx.Commit().Error
}
