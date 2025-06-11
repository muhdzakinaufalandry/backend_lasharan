package models


type Siswa struct {
    IDSiswa       int       `json:"id_siswa"`
	IDKelas       int       `json:"id_kelas"`
	IDUser        int       `json:"id_user"`
    NamaSiswa     string    `json:"nama_siswa"`
	Alamat     	  string    `json:"alamat"`
    TanggalLahir  string 	`json:"tanggal_lahir"`  // Tipe data time.Time untuk tanggal
	NISN           string 	`json:"nisn"`
}
