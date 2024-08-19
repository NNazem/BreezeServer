-- name: CreateMessage :one
INSERT INTO message(
                    username,
                    message_text,
                    group_id
) VALUES ($1, $2, $3) RETURNING *;

-- name: ListUserGroupMessage :many
SELECT *
FROM message
WHERE group_id = $1;

-- name: FetchLastMessage :one
SELECT  a.*
FROM message a
WHERE a.sent_datetime = (select max(b.sent_datetime) from message b where b.group_id = a.group_id) AND a.group_id = $1;