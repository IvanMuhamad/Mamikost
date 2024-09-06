package models

import (
	db "mamikost/db/sqlc"
	"time"
)

type CreateUserReq struct {
	UserName     *string `json:"user_name,omitempty" binding:"required"`
	UserPassword *string `json:"user_password,omitempty" binding:"required"`
}

type UserResponse struct {
	UserID       int32   `json:"user_id"`
	UserName     *string `json:"user_name"`
	UserPassword *string `json:"user_password"`
	UserPhone    *string `json:"user_phone"`
	UserToken    *string `json:"user_token"`
}

type UpdateCategoryParams struct {
	CateName string `json:"cate_name"`
	CateID   int32  `json:"cate_id"`
}

type CreateRentPropertyRequest struct {
	RepoName   string  `json:"repo_name" binding:"required"`
	RepoDesc   string  `json:"repo_desc"`
	RepoPrice  float64 `json:"repo_price" binding:"required"`
	RepoCateID int32   `json:"repo_cate_id" binding:"required"`
}

type UpdateRentPropertyRequest struct {
	RepoName   string  `json:"repo_name" binding:"required"`
	RepoDesc   string  `json:"repo_desc"`
	RepoPrice  float64 `json:"repo_price" binding:"required"`
	RepoCateID int32   `json:"repo_cate_id" binding:"required"`
}

type UploadImageRequest struct {
	FrimFilename string `json:"frim_filename" binding:"required"`
	FrimDefault  string `json:"frim_default" binding:"required"`
	FrimRepoID   int32  `json:"frim_repo_id" binding:"required"`
}

type CreateUserRequest struct {
	UserName     string `json:"user_name" binding:"required"`
	UserPassword string `json:"user_password" binding:"required"`
	UserEmail    string `json:"user_email"`
	UserPhone    string `json:"user_phone"`
}

type UpdateUserRequest struct {
	UserName  string `json:"user_name" binding:"required"`
	UserEmail string `json:"user_email" binding:"required"`
	UserPhone string `json:"user_phone" binding:"required"`
}

type CreateOrderRequest struct {
	OrpoPurchaseNo string  `json:"orpo_purchase_no"`
	OrpoTax        float64 `json:"orpo_tax"`
	OrpoSubtotal   float64 `json:"orpo_subtotal"`
	OrpoPatrxNo    *string `json:"orpo_patrx_no"`
	OrpoUserID     *int32  `json:"orpo_user_id"`
}

type UpdateOrderRequest struct {
	OrpoPurchaseNo string   `json:"orpo_purchase_no"`
	OrpoTax        *float64 `json:"orpo_tax"`
	OrpoSubtotal   *float64 `json:"orpo_subtotal"`
	OrpoPatrxNo    *string  `json:"orpo_patrx_no"`
	OrpoID         int32    `json:"orpo_id"`
}

type AddItemOrderRequest struct {
	OrpdQtyUnit int32   `json:"orpd_qty_unit"`
	OrpdPrice   float64 `json:"orpd_price"`
	OrpdOrpoID  *int32  `json:"orpd_orpo_id"`
	OrpdRepoID  *int32  `json:"orpd_repo_id"`
}

type CreateCartRequest struct {
	CartUserID    *int32    `json:"cart_user_id"`
	CartFrID      *int32    `json:"cart_fr_id"`
	CartStartDate time.Time `json:"cart_start_date"`
	CartEndDate   time.Time `json:"cart_end_date"`
	CartQty       int32     `json:"cart_qty"`
	CartPrice     float64   `json:"cart_price"`
	CartStatus    string    `json:"cart_status"`
}

type OrderRentProperty struct {
	OrpoID         int32     `json:"orpo_id"`
	OrpoPurchaseNo string    `json:"orpo_purchase_no"`
	OrpoTax        *float64  `json:"orpo_tax"`
	OrpoSubtotal   *float64  `json:"orpo_subtotal"`
	OrpoPatrxNo    *string   `json:"orpo_patrx_no"`
	OrpoModified   time.Time `json:"orpo_modified"`
	OrpoUserID     *int32    `json:"orpo_user_id"`
}

type AddOrderDetailRequest struct {
	OrpdQtyUnit    int32    `json:"orpd_qty_unit"`
	OrpdPrice      float64  `json:"orpd_price"`
	OrpdTotalPrice *float64 `json:"orpd_total_price"`
	OrpdOrpoID     *int32   `json:"orpd_orpo_id"`
	OrpdRepoID     *int32   `json:"orpd_repo_id"`
}

type UpdateOrderItemRequest struct {
	OrpdQtyUnit int32   `json:"orpd_qty_unit"`
	OrpdPrice   float64 `json:"orpd_price"`
	OrpdRepoID  *int32  `json:"orpd_repo_id"`
	OrpdID      int32   `json:"orpd_id"`
}

type UpdateCartRequest struct {
	CartFrID      *int32    `json:"cart_fr_id"`
	CartStartDate time.Time `json:"cart_start_date"`
	CartEndDate   time.Time `json:"cart_end_date"`
	CartQty       int32     `json:"cart_qty"`
	CartPrice     float64   `json:"cart_price"`
	CartStatus    *string   `json:"cart_status"`
	CartID        int32     `json:"cart_id"`
}

type UpdateCartQtyRequest struct {
	CartQty int32 `json:"cart_qty"`
	CartID  int32 `json:"cart_id"`
}

type CartResponse struct {
	CartID     int32                       `json:"cart_id"`
	UserID     int32                       `json:"user_id"`
	CartStatus *string                     `json:"cart_status"`
	TotalPrice float64                     `json:"total_price"`
	Rentals    []*CreateRentPropertyDetail `json:"rentals"`
}

type CreateRentPropertyDetail struct {
	RepoName   string  `json:"repo_name"`
	RepoDesc   *string `json:"repo_desc"`
	RepoPrice  float64 `json:"repo_price"`
	RepoCateID *int32  `json:"repo_cate_id"`
}

type OrderResponse struct {
	OrderID      int32                        `json:"order_id"`
	PurchaseNo   string                       `json:"purchase_no"`
	Tax          float64                      `json:"tax"`
	Subtotal     float64                      `json:"subtotal"`
	PatrxNo      *string                      `json:"patrx_no"`
	UserID       int32                        `json:"user_id"`
	OrderDetails []*db.GetAllItemsForOrderRow `json:"order_detail"`
}

type CreateCategoryRequest struct {
	CateName string `json:"cate_name"`
}
