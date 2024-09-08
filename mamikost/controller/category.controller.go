package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	storedb services.Store
}

func NewCategoryController(store services.Store) *CategoryController {
	return &CategoryController{
		storedb: store,
	}
}

// CreateCategory godoc
// @Summary CreateCategory
// @Description CreateCategory
// @Tags Category
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /category/ [post]
func (cate *CategoryController) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := cate.storedb.CreateCategory(c, req.CateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategoryById godoc
// @Summary GetCategoryById
// @Description GetCategoryById
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {} http.StatusNotFound
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /category/{id} [get]
func (cate *CategoryController) GetCategoryById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := models.Nullable(cate.storedb.GetCategoryByID(c, int32(id)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewError(err))
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, models.NewError(models.ErrCategoryNotFound))
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve a list of all available categories.
// @Tags Category
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /category/ [get]
func (cate *CategoryController) GetAllCategories(c *gin.Context) {
	categories, err := cate.storedb.GetAllCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary UpdateCategory
// @Description UpdateCategory
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.UpdateCategoryParams true "Update Category Request"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {} http.StatusUnprocessableEntity
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /category/{id} [put]
func (cate *CategoryController) UpdateCategory(c *gin.Context) {
	var payload *models.UpdateCategoryParams
	cateId, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(err))
		return
	}

	args := &db.UpdateCategoryParams{
		CateID:   int32(cateId),
		CateName: payload.CateName,
	}

	category, err := models.Nullable(cate, cate.storedb.UpdateCategory(c, *args))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewError(err))
		return
	}
	c.JSON(http.StatusOK, category)

}

// DeleteCategory godoc
// @Summary DeleteCategory
// @Description DeleteCategory
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /category/{id} [delete]
func (cate *CategoryController) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = cate.storedb.DeleteCategory(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
