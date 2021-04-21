package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/batijo/video-platform/backend/models"
	"github.com/twinj/uuid"

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

// func JwtVerify2(next http.Handler)http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tokenAuth, err := extractTokenMetadata(r)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		userID, err := fetchAuth(tokenAuth)
// 		if err != nil{
// 			log.Println(err)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), key, tk)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString, err := parseAuthHeader(r)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func isTokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("token is not valid")
	}
	return nil
}

func CreateToken(user models.User) (*models.Tokens, error) {
	var err error
	ts := &models.Tokens{
		AtExpires:   time.Now().Add(time.Minute * time.Duration(Conf.JWTExp)).Unix(),
		AccessUuid:  uuid.NewV4().String(),
		RtExpires:   time.Now().Add(time.Hour * time.Duration(Conf.JWTRef)).Unix(),
		RefreshUuid: uuid.NewV4().String(),
	}

	atk := &models.Token{
		UserID:     user.ID,
		Email:      user.Email,
		Admin:      user.Admin,
		AccessUuid: ts.AccessUuid,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: ts.AtExpires,
		},
	}
	atoken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), atk)
	ts.AccessToken, err = atoken.SignedString([]byte(Conf.JWTSecret))
	if err != nil {
		return nil, err
	}

	rtk := jwt.MapClaims{}
	rtk["refresh_uuid"] = ts.RefreshUuid
	rtk["user_id"] = user.ID
	rtk["exp"] = ts.RtExpires
	rtoken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), rtk)
	ts.RefreshToken, err = rtoken.SignedString([]byte(Conf.JWTSecret))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func CreateAuth(userID uint, ts *models.Tokens) error {
	at := time.Unix(ts.AtExpires, 0)
	rt := time.Unix(ts.RtExpires, 0)
	now := time.Now()
	ctx := context.Background()

	err := RedisCl.SetEX(ctx, ts.AccessUuid, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if err != nil {
		return err
	}
	errRefresh := RedisCl.SetEX(ctx, ts.RefreshUuid, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func extractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(models.Token)
	if ok && token.Valid {
		accessUuid := claims.AccessUuid
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserID:     claims.UserID,
		}, nil
	}
	return nil, errors.New("token is not valid")
}

func fetchAuth(ad *models.AccessDetails) (uint, error) {
	uid, err := RedisCl.Get(context.Background(), ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
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
