package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"			json:"id"`
	Url      string             `bson:"url"					json:"url"`
	Title    string             `bson:"title"				json:"title"`
	Date     time.Time          `bson:"date"				json:"date"`
	Tags     []string           `bson:"tags"				json:"tags"`
	MarkDown string             `bson:"markdown"				json:"markdown"`
	Html     string             `bson:"html"				json:"html"`
}
type FilterId struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
}
type FilterUrl struct {
	Url string `bson:"url"`
}
type FilterTitle struct {
	Title string `bson:"title"`
}
type FilterTag struct {
	Tag string `bson:"tags"`
}
