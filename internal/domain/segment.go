package domain

import "context"

type Segment struct {
	ID       int     `json:"id"`
	Slug     string  `json:"slug"`
	AutoProb float32 `json:"auto_prob"`
}

type CreateSegmentRequest struct {
	Slug        string  `json:"slug"                  binding:"required"`
	Probability float32 `json:"probability,omitempty" binding:"omitempty"`
}

type DeleteSegmentRequest struct {
	Slug string `form:"slug" binding:"required"`
}

type SegmentService interface {
	CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error
	DeleteSegment(ctx context.Context, segmentSlug string) error
}

type SegmentRepository interface {
	CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error
	DeleteSegment(ctx context.Context, segmentSlug string) error
}
