package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kwangjong/kwangjong.github.io/db"
)

const (
	BLOGPATH     = "/blog/"
	BLOGLISTPATH = "/blog/list"
	BLOGLISTALLPATH = "/blog/list/all"
	TAGSLISTPATH = "/tags/list"
	AUTHPATH     = "/auth"
)

var PostDB *db.DBCollection

func Get_Blog(w http.ResponseWriter, r *http.Request) {
	post_url := r.URL.Path[len(BLOGPATH):]
	if post_url == "" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	post, err := PostDB.Get(post_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post_json, err := json.Marshal(post)
	if err != nil {
		switch err.Error() {
		case "page not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(post_json)
}

func Post_Blog(w http.ResponseWriter, r *http.Request) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = PostDB.Insert(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Put_Blog(w http.ResponseWriter, r *http.Request) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s", post)

	err = PostDB.Update(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Delete_Blog(w http.ResponseWriter, r *http.Request) {
	post_url := r.URL.Path[len(BLOGPATH):]
	if post_url == "" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	err := PostDB.Delete(post_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Blog(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token")

	switch r.Method {
	case http.MethodGet:
		Get_Blog(w, r)
	case http.MethodPost:
		validateJwt(Post_Blog).ServeHTTP(w, r)
	case http.MethodPut:
		validateJwt(Put_Blog).ServeHTTP(w, r)
	case http.MethodDelete:
		validateJwt(Delete_Blog).ServeHTTP(w, r)
	default:
		return
	}
}

func BlogList(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	var err error
	var err_code int
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	skip, err := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
	numPost, err := strconv.ParseInt(r.URL.Query().Get("numPost"), 10, 64)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
		return
	}

	list, hasNext, err := PostDB.Read(skip, numPost)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
		return
	}

	list_json, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	out := fmt.Sprintf("{\"entries\": %s, \"hasNext\": %v}", list_json, hasNext)
	io.WriteString(w, out)
}

func BlogListAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	var err error
	var err_code int
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	tags, err := PostDB.Distinct("url")
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	tags_json, err := json.Marshal(tags)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tags_json)
}

func TagsList(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	var err error
	var err_code int
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	tags, err := PostDB.Distinct("tags")
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	tags_json, err := json.Marshal(tags)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tags_json)
}

func Run() {
	http.HandleFunc(BLOGPATH, Blog)
	http.HandleFunc(BLOGLISTPATH, BlogList)
	http.HandleFunc(BLOGLISTALLPATH, BlogListAll)
	http.HandleFunc(TAGSLISTPATH, TagsList)
	http.HandleFunc(AUTHPATH, getJwt)

	log.Printf("Starting server...\n")
	client, err := db.Connect_DB()
	if err != nil {
		panic(err)
	}

	PostDB, err = client.Collection("post")
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		log.Printf("Received kill signal, cleaning up")
		client.Close()
		os.Exit(0)
	}()

	server := &http.Server{
		Addr: ":8080",
		Handler: nil
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
