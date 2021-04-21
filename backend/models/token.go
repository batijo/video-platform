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

type Tokens struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserID     uint
}
