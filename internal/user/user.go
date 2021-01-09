package user

import (
	"context"
	"github.com/google/uuid"
	"log"
	"sort"
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

	result, err := mongodb.UsersCol.Find(context.TODO(), bson.M{})
	if err != nil {
		return -1, err
	}
	var users []*data.User
	if err = result.All(context.TODO(), &users); err != nil {
		return -1, nil
	}

	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Int > users[j].Int
	})

	return users[0].Int, nil
}

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

//CheckCounter on server start
func CheckCounter() {
	lastUser, err := getLastUser()
	if err != nil {
		panic("Cant find user info document!")
	}
	counter = lastUser
	log.Println("Last user int: ", counter)
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

	err := stats.NewUser()
	if err != nil {
		return -1, err
	}

	return counter, nil
}
