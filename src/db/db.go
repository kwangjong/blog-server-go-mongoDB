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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func Insert(post *Post) (primitive.ObjectID, error) {
	if mongo_client.client == nil {
		return primitive.NilObjectID, errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl
	
	result, err := coll.InsertOne(context.TODO(), *post)
	if err != nil {
		return primitive.NilObjectID, err
	}

	log.Printf("Inserted document with _id: %s\n", result.InsertedID)

	id, _:= result.InsertedID.(primitive.ObjectID)
	return id, nil
}

func Delete(id primitive.ObjectID) error {
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

func Update(id primitive.ObjectID, post *Post) (primitive.ObjectID, error) {
	if mongo_client.client == nil {
		return primitive.NilObjectID, errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	filter := FilterId{id}
	update := bson.D{{"$set", *post}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return primitive.NilObjectID, err
	}
	
	log.Printf("Updated document with _id: %s\n", result.UpsertedID)

	id, _ = result.UpsertedID.(primitive.ObjectID)
	return id, nil
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
	opts := options.Find().SetSort(bson.D{{"dateCreated", -1}}).SetSkip(skip).SetLimit(numPost)
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