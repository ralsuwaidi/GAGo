package ga

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Title string `json:"title,omitempty"`

	Body string `json:"body,omitempty"`
}

func MongoConnect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017/new_test"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
}

func InsertPost(title string, body string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	post := Post{title, body}
	collection := client.Database("test").Collection("posts")

	insertResult, err := collection.InsertOne(ctx, post)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)

}
