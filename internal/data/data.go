package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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
	UserID       int
	Time         float64 //relative to post content and shares JAKI WZOR?
	Localization string
	Lati         float64
	Long         float64
	Device       int
	Date         time.Time
}

//Requesting user struct
type Requesting struct {
	UserID int
	Lat    float64
	Long   float64
	Date   time.Time
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
