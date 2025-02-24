package database

import (
	"fmt"

	"kai-app/api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDB sets up the SQLite database connection and returns the instance
func InitializeDB() (*gorm.DB, error) {

	// dsn := "host=localhost user=root password=password dbname=postgres port=5432 sslmode=disable TimeZone=UTC"

	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("vulnerabilities.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// add migrations here:
	db.AutoMigrate(&models.ScanResult{},
		&models.Vulnerability{},
		&models.ScanSummary{},
		&models.ScanMetadata{})

	return db, nil
}

// ConnectDB establishes a connection to the SQLite database.
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("vulnerabilities.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// DisconnectDB closes the database connection.
func DisconnectDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	fmt.Println("Database disconnected successfully")
	return nil
}
