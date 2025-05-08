package models

type User struct {
	IDUser      int    `json:"id_user"`
	IDRole     	string `json:"id_role"`
	Username 	string `json:"username"`
	Password 	string `json:"password"`
}
