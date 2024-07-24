-- name: CreateMessage :one
INSERT INTO message(
                    username,
                    message_text,
                    group_id
) VALUES ($1, $2, $3) RETURNING *;

-- name: ListUserGroupMessage :many
SELECT *
FROM message
WHERE username = $1 AND group_id = $2;