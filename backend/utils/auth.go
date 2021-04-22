package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/batijo/video-platform/backend/models"
	"github.com/twinj/uuid"

	jwt "github.com/dgrijalva/jwt-go"
)

func GetUserID(r *http.Request) (uint, bool, error) {
	//Grab the token from the header
	token, err := parseAuthHeader(r)
	if err != nil {
		return 0, false, err
	}

	return GetUserIDFromToken(token)
}

func GetUserIDFromToken(token string) (uint, bool, error) {
	tk := &models.Token{}

	_, err := jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(Conf.JWTSecret), nil
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
		tokenAuth, err := ExtractTokenMetadata(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := fetchAuth(tokenAuth)
		if err != nil || userID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
		return []byte(Conf.JWTSecret), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		resp := models.Response{Status: false, Message: "Not authorised", Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return *tk, errors.New("Not authorised")
	}

	return *tk, nil
}

func verifyToken(r *http.Request) (*models.Token, error) {
	tokenString, err := parseAuthHeader(r)
	if err != nil {
		return nil, err
	}

	tk := &models.Token{}
	_, err = jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// }
		return []byte(Conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func CreateToken(user models.User) (*models.TokenDetails, error) {
	var err error

	atk := &models.Token{
		UserID:     user.ID,
		Email:      user.Email,
		Admin:      user.Admin,
		AccessUuid: uuid.NewV4().String(),
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(Conf.JWTExp)).Unix(),
		},
	}
	atoken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), atk)
	token, err := atoken.SignedString([]byte(Conf.JWTSecret))
	if err != nil {
		return nil, err
	}

	td := &models.TokenDetails{
		AccessToken: token,
		AccessUuid:  atk.AccessUuid,
		AtExpires:   atk.ExpiresAt,
	}

	return td, nil
}

func ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	accessUuid := token.AccessUuid
	return &models.AccessDetails{
		AccessUuid: accessUuid,
		UserID:     token.UserID,
	}, nil
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
