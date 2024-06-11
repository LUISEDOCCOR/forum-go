package routes

import (
	"encoding/json"
	"net/http"

	"github.com/LUISEDOCCOR/api/db"
	"github.com/LUISEDOCCOR/api/models"
	"github.com/LUISEDOCCOR/api/types"
	"github.com/LUISEDOCCOR/api/utils"
	"github.com/gorilla/mux"
)

func GetAllPostsPreview(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	db.DB.Limit(10).Find(&posts)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var Posts []models.Post
	db.DB.Find(&Posts)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Posts)
}
func GetMyPosts(w http.ResponseWriter, r *http.Request) {
	credentialsUser := r.Context().Value("credentialsUser").(types.CredentialsUser)
	var Posts []models.Post
	db.DB.Where("user_id = ?", credentialsUser.ID).Find(&Posts)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Posts)
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
	var postData types.PostData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postData)

	credentialsUser := r.Context().Value("credentialsUser").(types.CredentialsUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if postData.Title == "" || postData.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid data")
		json.NewEncoder(w).Encode(response)
		return
	}

	var post models.Post
	post.Author = credentialsUser.Name
	post.UserId = credentialsUser.ID
	post.Title = postData.Title
	post.Content = postData.Content

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
	credentialsUser := r.Context().Value("credentialsUser").(types.CredentialsUser)

	var post models.Post
	params := mux.Vars(r)
	postID := params["postID"]

	result := db.DB.Where("user_id = ?", credentialsUser.ID).Delete(&post, postID)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", result.Error.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "The post is not yours")
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
	credentialsUser := r.Context().Value("credentialsUser").(types.CredentialsUser)

	params := mux.Vars(r)
	postID := params["postID"]

	result := db.DB.Where("user_id = ?", credentialsUser.ID).First(&currentPost, postID)

	if currentPost.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "The post is not yours")
		json.NewEncoder(w).Encode(response)
		return
	}

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
