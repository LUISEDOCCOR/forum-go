package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LUISEDOCCOR/api/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func getJwtPassword() []byte {

	ok := godotenv.Load()
	if ok != nil {
		log.Fatal("I don have .env file")
	}

	jwtkey := os.Getenv("JWTKEY")
	return []byte(jwtkey)
}

func CreateToken(name string, id uint) string {
	jwtpassword := getJwtPassword()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["name"] = name
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * (24 * 30)).Unix() // 30 days
	tokenString, err := token.SignedString(jwtpassword)

	if err != nil {
		log.Fatal(err.Error())
	}

	return tokenString

}

func IsAuthorized(next http.Handler) http.Handler {
	jwtpassword := getJwtPassword()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
			//jwt start
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error in auth")
				}
				return jwtpassword, nil
			})
			//jwt end
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Println("claims")
				return
			}

			userId, ok := claims["id"].(float64)

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			name, ok := claims["name"].(string)

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			var credentials = types.CredentialsUser{
				ID:   userId,
				Name: name,
			}

			ctx := context.WithValue(r.Context(), "credentialsUser", credentials)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("I don see the token"))
		}
	})
}
