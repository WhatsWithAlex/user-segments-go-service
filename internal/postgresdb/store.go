package postgresdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	dbConn *pgxpool.Pool
}

func NewStore(dbConn *pgxpool.Pool) *Store {
	return &Store{
		dbConn:  dbConn,
		Queries: New(dbConn),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Join(err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
