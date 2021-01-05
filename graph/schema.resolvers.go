package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"QuicPos/graph/generated"
	"QuicPos/graph/model"
	"QuicPos/internal/data"
	"QuicPos/internal/ip"
	"QuicPos/internal/post"
	"QuicPos/internal/stats"
	"QuicPos/internal/storage"
	"QuicPos/internal/stripe"
	"QuicPos/internal/tensorflow"
	"QuicPos/internal/user"
	"context"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost, password string) (*model.PostOut, error) {
	if password == data.Pass {
		var postO data.Post
		var err error
		postO.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		postO.Text = input.Text
		userO, err := user.GetUser(input.UserID)
		if err != nil {
			return &model.PostOut{}, err
		}
		postO.User = userO
		postO.CreationTime = time.Now()
		postO.Image, err = storage.UploadFile(input.Image)
		if err != nil {
			return &model.PostOut{}, err
		}
		postO.InitialReview = false
		postO.Reports = nil
		postO.Shares = nil
		postO.Views = nil
		postO.Blocked = false
		postO.OutsideViews = nil
		postO.Money = 0
		postID, err := post.Save(postO)
		return &model.PostOut{ID: postID, Text: postO.Text, UserID: postO.User.UUID, Shares: len(postO.Shares), Views: len(postO.Views), CreationTime: postO.CreationTime.String(), InitialReview: postO.InitialReview, Image: postO.Image, Blocked: postO.Blocked, Money: postO.Money}, err
	}
	return &model.PostOut{}, errors.New("bad key")
}

func (r *mutationResolver) Review(ctx context.Context, input model.Review) (bool, error) {
	if input.Password == data.AdminPass {
		result, err := post.ReviewAction(input.New, input.PostID, input.Delete)
		return result, err
	}
	return false, errors.New("bad key")
}

func (r *mutationResolver) Share(ctx context.Context, input model.NewReportShare, password string) (bool, error) {
	if password == data.Pass {
		result, err := post.Share(input)
		return result, err
	}
	return false, errors.New("bad key")
}

func (r *mutationResolver) Report(ctx context.Context, input model.NewReportShare) (bool, error) {
	result, err := post.Report(input)
	return result, err
}

func (r *mutationResolver) View(ctx context.Context, input model.NewView, password string) (bool, error) {
	if password == data.Pass {
		//tu mam dokladne dane urzadzenia
		status, err := post.AddView(input, ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).IP)
		return status, err
	}
	return false, errors.New("bad key")
}

func (r *mutationResolver) Learning(ctx context.Context, input model.Learning, password string) (bool, error) {
	if password == data.Pass {
		err := stats.UpdateNets(input.Recommender, input.Detector)
		if err != nil {
			return false, err
		}
		err = tensorflow.InitModels()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("bad key")
}

func (r *mutationResolver) Payment(ctx context.Context, input model.Payment) (bool, error) {
	result, err := post.AddMoney(input)
	return result, err
}

func (r *mutationResolver) RemovePost(ctx context.Context, input model.Remove, password string) (bool, error) {
	if password == data.Pass {
		err := post.Remove(input.PostID, input.UserID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("bad password")
}

func (r *queryResolver) Post(ctx context.Context, userID string, normalMode bool, password string, ad bool) (*model.PostOut, error) {
	if password == data.Pass {
		if normalMode || ad {
			post, err := post.GetOne(userID, ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).IP, ad)
			return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.User.UUID, Shares: len(post.Shares), Views: len(post.Views) + len(post.OutsideViews), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked, Money: post.Money}, err
		}
		post, err := post.GetOneRandom()
		return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.User.UUID, Shares: len(post.Shares), Views: len(post.Views) + len(post.OutsideViews), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked, Money: post.Money}, err
	}
	return &model.PostOut{}, errors.New("bad key")
}

func (r *queryResolver) CreateUser(ctx context.Context, password string) (string, error) {
	if password == data.Pass {
		id, err := user.Create(ctx.Value(ip.IPCtxKey).(*ip.DeviceDetails).IP)
		return id, err
	}
	return "", errors.New("bad key")
}

func (r *queryResolver) ViewerPost(ctx context.Context, id string) (*model.PostOut, error) {
	post, err := post.GetByID(id, true)
	return &model.PostOut{ID: post.ID.String(), Text: post.Text, UserID: post.User.UUID, Shares: len(post.Shares), Views: len(post.Views) + len(post.OutsideViews), InitialReview: post.InitialReview, Image: post.Image, CreationTime: post.CreationTime.String(), Blocked: post.Blocked, Money: post.Money}, err
}

func (r *queryResolver) UnReviewed(ctx context.Context, password string, new bool) (*model.PostReview, error) {
	if password == data.AdminPass {
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
			UserID:        postReview.Post.User.UUID,
			Shares:        len(postReview.Post.Shares),
			Views:         len(postReview.Post.Views),
			InitialReview: postReview.Post.InitialReview,
			Image:         postReview.Post.Image,
			CreationTime:  postReview.Post.CreationTime.String(),
			Blocked:       postReview.Post.Blocked,
			Money:         postReview.Post.Money,
		}
		return &model.PostReview{Post: &post, Left: postReview.Left, Spam: float64(spam)}, err
	}
	return &model.PostReview{Post: &model.PostOut{}, Left: 0}, errors.New("bad password")
}

func (r *queryResolver) StorageIntegrity(ctx context.Context, password string) (string, error) {
	if password == data.Pass {
		deleted, err := storage.RemoveParentless()
		if err != nil {
			return "Error!", err
		}
		return strconv.Itoa(deleted) + " images deleted", nil
	}
	return "Bad password!", nil
}

func (r *queryResolver) GetStats(ctx context.Context, id string) (*model.Stats, error) {
	post, err := post.GetByID(id, false)
	if err != nil {
		return &model.Stats{}, err
	}

	var views []*model.View
	for _, view := range append(post.Views, post.OutsideViews...) {
		view := &model.View{
			Localization: view.Localization,
			Date:         view.Date.String(),
		}
		views = append(views, view)
	}

	return &model.Stats{Text: post.Text, Userid: post.User.UUID, Views: views, Money: float64(post.Money) / 100}, nil
}

func (r *queryResolver) GetStripeClient(ctx context.Context, amount float64) (string, error) {
	client, err := stripe.CreatePaymentIntent(amount)
	return client, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
