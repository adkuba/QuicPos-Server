package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is mongodb client
var Client *mongo.Client

//Ctx is context
var Ctx context.Context

//Cancel function for context
var Cancel func()

//PostsCol collection
var PostsCol *mongo.Collection

//InitDB starts connection with database
func InitDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:funia@quicpos.felpr.gcp.mongodb.net/quicpos?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	Cancel = cancel

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	posts := client.Database("quicpos").Collection("posts")

	PostsCol = posts
	Client = client
	Ctx = ctx
}

//DisconnectDB ends database connection
func DisconnectDB() {
	Client.Disconnect(Ctx)
	Cancel()
}
