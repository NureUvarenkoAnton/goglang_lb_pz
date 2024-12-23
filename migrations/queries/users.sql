-- name: GetUserById :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: GetUsers :many
SELECT * FROM users
WHERE 
(sqlc.narg(user_type) is NULL OR user_type = sqlc.narg(user_type)) AND
(sqlc.narg(is_banned) is NULL OR is_banned = sqlc.narg(is_banned)) AND
(sqlc.narg(is_deleted) is NULL OR is_deleted = sqlc.narg(is_deleted));

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (name, email, password, user_type, created_at) VALUES(?, ?, ?, ?, NOW());

-- name: UpdateUser :exec
UPDATE users
SET 
  name = ?,
  email = ?
WHERE
  id = ?;

-- name: SetDeleteState :exec
UPDATE users
SET is_deleted = ?
WHERE id = ?;

-- name: DeleteMarkedUsers :exec
DELETE FROM users
WHERE 
  is_deleted = true AND
  deleted_at > ?;

-- name: SetBanState :exec
UPDATE users
SET is_banned = ?
WHERE id = ?;
