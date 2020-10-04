package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/mongodb"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Post struct
type Post struct {
	Text          string
	UserID        string
	Reports       []string
	CreationTime  time.Time
	Image         string
	InitialReview bool
	Views         []*model.View
	Shares        int
}

//Output post struct
type Output struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text          string
	UserID        string
	Reports       []string
	Views         []*model.View
	Shares        int
	CreationTime  time.Time
	Image         string
	InitialReview bool
}

//Save saves post to database
func (post Post) Save() string {
	result, insertErr := mongodb.PostsCol.InsertOne(mongodb.Ctx, post)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	newID := result.InsertedID.(primitive.ObjectID).String()
	return newID
}

//GetOne gets one random post
func GetOne() Output {
	o1 := bson.D{{"$sample", bson.D{{"size", 1}}}}
	showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{o1})
	if err != nil {
		log.Fatal(err)
	}
	var showsLoaded []*Output
	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		panic(err)
	}
	return *showsLoaded[0]
}
