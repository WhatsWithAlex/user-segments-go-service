-- name: CreateOperation :exec
INSERT INTO
    operations (user_id, segment_slug, op_type)
VALUES
    ($2, $3, $1);

-- name: CreateOperationWithTS :exec
INSERT INTO
    operations (user_id, segment_slug, op_type, done_at)
VALUES
    ($2, $3, $1, $4);

-- name: GetOperationsByUserID :many
SELECT
    segment_slug,
    op_type,
    done_at
FROM
    operations
WHERE
    user_id = $1
    AND (
        done_at BETWEEN @from_ts
        AND @to_ts
    );