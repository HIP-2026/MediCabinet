-- name: CreateInventoryItem :one
INSERT INTO inventory_items (medication_id, quantity, expiry_date, location)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListInventoryItems :many
SELECT
    ii.id,
    ii.medication_id,
    m.name   AS medication_name,
    ii.quantity,
    ii.expiry_date,
    ii.location,
    ii.created_at
FROM inventory_items ii
JOIN medications m ON m.id = ii.medication_id
ORDER BY ii.created_at DESC;
