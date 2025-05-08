package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Konfigurasi database
const (
	DbUser     = "dzaki"  // Ganti dengan username PostgreSQL kamu
	DbPassword = "password" // Ganti dengan password PostgreSQL kamu
	DbName     = "lasharan" // Ganti dengan nama database PostgreSQL kamu
)

// ConnectToDB membuka koneksi ke PostgreSQL dan mengembalikan objek *sql.DB
func ConnectToDB() (*sql.DB, error) {
	// Membuat URL koneksi PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DbUser, DbPassword, DbName)

	// Membuka koneksi ke database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return nil, err
	}

	// Mengecek koneksi
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}
