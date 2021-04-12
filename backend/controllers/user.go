package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/batijo/video-platform/backend/utils/auth"
	"github.com/jinzhu/gorm"

	jwt "github.com/dgrijalva/jwt-go"
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
	expiresAt := time.Now().Add(time.Minute * time.Duration(utils.Conf.JWTExp)).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil || errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		resp := models.Response{Status: false, Message: "Invalid login credentials. Please try again", Error: errf.Error()}
		return resp, http.StatusUnauthorized
	}

	tk := &models.Token{
		UserID: user.ID,
		Email:  user.Email,
		Admin:  user.Admin,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(utils.Conf.JWTSecret))
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid login credentials. Please try again", Error: err.Error()}
		return resp, http.StatusInternalServerError
	}

	resp := models.Response{Status: true, Message: "logged in", Data: tokenString}
	// resp["token"] = tokenString // Store the token in the response
	// resp["user"] = user
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

	// Check if user trying to gain admin access
	if user.Admin {
		resp := models.Response{Status: false, Message: "You can not make yourself an admin"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		resp := models.Response{Status: false, Message: "Username, Email and Password must be provided"}
		w.WriteHeader(http.StatusBadRequest)
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

	createdUser := utils.DB.Create(user)
	if createdUser.Error != nil {
		log.Println(createdUser.Error)
		resp := models.Response{Status: false, Message: "Error ocured while creating user", Error: createdUser.Error.Error(), Data: createdUser}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "User created", Data: createdUser}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// FetchUsers function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var (
		users []models.User
		res   *gorm.DB
	)

	userId, admin, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if admin {
		res = utils.DB.Find(&users)
	} else {
		res = utils.DB.Where("id = ? OR public = ?", userId, true).Find(&users)
	}

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not fetch users", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: users}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateUser ...
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		user models.User
		id   = mux.Vars(r)["id"]
	)

	userId, admin, err := auth.GetUserID(r)
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

	if userId != user.ID || !admin {
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Forgot what I wanted to do here
	// if user.Email != "" {
	// 	if rows := utils.DB.Where("email = ?", user.Email).First(&user).RowsAffected; rows != 0 && userId != user.ID {

	// 	}
	// }
	// if user.Username != "" {
	// 	if rows := utils.DB.Where("username = ?", user.Email).First(&user).RowsAffected; rows != 0 {

	// 	}
	// }

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

	userId, admin, err := auth.GetUserID(r)
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

	userId, admin, err := auth.GetUserID(r)
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

	if !user.Public || !admin {
		if userId != user.ID {
			resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(resp)
			return
		}
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
