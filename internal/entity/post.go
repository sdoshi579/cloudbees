package entity

import (
	"github.com/google/uuid"
	"time"
)

type CreatePostRequest struct {
	Title       string
	Content     string
	Author      string
	PublishedOn time.Time
	Tags        []string
}

type UpdatePostRequest struct {
	Title   *string
	Content *string
	Author  *string
	Tags    []string
}

type PostDetail struct {
	ID          uuid.UUID
	Title       string
	Content     string
	Author      string
	PublishedOn time.Time
	Tags        []string
}
