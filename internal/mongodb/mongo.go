package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is mongodb client
var Client *mongo.Client

//Ctx is context
var Ctx context.Context

//Cancel function for context
var Cancel func()

//InitDB starts connection with database
func InitDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:funia@quicpos.felpr.gcp.mongodb.net/quicpos?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	Cancel = cancel

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	Ctx = ctx
}

//DisconnectDB ends database connection
func DisconnectDB() {
	Client.Disconnect(Ctx)
	Cancel()
}

//List collections
func List() {
	databases, err := Client.ListDatabaseNames(Ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(databases)
}
