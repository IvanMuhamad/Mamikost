-- name: CreateCart :one
INSERT INTO carts (cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_status)
VALUES ($1, $2, $3, $4, $5, $6, 'PENDING')
RETURNING *;

-- name: GetCartByUserID :many
SELECT c.cart_id, c.cart_start_date, c.cart_end_date, c.cart_qty, c.cart_price, (c.cart_qty * c.cart_price)::Double Precision AS cart_total_price, c.cart_status, u.user_name, rp.repo_name
FROM carts c
JOIN users u ON c.cart_user_id = u.user_id
JOIN rent_properties_images rpi ON c.cart_fr_id = rpi.frim_id
JOIN rent_properties rp ON rpi.frim_repo_id = rp.repo_id
WHERE c.cart_user_id = $1;

-- name: FindCartByUserandRentProperty :one
SELECT c.cart_id, c.cart_start_date, c.cart_end_date, c.cart_qty, c.cart_price, (c.cart_qty * c.cart_price)::Double Precision AS cart_total_price, c.cart_status, u.user_name, rp.repo_name
FROM carts c
JOIN users u ON c.cart_user_id = u.user_id
JOIN rent_properties_images rpi ON c.cart_fr_id = rpi.frim_id
JOIN rent_properties rp ON rpi.frim_repo_id = rp.repo_id WHERE u.user_id=$1 AND rp.repo_id=$2 LIMIT 1;

-- name: UpdateCart :one
UPDATE carts
SET cart_fr_id = $1, cart_start_date = $2, cart_end_date = $3, cart_qty = $4, cart_price = $5, cart_status = $6
WHERE cart_id = $7 
RETURNING *;

-- update_cart_item.sql
UPDATE carts
SET cart_qty = $1, cart_price = $2, cart_total_price = $1 * $2, cart_modified = CURRENT_TIMESTAMP
WHERE cart_user_id = $3
  AND cart_fr_id = $4;

-- name: UpdateCartQty :one
UPDATE carts
	SET cart_qty=$1
	WHERE cart_id=$2
	RETURNING *;

-- name: DeleteCart :exec
DELETE FROM carts
	WHERE cart_id=$1
    RETURNING *;

-- name: TransferCartItemsToOrder :exec
--INSERT INTO order_rent_properties_detail (orpd_qty_unit, orpd_price, orpd_total_price, orpd_orpo_id, orpd_repo_id)
--SELECT c.cart_qty, c.cart_price, c.cart_total_price, $1 AS orpd_orpo_id, rpi.frim_repo_id
--FROM carts c
--JOIN rent_properties_images rpi ON c.cart_fr_id = rpi.frim_id
--WHERE c.cart_user_id = $2;

-- name: UpdateOrderTotals :exec
--UPDATE order_rent_properties
--SET orpo_tax = $1, orpo_subtotal = $2, orpo_modified = NOW()
--WHERE orpo_id = $3;




