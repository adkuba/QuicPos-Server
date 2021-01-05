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

var counter = 0

func getLastUser() (int, error) {
	objectID, err := primitive.ObjectIDFromHex("5fe9a7132f3a5604036e87cd")
	if err != nil {
		return -1, err
	}

	result := mongodb.UsersCol.FindOne(context.TODO(), bson.M{"_id": objectID})
	var user data.UserModel
	result.Decode(&user)
	return user.LastUser, nil
}

//GetUser by id
func GetUser(id string) (data.User, error) {
	result := mongodb.UsersCol.FindOne(context.TODO(), bson.M{"uuid": id})
	var user data.User
	result.Decode(&user)
	return user, nil
}

//CheckCounter on server start
func CheckCounter() {
	lastUser, err := getLastUser()
	if err != nil {
		panic("Cant find user info document!")
	}
	counter = lastUser
}

func updateLastUser() error {
	objectID, _ := primitive.ObjectIDFromHex("5fe9a7132f3a5604036e87cd")
	_, err := mongodb.UsersCol.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"lastuser", counter}}},
		},
	)
	return err
}

//Create new user
func Create(ip string) (string, error) {
	uuid := uuid.New().String()
	uuid = strings.ReplaceAll(uuid, "-", "")
	intNum, err := getNextUser(ip)
	if err != nil {
		return "", err
	}
	user := &data.User{
		ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		UUID: uuid,
		Int:  intNum,
	}
	_, insertErr := mongodb.UsersCol.InsertOne(mongodb.Ctx, user)
	return uuid, insertErr
}

//GetNextUser id
func getNextUser(ip string) (int, error) {
	counter++

	err := updateLastUser()
	if err != nil {
		return -1, err
	}

	err = stats.NewUser()
	if err != nil {
		return -1, err
	}

	return counter, nil
}
