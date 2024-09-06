package controller

import (
	db "Mamikost2/db/sqlc"
	"Mamikost2/models"
	"Mamikost2/services"
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

func (ctrl *RentPropertyController) GetAllRentProperties(c *gin.Context) {
	rentProperties, err := ctrl.storedb.GetAllRentProperties(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rentProperties)
}

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
