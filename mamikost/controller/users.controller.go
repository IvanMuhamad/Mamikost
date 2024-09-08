package controller

import (
	"context"
	"database/sql"
	db "mamikost/db/sqlc"
	"mamikost/middleware"
	"mamikost/models"
	"mamikost/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	storedb services.Store
}

func NewUsersController(store services.Store) *UsersController {
	return &UsersController{
		storedb: store,
	}
}

func (ctrl *UsersController) CreateUser(c *gin.Context) {
	var payload db.CreateUserParams

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := ctrl.storedb.CreateUser(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (ctrl *UsersController) GetUserByUsername(c *gin.Context) {
	var payload db.FindUserByUsernameRow

	if err := c.ShouldBindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.storedb.FindUserByUsername(c, payload.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UsersController) GetUserByPhone(c *gin.Context) {
	var payload db.FindUserByPhoneRow

	if err := c.ShouldBindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.storedb.FindUserByUsername(c, payload.UserPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Signup godoc
// @Summary User signup
// @Description Registers a new user with the provided details.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User signup details"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {} http.StatusUnprocessableEntity
// @Failure 500 {} http.StatusInternalServerError
// @Router /signup [post]
func (uc *UsersController) Signup(c *gin.Context) {
	var payload models.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(err))
		return
	}

	args := &models.CreateUserRequest{
		UserName:     payload.UserName,
		UserPassword: payload.UserPassword,
		UserEmail:    payload.UserEmail,
		UserPhone:    payload.UserPhone,
	}
	user, err := uc.storedb.Signup(c, *args)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Signin godoc
// @Summary Signin
// @Description Signin
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User signin details"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {} http.StatusUnprocessableEntity
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /signin [post]
func (uc *UsersController) Signin(c *gin.Context) {
	var payload models.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.NewValidationError(err))
		return
	}

	args := &models.CreateUserRequest{
		UserName:     payload.UserName,
		UserPassword: payload.UserPassword,
	}
	user, err := uc.storedb.Signin(c, *args)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Signout godoc
// @Summary Signout
// @Description Signout
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer access token"
// @Success 204 "Signout Success"
// @Failure 400 {} http.StatusBadRequest
// @Security Bearer
// @Router /signout [post]
func (uc *UsersController) Signout(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")

	err := uc.storedb.Signout(c, accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProfile godoc
// @Summary GetProfile
// @Description GetProfile
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {} http.StatusInternalServerError
// @Security Bearer
// @Router /profile [get]
func (ctrl *UsersController) GetProfile(c *gin.Context) {
	var userID string

	token := middleware.GetJWTFromHeader(c)
	//token := ""
	if token != "" {
		userID = middleware.GetIDFromToken(token)
	}

	foundUser, err := ctrl.storedb.FindUserByUsername(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	profile := &models.UserResponse{
		UserID:    foundUser.UserID,
		UserName:  &foundUser.UserName,
		UserPhone: &foundUser.UserPhone,
	}

	c.JSON(http.StatusOK, profile)
}
