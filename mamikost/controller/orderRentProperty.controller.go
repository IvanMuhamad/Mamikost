package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderRentPropertyController struct {
	storedb services.Store
}

func NewOrderRentPropertyController(store services.Store) *OrderRentPropertyController {
	return &OrderRentPropertyController{
		storedb: store,
	}
}

func (ctrl *OrderRentPropertyController) CreateOrder(c *gin.Context) {
	var payload models.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.CreateOrderParams{
		OrpoPurchaseNo: payload.OrpoPurchaseNo,
		OrpoTax:        &payload.OrpoTax,
		OrpoSubtotal:   &payload.OrpoSubtotal,
		OrpoPatrxNo:    payload.OrpoPatrxNo,
		OrpoUserID:     payload.OrpoUserID,
	}

	orderID, err := ctrl.storedb.CreateOrderTx(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order",
			"error": err})
		return
	}

	orderdtl, err := ctrl.storedb.GetAllItemsForOrder(c, &orderID.OrpoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch order detail",
			"error": err})
	}

	response := models.OrderResponse{
		OrderID:      orderID.OrpoID,
		PurchaseNo:   orderID.OrpoPurchaseNo,
		Tax:          *orderID.OrpoTax,
		Subtotal:     *orderID.OrpoSubtotal,
		PatrxNo:      orderID.OrpoPatrxNo,
		UserID:       *orderID.OrpoUserID,
		OrderDetails: orderdtl,
	}

	//c.JSON(http.StatusCreated, gin.H{
	//"order": gin.H{
	//"order_id":     orderID.OrpoID,
	//"purchase_no":  orderID.OrpoPurchaseNo,
	//"tax":          orderID.OrpoTax,
	//"subtotal":     orderID.OrpoSubtotal,
	//"patrx_no":     orderID.OrpoPatrxNo,
	//"user_id":      orderID.OrpoUserID,
	//"order_detail": orderdtl,
	//}})

	c.JSON(http.StatusCreated, response)
}

func (ctrl *OrderRentPropertyController) AddItemToOrder(c *gin.Context) {
	var payload models.AddItemOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.AddItemOrderParams{
		OrpdQtyUnit: payload.OrpdQtyUnit,
		OrpdPrice:   payload.OrpdPrice,
		OrpdOrpoID:  payload.OrpdOrpoID,
		OrpdRepoID:  payload.OrpdRepoID,
	}

	if _, err := ctrl.storedb.AddItemOrder(c, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to order successfully"})
}

func (ctrl *OrderRentPropertyController) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctrl.storedb.FindOrderByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (ctrl *OrderRentPropertyController) UpdateOrder(c *gin.Context) {
	var payload models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.UpdateOrderParams{
		OrpoPurchaseNo: payload.OrpoPurchaseNo,
		OrpoTax:        payload.OrpoTax,
		OrpoSubtotal:   payload.OrpoSubtotal,
		OrpoPatrxNo:    payload.OrpoPatrxNo,
		OrpoID:         payload.OrpoID,
	}

	if _, err := ctrl.storedb.UpdateOrder(c, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder deletes an order by its ID
func (ctrl *OrderRentPropertyController) DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = ctrl.storedb.DeleteOrder(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
