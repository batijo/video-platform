package models

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Admin  bool   `json:"admin"`
	*jwt.StandardClaims
}
