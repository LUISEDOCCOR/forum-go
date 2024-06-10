package routes

import (
	"encoding/json"
	"net/http"

	"github.com/LUISEDOCCOR/api/db"
	"github.com/LUISEDOCCOR/api/models"
	"github.com/LUISEDOCCOR/api/utils"
	"github.com/gorilla/mux"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	db.DB.Find(&posts)
	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	params := mux.Vars(r)
	id := params["id"]

	db.DB.Find(&post, id)

	if post.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := utils.CreateResponse("error", "The post does not exist")
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.Title == "" || post.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid data")
		json.NewEncoder(w).Encode(response)
		return
	}

	result := db.DB.Create(&post)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", result.Error.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	response := utils.CreateResponse("success", "added successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	params := mux.Vars(r)
	postID := params["postID"]

	result := db.DB.Delete(&post, postID)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", result.Error.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Post not found")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := utils.CreateResponse("success", "remove successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	var currentPost models.Post

	params := mux.Vars(r)
	postID := params["postID"]

	result := db.DB.First(&currentPost, postID)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", result.Error.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.Title == "" || post.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid data")
		json.NewEncoder(w).Encode(response)
		return
	}

	currentPost.Title = post.Title
	currentPost.Content = post.Content

	db.DB.Save(&currentPost)

	response := utils.CreateResponse("success", "update successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
