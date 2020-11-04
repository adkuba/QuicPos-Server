package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"QuicPos/graph/generated"
	"QuicPos/graph/model"
	"QuicPos/internal/data"
	"QuicPos/internal/ip"
	"QuicPos/internal/post"
	"QuicPos/internal/storage"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.PostOut, error) {
	var postO data.Post
	postO.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	postO.Text = input.Text
	postO.UserID = input.UserID
	postO.CreationTime = time.Now()
	postO.Image = storage.UploadFile(input.Image)
	postO.InitialReview = false
	postO.Reports = nil
	postO.Shares = nil
	postO.Views = nil
	postO.Blocked = false
	postID, err := post.Save(postO)
	return &model.PostOut{ID: postID, Text: postO.Text, UserID: postO.UserID, Shares: len(postO.Shares), Views: len(postO.Views), CreationTime: postO.CreationTime.String(), InitialReview: postO.InitialReview, Image: postO.Image, Blocked: postO.Blocked}, err
}

func (r *mutationResolver) Review(ctx context.Context, input model.Review) (bool, error) {
	if input.Password == "funia" {
		result, err := post.ReviewAction(input.New, input.PostID, input.Delete)
		return result, err
	}
	return false, errors.New("bad password")
}

func (r *mutationResolver) Share(ctx context.Context, input model.NewReportShare) (bool, error) {
	result, err := post.Share(input)
	return result, err
}

func (r *mutationResolver) Report(ctx context.Context, input model.NewReportShare) (bool, error) {
	result, err := post.Report(input)
	return result, err
}

func (r *mutationResolver) View(ctx context.Context, input model.NewView) (bool, error) {
	//tu mam dokladne dane urzadzenia
	status, err := post.AddView(input, ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).IP)
	return status, err
}

func (r *queryResolver) Post(ctx context.Context, userID int, normalMode bool) (*model.PostOut, error) {
	//userID and normalMode to be used
	post, err := post.GetOne(userID, ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).IP)
	return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Shares: len(post.Shares), Views: len(post.Views), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked}, err
}

func (r *queryResolver) CreateUser(ctx context.Context) (int, error) {
	counter++
	return counter, nil
}

func (r *queryResolver) ViewerPost(ctx context.Context, id string) (*model.PostOut, error) {
	post, err := post.GetByID(id)
	return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Shares: len(post.Shares), Views: len(post.Views), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked}, err
}

func (r *queryResolver) UnReviewed(ctx context.Context, password string, new bool) (*model.PostReview, error) {
	if password == "funia" {
		var postReview data.OutputReview
		var spam float32
		var err error
		if new {
			postReview, spam, err = post.GetOneNew()
		} else {
			postReview, spam, err = post.GetOneReported()
		}
		post := model.PostOut{
			ID:            postReview.Post.ID.String(),
			Text:          postReview.Post.Text,
			UserID:        postReview.Post.UserID,
			Shares:        len(postReview.Post.Shares),
			Views:         len(postReview.Post.Views),
			InitialReview: postReview.Post.InitialReview,
			Image:         postReview.Post.Image,
			CreationTime:  postReview.Post.CreationTime.String(),
			Blocked:       postReview.Post.Blocked,
		}
		return &model.PostReview{Post: &post, Left: postReview.Left, Spam: float64(spam)}, err
	}
	return &model.PostReview{Post: &model.PostOut{}, Left: 0}, errors.New("bad password")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var counter = 0
