package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

//Exception struct
type Exception models.Exception

func GetUserID(r *http.Request) (uint, error) {
	var header = r.Header.Get("x-access-token") //Grab the token from the header
	header = strings.TrimSpace(header)

	tk := &models.Token{}

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Conf.JWTSecret), nil
	})
	if err != nil {
		return tk.UserID, err
	}
	return tk.UserID, nil
}

func AdminVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tk, err := jwtParser(w, r)
		if err != nil {
			return
		}

		if tk.Role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "No administration privilage"})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tk, err := jwtParser(w, r)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jwtParser(w http.ResponseWriter, r *http.Request) (models.Token, error) {

	tk := &models.Token{}

	var header = r.Header.Get("x-access-token") //Grab the token from the header
	header = strings.TrimSpace(header)

	if header == "" {
		//Token is missing, returns with error code 403 Unauthorized
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
		return *tk, errors.New("Missing auth token")
	}

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Conf.JWTSecret), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return *tk, errors.New("Not authorised")
	}

	return *tk, nil
}
