package devices

import (
	"QuicPos/internal/data"
	"QuicPos/internal/mongodb"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var counter = 0

//GetDevice id
func GetDevice(name string) (int, error) {
	number := exists(name)
	if number == -1 {
		newDevice := data.ViewModel{
			ID:     primitive.NewObjectIDFromTimestamp(time.Now()),
			Name:   name,
			Number: counter,
		}
		_, insertErr := mongodb.DevicesCol.InsertOne(mongodb.Ctx, newDevice)
		if insertErr != nil {
			return -1, insertErr
		}

		counter++
		err := updateCounter()
		if err != nil {
			return -1, err
		}

		return counter, nil
	}
	return number, nil
}

func updateCounter() error {

	infoID, _ := primitive.ObjectIDFromHex("5fd38ac0335016722636c1f2")
	_, err := mongodb.DevicesCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": infoID},
		bson.D{
			{"$set", bson.D{{"number", counter}}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//Init device manager
func Init() {
	result := mongodb.DevicesCol.FindOne(context.TODO(), bson.M{"name": "info-document"})
	var info data.ViewModel
	result.Decode(&info)
	counter = info.Number
}

func exists(name string) int {
	result := mongodb.DevicesCol.FindOne(context.TODO(), bson.M{"name": name})
	var device data.ViewModel
	result.Decode(&device)

	nullID, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	if device.ID == nullID {
		return -1
	}
	return device.Number
}
