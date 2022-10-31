package models

type Album struct {
	ID     int16   `json:"id" gorm:"primary_key"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
