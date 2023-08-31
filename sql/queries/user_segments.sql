-- name: AddUserSegment :exec
INSERT INTO
    user_segments (user_id, segment_id, remove_at)
VALUES
    ($1, $2, $3);

-- name: AddUserSegmentBySlug :exec
INSERT INTO
    user_segments (user_id, segment_id, remove_at)
SELECT
    $1,
    segments.id,
    $3
FROM
    segments
WHERE
    slug = $2;

-- name: GetUserSegments :many
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
    INNER JOIN segments ON t.segment_id = segments.id;

-- name: GetActiveUserSegments :many
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
    INNER JOIN segments ON t.segment_id = segments.id;

-- name: GetUsersBySegmentID :many
SELECT
    user_id
FROM
    user_segments
WHERE
    segment_id = $1;

-- name: GetUsers :many
SELECT
    DISTINCT user_id
FROM
    user_segments;

-- name: CountUserSegments :one
SELECT
    count(*)
FROM
    user_segments
WHERE
    user_id = $1;

-- name: DeleteUserSegment :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id = $2;

-- name: DeleteUserSegmentBySlug :exec
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
    );

-- name: DeleteUserSegments :exec
DELETE FROM
    user_segments
WHERE
    user_id = $1
    AND segment_id = ANY(@segment_ids :: int [ ]);

-- name: DeleteUserSegmentsBySlugs :exec
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
            slug = ANY(@segment_slugs :: char [ ])
    );

-- name: DeleteAllExpiredUserSegments :many
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
    remove_at;