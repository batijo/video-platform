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
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
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
func findOne(email, password string) map[string]interface{} {
	user := &models.User{}

	if err := utils.DB.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * time.Duration(utils.Conf.JWTExp)).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil || errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
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

	var resp = map[string]interface{}{"status": true, "message": "logged in"}
	resp["token"] = tokenString // Store the token in the response
	// resp["user"] = user
	return resp
}

// CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	// Check if user trying to gain admin access
	if user.Admin {
		json.NewEncoder(w).Encode("Nice try, but you can not make yourself an admin")
		return
	}

	if user.Email == "" || user.Password == "" {
		json.NewEncoder(w).Encode("Email and/or Password must be provided")
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Password Encryption failed",
		}
		json.NewEncoder(w).Encode(err)
	}

	user.Password = string(pass)

	createdUser := utils.DB.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
	}
	json.NewEncoder(w).Encode(createdUser)
}

// FetchUsers function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	utils.DB.Preload("auths").Find(&users)

	json.NewEncoder(w).Encode(users)
}

// UpdateUser ...
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	params := mux.Vars(r)
	var id = params["id"]
	utils.DB.First(&user, id)

	if user.Email == "" {
		json.NewEncoder(w).Encode("User not found")
		return
	}

	json.NewDecoder(r.Body).Decode(&user)

	// Check if user trying to gain admin access
	if user.Admin {
		json.NewEncoder(w).Encode("Nice try, but you can not make yourself an admin")
		return
	}

	// Create new hash if pasword is changed
	if user.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			err := ErrorResponse{
				Err: "Password Encryption failed",
			}
			json.NewEncoder(w).Encode(err)
		}

		user.Password = string(pass)
	}

	utils.DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser ...
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	utils.DB.First(&user, id)

	// For some reason if you try to delete user which does not exist it deletes all users
	if user.Email == "" {
		json.NewEncoder(w).Encode("User not found")
		return
	}

	utils.DB.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
	json.NewEncoder(w).Encode(&user)
}

// GetUser ...
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	utils.DB.Preload("Video").Preload("Video.AudioT").Preload("Video.SubtitleT").First(&user, id)
	json.NewEncoder(w).Encode(&user)
}

// GetUserByEmail ...
func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var email = params["email"]
	println(email)
	var user models.User
	utils.DB.Where("email = ?", email).First(&user)
	json.NewEncoder(w).Encode(&user)
}
