package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RentPropertyController struct {
	storedb services.Store
}

func NewRentPropertyController(store services.Store) *RentPropertyController {
	return &RentPropertyController{
		storedb: store,
	}
}

// CreateRentProperty godoc
// @Summary CreateRentProperty
// @Description CreateRentProperty
// @Tags RentProperty
// @Accept json
// @Produce json
// @Param property body models.CreateRentPropertyRequest true "Create Rent Property Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/ [post]
func (ctrl *RentPropertyController) CreateRentProperty(c *gin.Context) {
	var payload models.CreateRentPropertyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.CreateRentPropertyParams{
		RepoName:   payload.RepoName,
		RepoDesc:   &payload.RepoDesc,
		RepoPrice:  payload.RepoPrice,
		RepoCateID: &payload.RepoCateID,
	}

	rentPropertyID, err := ctrl.storedb.CreateRentProperty(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"rent_property_id": rentPropertyID})
}

// GetAllRentProperties godoc
// @Summary GetAllRentProperties
// @Description GetAllRentProperties
// @Tags RentProperty
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/ [get]
func (ctrl *RentPropertyController) GetAllRentProperties(c *gin.Context) {
	rentProperties, err := ctrl.storedb.GetAllRentProperties(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rentProperties)
}

// GetRentPropertyByID godoc
// @Summary GetRentPropertyByID
// @Description GetRentPropertyByID
// @Tags RentProperty
// @Accept json
// @Produce json
// @Param id path int true "Rent Property ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/{id} [get]
func (ctrl *RentPropertyController) GetRentPropertyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent property ID"})
		return
	}

	rentProperty, err := ctrl.storedb.GetRentPropertyByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rentProperty)
}

// UpdateRentProperty godoc
// @Summary UpdateRentProperty
// @Description UpdateRentProperty
// @Tags RentProperty
// @Accept json
// @Produce json
// @Param id path int true "Rent Property ID"
// @Param property body models.UpdateRentPropertyRequest true "Update Rent Property Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/{id} [put]
func (ctrl *RentPropertyController) UpdateRentProperty(c *gin.Context) {
	var payload models.UpdateRentPropertyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rentPropertyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent property ID"})
		return
	}

	repoDesc := &payload.RepoDesc
	repoCateID := &payload.RepoCateID

	params := db.UpdateRentPropertyParams{
		RepoID:     int32(rentPropertyID),
		RepoName:   payload.RepoName,
		RepoDesc:   repoDesc,
		RepoPrice:  payload.RepoPrice,
		RepoCateID: repoCateID,
	}

	err = ctrl.storedb.UpdateRentProperty(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rent property updated successfully"})
}

// DeleteRentProperty godoc
// @Summary DeleteRentProperty
// @Description DeleteRentProperty
// @Tags RentProperty
// @Accept json
// @Produce json
// @Param id path int true "Rent Property ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/{id} [delete]
func (ctrl *RentPropertyController) DeleteRentProperty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent property ID"})
		return
	}

	err = ctrl.storedb.DeleteRentProperty(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rent property deleted"})
}
