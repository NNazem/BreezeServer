-- name: CreateMessageGroup :one
INSERT INTO message_group (group_name)
VALUES ($1)
RETURNING *;