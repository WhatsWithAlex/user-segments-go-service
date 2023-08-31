package postgresdb

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateSegmentTXArgs struct {
	SegmentSlug string  `json:"segment_slug"`
	AutoProb    float32 `json:"auto_prob"`
}

func (store *Store) CreateSegmentTX(ctx context.Context, arg CreateSegmentTXArgs) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var txerr error

		segmentID, txerr := q.CreateSegment(ctx, CreateSegmentParams{
			Slug:     arg.SegmentSlug,
			AutoProb: arg.AutoProb,
		})
		if txerr != nil {
			return txerr
		}

		// If segment auto probability specified then add users to it with given probability
		if arg.AutoProb > 0.0 {
			userIDs, txerr := q.GetUsers(ctx)
			if txerr != nil {
				return txerr
			}

			for _, userID := range userIDs {
				if rand.Float32() <= arg.AutoProb {
					txerr = q.AddUserSegment(ctx, AddUserSegmentParams{
						UserID:    userID,
						SegmentID: segmentID,
						RemoveAt:  pgtype.Timestamptz{},
					})
					if txerr != nil {
						return txerr
					}

					// Add records to operations history for automatically added user segments.
					txerr = q.CreateOperation(ctx, CreateOperationParams{
						OpType:      OpAdd,
						UserID:      userID,
						SegmentSlug: arg.SegmentSlug,
					})
					if txerr != nil {
						return txerr
					}
				}
			}
		}

		return txerr
	})

	return err
}

func (store *Store) DeleteSegmentTX(ctx context.Context, segmentSlug string) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var txerr error

		segmentID, txerr := q.DeleteSegmentBySlug(ctx, segmentSlug)
		if txerr != nil {
			return txerr
		}

		userIDs, txerr := q.GetUsersBySegmentID(ctx, segmentID)
		if txerr != nil {
			return txerr
		}

		// Add records to operations history for all users in deleted segment.
		for _, userID := range userIDs {
			txerr = q.CreateOperation(ctx, CreateOperationParams{
				OpType:      OpRemove,
				UserID:      userID,
				SegmentSlug: segmentSlug,
			})
			if txerr != nil {
				return txerr
			}
		}

		return txerr
	})

	return err
}
