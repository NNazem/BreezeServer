-- name: CreateGroupMember :one
INSERT INTO group_member(
                         contact_id,
                         group_id
) VALUES ($1, $2) RETURNING *;