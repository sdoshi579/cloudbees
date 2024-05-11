package post

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sdoshi579/cloudbees/internal/entity"
	"github.com/sdoshi579/cloudbees/internal/repository/post"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=../../mockgen/service/post/post_service.go -source=./post_service.go Service
type Service interface {
	CreatePost(ctx context.Context, request entity.CreatePostRequest) (*entity.PostDetail, error)
	GetPost(ctx context.Context, id uuid.UUID) (*entity.PostDetail, error)
	UpdatePost(ctx context.Context, id uuid.UUID, request entity.UpdatePostRequest) (*entity.PostDetail, error)
	DeletePost(ctx context.Context, id uuid.UUID) (bool, error)
}

type serviceImplementation struct {
	repository post.Repository
	logger     *zap.Logger
}

type ServiceConfiguration func(r *serviceImplementation)

func NewService(configs ...ServiceConfiguration) Service {
	r := serviceImplementation{}
	for _, config := range configs {
		config(&r)
	}
	return &r
}

func WithLogger(logger *zap.Logger) ServiceConfiguration {
	return func(r *serviceImplementation) {
		r.logger = logger
	}
}

func WithRepository(repository post.Repository) ServiceConfiguration {
	return func(r *serviceImplementation) {
		r.repository = repository
	}
}

func (s *serviceImplementation) CreatePost(ctx context.Context, request entity.CreatePostRequest) (*entity.PostDetail, error) {
	return s.repository.CreatePost(ctx, request)
}

func (s *serviceImplementation) GetPost(ctx context.Context, id uuid.UUID) (*entity.PostDetail, error) {
	return s.repository.GetPost(ctx, id)
}

func (s *serviceImplementation) UpdatePost(ctx context.Context, id uuid.UUID,
	request entity.UpdatePostRequest) (*entity.PostDetail, error) {

	_, err := s.repository.GetPost(ctx, id)
	if err != nil {
		s.logger.Error("invalid post id for update", zap.Error(err), zap.Any("postID", id))
		return nil, errors.New("post is not available or is deleted")
	}
	return s.repository.UpdatePost(ctx, id, request)
}

func (s *serviceImplementation) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.DeletePost(ctx, id)
}
