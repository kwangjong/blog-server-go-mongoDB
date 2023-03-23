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
	MONGO_KEY = "/home/kwangjong/kwangjong.github.io/key.json"
	MONGO_URL = "mongodb+srv://cluster0.rtswz75.mongodb.net"
)

type Mongo_Client struct {
	client   *mongo.Client
	postColl *mongo.Collection
}

func Connect_DB(collection string) (*Mongo_Client, error) {
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

	return &Mongo_Client {
		client:    client,
		postColl:  client.Database("db").Collection(collection),
	}, nil
}

func (mongo_client *Mongo_Client) Close() error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}

	err := mongo_client.client.Disconnect(context.TODO())
	log.Println("Mongo db connection closed")
	
	return err
}

func (mongo_client *Mongo_Client) Insert(post *Post) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl
	
	_, err := coll.InsertOne(context.TODO(), *post)
	if err != nil {
		return err
	}

	log.Printf("Inserted document with url: %s\n", post.Url)

	return nil
}

func (mongo_client *Mongo_Client) Delete(url string) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl
	
	filter := FilterUrl{url}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	log.Printf("Deleted document with url: %s\n", url)
	return nil
}

func (mongo_client *Mongo_Client) Update(url string, post *Post) error {
	if mongo_client.client == nil {
		return errors.New("Mongo db not connected")
	}
	coll := mongo_client.postColl

	filter := FilterUrl{url}
	update := bson.D{{"$set", *post}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	
	log.Printf("Updated document with url: %s\n", url)

	return nil
}

func (mongo_client *Mongo_Client) Get(url string) (*Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
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
	log.Printf("Found documents with url: %v", url)
	return results[0], err
}

func (mongo_client *Mongo_Client) Find(filter interface{}, skip int64, numPost int64) ([]*Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
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

func (mongo_client *Mongo_Client) Read(skip int64, numPost int64) ([]*Post, error) {
	if mongo_client.client == nil {
		return nil, errors.New("Mongo db not connected")
	}

	coll := mongo_client.postColl
	opts := options.Find().SetSort(bson.D{{"date", -1}}).SetSkip(skip).SetLimit(numPost).SetProjection(bson.D{
			{"url", 1}, 
			{"title", 1},
			{"date", 1},
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

func (mongo_client *Mongo_Client) Distinct(fieldName string) ([]interface{}, error) {
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