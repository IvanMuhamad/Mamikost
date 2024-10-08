// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: carts.sql

package db

import (
	"context"

	"time"
)

const createCart = `-- name: CreateCart :one
INSERT INTO carts (cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_status)
VALUES ($1, $2, $3, $4, $5, $6, 'PENDING')
RETURNING cart_id, cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_total_price, cart_modified, cart_status, cart_cart_id
`

type CreateCartParams struct {
	CartUserID    *int32           `json:"cart_user_id"`
	CartFrID      *int32           `json:"cart_fr_id"`
	CartStartDate time.Time `json:"cart_start_date"`
	CartEndDate   time.Time      `json:"cart_end_date"`
	CartQty       int32            `json:"cart_qty"`
	CartPrice     float64          `json:"cart_price"`
}

func (q *Queries) CreateCart(ctx context.Context, arg CreateCartParams) (*Cart, error) {
	row := q.db.QueryRow(ctx, createCart,
		arg.CartUserID,
		arg.CartFrID,
		arg.CartStartDate,
		arg.CartEndDate,
		arg.CartQty,
		arg.CartPrice,
	)
	var i Cart
	err := row.Scan(
		&i.CartID,
		&i.CartUserID,
		&i.CartFrID,
		&i.CartStartDate,
		&i.CartEndDate,
		&i.CartQty,
		&i.CartPrice,
		&i.CartTotalPrice,
		&i.CartModified,
		&i.CartStatus,
		&i.CartCartID,
	)
	return &i, err
}

const deleteCart = `-- name: DeleteCart :exec
DELETE FROM carts
	WHERE cart_id=$1
    RETURNING cart_id, cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_total_price, cart_modified, cart_status, cart_cart_id
`

func (q *Queries) DeleteCart(ctx context.Context, cartID int32) error {
	_, err := q.db.Exec(ctx, deleteCart, cartID)
	return err
}

const findCartByUserandRentProperty = `-- name: FindCartByUserandRentProperty :one
SELECT c.cart_id, c.cart_start_date, c.cart_end_date, c.cart_qty, c.cart_price, (c.cart_qty * c.cart_price)::Double Precision AS cart_total_price, c.cart_status, u.user_name, rp.repo_name
FROM carts c
JOIN users u ON c.cart_user_id = u.user_id
JOIN rent_properties_images rpi ON c.cart_fr_id = rpi.frim_id
JOIN rent_properties rp ON rpi.frim_repo_id = rp.repo_id WHERE u.user_id=$1 AND rp.repo_id=$2 LIMIT 1
`

type FindCartByUserandRentPropertyParams struct {
	UserID int32 `json:"user_id"`
	RepoID int32 `json:"repo_id"`
}

type FindCartByUserandRentPropertyRow struct {
	CartID         int32            `json:"cart_id"`
	CartStartDate  time.Time        `json:"cart_start_date"`
	CartEndDate    time.Time      `json:"cart_end_date"`
	CartQty        int32            `json:"cart_qty"`
	CartPrice      float64          `json:"cart_price"`
	CartTotalPrice float64          `json:"cart_total_price"`
	CartStatus     *string          `json:"cart_status"`
	UserName       string           `json:"user_name"`
	RepoName       string           `json:"repo_name"`
}

func (q *Queries) FindCartByUserandRentProperty(ctx context.Context, arg FindCartByUserandRentPropertyParams) (*FindCartByUserandRentPropertyRow, error) {
	row := q.db.QueryRow(ctx, findCartByUserandRentProperty, arg.UserID, arg.RepoID)
	var i FindCartByUserandRentPropertyRow
	err := row.Scan(
		&i.CartID,
		&i.CartStartDate,
		&i.CartEndDate,
		&i.CartQty,
		&i.CartPrice,
		&i.CartTotalPrice,
		&i.CartStatus,
		&i.UserName,
		&i.RepoName,
	)
	return &i, err
}

const getCartByUserID = `-- name: GetCartByUserID :many
SELECT c.cart_id, c.cart_fr_id, c.cart_start_date, c.cart_end_date, c.cart_qty, c.cart_price, (c.cart_qty * c.cart_price)::Double Precision AS cart_total_price, c.cart_status, u.user_id, rp.repo_name
FROM carts c
JOIN users u ON c.cart_user_id = u.user_id
JOIN rent_properties_images rpi ON c.cart_fr_id = rpi.frim_id
JOIN rent_properties rp ON rpi.frim_repo_id = rp.repo_id
WHERE c.cart_user_id = $1
`

type GetCartByUserIDRow struct {
	CartID         int32            `json:"cart_id"`
	CartFrID	   int32			`json:"cart_fr_id"`
	CartStartDate  time.Time 		`json:"cart_start_date"`
	CartEndDate    time.Time      	`json:"cart_end_date"`
	CartQty        int32            `json:"cart_qty"`
	CartPrice      float64          `json:"cart_price"`
	CartTotalPrice float64          `json:"cart_total_price"`
	CartStatus     *string          `json:"cart_status"`
	UserID         int32           `json:"user_name"`
	RepoName       string           `json:"repo_name"`
}

func (q *Queries) GetCartByUserID(ctx context.Context, cartUserID *int32) ([]*GetCartByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getCartByUserID, cartUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetCartByUserIDRow
	for rows.Next() {
		var i GetCartByUserIDRow
		if err := rows.Scan(
			&i.CartID,
			&i.CartFrID,
			&i.CartStartDate,
			&i.CartEndDate,
			&i.CartQty,
			&i.CartPrice,
			&i.CartTotalPrice,
			&i.CartStatus,
			&i.UserID,
			&i.RepoName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCart = `-- name: UpdateCart :one
UPDATE carts
SET cart_fr_id = $1, cart_start_date = $2, cart_end_date = $3, cart_qty = $4, cart_price = $5, cart_status = $6, cart_total_price = cart_price * cart_qty
WHERE cart_id = $7 
RETURNING cart_id, cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_total_price, cart_modified, cart_status, cart_cart_id
`

type UpdateCartParams struct {
	CartFrID      *int32           `json:"cart_fr_id"`
	CartStartDate time.Time `json:"cart_start_date"`
	CartEndDate   time.Time      `json:"cart_end_date"`
	CartQty       int32            `json:"cart_qty"`
	CartPrice     float64          `json:"cart_price"`
	CartStatus    *string          `json:"cart_status"`
	CartID        int32            `json:"cart_id"`
}

func (q *Queries) UpdateCart(ctx context.Context, arg UpdateCartParams) (*Cart, error) {
	row := q.db.QueryRow(ctx, updateCart,
		arg.CartFrID,
		arg.CartStartDate,
		arg.CartEndDate,
		arg.CartQty,
		arg.CartPrice,
		arg.CartStatus,
		arg.CartID,
	)
	var i Cart
	err := row.Scan(
		&i.CartID,
		&i.CartUserID,
		&i.CartFrID,
		&i.CartStartDate,
		&i.CartEndDate,
		&i.CartQty,
		&i.CartPrice,
		&i.CartTotalPrice,
		&i.CartModified,
		&i.CartStatus,
		&i.CartCartID,
	)
	return &i, err
}

const updateCartQty = `-- name: UpdateCartQty :one
UPDATE carts
	SET cart_qty=$1
	WHERE cart_id=$2
	RETURNING cart_id, cart_user_id, cart_fr_id, cart_start_date, cart_end_date, cart_qty, cart_price, cart_total_price, cart_modified, cart_status, cart_cart_id
`

type UpdateCartQtyParams struct {
	CartQty int32 `json:"cart_qty"`
	CartID  int32 `json:"cart_id"`
}

func (q *Queries) UpdateCartQty(ctx context.Context, arg UpdateCartQtyParams) (*Cart, error) {
	row := q.db.QueryRow(ctx, updateCartQty, arg.CartQty, arg.CartID)
	var i Cart
	err := row.Scan(
		&i.CartID,
		&i.CartUserID,
		&i.CartFrID,
		&i.CartStartDate,
		&i.CartEndDate,
		&i.CartQty,
		&i.CartPrice,
		&i.CartTotalPrice,
		&i.CartModified,
		&i.CartStatus,
		&i.CartCartID,
	)
	return &i, err
}
