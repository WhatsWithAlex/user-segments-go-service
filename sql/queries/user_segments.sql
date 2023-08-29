-- name: AddUserSegment :exec
INSERT INTO user_segments (
    user_id,
    segment_id,
    remove_at
) VALUES (
    $1, $2, $3
);

-- name: AddUserSegments :copyfrom
INSERT INTO user_segments (
    user_id,
    segment_id,
    remove_at
) VALUES (
    $1, $2, $3
);

-- name: GetUserSegments :many
SELECT t.segment_id, slug, t.remove_at 
FROM (
    SELECT segment_id, remove_at FROM user_segments
    WHERE user_id = $1
) AS t
INNER JOIN segments
ON t.segment_id = segments.id;

-- name: CountUserSegments :one
SELECT count(*) FROM user_segments
WHERE user_id = $1;

-- name: DeleteUserSegment :exec
DELETE FROM user_segments
WHERE user_id = $1;

-- name: DeleteUserSegments :exec
DELETE FROM user_segments
WHERE user_id = ANY($1::int[]);