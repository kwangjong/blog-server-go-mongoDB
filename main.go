package main

import (
	"fmt"
	//"time"
	"github.com/kwangjong/kwangjong.github.io/db"
)

func main() {
	err := db.Connect_DB()
	if err != nil {
		panic(err)
	}

	// post := &db.Post{
	// 	Title: "test3 mongo",
	// 	Description: "testing mongo db insert",
	// 	By: "kj",
	// 	DateCreated: time.Now(),
	// 	LastUpdated: time.Now(),
	// 	Tags: []string{"test", "mongodb", "go"},
	// 	Body: "hello world",
	// }
	filter := struct{
		By	string
	}{"kj"}

	posts, err := db.Find(filter, 2, 2)
	if err != nil {
		panic(err)
	}

	fmt.Println(posts)

	defer db.Close()
}