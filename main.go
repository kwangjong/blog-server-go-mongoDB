package main

import (
	//"fmt"
	//"time"
	"github.com/kwangjong/kwangjong.github.io/db"
)

func main() {
	err := db.Connect_DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// posts := []*db.Post{
	// 	&db.Post{
	// 		Title: "mongo test1",
	// 		Tags: []string{"db", "test"},
	// 	},
	// 	&db.Post{
	// 		Title: "mongo test2",
	// 		Tags: []string{"foo", "test"},
	// 	}, 
	// 	&db.Post{
	// 		Title: "mongo test3",
	// 		Tags: []string{"db"},
	// 	},
	// }

	// for _, p := range posts {
	// 	_, err = db.Insert(p)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	db.Distinct("tags")
}