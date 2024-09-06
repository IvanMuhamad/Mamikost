package controller

import (
	db "Mamikost2/db/sqlc"
	"Mamikost2/models"
	"Mamikost2/services"
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

func (cate *CategoryController) GetAllCategories(c *gin.Context) {
	categories, err := cate.storedb.GetAllCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (cate *CategoryController) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := cate.storedb.GetCategoryByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

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
		/* 		if apiErr := models.ConvertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(apiErr))
		} */
		c.JSON(http.StatusInternalServerError, models.NewError(err))
		return
	}
	/* 	if category == nil {
		c.JSON(http.StatusNotFound, models.NewError(err))
		return
	} */
	c.JSON(http.StatusOK, category)

}

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
