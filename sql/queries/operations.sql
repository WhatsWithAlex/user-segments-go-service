-- name: CreateOperation :exec
INSERT INTO operations (
    user_id,
    segment_slug,
    op_type
) VALUES (
    $2, $3, $1
);

-- name: GetOperations :many
SELECT segment_slug, op_type, done_at FROM operations
WHERE user_id = $1 AND (done_at BETWEEN @ts_from AND @ts_to);