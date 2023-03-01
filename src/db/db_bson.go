package db

import "time"

type Post struct {
	Title 		 	string  	`bson:"title"`
	Description  	string  	`bson:"description"`
	By			 	string  	`bson:"by"`
	DateCreated		time.Time  	`bson:"dateCreated"`
	LastUpdated		time.Time	`bson:"lastUpdated"`
	Tags		 	[]string   	`bson:"tags"`
	Body       	 	string		`bson:"body"`
}