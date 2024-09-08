-- name: CreateOrder :one
INSERT INTO order_rent_properties (orpo_purchase_no, orpo_tax, orpo_subtotal, orpo_patrx_no, orpo_user_id) 
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: FindAllOrders :many
SELECT orpo_id, orpo_purchase_no, orpo_tax, orpo_subtotal, orpo_patrx_no, orpo_modified, orpo_user_id 
FROM order_rent_properties;

-- name: FindOrderByID :one
SELECT orpo_id, orpo_purchase_no, orpo_tax, orpo_subtotal, orpo_patrx_no, orpo_modified, orpo_user_id 
FROM order_rent_properties 
WHERE orpo_id = $1;

-- name: UpdateOrder :one
UPDATE order_rent_properties 
SET orpo_purchase_no = $1, 
    orpo_tax = $2, 
    orpo_subtotal = $3, 
    orpo_patrx_no = $4 
WHERE orpo_id = $5
RETURNING *;

-- name: UpdateOrderPatrxNo :one
UPDATE order_rent_properties 
SET orpo_patrx_no = $1 
WHERE orpo_id = $2
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM order_rent_properties 
WHERE orpo_id = $1;