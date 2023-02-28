package db

type Post struct {
	Title 		 string  	`bson:"title"`
	Description  string  	`bson:"description"`
	By			 string  	`bson:"by"`
	Timestamp    string  	`bson: "timestamp"`
	Tags		 []string   `bson:"tags"`
	Body       	 string		`bson:"body"`
}