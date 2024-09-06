// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type Cart struct {
	CartID         int32            `json:"cart_id"`
	CartUserID     *int32           `json:"cart_user_id"`
	CartFrID       *int32           `json:"cart_fr_id"`
	CartStartDate  time.Time `json:"cart_start_date"`
	CartEndDate    time.Time      `json:"cart_end_date"`
	CartQty        int32            `json:"cart_qty"`
	CartPrice      float64          `json:"cart_price"`
	CartTotalPrice *float64         `json:"cart_total_price"`
	CartModified   time.Time `json:"cart_modified"`
	CartStatus     *string          `json:"cart_status"`
	CartCartID     *int32           `json:"cart_cart_id"`
}

type Category struct {
	CateID   int32  `json:"cate_id"`
	CateName string `json:"cate_name"`
}

type OrderRentPropertiesDetail struct {
	OrpdID         int32    `json:"orpd_id"`
	OrpdQtyUnit    int32    `json:"orpd_qty_unit"`
	OrpdPrice      float64  `json:"orpd_price"`
	OrpdTotalPrice *float64 `json:"orpd_total_price"`
	OrpdOrpoID     *int32   `json:"orpd_orpo_id"`
	OrpdRepoID     *int32   `json:"orpd_repo_id"`
}

type OrderRentProperty struct {
	OrpoID         int32            `json:"orpo_id"`
	OrpoPurchaseNo string           `json:"orpo_purchase_no"`
	OrpoTax        *float64         `json:"orpo_tax"`
	OrpoSubtotal   *float64         `json:"orpo_subtotal"`
	OrpoPatrxNo    *string          `json:"orpo_patrx_no"`
	OrpoModified   time.Time `json:"orpo_modified"`
	OrpoUserID     *int32           `json:"orpo_user_id"`
}

type RentPropertiesImage struct {
	FrimID       int32   `json:"frim_id"`
	FrimFilename string  `json:"frim_filename"`
	FrimDefault  *string `json:"frim_default"`
	FrimRepoID   *int32  `json:"frim_repo_id"`
}

type RentProperty struct {
	RepoID       int32            `json:"repo_id"`
	RepoName     string           `json:"repo_name"`
	RepoDesc     *string          `json:"repo_desc"`
	RepoPrice    float64          `json:"repo_price"`
	RepoModified time.Time `json:"repo_modified"`
	RepoCateID   *int32           `json:"repo_cate_id"`
}

type User struct {
	UserID       int32   `json:"user_id"`
	UserName     string  `json:"user_name"`
	UserPassword string  `json:"user_password"`
	UserEmail    string  `json:"user_email"`
	UserPhone    string  `json:"user_phone"`
	UserToken    *string `json:"user_token"`
}
