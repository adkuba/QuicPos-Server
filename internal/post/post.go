package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/data"
	"QuicPos/internal/geoloc"
	"QuicPos/internal/mongodb"
	"QuicPos/internal/stats"
	"QuicPos/internal/tensorflow"
	"context"
	"math/rand"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Share post
func Share(newReport model.NewReportShare) (bool, error) {
	post, err := GetByID(newReport.PostID)
	if err != nil {
		return false, err
	}
	shares := post.Shares
	for _, sh := range shares {
		if sh == newReport.UserID {
			return true, err
		}
	}
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
	for _, rep := range reports {
		if rep == newReport.UserID {
			return true, err
		}
	}
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
	view := &data.View{
		UserID:       newView.UserID,
		Time:         newView.Time,
		Localization: loc,
		Lati:         lati,
		Long:         long,
		Device:       newView.DeviceDetails,
		Date:         time.Now(),
	}
	views = append(views, view)

	if newView.PostID != "5fa55095fd6ff21ede156479" {
		err = stats.NewView(*view)
		if err != nil {
			return false, err
		}
	}

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

//NewOutsideView to stats
func NewOutsideView(postid string) error {
	view := &data.View{
		Time: 0,
	}
	objectID, _ := primitive.ObjectIDFromHex(postid)

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post data.Post
	result.Decode(&post)

	_, err := mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"outsideViews", post.OutsideViews + 1}}},
		},
	)
	if err != nil {
		return err
	}

	err = stats.NewView(*view)
	return err
}

//Save saves post to database
func Save(post data.Post) (string, error) {
	result, insertErr := mongodb.PostsCol.InsertOne(mongodb.Ctx, post)
	if insertErr != nil {
		return "", insertErr
	}

	err := stats.NewPost()
	if err != nil {
		return "", err
	}

	newID := result.InsertedID.(primitive.ObjectID).String()
	return newID, nil
}

//GetOneRandom post
func GetOneRandom() (data.Post, error) {
	sample := bson.D{{"$sample", bson.D{{"size", 1}}}}
	cursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{sample})
	if err != nil {
		return data.Post{}, err
	}
	var showsLoaded []*data.Post
	if err := cursor.All(context.TODO(), &showsLoaded); err != nil {
		return data.Post{}, err
	}
	err = NewOutsideView(showsLoaded[0].ID.String()[10:34])
	if err != nil {
		return data.Post{}, err
	}
	return *showsLoaded[0], nil
}

//GetOne gets one random post
func GetOne(userID int, ip string) (data.Post, error) {
	reviewed := bson.D{{"$match", bson.M{"initialreview": true}}}
	notBlocked := bson.D{{"$match", bson.M{"blocked": false}}}
	notWatched := bson.D{{"$match", bson.M{"views": bson.M{"$not": bson.M{"$elemMatch": bson.M{"userid": userID}}}}}}
	sample := bson.D{{"$sample", bson.D{{"size", 10}}}}

	//sometimes get only posts with less than 10 views
	lessViews := bson.D{{"$match", bson.M{"views.9": bson.M{"$exists": false}}}}

	//<0, 49> 1 to 50 chance of hapenning
	number := rand.Intn(50)
	var showLoadedCursor *mongo.Cursor
	if number == 20 {
		//only less than 10 views posts choosing
		cursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{reviewed, notBlocked, notWatched, lessViews, sample})
		if err != nil {
			return data.Post{}, err
		}
		showLoadedCursor = cursor
	} else {
		//normal post choosing
		cursor, err := mongodb.PostsCol.Aggregate(context.TODO(), mongo.Pipeline{reviewed, notBlocked, notWatched, sample})
		if err != nil {
			return data.Post{}, err
		}
		showLoadedCursor = cursor
	}

	//convert
	var showsLoaded []*data.Post
	if err := showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		return data.Post{}, err
	}

	//predict
	_, lati, long, err := geoloc.GetLocalization(ip)
	if err != nil {
		return data.Post{}, err
	}
	requesting := data.Requesting{
		UserID: userID,
		Lat:    lati,
		Long:   long,
		Date:   time.Now(),
	}
	bestValue := -1
	best := -1
	start := time.Now()
	for idx, post := range showsLoaded {
		inter, err := tensorflow.Recommend(*post, requesting)
		if err != nil {
			return data.Post{}, err
		}
		results := inter.([][]float32)
		categoryBest := float32(0)
		categoryIndex := -1
		for cat, result := range results[0] {
			if categoryBest < result {
				categoryBest = result
				categoryIndex = cat
			}
		}
		if categoryIndex > bestValue {
			bestValue = categoryIndex
			best = idx
		}
	}
	err = stats.NewProcessing(float64(time.Now().UnixNano()-start.UnixNano()) / 1000000000)
	if err != nil {
		return data.Post{}, err
	}
	//log.Println("Result: ", result[0])

	return *showsLoaded[best], nil

}

//GetByID gets post by id
func GetByID(id string) (data.Post, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return data.Post{}, err
	}

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post data.Post
	result.Decode(&post)
	err = NewOutsideView(id)
	if err != nil {
		return data.Post{}, err
	}
	return post, nil
}

//GetOneNew get the oldest post without initial review
func GetOneNew() (data.OutputReview, float32, error) {

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"creationtime", -1}})

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"initialreview": false}, findOptions)
	if err != nil {
		return data.OutputReview{}, 0, err
	}

	var posts []*data.Post
	if err = result.All(context.TODO(), &posts); err != nil {
		return data.OutputReview{}, 0, nil
	}
	if len(posts) > 0 {
		//predict
		inter, err := tensorflow.Spam(*posts[0])
		if err != nil {
			return data.OutputReview{}, 0, err
		}
		spam := inter.([][]float32)
		return data.OutputReview{*posts[0], len(posts)}, spam[0][0], nil
	}
	return data.OutputReview{data.Post{}, 0}, 0, nil
}

//GetOneReported gets the most repoted post
func GetOneReported() (data.OutputReview, float32, error) {

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"reports": bson.M{"$ne": nil}})
	if err != nil {
		return data.OutputReview{}, 0, err
	}

	var posts []*data.Post
	if err = result.All(context.TODO(), &posts); err != nil {
		return data.OutputReview{}, 0, nil
	}

	sort.SliceStable(posts, func(i, j int) bool {
		return len(posts[i].Reports) > len(posts[j].Reports)
	})

	if len(posts) > 0 {
		//predict
		inter, err := tensorflow.Spam(*posts[0])
		if err != nil {
			return data.OutputReview{}, 0, err
		}
		spam := inter.([][]float32)
		return data.OutputReview{*posts[0], len(posts)}, spam[0][0], nil
	}
	return data.OutputReview{data.Post{}, 0}, 0, nil
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
