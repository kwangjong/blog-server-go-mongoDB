package server

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET = []byte("SUPER-SECRET-KEY")
var API_KEY = "1234"

func generateJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		log.Printf(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func validateJwt(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return SECRET, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

func getJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["Api-Key"][0] == API_KEY {
		token, err := generateJwt()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(token))
	}
}
