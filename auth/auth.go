package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtpassword = []byte("luisocegueda25#")

func CreateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * (24 * 30)).Unix()
	tokenString, err := token.SignedString(jwtpassword)

	if err != nil {
		log.Fatal(err.Error())
	}

	return tokenString

}

func IsAuthorized(next http.Handler) http.Handler {
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
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("I don see the token"))
		}
	})
}
