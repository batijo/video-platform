package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// ErrorResponse ...
type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid request", Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := findOne(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

// LogOut ...
func LogOut(w http.ResponseWriter, r *http.Request) {
	// TO DO ...
}

// FindOne ...
func findOne(email, password string) models.Response {
	user := &models.User{}

	if err := utils.DB.Where("Email = ?", email).First(user).Error; err != nil {
		resp := models.Response{Status: false, Message: "Email address not found", Error: err.Error()}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * time.Duration(utils.Conf.JWTExp)).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil || errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		resp := models.Response{Status: false, Message: "Invalid login credentials. Please try again", Error: errf.Error()}
		return resp
	}

	tk := &models.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Admin:  user.Admin,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(utils.Conf.JWTSecret))
	if err != nil {
		fmt.Println(err)
	}

	resp := models.Response{Status: true, Message: "logged in", Data: tokenString}
	// resp["token"] = tokenString // Store the token in the response
	// resp["user"] = user
	return resp
}

// CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	// Check if user trying to gain admin access
	if user.Admin {
		resp := models.Response{Status: false, Message: "You can not make yourself an admin"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Email == "" || user.Password == "" {
		resp := models.Response{Status: false, Message: "Email and/or Password must be provided"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		resp := models.Response{Status: false, Message: "Password Encryption failed", Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	user.Password = string(pass)

	createdUser := utils.DB.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
		resp := models.Response{Status: false, Message: "Error ocured while creating user", Error: errMessage.Error(), Data: createdUser}
		json.NewEncoder(w).Encode(resp)
	}
	resp := models.Response{Status: true, Message: "User created", Data: createdUser}
	json.NewEncoder(w).Encode(resp)
}

// FetchUsers function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	utils.DB.Preload("auths").Find(&users)

	resp := models.Response{Status: true, Message: "Success", Data: users}
	json.NewEncoder(w).Encode(resp)
}

// UpdateUser ...
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	params := mux.Vars(r)
	var id = params["id"]
	res := utils.DB.First(&user, id)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewDecoder(r.Body).Decode(&user)

	// Check if user trying to gain admin access
	if user.Admin {
		resp := models.Response{Status: false, Message: "You can not make yourself an admin"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Create new hash if pasword is changed
	if user.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			resp := models.Response{Status: false, Message: "Password Encryption failed", Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		user.Password = string(pass)
	}

	utils.DB.Save(&user)
	resp := models.Response{Status: true, Message: "User updated", Data: user}
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser ...
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	res := utils.DB.First(&user, id)

	// For some reason if you try to delete user which does not exist it deletes all users
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	utils.DB.Delete(&user)

	resp := models.Response{Status: true, Message: "User deleted", Data: user}
	json.NewEncoder(w).Encode(resp)
}

// GetUser ...
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	res := utils.DB.Preload("Video").Preload("Video.AudioT").Preload("Video.SubtitleT").First(&user, id)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "User not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: user}
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
