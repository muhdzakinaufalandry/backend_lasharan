package models

type User struct {
	IDUser      int    `json:"id_user"`
	IDRole     	int `json:"id_role"`
	Username 	string `json:"username"`
	Password 	string `json:"password"`
	TanggalRegistrasi string `json:"tanggal_registrasi"`
}
