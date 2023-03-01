package db

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MONGO_KEY = "key.json"
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

func Insert(post *Post) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl
	
	id, err := coll.InsertOne(context.TODO(), *post)
	if err != nil {
		return err
	}
	log.Printf("Inserted document with _id: %s\n", id)
	return nil
}

func Delete(id string) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	filter := bson.D{{"_id", object_id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	log.Printf("Deleted document with _id: %s\n", id)
	return nil
}

func Update(id string, post *Post) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", object_id}}
	update := bson.D{{"$set", *post}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	log.Printf("Updated document with _id: %s\n", id)
	return err
}

func Read(skip int64, numPost int64) (*[]Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"dateCreated", -1}}).SetSkip(skip).SetLimit(numPost)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var results []Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	log.Printf("Read %d documents%v", numPost)
	return &results, err
}

func Find(filter interface{}, skip int64, numPost int64) (*[]Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	opts := options.Find().SetSort(bson.D{{"dateCreated", -1}}).SetSkip(skip).SetLimit(numPost)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var results []Post
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	log.Printf("Read %d documents with filter: %v", numPost, filter)
	return &results, err
}