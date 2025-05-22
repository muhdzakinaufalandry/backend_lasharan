package models


type PenilaianResponse struct {
    PenilaianList []Penilaian `json:"penilaian"`
    TotalNilai    string      `json:"total"`
}