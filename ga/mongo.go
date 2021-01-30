package ga

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PushToMongo goes into each line of file
// and sends json object to mongo
func PushToMongo(filepath string) {

	// connect mongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect client to context
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// open file for reading
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// scan file content
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		jsonData := []byte(scanner.Text())

		var githubEvent GithubEvent

		err = json.Unmarshal(jsonData, &githubEvent)
		if err != nil {
			panic(err)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// post := Post{title, body}
		collection := client.Database("test").Collection("posts")

		_, err := collection.InsertOne(ctx, githubEvent)
		if err != nil {
			panic(err)

		}
	}

}

// GithubEvent struct to hold github events from json
type GithubEvent struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		DisplayLogin string `json:"display_login"`
		GravatarID   string `json:"gravatar_id"`
		URL          string `json:"url"`
		AvatarURL    string `json:"avatar_url"`
	} `json:"actor"`
	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
	Payload struct {
		Ref        string `json:"ref"`
		RefType    string `json:"ref_type"`
		PusherType string `json:"pusher_type"`
	} `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}
