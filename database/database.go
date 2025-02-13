package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error

	// Gunakan PostgreSQL atau MySQL sesuai kebutuhan
	dsn := "host=localhost user=postgres password= dbname=zoomdb port=5432 sslmode=disable" // Set variabel ini di .env
	dialect := "postgres"                                                                   // "postgres" atau "mysql"

	if dialect == "postgres" {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else if dialect == "mysql" {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		log.Fatal("Unsupported database dialect")
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected successfully")
}
