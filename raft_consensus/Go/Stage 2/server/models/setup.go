package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDatabase() {

	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"records_store",
		"password",
		os.Getenv("DB_NETWORK"),
		"5432",
		"records_store")
	database, err := gorm.Open("postgres", url)

	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&Album{})
	DB = database
}
