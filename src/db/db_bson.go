package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"			json:"id"`
	Url      string             `bson:"url,omitempty"			json:"url"`
	Title    string             `bson:"title,omitempty"			json:"title"`
	Date     time.Time          `bson:"date,omitempty"			json:"date"`
	Tags     []string           `bson:"tags,omitempty"			json:"tags"`
	MarkDown string             `bson:"markdown,omitempty"		json:"markdown"`
	Html     string             `bson:"html,omitempty"			json:"html"`
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
