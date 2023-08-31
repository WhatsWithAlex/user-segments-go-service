-- name: CreateSegment :one
INSERT INTO
    segments (slug, auto_prob)
VALUES
    ($1, $2) RETURNING id;

-- name: GetSegmentByID :one
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    id = $1;

-- name: GetSegmentBySlug :one
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    slug = $1;

-- name: GetAutoSegments :many
SELECT
    id,
    slug,
    auto_prob
FROM
    segments
WHERE
    auto_prob > 0.0;

-- name: UpdateSegment :exec
UPDATE
    segments
SET
    auto_prob = $2
WHERE
    id = $1;

-- name: DeleteSegment :exec
DELETE FROM
    segments
WHERE
    id = $1;

-- name: DeleteSegmentBySlug :one
DELETE FROM
    segments
WHERE
    slug = $1 RETURNING id;