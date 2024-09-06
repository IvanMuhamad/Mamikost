package services

import (
	db "Mamikost2/db/sqlc"
	"context"
	"log"
)

func (sm *StoreManager) CreateOrderTx(ctx context.Context, args db.CreateOrderParams) (*db.OrderRentProperty, error) {
	tx, err := sm.dbConn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	qtx := sm.Queries.WithTx(tx)

	// Populate cart list
	carts, err := qtx.GetCartByUserID(ctx, args.OrpoUserID)
	if err != nil {
		return nil, err
	}

	log.Println(carts)

	//create order
	newOrder, err := qtx.CreateOrder(ctx, args)
	if err != nil {
		return nil, err
	}

	var subtotal float64

	for _, cartItem := range carts {
		var totalPrice = cartItem.CartPrice * float64(cartItem.CartQty)
		subtotal += totalPrice
		orderDetailArgs := db.AddOrderDetailParams{
			OrpdOrpoID:     &newOrder.OrpoID,
			OrpdRepoID:     &cartItem.CartFrID,
			OrpdQtyUnit:    cartItem.CartQty,
			OrpdPrice:      cartItem.CartPrice,
			OrpdTotalPrice: &totalPrice,
		}

		_, err = qtx.AddOrderDetail(ctx, orderDetailArgs)
		if err != nil {
			return nil, err
		}
	}

	_, err = qtx.UpdateOrder(ctx, db.UpdateOrderParams{
		OrpoPurchaseNo: newOrder.OrpoPurchaseNo,
		OrpoTax:        newOrder.OrpoTax,
		OrpoSubtotal:   &subtotal,
		OrpoPatrxNo:    newOrder.OrpoPatrxNo,
		OrpoID:         newOrder.OrpoID,
	})
	if err != nil {
		return nil, err
	}

	for _, v := range carts {
		_, err = qtx.UpdateCart(ctx, db.UpdateCartParams{
			CartFrID:      &v.CartFrID,
			CartStartDate: v.CartStartDate,
			CartEndDate:   v.CartEndDate,
			CartQty:       v.CartQty,
			CartPrice:     v.CartPrice,
			CartStatus:    &[]string{"COMPLETED"}[0],
			CartID:        v.CartID,
		})
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	newOrder.OrpoSubtotal = &subtotal
	return newOrder, nil

}
