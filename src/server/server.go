package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

const (
	BLOGPATH        = "/blog"
	BLOGSLUGPATH    = "/blog/{slug}"
	BLOGLISTPATH    = "/blog/list"
	BLOGLISTALLPATH = "/blog/list/all"
	TAGSLISTALLPATH = "/tags/list/all"
	AUTHPATH        = "/auth"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Api-Key, Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Run() {
	log.Printf("Starting server...\n")
	blog := BlogHandler{}
	err := blog.setDB("post")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc(BLOGLISTPATH, blog.BlogList).Methods("GET")
	r.HandleFunc(BLOGLISTALLPATH, blog.BlogListAll).Methods("GET")
	r.HandleFunc(TAGSLISTALLPATH, blog.TagsListAll).Methods("GET")

	r.HandleFunc(BLOGSLUGPATH, blog.BlogGET).Methods("GET")
	r.HandleFunc(BLOGSLUGPATH, validateJwtHandler(blog.BlogDELETE)).Methods("DELETE")
	r.HandleFunc(BLOGPATH, validateJwtHandler(blog.BlogPOST)).Methods("POST")
	r.HandleFunc(BLOGPATH, validateJwtHandler(blog.BlogPUT)).Methods("PUT")

	r.HandleFunc(AUTHPATH, AuthGET).Methods("GET")
	r.HandleFunc(AUTHPATH, AuthPOST).Methods("POST")

	r.Use(loggingMiddleware)

	handler := corsMiddleware(r)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		log.Printf("Received kill signal, cleaning up")
		blog.closeDB()
		os.Exit(0)
	}()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
