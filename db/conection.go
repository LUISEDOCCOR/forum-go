package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DNS = "host=localhost user=luisedoccor password=654321 dbname=gorm port=4321"

func Conn() {
	var error error
	DB, error = gorm.Open(postgres.Open(DNS), &gorm.Config{})

	if error != nil {
		log.Fatal("Error 🙇‍♀️")
	} else {
		log.Println("The server is 🚀")
	}

}
