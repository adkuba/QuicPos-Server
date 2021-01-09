package post

import (
	"QuicPos/graph/model"
	"QuicPos/internal/data"
	"QuicPos/internal/devices"
	"QuicPos/internal/geoloc"
	"QuicPos/internal/mongodb"
	"QuicPos/internal/stats"
	"QuicPos/internal/tensorflow"
	"QuicPos/internal/user"
	"context"
	"errors"
	"math/rand"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//AddMoney to post
func AddMoney(payment model.Payment) (bool, error) {
	post, err := GetByID(payment.Postid, false)
	if err != nil {
		return false, err
	}

	objectID, _ := primitive.ObjectIDFromHex(payment.Postid)
	_, err = mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"money", post.Money + int(payment.Amount*100)}}},
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Remove post
func Remove(postID string, userUUID string) error {

	post, err := GetByID(postID, false)
	if err != nil {
		return err
	}

	if post.User.UUID == userUUID {
		objectID, _ := primitive.ObjectIDFromHex(postID)
		_, err := mongodb.PostsCol.DeleteOne(
			context.TODO(),
			bson.M{"_id": objectID},
		)
		return err
	}
	return errors.New("bad user request")
}

//Share post
func Share(newReport model.NewReportShare) (bool, error) {
	post, err := GetByID(newReport.PostID, false)
	if err != nil {
		return false, err
	}
	shares := post.Shares

	user, err := user.GetUser(newReport.UserID)
	if err != nil {
		return false, err
	}

	for _, sh := range shares {
		if sh == &user {
			return true, err
		}
	}
	shares = append(shares, &user)

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
	post, err := GetByID(newReport.PostID, false)
	if err != nil {
		return false, err
	}
	reports := post.Reports

	user, err := user.GetUser(newReport.UserID)
	if err != nil {
		return false, err
	}

	for _, rep := range reports {
		if rep == &user {
			return true, err
		}
	}
	reports = append(reports, &user)

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
	post, err := GetByID(newView.PostID, false)
	if err != nil {
		return false, err
	}
	views := post.Views
	loc, lati, long, err := geoloc.GetLocalization(ip)
	if err != nil {
		return false, err
	}
	deviceID, err := devices.GetDevice(newView.DeviceDetails)
	if err != nil {
		return false, err
	}

	user, err := user.GetUser(newView.UserID)
	if err != nil {
		return false, err
	}

	view := &data.View{
		User:         user,
		Time:         newView.Time,
		Localization: loc,
		Lati:         lati,
		Long:         long,
		Device:       deviceID,
		Date:         time.Now(),
	}
	views = append(views, view)

	err = stats.NewView(*view)
	if err != nil {
		return false, err
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
		Time:         0,
		Localization: "Private",
		Date:         time.Now(),
	}
	objectID, _ := primitive.ObjectIDFromHex(postid)

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post data.Post
	result.Decode(&post)

	var outside = post.OutsideViews
	outside = append(outside, view)

	_, err := mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"outsideviews", outside}}},
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
func GetOne(userID string, ip string, ad bool) (data.Post, error) {
	reviewed := bson.D{{"$match", bson.M{"initialreview": true}}}
	notBlocked := bson.D{{"$match", bson.M{"blocked": false}}}
	notWatched := bson.D{{"$match", bson.M{"views": bson.M{"$not": bson.M{"$elemMatch": bson.M{"userid": userID}}}}}}
	sample := bson.D{{"$sample", bson.D{{"size", 10}}}}
	lessViews := bson.D{{"$match", bson.M{"views.9": bson.M{"$exists": false}}}}
	ads := bson.D{{"$match", bson.M{"money": bson.M{"$gt": 0}}}}

	user, err := user.GetUser(userID)
	if err != nil {
		return data.Post{}, err
	}

	if user.Blocking == nil {
		user.Blocking = append(user.Blocking, "")
	}

	//user blocking
	notUserBlocked := bson.D{{"$match", bson.M{"user.uuid": bson.M{"$nin": user.Blocking}}}}

	//normal
	pipeline := mongo.Pipeline{reviewed, notUserBlocked, notBlocked, notWatched, sample}

	//sometimes get only posts with less than 10 views
	//<0, 49> 1 to 50 chance of hapenning
	number := rand.Intn(50)
	if number == 20 {
		pipeline = mongo.Pipeline{reviewed, notUserBlocked, notBlocked, notWatched, lessViews, sample}
	}

	//ad choosing
	if ad {
		pipeline = mongo.Pipeline{reviewed, notUserBlocked, notBlocked, ads, sample}
	}

	showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return data.Post{}, err
	}

	//convert
	var showsLoaded []*data.Post
	if err := showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		return data.Post{}, err
	}

	//no ads
	if len(showsLoaded) == 0 {
		pipeline = mongo.Pipeline{reviewed, notUserBlocked, notBlocked, notWatched, sample}
		showLoadedCursor, err := mongodb.PostsCol.Aggregate(context.TODO(), pipeline)
		if err != nil {
			return data.Post{}, err
		}
		if err := showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
			return data.Post{}, err
		}
		ad = false
	}

	//predict
	_, lati, long, err := geoloc.GetLocalization(ip)
	if err != nil {
		return data.Post{}, err
	}

	requesting := data.Requesting{
		User: user,
		Lat:  lati,
		Long: long,
		Date: time.Now(),
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

	//ad calculation
	if ad {
		err = reduceMoney(showsLoaded[best].ID, showsLoaded[best].Money)
		if err != nil {
			return data.Post{}, err
		}
	}

	return *showsLoaded[best], nil

}

//reduce ad view money
func reduceMoney(id primitive.ObjectID, money int) error {
	_, err := mongodb.PostsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"money", money - 1}}},
		},
	)
	return err
}

//GetByID gets post by id
func GetByID(id string, countView bool) (data.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return data.Post{}, err
	}

	result := mongodb.PostsCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var post data.Post
	result.Decode(&post)
	if countView {
		err = NewOutsideView(id)
		if err != nil {
			return data.Post{}, err
		}
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
