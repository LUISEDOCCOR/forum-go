package routes

import (
	"encoding/json"
	"net/http"

	"github.com/LUISEDOCCOR/api/auth"
	"github.com/LUISEDOCCOR/api/db"
	"github.com/LUISEDOCCOR/api/models"
	"github.com/LUISEDOCCOR/api/utils"
	"github.com/alexedwards/argon2id"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Invalid data")
		json.NewEncoder(w).Encode(response)
		return
	}

	hashPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.CreateResponse("error", "Server Error")
		json.NewEncoder(w).Encode(response)
		return
	}

	user.Password = hashPassword

	result := db.DB.Create(&user)
	if result.Error != nil {
		response := utils.CreateResponse("error", result.Error.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	token := auth.CreateToken()
	response := utils.CreateResponse("success", token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {

	type creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var credentials creds

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		response := utils.CreateResponse("error", "Invalid json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		response := utils.CreateResponse("error", "Invalid data")
		json.NewEncoder(w).Encode(response)
		return
	}

	var user models.User
	db.DB.Where("email = ?", credentials.Email).First(&user)

	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "User does not exist")
		json.NewEncoder(w).Encode(response)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(credentials.Password, user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.CreateResponse("error", "Server Error")
		json.NewEncoder(w).Encode(response)
		return
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.CreateResponse("error", "Wrong credentials")
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	token := auth.CreateToken()
	response := utils.CreateResponse("success", token)
	json.NewEncoder(w).Encode(response)
}
