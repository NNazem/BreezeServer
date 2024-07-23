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