package api

import (
	"encoding/json"
	"log"
	"net/http"
	"myapp/internal/db"
	"myapp/internal/models"
	"github.com/gorilla/mux"
    "database/sql"
	"fmt"
	"strconv"
	"strings"
)

// GetGuruHandler - Mendapatkan semua data guru
func GetGuruHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_guru, id_user, id_mapel, nama_guru, mata_pelajaran, nip, alamat, email, no_telp FROM guru")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var gurus []models.Guru
	for rows.Next() {
		var guru models.Guru
		if err := rows.Scan(&guru.IDGuru, &guru.IDUser, &guru.IDMapel, &guru.NamaGuru, &guru.MataPelajaran, &guru.NIP, &guru.Alamat, &guru.Email, &guru.NoTelp); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		gurus = append(gurus, guru)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gurus)
}

// CreateGuruHandler - Menambahkan data guru baru
func CreateGuruHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var guru models.Guru
	err = json.NewDecoder(r.Body).Decode(&guru)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO guru (id_user, id_mapel, nama_guru, mata_pelajaran, nip, alamat, email, no_telp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = database.Exec(query, guru.IDUser, guru.IDMapel, guru.NamaGuru, guru.MataPelajaran, guru.NIP, guru.Alamat, guru.Email, guru.NoTelp)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Guru berhasil ditambahkan:", guru)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Guru berhasil ditambahkan"))
}

// UpdateGuruHandler - Mengupdate data guru berdasarkan ID
func UpdateGuruHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var guru models.Guru
	if err := json.NewDecoder(r.Body).Decode(&guru); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE guru SET id_user=$1, id_mapel=$2, nama_guru=$3, mata_pelajaran=$4, nip=$5, alamat=$6, email=$7, no_telp=$8 WHERE id_guru=$9",
		guru.IDUser, guru.IDMapel, guru.NamaGuru, guru.MataPelajaran, guru.NIP, guru.Alamat, guru.Email, guru.NoTelp, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guru)
}

// DeleteGuruHandler - Menghapus data guru berdasarkan ID
func DeleteGuruHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting guru with ID:", id) // Log ID yang akan dihapus

	// Menghapus guru dari database berdasarkan id_guru
	result, err := database.Exec("DELETE FROM guru WHERE id_guru=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No guru found with ID:", id)
		http.Error(w, "Guru not found", http.StatusNotFound)
		return
	}

	log.Println("Guru successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Guru berhasil dihapus"))
}

// GetGuruByIDHandler - Mendapatkan data guru berdasarkan ID
func GetGuruByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data guru berdasarkan ID
	var guru models.Guru
	err = database.QueryRow("SELECT id_guru, id_user, id_mapel, nama_guru, mata_pelajaran, nip, alamat, email, no_telp FROM guru WHERE id_guru=$1", id).
		Scan(&guru.IDGuru, &guru.IDUser, &guru.IDMapel, &guru.NamaGuru, &guru.MataPelajaran, &guru.NIP, &guru.Alamat, &guru.Email, &guru.NoTelp)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Guru not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data guru dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guru)
}

// GetSiswaHandler - Mendapatkan semua data guru
func GetSiswaHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_siswa, id_user, id_kelas, nama_siswa, alamat, tanggal_lahir, nisn FROM siswa")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var siswas []models.Siswa
	for rows.Next() {
		var siswa models.Siswa
		if err := rows.Scan(&siswa.IDSiswa, &siswa.IDUser, &siswa.IDKelas, &siswa.NamaSiswa, &siswa.Alamat, &siswa.TanggalLahir, &siswa.NISN); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		siswas = append(siswas, siswa)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siswas)
}

// CreateSiswaHandler - Menambahkan data siswa baru
func CreateSiswaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var siswa models.Siswa
	err = json.NewDecoder(r.Body).Decode(&siswa)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO siswa (id_user, id_kelas, nama_siswa, alamat, tanggal_lahir, nisn) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = database.Exec(query, siswa.IDUser, siswa.IDKelas, siswa.NamaSiswa, siswa.Alamat, siswa.TanggalLahir, siswa.NISN)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Siswa berhasil ditambahkan:", siswa)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Siswa berhasil ditambahkan"))
}

// UpdateSiswaHandler - Mengupdate data siswa berdasarkan ID
func UpdateSiswaHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var siswa models.Siswa
	if err := json.NewDecoder(r.Body).Decode(&siswa); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE siswa SET id_user=$1, id_kelas=$2, nama_siswa=$3, alamat=$4, tanggal_lahir=$5, nisn=$6 WHERE id_siswa=$7",
		siswa.IDUser, siswa.IDKelas, siswa.NamaSiswa, siswa.Alamat, siswa.TanggalLahir, siswa.NISN, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(siswa)
}

// DeleteSiswaHandler - Menghapus data siswa berdasarkan ID
func DeleteSiswaHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting siswa with ID:", id) // Log ID yang akan dihapus

	// Menghapus guru dari database berdasarkan id_guru
	result, err := database.Exec("DELETE FROM siswa WHERE id_siswa=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No siswa found with ID:", id)
		http.Error(w, "Siswa not found", http.StatusNotFound)
		return
	}

	log.Println("Siswa successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Guru berhasil dihapus"))
}

// GetSiswayIDHandler - Mendapatkan data siswa berdasarkan ID
func GetSiswaByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data siswa berdasarkan ID
	var siswa models.Siswa
	err = database.QueryRow("SELECT id_siswa,id_user, id_kelas, nama_siswa, alamat, tanggal_lahir, nisn FROM siswa WHERE id_siswa=$1", id).
		Scan(&siswa.IDSiswa, &siswa.IDUser, &siswa.IDKelas, &siswa.NamaSiswa, &siswa.Alamat, &siswa.TanggalLahir, &siswa.NISN)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Siswa not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data guru dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siswa)
}

// GetKelasHandler - Mendapatkan semua data kelas
func GetKelasHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_kelas, id_guru, nama_kelas, tahun_ajaran FROM kelas")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var kelass []models.Kelas
	for rows.Next() {
		var kelas models.Kelas
		if err := rows.Scan(&kelas.IDKelas, &kelas.IDGuru, &kelas.NamaKelas, &kelas.TahunAjaran); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		kelass = append(kelass, kelas)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kelass)
}

// CreateKelasHandler - Menambahkan data kelas baru
func CreateKelasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var kelas models.Kelas
	err = json.NewDecoder(r.Body).Decode(&kelas)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO kelas (id_guru, nama_kelas, tahun_ajaran) VALUES ($1, $2, $3)"
	_, err = database.Exec(query, kelas.IDGuru, kelas.NamaKelas, kelas.TahunAjaran)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Kelas berhasil ditambahkan:", kelas)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Kelas berhasil ditambahkan"))
}

// UpdateKelasHandler - Mengupdate data kelas berdasarkan ID
func UpdateKelasHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var kelas models.Kelas
	if err := json.NewDecoder(r.Body).Decode(&kelas); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE kelas SET id_guru=$1, nama_kelas=$2, tahun_ajaran=$3 WHERE id_kelas=$4",
		kelas.IDGuru, kelas.NamaKelas, kelas.TahunAjaran, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kelas)
}

// DeleteKelasHandler - Menghapus data kelas berdasarkan ID
func DeleteKelasHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting kelas with ID:", id) // Log ID yang akan dihapus

	// Menghapus kelas dari database berdasarkan id_guru
	result, err := database.Exec("DELETE FROM kelas WHERE id_kelas=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No kelas found with ID:", id)
		http.Error(w, "Kelas not found", http.StatusNotFound)
		return
	}

	log.Println("Kelas successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Kelas berhasil dihapus"))
}

// GetKelasByIDHandler - Mendapatkan data kelas berdasarkan ID
func GetKelasByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data guru berdasarkan ID
	var kelas models.Kelas
	err = database.QueryRow("SELECT id_kelas, id_guru, nama_kelas, tahun_ajaran FROM kelas WHERE id_kelas=$1", id).
		Scan(&kelas.IDKelas, &kelas.IDGuru, &kelas.NamaKelas, &kelas.TahunAjaran)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Guru not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data guru dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kelas)
}

// GetMataPelajaranHandler - Mendapatkan semua data mata pelajaran
func GetMataPelajaranHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_mapel, id_kelas, nama_mata_pelajaran FROM mata_pelajaran")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mataPelajaran []models.MataPelajaran
	for rows.Next() {
		var mp models.MataPelajaran
		if err := rows.Scan(&mp.IDMapel, &mp.IDKelas, &mp.NamaMataPelajaran); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		mataPelajaran = append(mataPelajaran, mp)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mataPelajaran)
}

// CreateMataPelajaranHandler - Menambahkan data mata pelajaran baru
func CreateMataPelajaranHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var mataPelajaran models.MataPelajaran
	err = json.NewDecoder(r.Body).Decode(&mataPelajaran)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO mata_pelajaran (id_kelas, nama_mata_pelajaran) VALUES ($1,$2)"
	_, err = database.Exec(query, mataPelajaran.IDKelas,mataPelajaran.NamaMataPelajaran)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Mata Pelajaran berhasil ditambahkan:", mataPelajaran)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Mata Pelajaran berhasil ditambahkan"))
}

// UpdateMataPelajaranHandler - Mengupdate data mata pelajaran berdasarkan ID
func UpdateMataPelajaranHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var mataPelajaran models.MataPelajaran
	if err := json.NewDecoder(r.Body).Decode(&mataPelajaran); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		"UPDATE mata_pelajaran SET id_kelas=$1, nama_mata_pelajaran=$2 WHERE id_mapel=$3",
		 mataPelajaran.IDKelas, mataPelajaran.NamaMataPelajaran, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mataPelajaran)
}

// DeleteMataPelajaranHandler - Menghapus data mata pelajaran berdasarkan ID
func DeleteMataPelajaranHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting mata pelajaran with ID:", id) // Log ID yang akan dihapus

	// Menghapus mata pelajaran dari database berdasarkan id_mapel
	result, err := database.Exec("DELETE FROM mata_pelajaran WHERE id_mapel=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No mata pelajaran found with ID:", id)
		http.Error(w, "Mata Pelajaran not found", http.StatusNotFound)
		return
	}

	log.Println("Mata Pelajaran successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mata Pelajaran berhasil dihapus"))
}

// GetMataPelajaranByIDHandler - Mendapatkan data mata pelajaran berdasarkan ID
func GetMataPelajaranByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data mata pelajaran berdasarkan ID
	var mataPelajaran models.MataPelajaran
	err = database.QueryRow("SELECT id_mapel, id_kelas,nama_mata_pelajaran FROM mata_pelajaran WHERE id_mapel=$1", id).
		Scan(&mataPelajaran.IDMapel,&mataPelajaran.IDKelas, &mataPelajaran.NamaMataPelajaran)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Mata Pelajaran not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data mata pelajaran dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mataPelajaran)
}

// Handler
func GetMataPelajaranByKelasHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idKelas := vars["id"]

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query("SELECT id_mapel, id_kelas, nama_mata_pelajaran FROM mata_pelajaran WHERE id_kelas = $1", idKelas)
    if err != nil {
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var results []models.MataPelajaran
    for rows.Next() {
        var mp models.MataPelajaran
        if err := rows.Scan(&mp.IDMapel, &mp.IDKelas, &mp.NamaMataPelajaran); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        results = append(results, mp)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Login attempt:", creds.Username, creds.Password) // ✅ Log input login

	conn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var user models.User
	query := `SELECT id_user, id_role, username, password FROM "user" WHERE username=$1`
	err = conn.QueryRow(query, creds.Username).Scan(&user.IDUser, &user.IDRole, &user.Username, &user.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if creds.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"id_user": user.IDUser,
		"id_role": user.IDRole,
		"token":   "dummy-token", // atau JWT jika diperlukan
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// Handler
func GetKelasByGuru(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idGuruStr := vars["id_guru"]
    log.Println("id_guru dari URL:", idGuruStr) // debug log

    idGuru, err := strconv.Atoi(idGuruStr)
    if err != nil {
        http.Error(w, "id_guru harus berupa angka", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query("SELECT id_kelas, id_guru, nama_kelas, tahun_ajaran FROM kelas WHERE id_guru = $1", idGuru)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var kelasList []models.Kelas
    for rows.Next() {
        var k models.Kelas
        if err := rows.Scan(&k.IDKelas, &k.IDGuru, &k.NamaKelas, &k.TahunAjaran); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        kelasList = append(kelasList, k)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(kelasList)
}


// Handler untuk mendapatkan id_guru berdasarkan id_user
func GetGuruByUserIDHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idUserStr := vars["id_user"]

    // Konversi id_user dari string ke integer
    idUser, err := strconv.Atoi(idUserStr)
    if err != nil {
        http.Error(w, "id_user harus berupa angka", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    // Query untuk mengambil id_guru berdasarkan id_user
    var idGuru int
    err = dbConn.QueryRow("SELECT id_guru FROM guru WHERE id_user = $1", idUser).Scan(&idGuru)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Guru tidak ditemukan", http.StatusNotFound)
        } else {
            http.Error(w, "Query error", http.StatusInternalServerError)
        }
        return
    }

    // Kirim id_guru sebagai response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"id_guru": idGuru})
}

// Handler untuk mendapatkan id_siswa berdasarkan id_user
func GetSiswaByUserIDHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idUserStr := vars["id_user"]

    // Konversi id_user dari string ke integer
    idUser, err := strconv.Atoi(idUserStr)
    if err != nil {
        http.Error(w, "id_user harus berupa angka", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    // Query untuk mengambil id_siswa berdasarkan id_user
    var idSiswa int
    err = dbConn.QueryRow("SELECT id_siswa FROM siswa WHERE id_user = $1", idUser).Scan(&idSiswa)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Siswa tidak ditemukan", http.StatusNotFound)
        } else {
            http.Error(w, "Query error", http.StatusInternalServerError)
        }
        return
    }

    // Kirim id_siswa sebagai response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"id_siswa": idSiswa})
}

func GetKelasWithSubjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idKelasStr := vars["id_kelas"]

	idKelas, err := strconv.Atoi(idKelasStr)
	if err != nil {
		http.Error(w, "ID kelas tidak valid", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal konek database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Ambil data kelas
	var kelas models.Kelas
	err = dbConn.QueryRow(`
		SELECT id_kelas, id_guru, nama_kelas, tahun_ajaran
		FROM kelas
		WHERE id_kelas = $1
	`, idKelas).Scan(&kelas.IDKelas, &kelas.IDGuru, &kelas.NamaKelas, &kelas.TahunAjaran)
	if err != nil {
		http.Error(w, "Kelas tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil daftar mata pelajaran dari kelas ini
	rows, err := dbConn.Query(`
		SELECT id_mapel, id_kelas, nama_mata_pelajaran
		FROM mata_pelajaran
		WHERE id_kelas = $1
	`, idKelas)
	if err != nil {
		http.Error(w, "Gagal mengambil mata pelajaran", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mataPelajaranList []models.MataPelajaran
	for rows.Next() {
		var mp models.MataPelajaran
		if err := rows.Scan(&mp.IDMapel, &mp.IDKelas, &mp.NamaMataPelajaran); err == nil {
			mataPelajaranList = append(mataPelajaranList, mp)
		}
	}

	 // Ambil jumlah siswa di kelas ini
	 var jumlahSiswa int
	 err = dbConn.QueryRow(`SELECT COUNT(*) FROM siswa WHERE id_kelas = $1`, idKelas).Scan(&jumlahSiswa)
	 if err != nil {
		 log.Println("Error menghitung jumlah siswa:", err)
		 jumlahSiswa = 0 // fallback
	 }

	// Gabungkan respons
	response := map[string]interface{}{
		"id_kelas":         kelas.IDKelas,
		"id_guru":          kelas.IDGuru,
		"nama_kelas":       kelas.NamaKelas,
		"tahun_ajaran":     kelas.TahunAjaran,
		"mata_pelajaran":   mataPelajaranList,
		"jumlah_siswa":   jumlahSiswa,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func GetSiswaByKelas(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idKelasStr := vars["id_kelas"]
    idKelas, err := strconv.Atoi(idKelasStr)
    if err != nil {
        http.Error(w, "ID kelas tidak valid", http.StatusBadRequest)
        return
    }

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Gagal konek database", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    rows, err := dbConn.Query(`
        SELECT id_siswa, id_kelas, id_user, nama_siswa, alamat, tanggal_lahir 
        FROM siswa 
        WHERE id_kelas = $1`, idKelas)
    if err != nil {
        http.Error(w, "Gagal mengambil data siswa", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var siswaList []models.Siswa
    for rows.Next() {
        var siswa models.Siswa
        if err := rows.Scan(&siswa.IDSiswa, &siswa.IDKelas, &siswa.IDUser, &siswa.NamaSiswa, &siswa.Alamat, &siswa.TanggalLahir); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        siswaList = append(siswaList, siswa)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(siswaList)
}

func GetMataPelajaranBySiswaIDHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idSiswa := vars["id_siswa"]

    dbConn, err := db.ConnectToDB()
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    query := `
        SELECT mp.id_mapel, mp.id_kelas, mp.nama_mata_pelajaran
        FROM siswa s
        JOIN mata_pelajaran mp ON s.id_kelas = mp.id_kelas
        WHERE s.id_siswa = $1
    `
    rows, err := dbConn.Query(query, idSiswa)
    if err != nil {
        http.Error(w, "Query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var results []models.MataPelajaran
    for rows.Next() {
        var mp models.MataPelajaran
        if err := rows.Scan(&mp.IDMapel, &mp.IDKelas, &mp.NamaMataPelajaran); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        results = append(results, mp)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func GetSimpleSubjectDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idMapelStr := vars["id_mapel"]
	idMapel, err := strconv.Atoi(idMapelStr)
	if err != nil {
		http.Error(w, "ID mapel tidak valid", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	query := `
		SELECT mp.nama_mata_pelajaran, g.nama_guru, k.tahun_ajaran, k.id_kelas
		FROM mata_pelajaran mp
		JOIN kelas k ON mp.id_kelas = k.id_kelas
		JOIN guru g ON k.id_guru = g.id_guru
		WHERE mp.id_mapel = $1
	`

	var (
		namaMapel, namaGuru, tahunAjaran string
		idKelas                          int
	)

	err = dbConn.QueryRow(query, idMapel).Scan(&namaMapel, &namaGuru, &tahunAjaran, &idKelas)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	var jumlahSiswa int
	err = dbConn.QueryRow(`SELECT COUNT(*) FROM siswa WHERE id_kelas = $1`, idKelas).Scan(&jumlahSiswa)
	if err != nil {
		jumlahSiswa = 0
	}

	// Tambahkan id_mapel ke dalam response
	response := map[string]interface{}{
		"id_mapel":            idMapel,
		"nama_mata_pelajaran": namaMapel,
		"nama_guru":           namaGuru,
		"tahun_ajaran":        tahunAjaran,
		"jumlah_siswa":        jumlahSiswa,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetStudentsByMapelID(w http.ResponseWriter, r *http.Request) {
	// Ambil id_mapel dari query parameter
	vars := mux.Vars(r)
	idMapelStr := vars["id_mapel"]

	idMapel, err := strconv.Atoi(idMapelStr)
	if err != nil {
		http.Error(w, "ID mapel tidak valid", http.StatusBadRequest)
		return
	}

	// Koneksi ke database
	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Ambil id_kelas dari tabel mata_pelajaran
	var idKelas int
	err = dbConn.QueryRow(`SELECT id_kelas FROM mata_pelajaran WHERE id_mapel = $1`, idMapel).Scan(&idKelas)
	if err != nil {
		http.Error(w, "Mapel tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil data siswa berdasarkan id_kelas
	rows, err := dbConn.Query(`
		SELECT id_siswa, id_kelas, id_user, nama_siswa, alamat, tanggal_lahir 
		FROM siswa 
		WHERE id_kelas = $1
	`, idKelas)
	if err != nil {
		http.Error(w, "Gagal mengambil data siswa", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var siswaList []models.Siswa

	for rows.Next() {
		var s models.Siswa
		err := rows.Scan(&s.IDSiswa, &s.IDKelas, &s.IDUser, &s.NamaSiswa, &s.Alamat, &s.TanggalLahir)
		if err != nil {
			http.Error(w, "Gagal membaca data siswa", http.StatusInternalServerError)
			return
		}
		siswaList = append(siswaList, s)
	}

	// Encode hasil ke JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siswaList)
}


func GetPenilaianBySiswaAndMapelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idSiswa := r.URL.Query().Get("id_siswa")
	idMapel := r.URL.Query().Get("id_mapel")

	if idSiswa == "" || idMapel == "" {
		http.Error(w, "Missing id_siswa or id_mapel", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var idNilai int
	var totalNilai string

	err = dbConn.QueryRow(`
		SELECT id_nilai, total_nilai 
		FROM nilai 
		WHERE id_siswa = $1 AND id_mapel = $2
	`, idSiswa, idMapel).Scan(&idNilai, &totalNilai)

	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(models.PenilaianResponse{
				PenilaianList: []models.Penilaian{},
				TotalNilai:    "0",
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := dbConn.Query(`
		SELECT id_penilaian, nama_nilai, nilai, bobot 
		FROM penilaian 
		WHERE id_nilai = $1
	`, idNilai)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var penilaianList []models.Penilaian
	for rows.Next() {
		var idPenilaian int
		var namaNilai string
		var nilai int
		var bobot float64

		if err := rows.Scan(&idPenilaian, &namaNilai, &nilai, &bobot); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		penilaianList = append(penilaianList, models.Penilaian{
			IDPenilaian: idPenilaian,
			IDNilai:     idNilai, // ✅ tambahkan IDNilai dari hasil query sebelumnya
			NamaNilai:   namaNilai,
			Nilai:       nilai,
			Bobot:       fmt.Sprintf("%.2f%%", bobot*100),
			Range:       "0 - 100",
		})
	}

	response := models.PenilaianResponse{
		PenilaianList: penilaianList,
		TotalNilai:    totalNilai,
	}

	json.NewEncoder(w).Encode(response)
}



func CreatePenilaianHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var penilaian models.Penilaian
	if err := json.NewDecoder(r.Body).Decode(&penilaian); err != nil {
		http.Error(w, "Gagal membaca data", http.StatusBadRequest)
		return
	}

	// Pastikan id_nilai valid (FK constraint)
	var exists bool
	err = dbConn.QueryRow("SELECT EXISTS(SELECT 1 FROM nilai WHERE id_nilai = $1)", penilaian.IDNilai).Scan(&exists)
	if err != nil {
		http.Error(w, "Gagal mengecek id_nilai", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "id_nilai tidak ditemukan", http.StatusBadRequest)
		return
	}

	// Convert bobot string ke float (e.g. "20.00%" => 0.2)
	bobotFloat, err := strconv.ParseFloat(strings.TrimSuffix(penilaian.Bobot, "%"), 64)
	if err != nil {
		http.Error(w, "Format bobot salah", http.StatusBadRequest)
		return
	}
	bobotFloat = bobotFloat / 100

	// Insert dan ambil id_penilaian yang baru
	var idPenilaian int
	err = dbConn.QueryRow(`
		INSERT INTO penilaian (id_nilai, nama_nilai, nilai, bobot)
		VALUES ($1, $2, $3, $4)
		RETURNING id_penilaian
	`, penilaian.IDNilai, penilaian.NamaNilai, penilaian.Nilai, bobotFloat).Scan(&idPenilaian)

	if err != nil {
		http.Error(w, fmt.Sprintf("Insert error: %v", err), http.StatusInternalServerError)
		return
	}

	// Siapkan response JSON
	response := models.Penilaian{
		IDPenilaian: idPenilaian,
		IDNilai:     penilaian.IDNilai,
		NamaNilai:   penilaian.NamaNilai,
		Nilai:       penilaian.Nilai,
		Bobot:       fmt.Sprintf("%.2f%%", bobotFloat*100), // kembalikan ke format string persen
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}



func UpdatePenilaianHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"] // id_penilaian

	var penilaian models.Penilaian
	if err := json.NewDecoder(r.Body).Decode(&penilaian); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Konversi string bobot (misal: "20.00%") ke float64
	bobotFloat, err := strconv.ParseFloat(strings.TrimSuffix(penilaian.Bobot, "%"), 64)
	if err != nil {
		http.Error(w, "Format bobot salah", http.StatusBadRequest)
		return
	}
	bobotFloat = bobotFloat / 100 // Ubah jadi bentuk desimal: 20% → 0.2

	_, err = database.Exec(`
		UPDATE penilaian 
		SET nama_nilai=$1, nilai=$2, bobot=$3 
		WHERE id_penilaian=$4`,
		penilaian.NamaNilai, penilaian.Nilai, bobotFloat, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Penilaian berhasil diperbarui",
	})
}


func DeletePenilaianHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting penilaian with ID:", id)

	result, err := database.Exec("DELETE FROM penilaian WHERE id_penilaian=$1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No penilaian found with ID:", id)
		http.Error(w, "Penilaian not found", http.StatusNotFound)
		return
	}

	log.Println("Penilaian successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Penilaian berhasil dihapus"))
}

// GetPenilaianHandler - Mendapatkan semua data guru
func GetPenilaianHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_penilaian, id_nilai, nama_nilai, nilai, bobot FROM penilaian")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var penilaians []models.Penilaian
	for rows.Next() {
		var penilaian models.Penilaian
		if err := rows.Scan(&penilaian.IDPenilaian, &penilaian.IDNilai, &penilaian.NamaNilai, &penilaian.Nilai, &penilaian.Bobot); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		penilaians = append(penilaians, penilaian)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(penilaians)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query(`SELECT id_user, username, password, id_role, tanggal_registrasi FROM "user"`)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.IDUser, &user.Username, &user.Password, &user.IDRole, &user.TanggalRegistrasi); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	database, err := db.ConnectToDB()
	if err != nil {
		log.Println("DB connect error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO "user" (username, password, id_role, tanggal_registrasi) VALUES ($1, $2, $3, $4)`
	_, err = database.Exec(query, user.Username, user.Password, user.IDRole, user.TanggalRegistrasi)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("User berhasil ditambahkan:", user)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User berhasil ditambahkan"))
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, err = database.Exec(
		`UPDATE "user" SET username=$1, password=$2, id_role=$3, tanggal_registrasi=$4 WHERE id_user=$5`,
		user.Username, user.Password, user.IDRole, user.TanggalRegistrasi, id,
	)
	if err != nil {
		http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	id := mux.Vars(r)["id"]
	log.Println("Deleting user with ID:", id) // Log ID yang akan dihapus

	// Menghapus user dari database berdasarkan id_user
	result, err := database.Exec(`DELETE FROM "user" WHERE id_user=$1`, id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah baris data benar-benar terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No user found with ID:", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Println("User successfully deleted with ID:", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User berhasil dihapus"))
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari URL parameter
	id := mux.Vars(r)["id"]

	// Membuka koneksi ke database
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query untuk mendapatkan data user berdasarkan ID
	var user models.User
	err = database.QueryRow(`SELECT id_user, username, password, id_role, tanggal_registrasi FROM "user" WHERE id_user=$1`, id).
		Scan(&user.IDUser, &user.Username, &user.Password, &user.IDRole, &user.TanggalRegistrasi)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Mengirimkan data user dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}




















