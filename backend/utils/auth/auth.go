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

func GetUserID(r *http.Request) (uint, bool, error) {
	tk := &models.Token{}

	//Grab the token from the header
	token, err := parseAuthHeader(r)
	if err != nil {
		return tk.UserID, false, err
	}

	_, err = jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Conf.JWTSecret), nil
	})
	if err != nil {
		return tk.UserID, false, err
	}
	return tk.UserID, tk.Admin, nil
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

		ctx := context.WithValue(r.Context(), "admin", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var key string
		tk, err := jwtParser(w, r)
		if err != nil {
			return
		}

		if tk.Admin {
			key = "admin"
		} else {
			key = "user"
		}

		ctx := context.WithValue(r.Context(), key, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jwtParser(w http.ResponseWriter, r *http.Request) (models.Token, error) {
	tk := &models.Token{}

	//Grab the token from the header
	token, err := parseAuthHeader(r)
	if err != nil {
		//Token is missing, returns with error code 403 Unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		resp := models.Response{Status: false, Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return *tk, err
	}

	_, err = jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
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

func parseAuthHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("Missing 'Authorization' header")
	}

	header = strings.Replace(header, "\n", " ", -1)
	splitHeader := strings.Split(header, " ")
	switch len(splitHeader) {
	case 0:
		return "", errors.New("Missing auth token")
	case 1:
		return "", errors.New("Wrong header format")
	case 2:
		if strings.TrimSpace(splitHeader[0]) != "Bearer" {
			return "", errors.New("Wrong header format")
		} else {
			break
		}
	default:
		return "", errors.New("Wrong header format")
	}

	return splitHeader[1], nil
}
