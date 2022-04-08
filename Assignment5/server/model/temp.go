package model

type Temp struct {
	Npm  string `gorm:"primaryKey" json:"npm"`
	Nama string `json:"nama"`
}
