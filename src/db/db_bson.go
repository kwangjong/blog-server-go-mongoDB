package db

import (
	"time"
)

type Post struct {
	Url				string 				`bson:"url,omitempty"`
	Title 		 	string  			`bson:"title,omitempty"`
	Date			time.Time  			`bson:"date,omitempty"`
	Tags		 	[]string   			`bson:"tags,omitempty"`
	MarkDown       	string				`bson:"markdown,omitempty"`
	Html			string				`bson:"html,omitempty`
}

type FilterUrl 	  struct {Url    	string					`bson:"url"`}
type FilterTitle  struct {Title 	string  				`bson:"title"`}
type FilterTag    struct {Tag 		string 					`bson:"tags"`}