package main

import (
	"log"
	"net/http"

	"github.com/LUISEDOCCOR/api/auth"
	"github.com/LUISEDOCCOR/api/db"
	"github.com/LUISEDOCCOR/api/models"
	"github.com/LUISEDOCCOR/api/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.Conn()
	db.DB.AutoMigrate(models.User{})

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	//Users
	r.HandleFunc("/user/create", routes.CreateUser).Methods("POST")
	r.HandleFunc("/user/login", routes.Login).Methods("POST")

	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(auth.IsAuthorized)

	//Posts
	privateRouter.HandleFunc("/posts", routes.GetAllPosts).Methods("GET")

	err := http.ListenAndServe(":5000", r)
	log.Fatal(err)
}
