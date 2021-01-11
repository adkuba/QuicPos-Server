package devices

import (
	"QuicPos/internal/data"
	"QuicPos/internal/mongodb"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetDevice id
func GetDevice(name string) (string, error) {
	id := exists(name)
	if id == "" {
		newDevice := data.Device{
			ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
			Name: name,
		}
		_, insertErr := mongodb.DevicesCol.InsertOne(mongodb.Ctx, newDevice)
		if insertErr != nil {
			return "", insertErr
		}

		return newDevice.ID.String(), nil
	}
	return id, nil
}

func exists(name string) string {
	result := mongodb.DevicesCol.FindOne(context.TODO(), bson.M{"name": name})
	var device data.Device
	result.Decode(&device)

	nullID, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	if device.ID == nullID {
		return ""
	}
	return device.ID.String()
}
