package main

import (
	"log"
	"net/http"

	"github.com/LUISEDOCCOR/api/auth"
	"github.com/LUISEDOCCOR/api/db"
	"github.com/LUISEDOCCOR/api/models"
	"github.com/LUISEDOCCOR/api/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db.Conn()
	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.Post{})

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
	r.HandleFunc("/posts/preview", routes.GetAllPostsPreview).Methods("GET")                      //preview
	privateRouter.HandleFunc("/post/{id:[0-9]+}", routes.GetPost).Methods("GET")                  // 1 post
	privateRouter.HandleFunc("/posts", routes.GetAllPosts).Methods("GET")                         // posts when user is loged
	privateRouter.HandleFunc("/myposts", routes.GetMyPosts).Methods("GET")                        // get only my posts
	privateRouter.HandleFunc("/post/add", routes.CreatePost).Methods("POST")                      //Add
	privateRouter.HandleFunc("/post/edit/{postID:[0-9]+}", routes.UpdatePost).Methods("PUT")      //Edit
	privateRouter.HandleFunc("/post/delete/{postID:[0-9]+}", routes.DeletePost).Methods("DELETE") // Deltet

	err := http.ListenAndServe(":5000", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r))
	log.Fatal(err)
}
