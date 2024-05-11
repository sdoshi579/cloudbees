package post

import (
	"context"
	"github.com/google/uuid"
	"github.com/sdoshi579/cloudbees/internal/entity"
	"github.com/sdoshi579/cloudbees/internal/repository/ent"
	"github.com/sdoshi579/cloudbees/internal/repository/ent/post"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../../mockgen/repository/post/post_repository.go -source=./post_repository.go Repository
type Repository interface {
	CreatePost(ctx context.Context, request entity.CreatePostRequest) (*entity.PostDetail, error)
	GetPost(ctx context.Context, id uuid.UUID) (*entity.PostDetail, error)
	UpdatePost(ctx context.Context, id uuid.UUID, request entity.UpdatePostRequest) (*entity.PostDetail, error)
	DeletePost(ctx context.Context, id uuid.UUID) (bool, error)
}

type repositoryImplementation struct {
	entClient *ent.Client
	logger    *zap.Logger
}

type RepoConfiguration func(r *repositoryImplementation)

func NewRepository(configs ...RepoConfiguration) Repository {
	r := repositoryImplementation{}
	for _, config := range configs {
		config(&r)
	}
	return &r
}

func WithLogger(logger *zap.Logger) RepoConfiguration {
	return func(r *repositoryImplementation) {
		r.logger = logger
	}
}

func WithEntClient(client *ent.Client) RepoConfiguration {
	return func(r *repositoryImplementation) {
		r.entClient = client
	}
}

func (r *repositoryImplementation) CreatePost(ctx context.Context,
	request entity.CreatePostRequest) (*entity.PostDetail, error) {

	resp, err := r.entClient.Post.Create().SetTitle(request.Title).SetContent(request.Content).
		SetAuthor(request.Author).
		SetPublishedOn(request.PublishedOn).SetTags(request.Tags).Save(ctx)

	if err != nil {
		r.logger.Error("error in saving post", zap.Error(err), zap.Any("request", request))
		return nil, err
	}

	return decoratePostEntity(*resp), nil
}

func (r *repositoryImplementation) GetPost(ctx context.Context, id uuid.UUID) (*entity.PostDetail, error) {
	resp, err := r.entClient.Post.Query().Where(post.ID(id), post.IsDeleted(false)).Only(ctx)

	if err != nil {
		r.logger.Error("error in fetching post", zap.Error(err), zap.Any("postID", id))
		return nil, err
	}
	return decoratePostEntity(*resp), nil
}

func (r *repositoryImplementation) UpdatePost(ctx context.Context, id uuid.UUID,
	request entity.UpdatePostRequest) (*entity.PostDetail, error) {
	query := r.entClient.Post.UpdateOneID(id)

	if request.Title != nil {
		query.SetTitle(*request.Title)
	}
	if request.Content != nil {
		query.SetContent(*request.Content)
	}
	if request.Author != nil {
		query.SetAuthor(*request.Author)
	}
	if len(request.Tags) != 0 {
		query.SetTags(request.Tags)
	}

	resp, err := query.Save(ctx)
	if err != nil {
		r.logger.Error("error in updating post", zap.Error(err), zap.Any("request", request),
			zap.Any("postID", id))
		return nil, err
	}

	return decoratePostEntity(*resp), nil
}

func (r *repositoryImplementation) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	err := r.entClient.Post.UpdateOneID(id).SetIsDeleted(true).Exec(ctx)

	if err != nil {
		r.logger.Error("error in deleting post", zap.Error(err), zap.Any("postID", id))
		return false, err
	}
	return true, nil
}

func decoratePostEntity(postEnt ent.Post) *entity.PostDetail {
	return &entity.PostDetail{
		ID:          postEnt.ID,
		Author:      postEnt.Author,
		Tags:        postEnt.Tags,
		Title:       postEnt.Title,
		Content:     postEnt.Content,
		PublishedOn: postEnt.PublishedOn,
	}
}
