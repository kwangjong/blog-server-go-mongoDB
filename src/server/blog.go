package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kwangjong/kwangjong.github.io/db"
)

type BlogHandler struct {
	db     *db.DBCollection
	client *db.DBClient
}

func (b *BlogHandler) setDB(collection string) error {
	var err error

	b.client, err = db.Connect_DB()
	if err != nil {
		return err
	}

	b.db, err = b.client.Collection(collection)
	if err != nil {
		return err
	}

	return nil
}

func (b *BlogHandler) closeDB() {
	b.client.Close()
}

func (b *BlogHandler) BlogGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post_url := vars["slug"]
	if post_url == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	post, err := b.db.Get(post_url, validateJwt(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post_json, err := json.Marshal(post)
	if err != nil {
		switch err.Error() {
		case "404 page not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(post_json)
}

func (b *BlogHandler) BlogPOST(w http.ResponseWriter, r *http.Request) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = b.db.Insert(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *BlogHandler) BlogPUT(w http.ResponseWriter, r *http.Request) {
	var post db.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s", post)

	err = b.db.Update(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *BlogHandler) BlogDELETE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post_url := vars["slug"]
	if post_url == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	err := b.db.Delete(post_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *BlogHandler) BlogList(w http.ResponseWriter, r *http.Request) {
	var err error
	var err_code int

	tag := r.URL.Query().Get("tag")
	skip, err := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
	numPost, err := strconv.ParseInt(r.URL.Query().Get("numPost"), 10, 64)
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
		return
	}

	list, hasNext, err := b.db.Read(skip, numPost, validateJwt(r), tag)
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

func (b *BlogHandler) BlogListAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var err_code int

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	tags, err := b.db.Distinct("url", validateJwt(r))
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

func (b *BlogHandler) TagsListAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var err_code int

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	tags, err := b.db.Distinct("tags", validateJwt(r))
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
