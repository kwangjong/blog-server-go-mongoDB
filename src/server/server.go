package server

import (
	"log"
	"net/http"
	"github.com/kwangjong/kwangjong.github.io/db"
)

const (
	POSTPATH = "/post/"
)

var mongoDB db.Mongo_Client

func Post(w http.ResponseWriter, r *http.Request) {
	post_id := r.URL.Path[len(POSTPATH):]
	if post_id == "" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	log.Printf("%s: %s%s", r.Method, r.URL.Path)
	
	switch r.Method {
	case "GET":
		post, err := db.Get(post_id)
		if err != nil {
			log.Printf("Error: %s\n", err)
			http.Error(w, err, http.StatusInternalServerError)
		}
		return
	case "POST":
		return
	case "PUT":
		return
	case "DELETE":
		return
	}
}

// func Blog(w http.ResposeWriter, r *http.Request) {
	
// }

// func Tags() {

// }

func Run() {
	http.HandleFunc(POSTPATH, Post)

	log.Printf("Starting server...\n")
	mongoDB, err := db.Connect_DB("posts")
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		log.Printf("Received kill signal, cleaning up")
		mongoDB.Close()
		os.Exit(0)
	}()
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}


}