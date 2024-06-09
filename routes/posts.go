package routes

import (
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
