package post

import (
	"context"
	"errors"
	"github.com/google/uuid"
	postv1 "github.com/sdoshi579/cloudbees/gen/post/v1"
	"github.com/sdoshi579/cloudbees/internal/entity"
	"github.com/sdoshi579/cloudbees/internal/service/post"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RPCImplementation struct {
	postv1.UnimplementedPostServiceServer
	service post.Service
	logger  *zap.Logger
}

func NewRPCImplementation(postService post.Service, logger *zap.Logger) *RPCImplementation {
	return &RPCImplementation{
		service: postService,
		logger:  logger,
	}
}

func (r *RPCImplementation) Create(ctx context.Context, request *postv1.CreateRequest) (*postv1.CreateResponse, error) {
	entityRequest := entity.CreatePostRequest{
		Title:       request.Title,
		Content:     request.Content,
		Author:      request.Author,
		PublishedOn: request.PublishedOn.AsTime(),
		Tags:        request.Tags,
	}

	resp, err := r.service.CreatePost(ctx, entityRequest)

	if err != nil {
		r.logger.Error("error in creating post", zap.Error(err), zap.Any("request", request))
		return &postv1.CreateResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &postv1.CreateResponse{
		Success:     true,
		Id:          resp.ID.String(),
		Title:       resp.Title,
		Content:     resp.Content,
		Author:      resp.Author,
		PublishedOn: timestamppb.New(resp.PublishedOn),
		Tags:        resp.Tags,
	}, nil
}
func (r *RPCImplementation) Get(ctx context.Context, request *postv1.GetRequest) (*postv1.GetResponse, error) {
	postID, err := uuid.Parse(request.Id)
	if err != nil {
		r.logger.Error("error in parsing post id", zap.Error(err), zap.Any("request", request))
		return nil, errors.New("invalid post id")
	}
	resp, err := r.service.GetPost(ctx, postID)

	if err != nil {
		r.logger.Error("error in fetching post", zap.Error(err), zap.Any("request", request))
		return &postv1.GetResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &postv1.GetResponse{
		Success:     true,
		Id:          resp.ID.String(),
		Title:       resp.Title,
		Content:     resp.Content,
		Author:      resp.Author,
		PublishedOn: timestamppb.New(resp.PublishedOn),
		Tags:        resp.Tags,
	}, nil
}
func (r *RPCImplementation) Update(ctx context.Context, request *postv1.UpdateRequest) (*postv1.UpdateResponse, error) {
	postID, err := uuid.Parse(request.Id)
	if err != nil {
		r.logger.Error("error in parsing post id", zap.Error(err), zap.Any("request", request))
		return nil, errors.New("invalid post id")
	}

	entityRequest := entity.UpdatePostRequest{
		Title:   request.Title,
		Content: request.Content,
		Author:  request.Author,
		Tags:    request.Tags,
	}

	resp, err := r.service.UpdatePost(ctx, postID, entityRequest)

	if err != nil {
		r.logger.Error("error in updating post", zap.Error(err), zap.Any("request", request))
		return &postv1.UpdateResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &postv1.UpdateResponse{
		Success:     true,
		Id:          resp.ID.String(),
		Title:       resp.Title,
		Content:     resp.Content,
		Author:      resp.Author,
		PublishedOn: timestamppb.New(resp.PublishedOn),
		Tags:        resp.Tags,
	}, nil
}
func (r *RPCImplementation) Delete(ctx context.Context, request *postv1.DeleteRequest) (*postv1.DeleteResponse, error) {
	postID, err := uuid.Parse(request.Id)
	if err != nil {
		r.logger.Error("error in parsing post id", zap.Error(err), zap.Any("request", request))
		return nil, errors.New("invalid post id")
	}
	resp, err := r.service.DeletePost(ctx, postID)

	if err != nil || !resp {
		r.logger.Error("error in fetching post", zap.Error(err), zap.Any("request", request))
		return &postv1.DeleteResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &postv1.DeleteResponse{
		Success: true,
	}, nil
}
