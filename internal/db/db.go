package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Konfigurasi database
const (
	DbHost     = "mydb-lasharan.cnyoam46yv1l.ap-southeast-3.rds.amazonaws.com"
	DbPort     = 5432
	DbUser     = "postgres"
	DbPassword = "AyamBangkok123"
	DbName     = "lasharan"
)

func ConnectToDB() (*sql.DB, error) {
	// Format connection string lengkap
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		DbHost, DbPort, DbUser, DbPassword, DbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}
