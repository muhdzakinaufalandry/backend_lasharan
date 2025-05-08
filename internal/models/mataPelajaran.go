package models


type MataPelajaran struct {
    IDMapel       		  int       `json:"id_mapel"`
    IDKelas               int       `json:"id_kelas"`
    NamaMataPelajaran     string    `json:"nama_mata_pelajaran"`
}
