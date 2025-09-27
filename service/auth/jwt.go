package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zinx110/golang-backend-rest/config"
	"github.com/zinx110/golang-backend-rest/types"
	"github.com/zinx110/golang-backend-rest/utils"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the toekn from the user requrest
		tokenString := getTokenFromRequest(r)
		// validate the jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		// if is we need to fetch the userid from db (id from the token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert user id from token to int: %v", err)
			permissionDenied(w)
			return
		}
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id from token: %v", err)
			permissionDenied(w)
			return
		}
		// set context "userILD"
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	return tokenAuth
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) (int, error) {
	userId, ok := ctx.Value(UserKey).(int)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}
	return userId, nil
}
