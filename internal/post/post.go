package post

import (
	"QuicPos/internal/mongodb"
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//View details struct
type View struct {
	UserID string
	Time   float32
}

//Post struct
type Post struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text          string
	UserID        string
	Reports       []string
	Views         []*View
	Shares        []string
	CreationTime  time.Time
	Image         string
	InitialReview bool
	Blocked       bool
}

//OutputReview struct with number of posts left
type OutputReview struct {
	Post Post
	Left int
}

//Save saves post to database
func (post Post) Save() (string, error) {
	result, insertErr := mongodb.PostsCol.InsertOne(mongodb.Ctx, post)
	if insertErr != nil {
		return "", insertErr
	}
	newID := result.InsertedID.(primitive.ObjectID).String()
	return newID, nil
}

//GetOne gets one random post
func GetOne() (Post, error) {
	o1 := bson.D{{"$sample", bson.D{{"size", 1}}}}
	showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{o1})
	if err != nil {
		return Post{}, err
	}
	var showsLoaded []*Post
	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		return Post{}, err
	}
	return *showsLoaded[0], nil
}

//GetByID gets post by id
func GetByID(id string) (Post, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Post{}, err
	}

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post Post
	result.Decode(&post)
	return post, nil
}

//GetOneNew get the oldest post without initial review
func GetOneNew() (OutputReview, error) {

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"creationtime", -1}})

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"initialreview": false}, findOptions)
	if err != nil {
		return OutputReview{}, err
	}

	var posts []*Post
	if err = result.All(context.TODO(), &posts); err != nil {
		return OutputReview{}, nil
	}
	if len(posts) > 0 {
		return OutputReview{*posts[0], len(posts)}, nil
	}
	return OutputReview{Post{}, 0}, nil
}

//GetOneReported gets the most repoted post
func GetOneReported() (OutputReview, error) {

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"reports": bson.M{"$ne": nil}})
	if err != nil {
		return OutputReview{}, err
	}

	var posts []*Post
	if err = result.All(context.TODO(), &posts); err != nil {
		return OutputReview{}, nil
	}

	sort.SliceStable(posts, func(i, j int) bool {
		return len(posts[i].Reports) > len(posts[j].Reports)
	})

	if len(posts) > 0 {
		return OutputReview{*posts[0], len(posts)}, nil
	}
	return OutputReview{Post{}, 0}, nil
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
