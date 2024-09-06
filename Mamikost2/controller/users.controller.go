package controller

import (
	db "Mamikost2/db/sqlc"
	"Mamikost2/middleware"
	"Mamikost2/models"
	"Mamikost2/services"
	"context"
	"database/sql"
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

func (uc *UsersController) Sigin(c *gin.Context) {
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

func (uc *UsersController) Signout(c *gin.Context) {
	accessToken := c.Request.Header.Get("Authorization")

	err := uc.storedb.Signout(c, accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

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
