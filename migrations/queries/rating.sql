-- name: RatingsByRaterId :many
SELECT * FROM ratings
WHERE rater_id = ?;

-- name: RatingsByRateeId :many
SELECT * FROM ratings
WHERE ratee_id = ?;

-- name: RatingByIds :one
SELECT * FROM ratings
WHERE ratee_id = ? AND rater_id = ?;

-- name: AddRating :exec
INSERT INTO ratings
  (rater_id, ratee_id, value)
VALUES
  (?, ?, ?);
