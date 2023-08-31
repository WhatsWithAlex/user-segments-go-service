package domain

import "context"

type Segment struct {
	ID       int     `json:"id"`
	Slug     string  `json:"slug"`
	AutoProb float32 `json:"auto_prob"`
}

type SegmentService interface {
	CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error
	DeleteSegment(ctx context.Context, segmentSlug string) error
}

type SegmentRepository interface {
	CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error
	DeleteSegment(ctx context.Context, segmentSlug string) error
}
