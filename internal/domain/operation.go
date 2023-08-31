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

type GetOperationsRequest struct {
	UserID int `form:"user_id" binding:"required"`
	Month  int `form:"month"   binding:"required,gte=1,lte=12"`
	Year   int `form:"year"    binding:"required,gte=1970,lte=9999"`
}

type GetOperationsResponse struct {
	FileURL string `json:"file_url"`
}

type OperationService interface {
	GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) (string, error)
}

type OperationRepository interface {
	CreateOperation(ctx context.Context, userID int, opType Op, segmentSlug string, doneAt sql.NullTime) error
	GetUserOperations(ctx context.Context, userID int, fromTS sql.NullTime, toTS sql.NullTime) ([]Operation, error)
}
