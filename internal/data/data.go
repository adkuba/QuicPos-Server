package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//Pass to access
var Pass = ""

//AdminPass to access
var AdminPass = ""

//Day for stats
type Day struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Date           time.Time
	NewUsers       int
	NewPosts       int
	Views          int
	WatchTime      float64
	ProcessingTime float64
	Requests       int
	Recommender    float64
	Detector       float64
}

//View details struct
type View struct {
	User         User
	Time         float64 //relative to post content and shares JAKI WZOR?
	Localization string
	Lati         float64
	Long         float64
	Device       int
	Date         time.Time
}

//ViewModel struct for database
type ViewModel struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name   string
	Number int
}

//UserModel struct for database
type UserModel struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	LastUser int
}

//Requesting user struct
type Requesting struct {
	User User
	Lat  float64
	Long float64
	Date time.Time
}

//User struct
type User struct {
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UUID string
	Int  int
}

//Post struct
type Post struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text          string
	User          User
	Reports       []*User
	Views         []*View
	Shares        []*User
	CreationTime  time.Time
	Image         string
	InitialReview bool
	Blocked       bool
	OutsideViews  []*View
	Money         int
}

//OutputReview struct with number of posts left
type OutputReview struct {
	Post Post
	Left int
}
