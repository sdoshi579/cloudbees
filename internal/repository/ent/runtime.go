// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/sdoshi579/cloudbees/internal/repository/ent/post"
	"github.com/sdoshi579/cloudbees/internal/repository/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	postFields := schema.Post{}.Fields()
	_ = postFields
	// postDescIsDeleted is the schema descriptor for is_deleted field.
	postDescIsDeleted := postFields[6].Descriptor()
	// post.DefaultIsDeleted holds the default value on creation for the is_deleted field.
	post.DefaultIsDeleted = postDescIsDeleted.Default.(bool)
	// postDescCreatedAt is the schema descriptor for created_at field.
	postDescCreatedAt := postFields[7].Descriptor()
	// post.DefaultCreatedAt holds the default value on creation for the created_at field.
	post.DefaultCreatedAt = postDescCreatedAt.Default.(func() time.Time)
	// postDescUpdatedAt is the schema descriptor for updated_at field.
	postDescUpdatedAt := postFields[8].Descriptor()
	// post.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	post.DefaultUpdatedAt = postDescUpdatedAt.Default.(func() time.Time)
	// post.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	post.UpdateDefaultUpdatedAt = postDescUpdatedAt.UpdateDefault.(func() time.Time)
	// postDescID is the schema descriptor for id field.
	postDescID := postFields[0].Descriptor()
	// post.DefaultID holds the default value on creation for the id field.
	post.DefaultID = postDescID.Default.(func() uuid.UUID)
}