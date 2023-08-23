package db

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MONGO_KEY = "/home/kwangjong/blog-server-go-mongoDB/key.json"
	MONGO_URL = "mongodb+srv://cluster0.rtswz75.mongodb.net"
)

type DBClient struct {
	client *mongo.Client
}

type DBCollection struct {
	collection *mongo.Collection
}

func Connect_DB() (*DBClient, error) {
	key_file, err := os.Open(MONGO_KEY)
	if err != nil {
		return nil, err
	}
	defer key_file.Close()

	key_byte, _ := ioutil.ReadAll(key_file)

	var credential options.Credential
	json.Unmarshal([]byte(key_byte), &credential)

	opts := options.Client().ApplyURI(MONGO_URL).
		SetAuth(credential)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	log.Println("Mongo db connection established")

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}
	log.Println("Mongo db successfully pinged")

	return &DBClient{client}, nil
}

func (db *DBClient) Collection(collection string) (*DBCollection, error) {
	if db.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	return &DBCollection{db.client.Database("db").Collection(collection)}, nil
}

func (db *DBClient) Close() error {
	if db.client == nil {
		return errors.New("Mongo db not connected")
	}

	err := db.client.Disconnect(context.TODO())
	log.Println("Mongo db connection closed")

	return err
}

func (db_coll *DBCollection) Insert(post *Post) error {
	coll := db_coll.collection

	_, err := coll.InsertOne(context.TODO(), *post)
	if err != nil {
		return err
	}

	log.Printf("Inserted document with url: %s\n", post.Url)

	return nil
}

func (db_coll *DBCollection) Delete(url string) error {
	coll := db_coll.collection

	filter := FilterUrl{url}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	log.Printf("Deleted document with url: %s\n", url)
	return nil
}

func (db_coll *DBCollection) Update(post *Post) error {
	coll := db_coll.collection

	filter := FilterId{post.Id}
	update := bson.D{{"$set", *post}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	log.Printf("Updated document with _id: %s\n", post.Id)

	return nil
}

func (db_coll *DBCollection) Get(url string) (*Post, error) {
	coll := db_coll.collection

	filter := FilterUrl{url}
	opts := options.Find()
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var results []*Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("page not found")
	}
	log.Printf("Found documents with url: %v", url)
	return results[0], err
}

func (db_coll *DBCollection) Find(filter interface{}, skip int64, numPost int64) ([]*Post, error) {
	coll := db_coll.collection

	opts := options.Find().SetSort(bson.D{{"date", -1}}).SetSkip(skip).SetLimit(numPost).SetProjection(bson.D{
		{"url", 1},
		{"title", 1},
		{"author", 1},
		{"date", 1},
		{"tags", 1},
	})
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

func (db_coll *DBCollection) Read(skip int64, numPost int64) ([]*Post, bool, error) {
	coll := db_coll.collection

	opts := options.Find().SetSort(bson.D{{"date", -1}}).SetSkip(skip).SetLimit(numPost + 1).SetProjection(bson.D{
		{"url", 1},
		{"title", 1},
		{"date", 1},
		{"tags", 1},
	})

	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, false, err
	}

	var results []*Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, false, err
	}
	log.Printf("Read %d documents\n", len(results))

	if len(results) < int(numPost)+1 {
		return results, false, nil
	}

	return results[:len(results)-1], true, nil
}

func (db_coll *DBCollection) Distinct(fieldName string) ([]interface{}, error) {
	results := []interface{}{}
	coll := db_coll.collection

	results, err := coll.Distinct(context.TODO(), fieldName, bson.D{})
	if err != nil {
		return results, err
	}
	log.Printf("Found %d distinct %s\n", len(results), fieldName)
	return results, nil
}
