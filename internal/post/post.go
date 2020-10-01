package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/mongodb"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Post struct
type Post struct {
	Text   string
	UserID string
	Views  []*model.View
	Shares int
}

//Output post struct
type Output struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text   string
	UserID string
	Views  []*model.View
	Shares int
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
	var result Output
	err := mongodb.PostsCol.FindOne(context.TODO(), bson.D{}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
