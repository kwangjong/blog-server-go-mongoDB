package server

import (
	"encoding/json"
	"errors"
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
	TAGSLISTPATH = "/tags/list"
	CERTFILEPATH = "/home/kwangjong/107106.xyz-ssl-bundle/domain.cert.pem"
	KEYFILEPATH  = "/home/kwangjong/107106.xyz-ssl-bundle/private.key.pem"
)

var PostDB *db.DBCollection

func Get_Blog(w http.ResponseWriter, r *http.Request) (error, int) {
	post_url := r.URL.Path[len(BLOGPATH):]
	if post_url == "" {
		return errors.New("page not found"), http.StatusNotFound
	}

	post, err := PostDB.Get(post_url)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	post_json, err := json.Marshal(post)
	if err != nil {
		switch err.Error() {
		case "page not found":
			return err, http.StatusNotFound
		default:
			return err, http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(post_json))

	return nil, 0
}

func Post_Blog(w http.ResponseWriter, r *http.Request) (error, int) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return err, http.StatusBadRequest
	}

	err = PostDB.Insert(&post)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, 0
}

func Put_Blog(w http.ResponseWriter, r *http.Request) (error, int) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return err, http.StatusBadRequest
	}

	log.Printf("%s", post)

	err = PostDB.Update(&post)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, 0
}

func Blog(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL.Path)

	var err error
	var err_code int
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case http.MethodGet:
		err, err_code = Get_Blog(w, r)
	case http.MethodPost:
		err, err_code = Post_Blog(w, r)
	default:
		return
	}
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
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
	}

	list, hasNext, err := PostDB.Read(skip, numPost)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	list_json, err := json.Marshal(list)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}

	w.Header().Set("Content-Type", "application/json")
	out := fmt.Sprintf("{\"entries\": %s, \"hasNext\": %v}", list_json, hasNext)
	io.WriteString(w, out)
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
	io.WriteString(w, string(tags_json))
}

func Run() {
	http.HandleFunc(BLOGPATH, Blog)
	http.HandleFunc(BLOGLISTPATH, BlogList)
	http.HandleFunc(TAGSLISTPATH, TagsList)

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

	if err := http.ListenAndServeTLS(":443", CERTFILEPATH, KEYFILEPATH, nil); err != nil {
		log.Fatal(err)
	}

}
