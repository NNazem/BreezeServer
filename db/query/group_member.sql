-- name: CreateGroupMember :one
INSERT INTO group_member(
                         contact_id,
                         group_id
) VALUES ($1, $2) RETURNING *;

-- name: GetGroupId :one
SELECT a.group_id
FROM group_member as a JOIN group_member as b ON a.group_id = b.group_id
WHERE a.contact_id != b.contact_id AND a.contact_id = $1 AND b.contact_id = $2;
