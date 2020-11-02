package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/geoloc"
	"QuicPos/internal/mongodb"
	"context"
	"log"
	"math/rand"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//View details struct
type View struct {
	UserID       int
	Time         float64 //relative to post content and shares JAKI WZOR?
	Localization string
	Lati         float64
	Long         float64
	Device       int
	Date         time.Time
}

//Post struct
type Post struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text          string
	UserID        int
	Reports       []int
	Views         []*View
	Shares        []int
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

//Share post
func Share(newReport model.NewReportShare) (bool, error) {
	post, err := GetByID(newReport.PostID)
	if err != nil {
		return false, err
	}
	shares := post.Shares
	shares = append(shares, newReport.UserID)

	objectID, _ := primitive.ObjectIDFromHex(newReport.PostID)
	_, err = mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"shares", shares}}},
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Report post
func Report(newReport model.NewReportShare) (bool, error) {
	post, err := GetByID(newReport.PostID)
	if err != nil {
		return false, err
	}
	reports := post.Reports
	reports = append(reports, newReport.UserID)

	objectID, _ := primitive.ObjectIDFromHex(newReport.PostID)
	_, err = mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"reports", reports}}},
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

//AddView to post
func AddView(newView model.NewView, ip string) (bool, error) {
	post, err := GetByID(newView.PostID)
	if err != nil {
		return false, err
	}
	views := post.Views
	loc, lati, long, err := geoloc.GetLocalization(ip)
	if err != nil {
		return false, err
	}
	view := &View{
		UserID:       newView.UserID,
		Time:         newView.Time,
		Localization: loc,
		Lati:         lati,
		Long:         long,
		Device:       newView.DeviceDetails,
		Date:         time.Now(),
	}
	views = append(views, view)

	objectID, _ := primitive.ObjectIDFromHex(newView.PostID)
	_, err = mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"views", views}}},
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
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
func GetOne(userID int) (Post, error) {
	reviewed := bson.D{{"$match", bson.M{"initialreview": true}}}
	notBlocked := bson.D{{"$match", bson.M{"blocked": false}}}
	notWatched := bson.D{{"$match", bson.M{"views": bson.M{"$ne": bson.M{"userid": userID}}}}}
	sample := bson.D{{"$sample", bson.D{{"size", 20}}}}

	//sometimes get only posts with less than 10 views
	lessViews := bson.D{{"$match", bson.M{"views.9": bson.M{"$exists": false}}}}

	//<0, 49> 1 to 50 chance of hapenning
	number := rand.Intn(50)
	if number == 20 {
		//only less than 10 views posts choosing
		showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{reviewed, notBlocked, notWatched, lessViews, sample})
		if err != nil {
			return Post{}, err
		}
		var showsLoaded []*Post
		if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
			return Post{}, err
		}
		for _, s := range showsLoaded {
			log.Println(s.Text)
		}
		return *showsLoaded[0], nil
	}

	//normal post choosing
	showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{reviewed, notBlocked, notWatched, sample})
	if err != nil {
		return Post{}, err
	}
	var showsLoaded []*Post
	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		return Post{}, err
	}
	for _, s := range showsLoaded {
		log.Println(s.Text)
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
