package services

import (
	db "Mamikost2/db/sqlc"
	"Mamikost2/middleware"
	"Mamikost2/models"
	"context"
)

func (sm *StoreManager) Signup(ctx context.Context, userReq models.CreateUserRequest) (*models.UserResponse, *models.Error) {
	//1. is user already signup
	foundUser, _ := sm.FindUserByUsername(ctx, userReq.UserName)

	if foundUser.UserID != 0 {
		return &models.UserResponse{}, models.NewError(models.ErrUserAlreadyExist)
	}

	// user not exist
	argCreateUser := db.CreateUserParams{
		UserName:     userReq.UserName,
		UserPassword: userReq.UserPassword,
		UserEmail:    userReq.UserEmail,
		UserPhone:    userReq.UserPhone,
	}
	newUser, err := sm.CreateUser(ctx, argCreateUser)
	if err != nil {
		return &models.UserResponse{}, models.NewError(err)
	}

	response := &models.UserResponse{
		UserID:       newUser.UserID,
		UserPassword: &newUser.UserPassword,
		UserPhone:    &newUser.UserPhone,
	}
	return response, nil
}

func (sm *StoreManager) Signin(ctx context.Context, userReq models.CreateUserRequest) (*models.UserResponse, *models.Error) {
	args := &db.FindUserByUserPasswordParams{
		UserName:     userReq.UserName,
		UserPassword: userReq.UserPassword,
	}

	foundUser, _ := sm.FindUserByUserPassword(ctx, *args)

	if foundUser.UserID == 0 {
		return &models.UserResponse{}, models.NewError(models.ErrInvalidUserPassword)
	}

	accessToken, error := middleware.GenerateJWT(foundUser.UserName)

	if error != nil {
		return &models.UserResponse{},
			models.NewError(models.ErrFailedGenerateToken)
	}

	argsUpdateToken := &db.UpdateTokenParams{
		UserID:    foundUser.UserID,
		UserToken: &accessToken,
	}

	_, err := sm.UpdateToken(ctx, *argsUpdateToken)

	if err != nil {
		return &models.UserResponse{},
			models.NewError(models.ErrFailedGenerateToken)
	}

	response := &models.UserResponse{
		UserID:       foundUser.UserID,
		UserName:     &foundUser.UserName,
		UserPassword: &foundUser.UserPassword,
		UserPhone:    &foundUser.UserPhone,
		UserToken:    &accessToken,
	}

	return response, nil
}

func (sm *StoreManager) Signout(ctx context.Context, accessToken string) *models.Error {
	err := sm.DeleteToken(ctx, &accessToken)
	if err != nil {
		return models.NewError(models.ErrUpdateToken)
	}

	return nil
}
