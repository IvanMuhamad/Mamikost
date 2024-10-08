package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderRentPropertyDetailsController struct {
	storedb services.Store
}

func NewOrderRentPropertyDetailsController(store services.Store) *OrderRentPropertyDetailsController {
	return &OrderRentPropertyDetailsController{
		storedb: store,
	}
}

// AddOrderDetail godoc
// @Summary AddOrderDetail
// @Description AddOrderDetail
// @Tags OrderDetails
// @Accept json
// @Produce json
// @Param body body models.AddOrderDetailRequest true "Order detail request body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /order/detail/ [post]
func (ctrl *OrderRentPropertyDetailsController) AddOrderDetail(c *gin.Context) {
	var payload models.AddOrderDetailRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.AddOrderDetailParams{
		OrpdQtyUnit:    payload.OrpdQtyUnit,
		OrpdPrice:      payload.OrpdPrice,
		OrpdTotalPrice: payload.OrpdTotalPrice,
		OrpdOrpoID:     payload.OrpdOrpoID,
		OrpdRepoID:     payload.OrpdRepoID,
	}

	orderDetail, err := ctrl.storedb.AddOrderDetail(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order detail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_detail_id": orderDetail.OrpdID})
}

func (ctrl *OrderRentPropertyDetailsController) GetAllItemsForOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	orderIDInt32 := int32(orderID)
	items, err := ctrl.storedb.GetAllItemsForOrder(c, &orderIDInt32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (ctrl *OrderRentPropertyDetailsController) UpdateOrderItem(c *gin.Context) {
	var payload models.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	params := db.UpdateOrderItemParams{
		OrpdQtyUnit: payload.OrpdQtyUnit,
		OrpdPrice:   payload.OrpdPrice,
		OrpdRepoID:  payload.OrpdRepoID,
		OrpdID:      int32(orderItemID),
	}

	err = ctrl.storedb.UpdateOrderItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item updated successfully"})
}

// RemoveItemFromOrder godoc
// @Summary RemoveItemFromOrder
// @Description RemoveItemFromOrder
// @Tags OrderDetails
// @Accept json
// @Produce json
// @Param id path int true "Order item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /order/detail/{id} [delete]
func (ctrl *OrderRentPropertyDetailsController) RemoveItemFromOrder(c *gin.Context) {
	orderItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	err = ctrl.storedb.RemoveItemFromOrder(c, int32(orderItemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from order successfully"})
}

func (ctrl *OrderRentPropertyDetailsController) UpdateOrderSubtotal(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctrl.storedb.UpdateOrderSubtotal(c, int32(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order subtotal"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (ctrl *OrderRentPropertyDetailsController) UpdateOrderTotalAndTax(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctrl.storedb.UpdateOrderTotalAndTax(c, int32(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order total and tax"})
		return
	}

	c.JSON(http.StatusOK, order)
}
