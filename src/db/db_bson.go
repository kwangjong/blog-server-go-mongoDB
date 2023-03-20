package db

import (
	"time"
)

type Post struct {
	ID				string 				`bson:"id,omitempty"`
	Title 		 	string  			`bson:"title,omitempty"`
	Description  	string  			`bson:"description,omitempty"`
	Author			string  			`bson:"author,omitempty"`
	DateCreated		time.Time  			`bson:"dateCreated,omitempty"`
	LastUpdated		time.Time			`bson:"lastUpdated,omitempty"`
	Tags		 	[]string   			`bson:"tags,omitempty"`
	MarkDown       	string				`bson:"markdown,omitempty"`
	Html			string				`bson:"html,omitempty`
}

type FilterId 	  struct {ID    	string					`bson:"id"`}
type FilterTitle  struct {Title 	string  				`bson:"title"`}
type FilterAuthor struct {Author	string  				`bson:"by"`}
type FilterTag    struct {Tag 		string 					`bson:"tags"`}