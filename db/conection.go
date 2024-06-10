package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conn() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("I don have .env file")
		return
	}

	dbHost := os.Getenv("HOST")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("PASSWORD")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("PORT")

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	fmt.Println(dns)

	var error error
	DB, error = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if error != nil {
		log.Fatal("Error 🙇‍♀️")
	} else {
		log.Println("The server is 🚀")
	}

}
