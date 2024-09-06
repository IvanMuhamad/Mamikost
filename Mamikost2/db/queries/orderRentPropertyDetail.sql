-- name: AddOrderDetail :one
INSERT INTO order_rent_properties_detail (orpd_qty_unit, orpd_price, orpd_total_price, orpd_orpo_id, orpd_repo_id) 
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: AddItemOrder :one
INSERT INTO order_rent_properties_detail (orpd_qty_unit, orpd_price, orpd_orpo_id, orpd_repo_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllItemsForOrder :many
SELECT ord.orpd_id, ord.orpd_qty_unit, ord.orpd_price, ord.orpd_total_price, rp.repo_name
FROM order_rent_properties_detail ord
JOIN rent_properties rp ON ord.orpd_repo_id = rp.repo_id
WHERE ord.orpd_orpo_id = $1;

-- name: UpdateOrderItem :exec
UPDATE order_rent_properties_detail
SET orpd_qty_unit = $1, orpd_price = $2, orpd_repo_id = $3
WHERE orpd_id = $4;

-- name: RemoveItemFromOrder :exec
DELETE FROM order_rent_properties_detail
WHERE orpd_id = $1;

-- name: UpdateOrderSubtotal :one
UPDATE order_rent_properties
SET orpo_subtotal = (
        SELECT COALESCE(SUM(orpd_total_price), 0) 
        FROM order_rent_properties_detail ord
        WHERE ord.orpd_orpo_id = order_rent_properties.orpo_id
    ),
    orpo_modified = CURRENT_TIMESTAMP
WHERE orpo_id = $1
RETURNING *;

-- name: UpdateOrderTotalAndTax :one
UPDATE order_rent_properties
SET orpo_subtotal = (
        SELECT COALESCE(SUM(orpd_total_price), 0) 
        FROM order_rent_properties_detail ord
        WHERE ord.orpd_orpo_id = order_rent_properties.orpo_id
    ),
    orpo_tax = (
        SELECT COALESCE(SUM(orpd_total_price) * 0.1, 0)
        FROM order_rent_properties_detail ord
        WHERE ord.orpd_orpo_id = order_rent_properties.orpo_id
    ),
    orpo_total_price = (
        SELECT COALESCE(SUM(orpd_total_price), 0) 
        FROM order_rent_properties_detail ord
        WHERE ord.orpd_orpo_id = order_rent_properties.orpo_id
    ) + (
        SELECT COALESCE(SUM(orpd_total_price) * 0.1, 0)
        FROM order_rent_properties_detail ord
        WHERE ord.orpd_orpo_id = order_rent_properties.orpo_id
    ),
    orpo_modified = CURRENT_TIMESTAMP
WHERE orpo_id = $1
RETURNING *;