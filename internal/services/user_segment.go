package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
)

type userSegmentService struct {
	userSegmentRepository domain.UserSegmentRepository
	contextTimeout        time.Duration
}

func NewUserSegmentService(userSegmentRepository domain.UserSegmentRepository, timeout time.Duration) domain.UserSegmentService {
	return &userSegmentService{
		userSegmentRepository: userSegmentRepository,
		contextTimeout:        timeout,
	}
}

func (uss *userSegmentService) GetActiveUserSegments(ctx context.Context, userID int) (segments []string, err error) {
	c, cancel := context.WithTimeout(ctx, uss.contextTimeout)
	defer cancel()
	return uss.userSegmentRepository.GetActiveUserSegments(c, userID)
}

func (uss *userSegmentService) UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error {
	c, cancel := context.WithTimeout(ctx, uss.contextTimeout)
	defer cancel()
	return uss.userSegmentRepository.UpdateUserSegments(c, userID, segmentSlugsAdd, segmentSlugsRemove, removeAt)
}
