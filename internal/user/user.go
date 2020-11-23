package user

import (
	"QuicPos/graph/model"
	"QuicPos/internal/post"
	"QuicPos/internal/stats"
)

var counter = 0

//CheckCounter on server start
func CheckCounter() {
	initPost, err := post.GetByID("5fa55095fd6ff21ede156479", false)
	if err != nil {
		panic("Cant find init post!")
	}
	maxUser := -1
	for _, view := range initPost.Views {
		if view.UserID > maxUser {
			maxUser = view.UserID
		}
	}
	if maxUser != -1 {
		counter = maxUser
	}
}

//GetNextUser id
func GetNextUser(ip string) (int, error) {
	counter++
	newView := model.NewView{
		PostID:        "5fa55095fd6ff21ede156479",
		UserID:        counter,
		Time:          0,
		DeviceDetails: 0,
	}
	_, err := post.AddView(newView, ip)
	if err != nil {
		return -1, err
	}

	err = stats.NewUser()
	if err != nil {
		return -1, err
	}

	return counter, nil
}
