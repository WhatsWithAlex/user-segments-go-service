package repository

import (
	"context"
	"database/sql"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
	"github.com/jackc/pgx/v5/pgtype"
)

type operationRepository struct {
	database postgresdb.Store
}

func NewOperationRepository(db postgresdb.Store) domain.OperationRepository {
	return &operationRepository{
		database: db,
	}
}

func (or *operationRepository) CreateOperation(ctx context.Context, userID int, opType domain.Op, segmentSlug string, doneAt sql.NullTime) error {
	var doneAtPg *pgtype.Timestamptz
	err := doneAtPg.Scan(doneAt.Time)
	if err != nil {
		return err
	}
	return or.database.CreateOperation(ctx, postgresdb.CreateOperationParams{
		UserID:      int32(userID),
		OpType:      postgresdb.Op(opType),
		SegmentSlug: segmentSlug,
		DoneAt:      *doneAtPg,
	})
}

func (or *operationRepository) GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) (operations []domain.Operation, err error) {
	var fromTSPg, toTSPg *pgtype.Timestamptz
	err = fromTSPg.Scan(fromTS.Time)
	if err != nil {
		return
	}
	err = toTSPg.Scan(toTS.Time)
	if err != nil {
		return
	}

	operationsPg, err := or.database.GetOperationsByUserID(ctx, postgresdb.GetOperationsByUserIDParams{
		UserID: int32(userID),
		FromTS: *fromTSPg,
		ToTS:   *toTSPg,
	})
	if err != nil {
		return
	}

	for _, operationPg := range operationsPg {
		operations = append(operations, domain.Operation{
			OpType:      domain.Op(operationPg.OpType),
			SegmentSlug: operationPg.SegmentSlug,
			DoneAt:      operationPg.DoneAt.Time,
		})
	}
	return
}
