-- name: GetContact :one
SELECT * FROM contact
WHERE username = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateContact :one
UPDATE contact
SET first_name = $1 AND last_name = $2 AND profile_photo = $3
WHERE username = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contact
WHERE username = $1;

-- name: CreateContact :one
INSERT INTO contact (first_name, last_name, profile_photo, phone_number, username, hashed_password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetContactList :many
SELECT c.contact_id, c.first_name, c.last_name, c.profile_photo, c.phone_number, c.username, a.group_id
FROM contact as c
JOIN group_member as a ON c.contact_id = a.contact_id
JOIN group_member as b on a.group_id = b.group_id
WHERE a.contact_id != b.contact_id AND a.contact_id != $1;

-- name: SearchContact :many
SELECT c.contact_id, c.first_name, c.last_name, c.profile_photo, c.phone_number, c.username
FROM contact as c
WHERE LOWER(c.username) LIKE '%' || LOWER($1) || '%';