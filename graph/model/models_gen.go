// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Block struct {
	ReqUser   string `json:"reqUser"`
	BlockUser string `json:"blockUser"`
}

type Learning struct {
	Recommender float64 `json:"recommender"`
	Detector    float64 `json:"detector"`
}

type NewPost struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
	Image  string `json:"image"`
}

type NewReportShare struct {
	UserID string `json:"userID"`
	PostID string `json:"postID"`
}

type NewView struct {
	PostID        string  `json:"postID"`
	UserID        string  `json:"userId"`
	Time          float64 `json:"time"`
	DeviceDetails string  `json:"deviceDetails"`
}

type Payment struct {
	Amount float64 `json:"amount"`
	Postid string  `json:"postid"`
}

type PostOut struct {
	ID            string `json:"ID"`
	Text          string `json:"text"`
	UserID        string `json:"userId"`
	Views         int    `json:"views"`
	Shares        int    `json:"shares"`
	CreationTime  string `json:"creationTime"`
	InitialReview bool   `json:"initialReview"`
	Image         string `json:"image"`
	Blocked       bool   `json:"blocked"`
	Money         int    `json:"money"`
}

type PostReview struct {
	Post *PostOut `json:"post"`
	Left int      `json:"left"`
	Spam float64  `json:"spam"`
}

type Remove struct {
	PostID string `json:"postID"`
	UserID string `json:"userID"`
}

type Review struct {
	PostID   string `json:"postID"`
	Type     int    `json:"type"`
	Delete   bool   `json:"delete"`
	Password string `json:"password"`
}

type Stats struct {
	Userid string  `json:"userid"`
	Text   string  `json:"text"`
	Views  []*View `json:"views"`
	Money  float64 `json:"money"`
}

type View struct {
	Localization string `json:"localization"`
	Date         string `json:"date"`
}
