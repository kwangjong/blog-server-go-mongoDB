package db

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MONGO_KEY = "/home/kwangjong/server/key.json"
	MONGO_URL = "mongodb+srv://cluster0.rtswz75.mongodb.net"
)

type Mongo_Struct struct {
	client   *mongo.Client
	postColl *mongo.Collection
}

var mongo_client Mongo_Struct

func Connect_DB() error {
	key_file, err := os.Open(MONGO_KEY)
    if err != nil {
        return err
    }
    defer key_file.Close()

    key_byte, _ := ioutil.ReadAll(key_file)

    var credential options.Credential
    json.Unmarshal([]byte(key_byte), &credential)

	opts := options.Client().ApplyURI(MONGO_URL).
	SetAuth(credential)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	log.Println("Mongo db connection established")

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return err
	}
	log.Println("Mongo db successfully pinged")

	mongo_client = Mongo_Struct {
		client:    client,
		postColl:  client.Database("db").Collection("posts"),
	}
	return nil
}

func Close() error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}

	err := mongo_client.client.Disconnect(context.TODO())
	log.Println("Mongo db connection closed")
	
	return err
}

func Insert(post *Post) (string, error) {
	if mongo_client.client == nil {
		return "", errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	if post.DateCreated.IsZero() {
		post.DateCreated = time.Now()
		post.LastUpdated = time.Now()
	} else if post.LastUpdated.IsZero() {
		post.LastUpdated = time.Now()
	}

	post.ID = fmt.Sprintf("%s-%s", post.DateCreated.Format("2006-01-02"), strings.Replace(post.Title, " ", "-", -1))
	
	_, err := coll.InsertOne(context.TODO(), *post)
	if err != nil {
		return "", err
	}

	log.Printf("Inserted document with id: %s\n", post.ID)

	return post.ID, nil
}

func Delete(id string) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl
	
	filter := FilterId{id}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	log.Printf("Deleted document with filter: %s\n", filter)
	return nil
}

func Update(id string, post *Post) (string, error) {
	if mongo_client.client == nil {
		return "", errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	if post.LastUpdated.IsZero() {
		post.LastUpdated = time.Now()
	}

	if post.Title != "" {
		orig_title := id[11:]
		new_title := strings.Replace(post.Title, " ", "-", -1)

		if orig_title != new_title {
			post.ID = fmt.Sprintf("%s-%s", id[:10], new_title)
		}
	}

	filter := FilterId{id}
	update := bson.D{{"$set", *post}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", err
	}
	
	log.Printf("Updated document with id: %s\n", post.ID)

	return post.ID, nil
}

func Find(filter interface{}, skip int64, numPost int64) ([]*Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	opts := options.Find().SetSort(bson.D{{"dateCreated", -1}}).SetSkip(skip).SetLimit(numPost)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var results []*Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d documents with filter: %v", len(results), filter)
	return results, err
}

func Read(skip int64, numPost int64) ([]*Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	opts := options.Find().SetSort(bson.D{{"dateCreated", -1}}).SetSkip(skip).SetLimit(numPost).SetProjection(bson.D{
			{"id", 1}, 
			{"title", 1},
			{"description", 1},
			{"author", 1},
			{"dateCreated", 1},
			{"lastUpdated", 1},
			{"tags", 1},
		})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []*Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	log.Printf("Read %d documents\n", len(results))
	return results, err
}

func Distinct(fieldName string) ([]interface{}, error) {
	results := []interface{}{}
	if mongo_client.client == nil {
		return results, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	results, err := coll.Distinct(context.TODO(), fieldName, bson.D{})
	if err != nil {
    	return results, err
	}
	log.Printf("Found %d distinct %s\n", len(results), fieldName)
	return results, nil
}