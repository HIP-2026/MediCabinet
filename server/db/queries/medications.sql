-- name: CreateMedication :one
INSERT INTO medications (name, barcode, active_ingredient)
VALUES ($1, $2, $3)
RETURNING *;
