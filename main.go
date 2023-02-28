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
		Title: "test mongo",
		Description: "testing mongo db insert",
		By: "kj",
		Timestamp: "1234",
		Tags: []string{"test", "mongodb", "go"},
		Body: "hello world",
	}

	
	err = db.Insert(post)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}