package stats

import (
	"QuicPos/internal/data"
	"QuicPos/internal/mongodb"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//check and create day object if not exists
func checkDay() (primitive.ObjectID, error) {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"date": today})
	var stat data.Day
	result.Decode(&stat)

	nullID, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	if stat.ID == nullID {
		newStat := data.Day{
			ID:             primitive.NewObjectIDFromTimestamp(time.Now()),
			Date:           today,
			NewUsers:       0,
			NewPosts:       0,
			Views:          0,
			WatchTime:      0,
			ProcessingTime: 0,
			Requests:       0,
			Recommender:    0,
			Detector:       0,
		}
		result, insertErr := mongodb.StatsCol.InsertOne(mongodb.Ctx, newStat)
		if insertErr != nil {
			return nullID, insertErr
		}
		return result.InsertedID.(primitive.ObjectID), nil
	}

	return stat.ID, nil
}

//NewUser to stats
func NewUser() error {
	statID, err := checkDay()

	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"_id": statID})
	var stat data.Day
	result.Decode(&stat)

	_, err = mongodb.StatsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": statID},
		bson.D{
			{"$set", bson.D{{"newusers", stat.NewUsers + 1}}},
		},
	)
	return err
}

//UpdateNets to mongo
func UpdateNets(recommender float64, detector float64) error {
	statID, err := checkDay()

	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"_id": statID})
	var stat data.Day
	result.Decode(&stat)

	_, err = mongodb.StatsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": statID},
		bson.D{
			{"$set", bson.D{{"recommender", recommender}}},
			{"$set", bson.D{{"detector", detector}}},
		},
	)
	return err
}

//NewProcessing adds time spend on recommending post by the net
func NewProcessing(time float64) error {
	statID, err := checkDay()

	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"_id": statID})
	var stat data.Day
	result.Decode(&stat)

	_, err = mongodb.StatsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": statID},
		bson.D{
			{"$set", bson.D{{"processingtime", stat.ProcessingTime + time}}},
			{"$set", bson.D{{"requests", stat.Requests + 1}}},
		},
	)
	return err
}

//NewView to stats
func NewView(view data.View) error {
	statID, err := checkDay()

	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"_id": statID})
	var stat data.Day
	result.Decode(&stat)

	_, err = mongodb.StatsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": statID},
		bson.D{
			{"$set", bson.D{{"views", stat.Views + 1}}},
			{"$set", bson.D{{"watchtime", stat.WatchTime + view.Time}}},
		},
	)
	return err
}

//NewPost to stats
func NewPost() error {
	statID, err := checkDay()

	result := mongodb.StatsCol.FindOne(context.TODO(), bson.M{"_id": statID})
	var stat data.Day
	result.Decode(&stat)

	_, err = mongodb.StatsCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": statID},
		bson.D{
			{"$set", bson.D{{"newposts", stat.NewPosts + 1}}},
		},
	)
	return err
}
