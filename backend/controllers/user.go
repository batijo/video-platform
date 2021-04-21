package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid request", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Email == "" || user.Password == "" {
		resp := models.Response{Status: false, Message: "Email and Password must be provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, status := findOne(user.Email, user.Password)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// LogOut ...
func LogOut(w http.ResponseWriter, r *http.Request) {
	// TO DO ...
}

// FindOne ...
func findOne(email, password string) (models.Response, int) {
	user := &models.User{}

	if err := utils.DB.Where("email = ?", email).First(user).Error; err != nil {
		resp := models.Response{Status: false, Message: "Email address not found", Error: err.Error()}
		return resp, http.StatusUnauthorized
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		resp := models.Response{Status: false, Message: "Invalid login credentials. Please try again", Error: err.Error()}
		return resp, http.StatusUnauthorized
	}

	ts, err := utils.CreateToken(*user)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid login credentials. Please try again", Error: err.Error()}
		return resp, http.StatusUnauthorized
	}
	err = utils.CreateAuth(user.ID, ts)
	if err != nil {
		resp := models.Response{Status: false, Message: "Failed to cache token", Error: err.Error()}
		return resp, http.StatusInternalServerError
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	resp := models.Response{Status: true, Message: "logged in", Data: tokens}

	return resp, http.StatusOK
}

// CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid request", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		resp := models.Response{Status: false, Message: "Username, Email and Password must be provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := utils.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&user).Error; err == nil {
		resp := models.Response{Status: false, Message: "User with the same email or username already exist"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Check if user trying to gain admin access
	if user.Admin {
		resp := models.Response{Status: false, Message: "You can not make yourself an admin"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		resp := models.Response{Status: false, Message: "Password Encryption failed", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user.Password = string(pass)

	createdUser := utils.DB.Create(&user)
	if createdUser.Error != nil {
		log.Println(createdUser.Error)
		resp := models.Response{Status: false, Message: "Error ocured while creating user", Error: createdUser.Error.Error()}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "User created", Data: createdUser.Value}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// FetchUsers function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var (
		users []models.User
		res   *gorm.DB
	)

	userId, admin, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res = utils.DB.Find(&users)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not fetch users", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if admin {
		serializeUsers(&users, admin, userId)
	} else {
		serializeUsers(&users, false, userId)
	}

	resp := models.Response{Status: true, Message: "Success", Data: users}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func serializeUsers(users *[]models.User, full bool, userID uint) {
	for i, u := range *users {
		if u.ID == userID {
			(*users)[i].Serialize(true)
		} else {
			(*users)[i].Serialize(full)
		}
	}
}

// UpdateUser ...
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		user models.User
		id   = mux.Vars(r)["id"]
	)

	userId, admin, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.First(&user, id)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch true {
	case userId == user.ID:
		break
	case admin:
		break
	default:
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid request", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Check if user trying to gain admin access
	if user.Admin {
		resp := models.Response{Status: false, Message: "You can not make yourself an admin"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Create new hash if pasword is changed
	if user.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			resp := models.Response{Status: false, Message: "Password Encryption failed", Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		user.Password = string(pass)
	}

	res = utils.DB.Save(&user)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not save user", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user.Serialize(true)
	resp := models.Response{Status: true, Message: "User updated", Data: user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser ...
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		id   = mux.Vars(r)["id"]
		user models.User
	)

	userId, admin, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.First(&user, id)
	// For some reason if you try to delete user which does not exist it deletes all users
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if userId != user.ID || !admin {
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res = utils.DB.Delete(&user)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not delete user", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user.Serialize(true)
	resp := models.Response{Status: true, Message: "User deleted", Data: user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetUser ...
func GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		id   = mux.Vars(r)["id"]
		user models.User
	)

	userId, admin, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.Preload("Video").Preload("Video.AudioT").Preload("Video.SubtitleT").First(&user, id)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if userId == user.ID || admin {
		user.Serialize(true)
	} else {
		user.Serialize(false)
	}

	resp := models.Response{Status: true, Message: "Success", Data: user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetUserByEmail ...
func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var email = params["email"]
	var user models.User
	res := utils.DB.Where("email = ?", email).First(&user)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: user}
	json.NewEncoder(w).Encode(resp)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	var username = mux.Vars(r)["username"]
	var user models.User
	res := utils.DB.Where("username = ?", username).First(&user)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: user}
	json.NewEncoder(w).Encode(resp)
}
