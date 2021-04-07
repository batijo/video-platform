package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

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

		if tk.Admin {
			w.WriteHeader(http.StatusForbidden)
			resp := models.Response{Status: false, Message: "No administration privilage"}
			json.NewEncoder(w).Encode(resp)
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
		resp := models.Response{Status: false, Message: "Missing auth token"}
		json.NewEncoder(w).Encode(resp)
		return *tk, errors.New("Missing auth token")
	}

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Conf.JWTSecret), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		resp := models.Response{Status: false, Message: "Not authorised", Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return *tk, errors.New("Not authorised")
	}

	return *tk, nil
}
