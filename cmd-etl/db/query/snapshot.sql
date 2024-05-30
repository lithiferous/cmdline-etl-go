-- name: CreateSnapshot :one
INSERT INTO snapshots (
	user_name,
	store_name,
	credit_limit,
	snapshot_at
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: ListSnapshots :many
SELECT * FROM snapshots
ORDER BY snapshot_at DESC;
