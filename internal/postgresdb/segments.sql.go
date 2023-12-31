// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: segments.sql

package postgresdb

import (
	"context"
)

const createSegment = `-- name: CreateSegment :one
INSERT INTO
    segments (slug, auto_prob)
VALUES
    ($1, $2) RETURNING id
`

type CreateSegmentParams struct {
	Slug     string  `json:"slug"`
	AutoProb float32 `json:"auto_prob"`
}

func (q *Queries) CreateSegment(ctx context.Context, arg CreateSegmentParams) (int32, error) {
	row := q.db.QueryRow(ctx, createSegment, arg.Slug, arg.AutoProb)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteSegment = `-- name: DeleteSegment :exec
DELETE FROM
    segments
WHERE
    id = $1
`

func (q *Queries) DeleteSegment(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteSegment, id)
	return err
}

const deleteSegmentBySlug = `-- name: DeleteSegmentBySlug :one
DELETE FROM
    segments
WHERE
    slug = $1 RETURNING id
`

func (q *Queries) DeleteSegmentBySlug(ctx context.Context, slug string) (int32, error) {
	row := q.db.QueryRow(ctx, deleteSegmentBySlug, slug)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getAutoSegments = `-- name: GetAutoSegments :many
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    auto_prob > 0.0
`

func (q *Queries) GetAutoSegments(ctx context.Context) ([]Segment, error) {
	rows, err := q.db.Query(ctx, getAutoSegments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Segment{}
	for rows.Next() {
		var i Segment
		if err := rows.Scan(&i.ID, &i.Slug, &i.AutoProb); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSegmentByID = `-- name: GetSegmentByID :one
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    id = $1
`

func (q *Queries) GetSegmentByID(ctx context.Context, id int32) (Segment, error) {
	row := q.db.QueryRow(ctx, getSegmentByID, id)
	var i Segment
	err := row.Scan(&i.ID, &i.Slug, &i.AutoProb)
	return i, err
}

const getSegmentBySlug = `-- name: GetSegmentBySlug :one
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    slug = $1
`

func (q *Queries) GetSegmentBySlug(ctx context.Context, slug string) (Segment, error) {
	row := q.db.QueryRow(ctx, getSegmentBySlug, slug)
	var i Segment
	err := row.Scan(&i.ID, &i.Slug, &i.AutoProb)
	return i, err
}

const updateSegment = `-- name: UpdateSegment :exec
UPDATE
    segments
SET
    auto_prob = $2
WHERE
    id = $1
`

type UpdateSegmentParams struct {
	ID       int32   `json:"id"`
	AutoProb float32 `json:"auto_prob"`
}

func (q *Queries) UpdateSegment(ctx context.Context, arg UpdateSegmentParams) error {
	_, err := q.db.Exec(ctx, updateSegment, arg.ID, arg.AutoProb)
	return err
}
