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
		Url: "2023-03-23-Test-Insert-Delete",
		Title: "Test Insert Delete",
		Date: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		MarkDown: "hello world",
		Html: "<p>hello world</p>",
	}

	err = client.Insert(post)
	if err != nil {
		t.Error(err)
	}

	err = client.Delete(post.Url)
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
			Url: "0",
			Title: "mongo test1",
			Tags: []string{"db", "test"},
		},
		&db.Post{
			Url: "1",
			Title: "mongo test2",
			Tags: []string{"foo", "test"},
		}, 
		&db.Post{
			Url: "2",
			Title: "mongo test3",
			Tags: []string{"db"},
		},
	}

	for _, p := range posts {
		err = client.Insert(p)
		if err != nil {
			t.Error(err)
		}
	}

	result, err := client.Get("0")
	if err != nil {
		t.Error(err)
	}

	if result.Title != "mongo test1" {
		t.Errorf("Expected: %s Received: %s\n", "mongo test1", result.Title)
	}
	
	for _, i := range []string{"0", "1", "2"} {
		err = client.Delete(i)
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

	dates := [3]time.Time{}
	dates[0], _ = time.Parse("2006-Jan-02", "2014-Feb-01")
	dates[1], _ = time.Parse("2006-Jan-02", "2014-Feb-02")
	dates[2], _ = time.Parse("2006-Jan-02", "2014-Feb-03")

	posts := []*db.Post{
		&db.Post{
			Url: "0",
			Title: "mongo test1",
			Date: dates[0],
		},
		&db.Post{
			Url: "1",
			Title: "mongo test2",
			Date: dates[1],
		}, 
		&db.Post{
			Url: "2",
			Title: "mongo test3",
			Date: dates[2],
		},
	}

	for _, p := range posts {
		err = client.Insert(p)
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

	for _, i := range []string{"0", "1", "2"} {
		err = client.Delete(i)
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

	dates := [3]time.Time{}
	dates[0], _ = time.Parse("2006-Jan-02", "2014-Feb-01")
	dates[1], _ = time.Parse("2006-Jan-02", "2014-Feb-02")
	dates[2], _ = time.Parse("2006-Jan-02", "2014-Feb-03")

	posts := []*db.Post{
		&db.Post{
			Url: "0",
			Title: "mongo test1",
			Date: dates[0],
			Tags: []string{"db", "test"},
		},
		&db.Post{
			Url: "1",
			Title: "mongo test2",
			Date: dates[1],
			Tags: []string{"foo", "test"},
		}, 
		&db.Post{
			Url: "2",
			Title: "mongo test3",
			Date: dates[2],
			Tags: []string{"db"},
		},
	}

	for _, p := range posts {
		err = client.Insert(p)
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
	
	for _, i := range []string{"0", "1", "2"} {
		err = client.Delete(i)
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
		Url: "test-url",
		Title: "mongo Test_Update",
		Date: time.Now(),
		Tags: []string{"test", "mongodb", "go"},
		MarkDown: "hello world",
	}

	err = client.Insert(post)
	if err != nil {
		t.Error(err)
	}

	post.MarkDown = "hello mongo"

	err = client.Update(post.Url, post)
	if err != nil {
		t.Error(err)
	}

	results, err := client.Find(db.FilterUrl{post.Url}, 0, 3)
	if err != nil {
		t.Error(err)
	}

	if results[0].MarkDown == "hello mongo" {
		t.Errorf("Expected: %s Received: %s\n", "hello mongo", results[0].MarkDown)
	}

	err = client.Delete(post.Url)
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
			Url: "0",
			Title: "mongo test1",
			Tags: []string{"db", "test"},
		},
		&db.Post{
			Url: "1",
			Title: "mongo test2",
			Tags: []string{"foo", "test"},
		}, 
		&db.Post{
			Url: "2",
			Title: "mongo test3",
			Tags: []string{"db"},
		},
	}

	expected := []string{"db", "foo", "test"}

	for _, p := range posts {
		err = client.Insert(p)
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
	
	for _, i := range []string{"0", "1", "2"} {
		err = client.Delete(i)
		if err != nil {
			t.Error(err)
		}
	}
}