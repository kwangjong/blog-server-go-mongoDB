package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var SECRET []byte
var API_KEY string

func LoadSecret() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	SECRET = []byte(os.Getenv("JWT_SECRET"))
	API_KEY = os.Getenv("API_KEY")
}
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

func validateJwtHandler(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("not authorized")
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

func validateJwt(r *http.Request) bool {
	if r.Header["Token"] != nil {
		token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("not authorized")
			}
			return SECRET, nil
		})

		if err == nil && token.Valid {
			return true
		}
	}

	return false
}

func getJwt(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Api-Key, Token")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodPost {
		validateJwtHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("authorized"))
		}).ServeHTTP(w, r)
		return
	}
	
	log.Printf("%s: %s", r.Header["Api-Key"], API_KEY)
	_, ok := r.Header["Api-Key"]
	if ok && r.Header["Api-Key"][0] == API_KEY {
		token, err := generateJwt()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(token))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not authorized"))
	}
}
