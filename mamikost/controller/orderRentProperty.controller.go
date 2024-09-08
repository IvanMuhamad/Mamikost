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

// CreateOrder godoc
// @Summary CreateOrder
// @Description CreateOrder
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order creation request payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /order/ [post]
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

	c.JSON(http.StatusCreated, response)
}

// GetOrderByID godoc
// @Summary GetOrderByID
// @Description GetOrderByID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /order/{id} [get]
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

// DeleteOrder godoc
// @Summary DeleteOrder
// @Description DeleteOrder
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /order/{id} [delete]
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
