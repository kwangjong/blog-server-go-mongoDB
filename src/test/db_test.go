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
		Author: "kj",
		DateCreated: time.Now(),
		LastUpdated: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		MarkDown: "hello world",
		Html: "<p>hello world</p>",
	}

	id, err := db.Insert(post)
	if err != nil {
		t.Error(err)
	}

	err = db.Delete(id)
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
	ids := [3]string{}
	for i, p := range posts {
		ids[i], err = db.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := db.Read(0, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range results {
		if p.Title != posts[2-i].Title {
		 	t.Errorf("Expected: %s Received: %s\n", posts[2-i].Title, p.Title)
		}
	}

	results, err = db.Read(1, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range results {
		if p.Title != posts[1-i].Title {
		 	t.Errorf("Expected: %s Received: %s\n", posts[1-i].Title, p.Title)
		}
	}

	for _, id := range ids {
		err = db.Delete(id)
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

	ids := [3]string{}
	for i, p := range posts {
		ids[i], err = db.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := db.Find(db.FilterTag{"db"}, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test3" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test3", results[0].Title)
	}
	
	if results[1].Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[1].Title)
	}

	results, err = db.Find(db.FilterTag{"test"}, 1, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[0].Title)
	}
	
	for _, id := range ids {
		err = db.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Update(t *testing.T) {
	err := db.Connect_DB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	post := &db.Post{
		Title: "mongo Test_Update",
		Description: "testing mongo db update",
		Author: "kj",
		DateCreated: time.Now(),
		LastUpdated: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		MarkDown: "hello world",
	}

	id, err := db.Insert(post)
	if err != nil {
		t.Error(err)
	}

	post.Author = "mongo"

	id, err = db.Update(id, post)
	if err != nil {
		t.Error(err)
	}

	results, err := db.Find(db.FilterId{id}, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Author != "mongo" {
		t.Errorf("Expected: %s Received: %s\n", "mongo", results[0].Author)
	}

	err = db.Delete(id)
	if err != nil {
		t.Error(err)
	}
}

func Test_Distinct(t *testing.T) {
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

	expected := []string{"db", "foo", "test"}

	ids := [3]string{}
	for i, p := range posts {
		ids[i], err = db.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := db.Distinct("tags")
	if err != nil {
		t.Error(err)
	}

	for i, tag := range results {
		if tag != expected[i] {
			t.Errorf("Expected: %s Received: %s\n", expected[i], tag)
		}
	}
	
	for _, id := range ids {
		err = db.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}