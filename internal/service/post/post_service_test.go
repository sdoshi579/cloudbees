package post

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sdoshi579/cloudbees/internal/entity"
	mockpostrepository "github.com/sdoshi579/cloudbees/internal/mockgen/repository/post"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func Test_serviceImplementation_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockpostrepository.NewMockRepository(ctrl)

	mockRepo.EXPECT().CreatePost(gomock.Any(), entity.CreatePostRequest{
		Title: "success post",
	}).MaxTimes(1).Return(&entity.PostDetail{
		Title: "success post",
	}, nil)

	mockRepo.EXPECT().CreatePost(gomock.Any(), entity.CreatePostRequest{
		Title: "failed post",
	}).MaxTimes(1).Return(nil, errors.New("error in creating post"))
	type args struct {
		ctx     context.Context
		request entity.CreatePostRequest
	}
	tests := []struct {
		name string
		args args
		want *entity.PostDetail
		err  error
	}{
		{
			name: "success in creating post",
			args: args{
				ctx: context.Background(),
				request: entity.CreatePostRequest{
					Title: "success post",
				},
			},
			want: &entity.PostDetail{
				Title: "success post",
			},
			err: nil,
		},
		{
			name: "error in creating post",
			args: args{
				ctx: context.Background(),
				request: entity.CreatePostRequest{
					Title: "failed post",
				},
			},
			want: nil,
			err:  errors.New("error in creating post"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serviceImplementation{
				repository: mockRepo,
				logger:     zap.NewExample(),
			}
			got, err := s.CreatePost(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePost() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("CreatePost() error got = %v, want %v", err, tt.err)
			}
		})
	}
}

func Test_serviceImplementation_GetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockpostrepository.NewMockRepository(ctrl)

	successPostID := uuid.New()
	failPostID := uuid.New()
	mockRepo.EXPECT().GetPost(gomock.Any(), successPostID).MaxTimes(1).Return(&entity.PostDetail{
		Title: "success post",
	}, nil)

	mockRepo.EXPECT().GetPost(gomock.Any(), failPostID).MaxTimes(1).Return(nil, errors.New("error in fetching post"))
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want *entity.PostDetail
		err  error
	}{
		{
			name: "success in fetching post",
			args: args{
				ctx: context.Background(),
				id:  successPostID,
			},
			want: &entity.PostDetail{
				Title: "success post",
			},
			err: nil,
		},
		{
			name: "error in fetching post",
			args: args{
				ctx: context.Background(),
				id:  failPostID,
			},
			want: nil,
			err:  errors.New("error in fetching post"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serviceImplementation{
				repository: mockRepo,
				logger:     zap.NewExample(),
			}
			got, err := s.GetPost(tt.args.ctx, tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPost() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetPost() error got = %v, want %v", err, tt.err)
			}
		})
	}
}

func Test_serviceImplementation_UpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockpostrepository.NewMockRepository(ctrl)

	successPostID := uuid.New()
	failPostID := uuid.New()
	postNotFoundPostID := uuid.New()
	mockRepo.EXPECT().UpdatePost(gomock.Any(), successPostID, gomock.Any()).MaxTimes(1).Return(&entity.PostDetail{
		Title: "success post",
	}, nil)
	mockRepo.EXPECT().UpdatePost(gomock.Any(), failPostID, gomock.Any()).
		MaxTimes(1).Return(nil, errors.New("error in updating post"))

	mockRepo.EXPECT().GetPost(gomock.Any(), postNotFoundPostID).MaxTimes(1).
		Return(nil, errors.New("post not found or is deleted"))

	mockRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).MaxTimes(2).
		Return(&entity.PostDetail{}, nil)

	type args struct {
		ctx     context.Context
		id      uuid.UUID
		request entity.UpdatePostRequest
	}
	tests := []struct {
		name string
		args args
		want *entity.PostDetail
		err  error
	}{
		{
			name: "success in updating post",
			args: args{
				ctx: context.Background(),
				id:  successPostID,
			},
			want: &entity.PostDetail{
				Title: "success post",
			},
			err: nil,
		},
		{
			name: "error in updating post",
			args: args{
				ctx: context.Background(),
				id:  failPostID,
			},
			want: nil,
			err:  errors.New("error in updating post"),
		},
		{
			name: "post is not available",
			args: args{
				ctx: context.Background(),
				id:  postNotFoundPostID,
			},
			want: nil,
			err:  errors.New("post is not available or is deleted"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serviceImplementation{
				repository: mockRepo,
				logger:     zap.NewExample(),
			}
			got, err := s.UpdatePost(tt.args.ctx, tt.args.id, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdatePost() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("UpdatePost() error got = %v, want %v", err, tt.err)
			}
		})
	}
}

func Test_serviceImplementation_DeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockpostrepository.NewMockRepository(ctrl)

	successPostID := uuid.New()
	failPostID := uuid.New()
	mockRepo.EXPECT().DeletePost(gomock.Any(), successPostID).MaxTimes(1).Return(true, nil)

	mockRepo.EXPECT().DeletePost(gomock.Any(), failPostID).MaxTimes(1).Return(false, errors.New("error in deleting post"))

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want bool
		err  error
	}{
		{
			name: "success in deleting post",
			args: args{
				ctx: context.Background(),
				id:  successPostID,
			},
			want: true,
			err:  nil,
		},
		{
			name: "error in deleting post",
			args: args{
				ctx: context.Background(),
				id:  failPostID,
			},
			want: false,
			err:  errors.New("error in deleting post"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serviceImplementation{
				repository: mockRepo,
				logger:     zap.NewExample(),
			}
			got, err := s.DeletePost(tt.args.ctx, tt.args.id)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("DeletePost() error got = %v, want %v", err, tt.err)
			}
			if got != tt.want {
				t.Errorf("DeletePost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
