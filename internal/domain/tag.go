package domain

import (
	"time"
	"context"
)

type TagStatus string

const (
	TagStatusPending  TagStatus = "pending"
	TagStatusApproved TagStatus = "approved"
	TagStatusRejected TagStatus = "rejected"
)

type Tag struct {
	ID        string
	Name      string
	Status    TagStatus
	CreatedBy string
	CreatedAt time.Time
}

//=============================================================================//
//                          Tags Interface                                     //
//=============================================================================//
type TagRepository interface {
	Create(ctx context.Context, tag *Tag) (*Tag, error)
	Update(ctx context.Context, tag *Tag) (*Tag, error)
	GetByID(ctx context.Context, id string) (*Tag, error)
	GetByName(ctx context.Context, name string) (*Tag, error)
	List(ctx context.Context, filter TagFilter) ([]Tag, error)
	Delete(ctx context.Context, id string) error
}

type TagUsecase interface {
	CreateTag(ctx context.Context, userID string, name string) (*Tag, error)
	ApproveTag(ctx context.Context, tagID string) (*Tag, error)
	RejectTag(ctx context.Context, tagID string) (*Tag, error)
	ListTags(ctx context.Context, status TagStatus) ([]Tag, error)
	DeleteTag(ctx context.Context, tagID string) error
	IsTagApproved(name string) bool
	ValidateTags(tags []string) error
}

type TagFilter struct {
	Status TagStatus
}