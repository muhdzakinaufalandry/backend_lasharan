package models

type Guru struct {
	IDGuru        int    `json:"id_guru"`
	IDUser        int    `json:"id_user"`
	IDMapel       int    `json:"id_mapel"`
	NamaGuru      string `json:"nama_guru"`
	MataPelajaran string `json:"mata_pelajaran"`
	NIP           string `json:"nip"`
	Alamat        string `json:"alamat"`
	Email         string `json:"email"`
	NoTelp        string `json:"no_telp"`
}
