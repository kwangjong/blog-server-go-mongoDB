package db

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID				primitive.ObjectID 	`bson:"_id,omitempty"`
	Title 		 	string  			`bson:"title"`
	Description  	string  			`bson:"description"`
	Author			string  			`bson:"author"`
	DateCreated		time.Time  			`bson:"dateCreated"`
	LastUpdated		time.Time			`bson:"lastUpdated"`
	Tags		 	[]string   			`bson:"tags"`
	Body       	 	string				`bson:"body"`
}

type FilterId 	  struct {ID    	primitive.ObjectID		`bson:"_id,omitempty"`}
type FilterTitle  struct {Title 	string  				`bson:"title"`}
type FilterAuthor struct {Author	string  				`bson:"by"`}
type FilterTag    struct {Tag 		string 					`bson:"tags"`}

func ObjectIDFromHex(id string) primitive.ObjectID {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return ObjId
}