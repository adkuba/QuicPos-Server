package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"QuicPos/graph/generated"
	"QuicPos/graph/model"
	"QuicPos/internal/post"
	"QuicPos/internal/storage"
	"context"
	"fmt"
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

func (r *queryResolver) Post(ctx context.Context, userID string) (*model.Post, error) {
	post := post.GetOne()
	return &model.Post{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Reports: post.Reports, Shares: post.Shares, Views: post.Views, InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String()}, nil
}

func (r *queryResolver) OpinionPost(ctx context.Context) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CreateUser(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
