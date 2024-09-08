package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	storedb services.Store
}

func NewCartController(store services.Store) *CartController {
	return &CartController{
		storedb: store,
	}
}

func (ctrl *CartController) CreateCart(c *gin.Context) {
	var payload models.CreateCartRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.CreateCartParams{
		CartUserID:    payload.CartUserID,
		CartFrID:      payload.CartFrID,
		CartStartDate: payload.CartStartDate,
		CartEndDate:   payload.CartEndDate,
		CartQty:       payload.CartQty,
		CartPrice:     payload.CartPrice,
	}

	cartID, err := ctrl.storedb.CreateCart(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"cart_id": cartID})
}

// GetCartByUserID godoc
// @Summary GetCartByUserID
// @Description GetCartByUserID
// @Tags Cart
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /cart/{user_id} [get]

func (ctrl *CartController) GetCartByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userIDInt32 := int32(userID)
	carts, err := ctrl.storedb.GetCartByUserID(c, &userIDInt32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

func (ctrl *CartController) UpdateCart(c *gin.Context) {
	var payload models.UpdateCartRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	params := db.UpdateCartParams{
		CartFrID:      payload.CartFrID,
		CartStartDate: payload.CartStartDate,
		CartEndDate:   payload.CartEndDate,
		CartQty:       payload.CartQty,
		CartPrice:     payload.CartPrice,
		CartStatus:    payload.CartStatus,
		CartID:        int32(cartID),
	}

	_, err = ctrl.storedb.UpdateCart(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully"})
}

func (ctrl *CartController) UpdateCartQty(c *gin.Context) {
	var payload models.UpdateCartQtyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	params := db.UpdateCartQtyParams{
		CartQty: payload.CartQty,
		CartID:  int32(cartID),
	}

	updatedCart, err := ctrl.storedb.UpdateCartQty(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart quantity"})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

// DeleteCart godoc
// @Summary DeleteCart
// @Description DeleteCart
// @Tags Cart
// @Param cart_id path int true "Cart ID"
// @Success 204 "Cart deleted successfully"
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /cart/{cart_id} [delete]

func (ctrl *CartController) DeleteCart(c *gin.Context) {
	cartID, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	err = ctrl.storedb.DeleteCart(c, int32(cartID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// AddToCart godoc
// @Summary AddToCart
// @Description AddToCart
// @Tags Cart
// @Accept json
// @Produce json
// @Param cart body models.CreateCartRequest true "Add to Cart Request"
// @Success 201 {object} map[string]interface{}
// @Failure 422 {} http.StatusUnprocessableEntity
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /cart/ [post]

func (ctrl *CartController) AddToCart(c *gin.Context) {
	var payload models.CreateCartRequest

	// Bind the JSON payload to the struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(err))
		return
	}

	// Find existing cart item for user and property
	argsFindCart := db.FindCartByUserandRentPropertyParams{
		UserID: *payload.CartUserID,
		RepoID: *payload.CartFrID,
	}

	cartItem, err := models.Nullable(ctrl.storedb.FindCartByUserandRentProperty(c, argsFindCart))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewError(err))
		return
	}

	var response models.CartResponse
	var cart = &db.Cart{}

	if cartItem == nil || cartItem.CartID == 0 {
		// Create new cart
		argsCreateCart := db.CreateCartParams{
			CartUserID:    payload.CartUserID,
			CartFrID:      payload.CartFrID,
			CartStartDate: payload.CartStartDate,
			CartEndDate:   payload.CartEndDate,
			CartQty:       payload.CartQty,
			CartPrice:     payload.CartPrice,
		}
		cart, err = ctrl.storedb.CreateCart(c, argsCreateCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewError(err))
			return
		}

	} else {
		// Update existing cart
		argsUpdateCart := db.UpdateCartParams{
			CartFrID:      payload.CartFrID,
			CartStartDate: payload.CartStartDate,
			CartEndDate:   payload.CartEndDate,
			CartID:        cartItem.CartID,
			CartQty:       payload.CartQty,
			CartPrice:     payload.CartPrice,
			CartStatus:    &payload.CartStatus,
		}
		cart, err = ctrl.storedb.UpdateCart(c, argsUpdateCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewError(err))
			return
		}
	}
	// Fetch all cart items for the user
	carts, err := ctrl.storedb.GetCartByUserID(c, cart.CartUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrDataNotFound)
		return
	}

	// Build the cart response
	response.CartID = int32(carts[0].CartID)
	response.UserID = int32(carts[0].UserID)
	response.CartStatus = carts[0].CartStatus
	response.TotalPrice = 0

	for _, v := range carts {
		detail := &models.CreateRentPropertyDetail{
			RepoName:  v.RepoName,
			RepoPrice: v.CartPrice,
		}
		response.TotalPrice += v.CartPrice * float64(v.CartQty)
		response.Rentals = append(response.Rentals, detail)
	}

	c.JSON(http.StatusCreated, response)
}
