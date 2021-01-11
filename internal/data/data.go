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
	User         string
	Time         float64 //relative to post content and shares JAKI WZOR?
	Localization string
	IP           string
	Device       string
	Date         time.Time
}

//Device struct for database
type Device struct {
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string
}

//Requesting user struct
type Requesting struct {
	User string
	IP   string
	Date time.Time
}

//User struct
type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UUID     string
	Blocking []string
}

//Post struct
type Post struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text          string
	User          string
	Reports       []*string
	Views         []*View
	Shares        []*string
	CreationTime  time.Time
	Image         string
	InitialReview bool
	Blocked       bool
	OutsideViews  []*View
	Money         int
	ImageFeatures []float32
	HumanReview   bool
}

//OutputReview struct with number of posts left
type OutputReview struct {
	Post Post
	Left int
}
