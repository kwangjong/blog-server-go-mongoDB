package test

import (
	"testing"
	"time"
	"github.com/kwangjong/kwangjong.github.io/db"
)

func Test_Connect_DB_Close(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}
}

func Test_Insert_Delete(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

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

	id, err := client.Insert(post)
	if err != nil {
		t.Error(err)
	}

	err = client.Delete(id)
	if err != nil {
		t.Error(err)
	}
}

func Test_Get(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

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
		ids[i], err = client.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	result, err := client.Get(db.FilterId{ids[0]})
	if err != nil {
		t.Error(err)
	}

	if result.Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", result.Title)
	}
	
	for _, id := range ids {
		err = client.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Read(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

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
		ids[i], err = client.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := client.Read(0, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range posts {
		if p.Title != results[len(results)-1-i].Title {
		 	t.Errorf("Expected: %s Received: %s\n",p.Title, results[len(results)-1-i].Title)
		}
	}

	results, err = client.Read(1, 3)
	if err != nil {
		t.Error(err)
	}

	for i, p := range posts[:2] {
		if p.Title != results[len(results)-1-i].Title {
		 	t.Errorf("Expected: %s Received: %s\n", p.Title, results[len(results)-1-i].Title)
		}
	}

	for _, id := range ids {
		err = client.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Find(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

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
		ids[i], err = client.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := client.Find(db.FilterTag{"db"}, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test3" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test3", results[0].Title)
	}
	
	if results[1].Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[1].Title)
	}

	results, err = client.Find(db.FilterTag{"test"}, 1, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", results[0].Title)
	}
	
	for _, id := range ids {
		err = client.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_Update(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	post := &db.Post{
		Title: "mongo Test_Update",
		Description: "testing mongo db update",
		Author: "kj",
		DateCreated: time.Now(),
		LastUpdated: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		MarkDown: "hello world",
	}

	id, err := client.Insert(post)
	if err != nil {
		t.Error(err)
	}

	post.Author = "mongo"

	id, err = client.Update(id, post)
	if err != nil {
		t.Error(err)
	}

	results, err := client.Find(db.FilterId{id}, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].Author != "mongo" {
		t.Errorf("Expected: %s Received: %s\n", "mongo", results[0].Author)
	}

	err = client.Delete(id)
	if err != nil {
		t.Error(err)
	}
}

func Test_Distinct(t *testing.T) {
	client, err := db.Connect_DB("test")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

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
		ids[i], err = client.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	results, err := client.Distinct("tags")
	if err != nil {
		t.Error(err)
	}

	for i, tag := range results {
		if tag != expected[i] {
			t.Errorf("Expected: %s Received: %s\n", expected[i], tag)
		}
	}
	
	for _, id := range ids {
		err = client.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}