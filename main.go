package main

import (
	"github.com/kwangjong/kwangjong.github.io/db"
)

func main() {
	err := db.Connect_DB()
	if err != nil {
		panic(err)
	}

	post := &db.Post{
		Title: "test2 mongo",
		Description: "testing mong2so db insert",
		By: "kj",
		Timestamp: "1234",
		Tags: []string{"test", "mongodb", "go"},
		Body: "hello world",
	}

	
	err = db.Update("63fdc604f30b565e9e716976", post)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}