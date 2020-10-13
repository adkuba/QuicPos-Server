package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"QuicPos/graph/generated"
	"QuicPos/graph/model"
	"QuicPos/internal/post"
	"QuicPos/internal/storage"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	var post post.Post
	post.Text = input.Text
	post.UserID = input.UserID
	post.CreationTime = time.Now()
	post.Image = storage.UploadFile(input.Image)
	post.InitialReview = false
	post.Reports = nil
	post.Shares = 0
	post.Views = nil
	postID := post.Save()
	return &model.Post{ID: postID, Text: post.Text, UserID: post.UserID, Reports: post.Reports, Shares: post.Shares, Views: post.Views, CreationTime: post.CreationTime.String(), InitialReview: post.InitialReview, Image: post.Image}, nil
}

func (r *mutationResolver) Review(ctx context.Context, input model.Review) (bool, error) {
	if input.Password == "funia" {
		result, err := post.ReviewAction(input.New, input.PostID, input.Delete)
		return result, err
	}
	return false, errors.New("bad password")
}

func (r *queryResolver) Post(ctx context.Context, userID string, normalMode bool) (*model.Post, error) {
	//userID and normalMode to be used
	post := post.GetOne()
	return &model.Post{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Reports: post.Reports, Shares: post.Shares, Views: post.Views, InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String()}, nil
}

func (r *queryResolver) CreateUser(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

func (r *queryResolver) ViewerPost(ctx context.Context, id string) (*model.Post, error) {
	post := post.GetByID(id)
	return &model.Post{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Reports: post.Reports, Shares: post.Shares, Views: post.Views, InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String()}, nil
}

func (r *queryResolver) UnReviewed(ctx context.Context, password string, new bool) (*model.PostReview, error) {
	if password == "funia" {
		var postReview post.OutputReview
		if new {
			postReview = post.GetOneNew()
		} else {
			postReview = post.GetOneReported()
		}
		post := model.Post{
			ID:            postReview.Post.ID.String(),
			Text:          postReview.Post.Text,
			UserID:        postReview.Post.UserID,
			Reports:       postReview.Post.Reports,
			Shares:        postReview.Post.Shares,
			Views:         postReview.Post.Views,
			InitialReview: postReview.Post.InitialReview,
			Image:         postReview.Post.Image,
			CreationTime:  postReview.Post.CreationTime.String(),
		}
		return &model.PostReview{Post: &post, Left: postReview.Left}, nil
	}
	return &model.PostReview{Post: &model.Post{}, Left: 0}, errors.New("bad password")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
