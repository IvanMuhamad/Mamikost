package controller

import (
	db "Mamikost2/db/sqlc"
	"Mamikost2/models"
	"Mamikost2/services"
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

// UploadImage handles the upload of a new image
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

// GetAllImages retrieves all images
func (ctrl *RentPropertiesImagesController) GetAllImages(c *gin.Context) {
	images, err := ctrl.storedb.GetAllImages(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

// GetImageByID retrieves a specific image by ID
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

// DeleteImage deletes a specific image by ID
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
