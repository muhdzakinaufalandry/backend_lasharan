package models


type Kelas struct {
    IDKelas      int       `json:"id_kelas"`
	IDGuru       int       `json:"id_guru"`
    NamaKelas    string    `json:"nama_kelas"`
	TahunAjaran  string    `json:"tahun_ajaran"`
	JumlahSiswa  int    `json:"jumlah_siswa"`
}
