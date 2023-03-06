package test

import (
	"testing"
	"time"
	"github.com/kwangjong/kwangjong.github.io/db"
)

func Test_Connect_DB_Close(t *testing.T) {
	err := db.Connect_DB()
	if err != nil {
		t.Error(err)
	}

	err = db.Close()
	if err != nil {
		t.Error(err)
	}
}

func Test_Insert_Delete(t *testing.T) {
	err := db.Connect_DB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	post := &db.Post{
		Title: "mongo Test_Insert_Delete",
		Description: "testing mongo db insert",
		By: "kj",
		DateCreated: time.Now(),
		LastUpdated: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		Body: "hello world",
	}

	err = db.Insert(post)
	if err != nil {
		t.Error(err)
	}

	filter_by_title := struct{
		title	string
	}{"mongo Test_Insert_Delete"}

	err = db.Delete(filter_by_title)
	if err != nil {
		t.Error(err)
	}
}

func Test_Read(t *testing.T) {
	err := db.Connect_DB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	posts := []*db.Post{
		&db.Post{
			Title: "mongo test1",
		},
		&db.Post{
			Title: "mongo test2",
		}, 
		&db.Post{
			Title: "mongo test3",
		},
	}

	for _, p := range posts {
		err = db.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := db.Load(0, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range posts {
		if results[i].Title != p.Title {
			t.Errorf("Expected: %s Received: %s\n", p.Title, results[i].Title)
		}
	}

	results, err = db.Load(1, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range posts[1:] {
		if results[i].Title != p.Title {
			t.Errorf("Expected: %s Received: %s\n", p.Title, results[i].Title)
		}
	}

	for _, p := range posts {
		err = db.Delete(p)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Find(t *testing.T) {
	err := db.Connect_DB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	posts := []*db.Post{
		&db.Post{
			Title: "mongo test1",
			Tags: []string{"db", "test"},
		},
		&db.Post{
			Title: "mongo test2",
			Tags: []string{"foo", "test"},
		}, 
		&db.Post{
			Title: "mongo test3",
			Tags: []string{"db"},
		},
	}

	for _, p := range posts {
		err = db.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	filter_db := struct{
		Tags	string
	}{"db"}

	filter_test := struct{
		Tags	string
	}{"test"}

	results, err := db.Find(filter_db, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[0].Title)
	}
	
	if results[1].Title != "mongo test3" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test3", results[1].Title)
	}

	results, err = db.Find(filter_test, 1, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test2" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[0].Title)
	}
	
	for _, p := range posts {
		err = db.Delete(p)
		if err != nil {
			t.Error(err)
		}
	}
}