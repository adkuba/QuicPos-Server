package mongodb

import (
	"QuicPos/internal/data"
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

//StatsCol collection
var StatsCol *mongo.Collection

//UsersCol collection
var UsersCol *mongo.Collection

//DevicesCol collection
var DevicesCol *mongo.Collection

//InitDB starts connection with database
func InitDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI(data.MongoSRV))
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
	stats := client.Database("quicpos").Collection("stats")
	devices := client.Database("quicpos").Collection("devices")
	users := client.Database("quicpos").Collection("users")

	PostsCol = posts
	StatsCol = stats
	DevicesCol = devices
	UsersCol = users
	Client = client
	Ctx = ctx
}

//DisconnectDB ends database connection
func DisconnectDB() {
	Client.Disconnect(Ctx)
	Cancel()
}
