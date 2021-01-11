package user

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"time"

	"QuicPos/internal/data"
	"QuicPos/internal/mongodb"
	"QuicPos/internal/stats"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetUser by id
func GetUser(id string) (data.User, error) {
	result := mongodb.UsersCol.FindOne(context.TODO(), bson.M{"uuid": id})
	var user data.User
	result.Decode(&user)
	return user, nil
}

//Block user by user
func Block(requestingUser string, blockUser string) error {

	user, err := GetUser(requestingUser)
	if err != nil {
		return err
	}

	user.Blocking = append(user.Blocking, blockUser)

	_, err = mongodb.UsersCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.D{
			{"$set", bson.D{{"blocking", user.Blocking}}},
		},
	)
	return err
}

//Create new user
func Create(ip string) (string, error) {
	uuid := uuid.New().String()
	uuid = strings.ReplaceAll(uuid, "-", "")

	user := &data.User{
		ID:       primitive.NewObjectIDFromTimestamp(time.Now()),
		UUID:     uuid,
		Blocking: nil,
	}
	_, insertErr := mongodb.UsersCol.InsertOne(mongodb.Ctx, user)
	if insertErr != nil {
		return "", insertErr
	}

	err := stats.NewUser()

	return uuid, err
}
