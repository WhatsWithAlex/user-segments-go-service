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

type UserSegmentService interface {
	UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
}

type UserSegmentRepository interface {
	UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error
	GetActiveUserSegments(ctx context.Context, userID int) ([]string, error)
	DeleteAllExpiredUserSegments(ctx context.Context) error
}
