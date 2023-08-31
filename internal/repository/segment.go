package repository

import (
	"context"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
)

type segmentRepository struct {
	database postgresdb.Store
}

func NewSegmentRepository(db postgresdb.Store) domain.SegmentRepository {
	return &segmentRepository{
		database: db,
	}
}

func (sr *segmentRepository) CreateSegment(ctx context.Context, segmentSlug string, autoProb float32) error {
	err := sr.database.CreateSegmentTX(ctx, postgresdb.CreateSegmentTXArgs{
		SegmentSlug: segmentSlug,
		AutoProb:    autoProb,
	})
	return err
}

func (sr *segmentRepository) DeleteSegment(ctx context.Context, segmentSlug string) error {
	err := sr.database.DeleteSegmentTX(ctx, segmentSlug)
	return err
}
