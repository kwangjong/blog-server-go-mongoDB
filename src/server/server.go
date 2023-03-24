package server

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"log"
	"net/http"
	"encoding/json"
	"github.com/kwangjong/kwangjong.github.io/db"
)

const (
	POSTPATH = "/blog/"
)

var PostDB	*db.DBCollection

func GetPost(w http.ResponseWriter, r *http.Request, url string) (error, int)  {
	post, err := PostDB.Get(url)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	post_json, err := json.Marshal(post)
	if err != nil {
		switch(err.Error()) {
		case "page not found":
			return err, http.StatusNotFound
		default:
			return err, http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(post_json));

	return nil, 0;
}

func Post(w http.ResponseWriter, r *http.Request) {
	post_url := r.URL.Path[len(POSTPATH):]
	if post_url == "" {
		log.Printf("Error: page not found\n")
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	log.Printf("%s: %s", r.Method, r.URL.Path)
	
	var err error
	var err_code int
	switch r.Method {
	case http.MethodGet:
		err, err_code = GetPost(w, r, post_url)
	default:
		return
	}
	if err != nil {
		log.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), err_code)
	}
}

// func Blog(w http.ResposeWriter, r *http.Request) {
	
// }

// func Tags() {

// }

func Run() {
	http.HandleFunc(POSTPATH, Post)

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
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}


}