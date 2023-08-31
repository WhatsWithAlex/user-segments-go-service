package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
	"github.com/jackc/pgx/v5/pgtype"
)

type userSegmentRepository struct {
	database *postgresdb.Store
}

func NewUserSegmentRepository(db *postgresdb.Store) domain.UserSegmentRepository {
	return &userSegmentRepository{
		database: db,
	}
}

func (usr *userSegmentRepository) GetActiveUserSegments(ctx context.Context, userID int) (segments []string, err error) {
	return usr.database.GetActiveUserSegmentsTX(ctx, int32(userID))
}

func (usr *userSegmentRepository) UpdateUserSegments(ctx context.Context, userID int, segmentSlugsAdd, segmentSlugsRemove []string, removeAt sql.NullTime) error {
	var removeAtPg pgtype.Timestamptz

	if removeAt.Valid && removeAtPg.Scan(removeAt.Time) != nil {
		return errors.New("removeAt scan error")
	}

	err := usr.database.UpdateUserSegmentsTX(ctx, postgresdb.UpdateUserSegmentsTXArgs{
		UserID:             int32(userID),
		SegmentSlugsAdd:    segmentSlugsAdd,
		SegmentSlugsRemove: segmentSlugsRemove,
		RemoveAt:           removeAtPg,
	})
	return err
}

func (usr *userSegmentRepository) DeleteAllExpiredUserSegments(ctx context.Context) error {
	return usr.database.DeleteAllExpiredUserSegmentsTX(ctx)
}
