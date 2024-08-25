-- name: CreateMessageGroup :one
INSERT INTO message_group (group_name)
VALUES ($1)
RETURNING *;

-- name: DeleteMessageGroup :exec
DELETE FROM message_group WHERE group_id = $1;