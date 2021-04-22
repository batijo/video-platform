package models

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Admin      bool   `json:"admin"`
	AccessUuid string `json:"access_uuid"`
	*jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUuid string
	UserID     uint
}
