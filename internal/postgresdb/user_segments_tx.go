package postgresdb

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUserSegmentsTXArgs struct {
	UserID             int32              `json:"user_id"`
	SegmentSlugsRemove []string           `json:"segment_slugs_remove"`
	SegmentSlugsAdd    []string           `json:"segment_slugs_add"`
	RemoveAt           pgtype.Timestamptz `json:"remove_at"`
}

func (store *Store) UpdateUserSegmentsTX(ctx context.Context, arg UpdateUserSegmentsTXArgs) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var txerr error

		cnt, txerr := q.CountUserSegments(ctx, arg.UserID)
		if txerr != nil {
			return txerr
		}

		for _, slug := range arg.SegmentSlugsAdd {
			txerr := q.AddUserSegmentBySlug(ctx, AddUserSegmentBySlugParams{
				UserID:   arg.UserID,
				Slug:     slug,
				RemoveAt: arg.RemoveAt,
			})
			if txerr != nil {
				return txerr
			}

			// Add records to operations history for added user segments.
			txerr = q.CreateOperation(ctx, CreateOperationParams{
				OpType:      OpAdd,
				UserID:      arg.UserID,
				SegmentSlug: slug,
			})
			if txerr != nil {
				return txerr
			}
		}

		txerr = q.DeleteUserSegmentsBySlugs(ctx, DeleteUserSegmentsBySlugsParams{
			UserID:       arg.UserID,
			SegmentSlugs: arg.SegmentSlugsRemove,
		})
		if txerr != nil {
			return txerr
		}

		for _, slug := range arg.SegmentSlugsRemove {
			// Add records to operations history for removed user segments.
			txerr = q.CreateOperation(ctx, CreateOperationParams{
				OpType:      OpRemove,
				UserID:      arg.UserID,
				SegmentSlug: slug,
			})
			if txerr != nil {
				return txerr
			}
		}

		// If user is new then automatically add him in segments with given probabilities.
		if cnt == 0 {
			autoSegments, txerr := q.GetAutoSegments(ctx)
			if txerr != nil {
				return txerr
			}
			for _, segment := range autoSegments {
				if rand.Float32() <= segment.AutoProb {
					txerr := q.AddUserSegmentBySlug(ctx, AddUserSegmentBySlugParams{
						UserID:   arg.UserID,
						Slug:     segment.Slug,
						RemoveAt: pgtype.Timestamptz{},
					})
					if txerr != nil {
						return txerr
					}

					// Add records to operations history for automatically added user segments.
					txerr = q.CreateOperation(ctx, CreateOperationParams{
						OpType:      OpAdd,
						UserID:      arg.UserID,
						SegmentSlug: segment.Slug,
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

func (store *Store) GetActiveUserSegmentsTX(ctx context.Context, userID int32) (segmentSlugs []string, err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		var txerr error

		segmentSlugs, txerr = q.GetActiveUserSegments(ctx, userID)
		if txerr != nil {
			return txerr
		}

		return checkExpiredAndUpdate(ctx, q)
	})

	return segmentSlugs, err
}

func (store *Store) DeleteAllExpiredUserSegmentsTX(ctx context.Context) error {
	err := store.execTx(ctx, func(q *Queries) error {
		return checkExpiredAndUpdate(ctx, q)
	})

	return err
}

func checkExpiredAndUpdate(ctx context.Context, q *Queries) error {
	expiredUserSegments, err := q.DeleteAllExpiredUserSegments(ctx)
	if err != nil {
		return err
	}

	// Add records to operations history for expired user segments.
	for _, expiredUserSegment := range expiredUserSegments {
		err = q.CreateOperationWithTS(ctx, CreateOperationWithTSParams{
			OpType:      OpRemove,
			UserID:      expiredUserSegment.UserID,
			SegmentSlug: expiredUserSegment.Slug,
			DoneAt:      expiredUserSegment.RemoveAt,
		})
		if err != nil {
			return err
		}
	}

	return err
}
