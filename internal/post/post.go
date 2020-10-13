package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/mongodb"
	"context"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//OutputReview struct with number of posts left
type OutputReview struct {
	Post Output
	Left int
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

//GetByID gets post by id
func GetByID(id string) Output {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid id")
	}

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post Output
	result.Decode(&post)
	return post
}

//GetOneNew get the oldest post without initial review
func GetOneNew() OutputReview {

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"creationtime", -1}})

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"initialreview": false}, findOptions)
	if err != nil {
		log.Fatal("Can't find")
	}

	var posts []*Output
	if err = result.All(context.TODO(), &posts); err != nil {
		panic(err)
	}
	if len(posts) > 0 {
		return OutputReview{*posts[0], len(posts)}
	}
	return OutputReview{Output{}, 0}
}

//GetOneReported gets the most repoted post
func GetOneReported() OutputReview {

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"reports": bson.M{"$ne": nil}})
	if err != nil {
		log.Fatal("Can't find")
	}

	var posts []*Output
	if err = result.All(context.TODO(), &posts); err != nil {
		panic(err)
	}

	sort.SliceStable(posts, func(i, j int) bool {
		return len(posts[i].Reports) > len(posts[j].Reports)
	})

	if len(posts) > 0 {
		return OutputReview{*posts[0], len(posts)}
	}
	return OutputReview{Output{}, 0}
}

//ReviewAction reviews post
func ReviewAction(new bool, id string, delete bool) (bool, error) {

	objectID, _ := primitive.ObjectIDFromHex(id)
	// INITIAL REVIEW
	if new {
		_, err := mongodb.PostsCol.UpdateOne(
			context.TODO(),
			bson.M{"_id": objectID},
			bson.D{
				{"$set", bson.D{{"initialreview", true}}},
			},
		)
		if err != nil {
			return false, err
		}
		if delete {
			_, err := mongodb.PostsCol.UpdateOne(
				context.TODO(),
				bson.M{"_id": objectID},
				bson.D{
					{"$set", bson.D{{"blocked", true}}},
				},
			)
			if err != nil {
				return false, err
			}
		}
	} else {
		_, err := mongodb.PostsCol.UpdateOne(
			context.TODO(),
			bson.M{"_id": objectID},
			bson.D{
				{"$set", bson.D{{"reports", nil}}},
			},
		)
		if err != nil {
			return false, err
		}
		if delete {
			_, err := mongodb.PostsCol.UpdateOne(
				context.TODO(),
				bson.M{"_id": objectID},
				bson.D{
					{"$set", bson.D{{"blocked", true}}},
				},
			)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
