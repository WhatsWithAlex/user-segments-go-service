package domain

import (
	"context"
	"database/sql"
	"time"
)

type UserSegment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	SegmentID int       `json:"segment_id"`
	RemoveAt  time.Time `json:"remove_at"`
}

type GetUserSegmentsRequest struct {
	UserID int `form:"user_id" binding:"required"`
}

type GetUserSegmentsResponse struct {
	Segments []string `json:"segments"`
}

type UpdateUserSegmentsRequest struct {
	UserID         int       `json:"user_id"             binding:"required,gte=0"`
	AddSegments    []string  `json:"add_segments"        binding:"required_without=RemoveSegments"`
	RemoveSegments []string  `json:"remove_segments"     binding:"required_without=AddSegments"`
	RemoveAt       *time.Time `json:"remove_at,omitempty" binding:"omitempty"`
}

type UserSegmentService interface {
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
	UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error
}

type UserSegmentRepository interface {
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
	UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error
	DeleteAllExpiredUserSegments(ctx context.Context) error
}
