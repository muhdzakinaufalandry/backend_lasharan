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
)

// GetGuruHandler - Mendapatkan semua data guru
func GetGuruHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.ConnectToDB()
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id_guru, id_user, id_mapel, nama_guru, mata_pelajaran FROM guru")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var gurus []models.Guru
	for rows.Next() {
		var guru models.Guru
		if err := rows.Scan(&guru.IDGuru, &guru.IDUser, &guru.IDMapel, &guru.NamaGuru, &guru.MataPelajaran); err != nil {
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

	query := "INSERT INTO guru (id_user, id_mapel, nama_guru, mata_pelajaran) VALUES ($1, $2, $3, $4)"
	_, err = database.Exec(query, guru.IDUser, guru.IDMapel, guru.NamaGuru, guru.MataPelajaran)
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
		"UPDATE guru SET id_user=$1, id_mapel=$2, nama_guru=$3, mata_pelajaran=$4 WHERE id_guru=$5",
		guru.IDUser, guru.IDMapel, guru.NamaGuru, guru.MataPelajaran, id,
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
	err = database.QueryRow("SELECT id_guru, id_user, id_mapel, nama_guru, mata_pelajaran FROM guru WHERE id_guru=$1", id).
		Scan(&guru.IDGuru, &guru.IDUser, &guru.IDMapel, &guru.NamaGuru, &guru.MataPelajaran)
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

	rows, err := database.Query("SELECT id_siswa, id_user, id_kelas, nama_siswa, alamat, tanggal_lahir FROM siswa")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var siswas []models.Siswa
	for rows.Next() {
		var siswa models.Siswa
		if err := rows.Scan(&siswa.IDSiswa, &siswa.IDUser, &siswa.IDKelas, &siswa.NamaSiswa, &siswa.Alamat, &siswa.TanggalLahir); err != nil {
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

	query := "INSERT INTO siswa (id_user, id_kelas, nama_siswa, alamat, tanggal_lahir) VALUES ($1, $2, $3, $4, $5)"
	_, err = database.Exec(query, siswa.IDUser, siswa.IDKelas, siswa.NamaSiswa, siswa.Alamat, siswa.TanggalLahir)
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
		"UPDATE siswa SET id_user=$1, id_kelas=$2, nama_siswa=$3, alamat=$4, tanggal_lahir=$5 WHERE id_siswa=$6",
		siswa.IDUser, siswa.IDKelas, siswa.NamaSiswa, siswa.Alamat, siswa.TanggalLahir, id,
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
	err = database.QueryRow("SELECT id_siswa,id_user, id_kelas, nama_siswa, alamat, tanggal_lahir FROM siswa WHERE id_siswa=$1", id).
		Scan(&siswa.IDSiswa, &siswa.IDUser, &siswa.IDKelas, &siswa.NamaSiswa, &siswa.Alamat, &siswa.TanggalLahir)
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












