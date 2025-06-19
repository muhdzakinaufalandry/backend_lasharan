package main

import (
	"log"
	"net/http"
	"myapp/internal/api"  // Pastikan path impor sesuai dengan folder proyek kamu
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"myapp/config"
)

func main() {
	config.InitS3()
	// Membuat router
	r := mux.NewRouter()

	// Menangani route untuk /guru/{id}
	r.HandleFunc("/guru/{id}", api.GetGuruByIDHandler).Methods("GET")
	r.HandleFunc("/siswa/{id}", api.GetSiswaByIDHandler).Methods("GET")
	r.HandleFunc("/kelas/{id}", api.GetKelasByIDHandler).Methods("GET")
	r.HandleFunc("/matapelajaran/{id}", api.GetMataPelajaranByIDHandler).Methods("GET")

	// Menghubungkan route dengan handler
	r.HandleFunc("/guru", api.CreateGuruHandler).Methods("POST")
	r.HandleFunc("/guru", api.GetGuruHandler).Methods("GET")
	r.HandleFunc("/guru/{id}", api.UpdateGuruHandler).Methods("PUT")
	r.HandleFunc("/guru/{id}", api.DeleteGuruHandler).Methods("DELETE")


	r.HandleFunc("/siswa", api.GetSiswaHandler).Methods("GET")
	r.HandleFunc("/siswa", api.CreateSiswaHandler).Methods("POST")
	r.HandleFunc("/siswa/{id}", api.UpdateSiswaHandler).Methods("PUT")
	r.HandleFunc("/siswa/{id}", api.DeleteSiswaHandler).Methods("DELETE")

	r.HandleFunc("/kelas", api.GetKelasHandler).Methods("GET")
	r.HandleFunc("/kelas", api.CreateKelasHandler).Methods("POST")
	r.HandleFunc("/kelas/{id}", api.UpdateKelasHandler).Methods("PUT")
	r.HandleFunc("/kelas/{id}", api.DeleteKelasHandler).Methods("DELETE")

	r.HandleFunc("/matapelajaran", api.GetMataPelajaranHandler).Methods("GET")
	r.HandleFunc("/matapelajaran", api.CreateMataPelajaranHandler).Methods("POST")
	r.HandleFunc("/matapelajaran/{id}", api.UpdateMataPelajaranHandler).Methods("PUT")
	r.HandleFunc("/matapelajaran/{id}", api.DeleteMataPelajaranHandler).Methods("DELETE")

	r.HandleFunc("/matapelajaran/bykelas/{id}", api.GetMataPelajaranByKelasHandler).Methods("GET")

	r.HandleFunc("/kelas/guru/{id_guru}", api.GetKelasByGuru).Methods("GET")

	r.HandleFunc("/guru/user/{id_user}", api.GetGuruByUserIDHandler).Methods("GET")
	r.HandleFunc("/siswa/user/{id_user}", api.GetSiswaByUserIDHandler).Methods("GET")

	r.HandleFunc("/matapelajaran/siswa/{id_siswa}", api.GetMataPelajaranBySiswaIDHandler).Methods("GET")

	r.HandleFunc("/kelass/{id_kelas}", api.GetKelasWithSubjects).Methods("GET")

	r.HandleFunc("/siswaa/{id_kelas}", api.GetSiswaByKelas).Methods("GET")

	r.HandleFunc("/mapel/simple-detail/{id_mapel}", api.GetSimpleSubjectDetailHandler).Methods("GET")

	r.HandleFunc("/siswa/by-mapel/{id_mapel}", api.GetStudentsByMapelID).Methods("GET")

	r.HandleFunc("/nilai-detail", api.GetPenilaianBySiswaAndMapelHandler).Methods("GET")

	r.HandleFunc("/penilaian", api.CreatePenilaianHandler).Methods("POST")
	r.HandleFunc("/penilaian", api.GetPenilaianHandler).Methods("GET")
	r.HandleFunc("/penilaian/{id}", api.UpdatePenilaianHandler).Methods("PUT")
	r.HandleFunc("/penilaian/{id}", api.DeletePenilaianHandler).Methods("DELETE")

	r.HandleFunc("/user", api.CreateUserHandler).Methods("POST")
	r.HandleFunc("/user", api.GetUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", api.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", api.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/user/{id}", api.GetUserByIDHandler).Methods("GET")

	
	r.HandleFunc("/upload-foto-guru", api.UploadFotoGuruHandler).Methods("POST")
	r.HandleFunc("/upload-foto-siswa", api.UploadFotoSiswaHandler).Methods("POST")

	

	r.HandleFunc("/login", api.LoginHandler)



	// Menambahkan CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Ganti dengan URL frontend Anda
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Menambahkan middleware CORS
	handler := c.Handler(r)

	// Menjalankan server di port 8080
	log.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
