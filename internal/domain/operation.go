package domain

import (
	"context"
	"database/sql"
	"time"
)

type Op string

const (
	OpRemove Op = "remove"
	OpAdd    Op = "add"
)

type Operation struct {
	SegmentSlug string
	DoneAt      time.Time
	OpType      Op
}

type OperationService interface {
	GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) ([]Operation, error)
}

type OperationRepository interface {
	CreateOperation(ctx context.Context, userID int, opType Op, segmentSlug string, doneAt sql.NullTime) error
	GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) ([]Operation, error)
}
