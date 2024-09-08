package controller

import (
	db "mamikost/db/sqlc"
	"mamikost/models"
	"mamikost/services"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RentPropertiesImagesController struct {
	storedb services.Store
}

type UploadImageDto struct {
	FrimFilename *SingleFileUpload //`form:"frim_filename" binding:"required"`
	FrimDefault  string            `form:"frim_default" binding:"required"`
	FrimRepoID   int32             `form:"frim_repo_id" binding:"required"`
}

type SingleFileUpload struct {
	FrimFilename *multipart.FileHeader `form:"frim_filename" binding:"required"`
}

func NewRentPropertiesImagesController(store services.Store) *RentPropertiesImagesController {
	return &RentPropertiesImagesController{
		storedb: store,
	}
}

// UploadImage godoc
// @Summary UploadImage
// @Description UploadImage
// @Tags RentPropertiesImages
// @Accept multipart/form-data
// @Produce json
// @Param frim_filename formData file true "Upload image file"
// @Param frim_default formData boolean false "Default Image"
// @Param frim_repo_id formData int32 false "Rent property ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 422 {} http.StatusUnprocessableEntity
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/image/ [post]
func (ctrl *RentPropertiesImagesController) UploadImage(c *gin.Context) {
	var payload UploadImageDto
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(err))
		return
	}

	fileUpload, err := c.FormFile("frim_filename")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	// Retrieve file information
	extension := filepath.Ext(fileUpload.Filename)
	// Generate random file name for the new uploaded file
	newFileName := uuid.New().String() + extension

	// Save the uploaded file
	if err := c.SaveUploadedFile(fileUpload, "./public/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	params := db.UploadImageParams{
		FrimFilename: newFileName,
		FrimDefault:  &payload.FrimDefault,
		FrimRepoID:   &payload.FrimRepoID,
	}

	imageID, err := ctrl.storedb.UploadImage(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"image_id": imageID})
}

// GetAllImages godoc
// @Summary GetAllImages
// @Description GetAllImages
// @Tags RentPropertiesImages
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/image/ [get]
func (ctrl *RentPropertiesImagesController) GetAllImages(c *gin.Context) {
	images, err := ctrl.storedb.GetAllImages(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

// GetImageByID godoc
// @Summary GetImageByID
// @Description GetImageByID
// @Tags RentPropertiesImages
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/image/{id} [get]
func (ctrl *RentPropertiesImagesController) GetImageByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := ctrl.storedb.GetImageByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// DeleteImage godoc
// @Summary DeleteImage
// @Description DeleteImage
// @Tags RentPropertiesImages
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {} http.StatusBadRequest
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /repro/image/{id} [delete]
func (ctrl *RentPropertiesImagesController) DeleteImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	err = ctrl.storedb.DeleteImageByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}
