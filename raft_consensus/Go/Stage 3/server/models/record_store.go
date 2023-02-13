package models

type Album struct {
	ID     int16   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
