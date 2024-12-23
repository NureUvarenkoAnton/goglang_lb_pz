-- name: CreateWalk :exec
INSERT INTO walks
  (owner_id, walker_id, pet_id, start_time, state)
VALUES
  (?, ?, ?, ?, 'pending');

-- name: UpdateWalkState :exec
UPDATE walks
SET 
  state = ?,
  finish_time = ?
WHERE
  id = ?;

-- name: GetWalkInfoByParams :many
SELECT * FROM walk_info
WHERE
(sqlc.narg(owner_id) IS NULL OR owner_id = sqlc.narg(owner_id)) AND
(sqlc.narg(walker_id) IS NULL OR walker_id = sqlc.narg(walker_id)) AND
(sqlc.narg(pet_id) IS NULL OR pet_id = sqlc.narg(pet_id)) AND
(sqlc.narg(walk_state) IS NULL OR state = sqlc.narg(walk_state));

-- name: GetWalkInfoByWalkId :one
SELECT * FROM walk_info
WHERE walk_id = ?;

-- name: GetWalksByParams :many
SELECT * FROM walks
WHERE
(sqlc.narg(owner_id) IS NULL OR owner_id = sqlc.narg(owner_id)) AND
(sqlc.narg(walker_id) IS NULL OR walker_id = sqlc.narg(walker_id)) AND
(sqlc.narg(pet_id) IS NULL OR pet_id = sqlc.narg(pet_id)) AND
(sqlc.narg(walk_state) IS NULL OR state = sqlc.narg(walk_state));

-- name: GetWalksByOwnerAndWalkerIds :many
SELECT * FROM walks
WHERE walker_id = ? AND 
      owner_id = ?;

-- name: GetWalksByWalkerId :many
SELECT * FROM walks
WHERE walker_id = ?;

-- name: GetWalksByOwnerId :many
SELECT * FROM walks
WHERE owner_id = ?;

-- name: GetWalkById :one
SELECT * FROM walks
WHERE id = ?;
