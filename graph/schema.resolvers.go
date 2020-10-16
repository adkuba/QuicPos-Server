package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"QuicPos/graph/generated"
	"QuicPos/graph/model"
	"QuicPos/internal/ip"
	"QuicPos/internal/post"
	"QuicPos/internal/storage"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.PostOut, error) {
	var post post.Post
	post.Text = input.Text
	post.UserID = input.UserID
	post.CreationTime = time.Now()
	post.Image = storage.UploadFile(input.Image)
	post.InitialReview = false
	post.Reports = nil
	post.Shares = nil
	post.Views = nil
	post.Blocked = false
	postID, err := post.Save()
	return &model.PostOut{ID: postID, Text: post.Text, UserID: post.UserID, Shares: len(post.Shares), Views: len(post.Views), CreationTime: post.CreationTime.String(), InitialReview: post.InitialReview, Image: post.Image, Blocked: post.Blocked}, err
}

func (r *mutationResolver) Review(ctx context.Context, input model.Review) (bool, error) {
	if input.Password == "funia" {
		result, err := post.ReviewAction(input.New, input.PostID, input.Delete)
		return result, err
	}
	return false, errors.New("bad password")
}

func (r *mutationResolver) Share(ctx context.Context, input string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Report(ctx context.Context, input string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) View(ctx context.Context, input model.NewView) (bool, error) {
	//tu mam dokladne dane urzadzenia
	log.Println(ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).UserAgent)
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Post(ctx context.Context, userID string, normalMode bool) (*model.PostOut, error) {
	//userID and normalMode to be used
	post, err := post.GetOne()
	return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Shares: len(post.Shares), Views: len(post.Views), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked}, err
}

func (r *queryResolver) CreateUser(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

func (r *queryResolver) ViewerPost(ctx context.Context, id string) (*model.PostOut, error) {
	post, err := post.GetByID(id)
	return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.UserID, Shares: len(post.Shares), Views: len(post.Views), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked}, err
}

func (r *queryResolver) UnReviewed(ctx context.Context, password string, new bool) (*model.PostReview, error) {
	if password == "funia" {
		var postReview post.OutputReview
		var err error
		if new {
			postReview, err = post.GetOneNew()
		} else {
			postReview, err = post.GetOneReported()
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
		return &model.PostReview{Post: &post, Left: postReview.Left}, err
	}
	return &model.PostReview{Post: &model.PostOut{}, Left: 0}, errors.New("bad password")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
