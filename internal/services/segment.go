package services

import (
	"context"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
)

type segmentService struct {
	segmentRepository domain.SegmentRepository
	contextTimeout    time.Duration
}

func NewSegmentService(segmentRepository domain.SegmentRepository, timeout time.Duration) domain.SegmentService {
	return &segmentService{
		segmentRepository: segmentRepository,
		contextTimeout:    timeout,
	}
}

func (ss *segmentService) CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error {
	c, cancel := context.WithTimeout(ctx, ss.contextTimeout)
	defer cancel()
	return ss.segmentRepository.CreateSegment(c, segmentSlug, autoProb)
}

func (ss *segmentService) DeleteSegment(ctx context.Context, segmentSlug string) error {
	c, cancel := context.WithTimeout(ctx, ss.contextTimeout)
	defer cancel()
	return ss.segmentRepository.DeleteSegment(c, segmentSlug)
}
