// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user_segments.sql

package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addUserSegment = `-- name: AddUserSegment :exec
INSERT INTO
    user_segments (user_id, segment_id, remove_at)
VALUES
    ($1, $2, $3) ON CONFLICT(user_id, segment_id) DO NOTHING
`

type AddUserSegmentParams struct {
	UserID    int32              `json:"user_id"`
	SegmentID int32              `json:"segment_id"`
	RemoveAt  pgtype.Timestamptz `json:"remove_at"`
}

func (q *Queries) AddUserSegment(ctx context.Context, arg AddUserSegmentParams) error {
	_, err := q.db.Exec(ctx, addUserSegment, arg.UserID, arg.SegmentID, arg.RemoveAt)
	return err
}

const addUserSegmentBySlug = `-- name: AddUserSegmentBySlug :exec
INSERT INTO
    user_segments (user_id, segment_id, remove_at)
SELECT
    $1,
    segments.id,
    $3
FROM
    segments
WHERE
    slug = $2 ON CONFLICT(user_id, segment_id) DO NOTHING
`

type AddUserSegmentBySlugParams struct {
	UserID   int32              `json:"user_id"`
	Slug     string             `json:"slug"`
	RemoveAt pgtype.Timestamptz `json:"remove_at"`
}

func (q *Queries) AddUserSegmentBySlug(ctx context.Context, arg AddUserSegmentBySlugParams) error {
	_, err := q.db.Exec(ctx, addUserSegmentBySlug, arg.UserID, arg.Slug, arg.RemoveAt)
	return err
}

const countUserSegments = `-- name: CountUserSegments :one
SELECT
    count(*)
FROM
    user_segments
WHERE
    user_id = $1
`

func (q *Queries) CountUserSegments(ctx context.Context, userID int32) (int64, error) {
	row := q.db.QueryRow(ctx, countUserSegments, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteAllExpiredUserSegments = `-- name: DeleteAllExpiredUserSegments :many
DELETE FROM
    user_segments
WHERE
    remove_at < now() RETURNING user_id,
    (
        SELECT
            slug
        FROM
            segments
        WHERE
            id = user_segments.segment_id
    ),
    remove_at
`

type DeleteAllExpiredUserSegmentsRow struct {
	UserID   int32              `json:"user_id"`
	Slug     string             `json:"slug"`
	RemoveAt pgtype.Timestamptz `json:"remove_at"`
}

func (q *Queries) DeleteAllExpiredUserSegments(ctx context.Context) ([]DeleteAllExpiredUserSegmentsRow, error) {
	rows, err := q.db.Query(ctx, deleteAllExpiredUserSegments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DeleteAllExpiredUserSegmentsRow{}
	for rows.Next() {
		var i DeleteAllExpiredUserSegmentsRow
		if err := rows.Scan(&i.UserID, &i.Slug, &i.RemoveAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteUserSegment = `-- name: DeleteUserSegment :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id = $2
`

type DeleteUserSegmentParams struct {
	UserID    int32 `json:"user_id"`
	SegmentID int32 `json:"segment_id"`
}

func (q *Queries) DeleteUserSegment(ctx context.Context, arg DeleteUserSegmentParams) error {
	_, err := q.db.Exec(ctx, deleteUserSegment, arg.UserID, arg.SegmentID)
	return err
}

const deleteUserSegmentBySlug = `-- name: DeleteUserSegmentBySlug :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id IN (
        SELECT
            id
        FROM
            segments
        WHERE
            slug = $2
    )
`

type DeleteUserSegmentBySlugParams struct {
	UserID int32  `json:"user_id"`
	Slug   string `json:"slug"`
}

func (q *Queries) DeleteUserSegmentBySlug(ctx context.Context, arg DeleteUserSegmentBySlugParams) error {
	_, err := q.db.Exec(ctx, deleteUserSegmentBySlug, arg.UserID, arg.Slug)
	return err
}

const deleteUserSegments = `-- name: DeleteUserSegments :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id = ANY($2 :: int [ ])
`

type DeleteUserSegmentsParams struct {
	UserID     int32   `json:"user_id"`
	SegmentIDs []int32 `json:"segment_ids"`
}

func (q *Queries) DeleteUserSegments(ctx context.Context, arg DeleteUserSegmentsParams) error {
	_, err := q.db.Exec(ctx, deleteUserSegments, arg.UserID, arg.SegmentIDs)
	return err
}

const deleteUserSegmentsBySlugs = `-- name: DeleteUserSegmentsBySlugs :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id IN (
        SELECT
            id
        FROM
            segments
        WHERE
            slug = ANY($2 :: char [ ])
    )
`

type DeleteUserSegmentsBySlugsParams struct {
	UserID       int32    `json:"user_id"`
	SegmentSlugs []string `json:"segment_slugs"`
}

func (q *Queries) DeleteUserSegmentsBySlugs(ctx context.Context, arg DeleteUserSegmentsBySlugsParams) error {
	_, err := q.db.Exec(ctx, deleteUserSegmentsBySlugs, arg.UserID, arg.SegmentSlugs)
	return err
}

const getActiveUserSegments = `-- name: GetActiveUserSegments :many
SELECT
    slug
FROM
    (
        SELECT
            segment_id,
            remove_at
        FROM
            user_segments
        WHERE
            user_id = $1
            AND (
                remove_at > now()
                OR remove_at is NULL
            )
    ) AS t
    INNER JOIN segments ON t.segment_id = segments.id
`

func (q *Queries) GetActiveUserSegments(ctx context.Context, userID int32) ([]string, error) {
	rows, err := q.db.Query(ctx, getActiveUserSegments, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		items = append(items, slug)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSegments = `-- name: GetUserSegments :many
SELECT
    t.segment_id,
    slug,
    t.remove_at
FROM
    (
        SELECT
            segment_id,
            remove_at
        FROM
            user_segments
        WHERE
            user_id = $1
    ) AS t
    INNER JOIN segments ON t.segment_id = segments.id
`

type GetUserSegmentsRow struct {
	SegmentID int32              `json:"segment_id"`
	Slug      string             `json:"slug"`
	RemoveAt  pgtype.Timestamptz `json:"remove_at"`
}

func (q *Queries) GetUserSegments(ctx context.Context, userID int32) ([]GetUserSegmentsRow, error) {
	rows, err := q.db.Query(ctx, getUserSegments, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserSegmentsRow{}
	for rows.Next() {
		var i GetUserSegmentsRow
		if err := rows.Scan(&i.SegmentID, &i.Slug, &i.RemoveAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsers = `-- name: GetUsers :many
SELECT
    DISTINCT user_id
FROM
    user_segments
`

func (q *Queries) GetUsers(ctx context.Context) ([]int32, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var user_id int32
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersBySegmentID = `-- name: GetUsersBySegmentID :many
SELECT
    user_id
FROM
    user_segments
WHERE
    segment_id = $1
`

func (q *Queries) GetUsersBySegmentID(ctx context.Context, segmentID int32) ([]int32, error) {
	rows, err := q.db.Query(ctx, getUsersBySegmentID, segmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var user_id int32
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersBySegmentSlug = `-- name: GetUsersBySegmentSlug :many
SELECT
    user_id
FROM
    user_segments
WHERE
    segment_id IN (
        SELECT
            id
        FROM
            segments
        WHERE
            slug = $1
    )
`

func (q *Queries) GetUsersBySegmentSlug(ctx context.Context, slug string) ([]int32, error) {
	rows, err := q.db.Query(ctx, getUsersBySegmentSlug, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var user_id int32
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
