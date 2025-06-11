package models

type Penilaian struct {
	IDPenilaian int     `json:"id_penilaian"`
	IDNilai     int     `json:"id_nilai"`
	IDMapel     int     `json:"id_mapel"`
	IDSiswa     int     `json:"id_siswa"`
	NamaNilai   string  `json:"nama_nilai"`
	Nilai       int `json:"nilai"`
	Bobot       string `json:"bobot"`
    Range string  `json:"range"`
}