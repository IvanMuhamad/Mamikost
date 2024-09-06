package controller

import "Mamikost2/services"

type ControllerManager struct {
	*CategoryController
	*UsersController
	*RentPropertyController
	*RentPropertiesImagesController
	*OrderRentPropertyController
	*OrderRentPropertyDetailsController
	*CartController
}

func NewControllerManager(store services.Store) *ControllerManager {
	return &ControllerManager{
		CategoryController:             NewCategoryController(store),
		UsersController:                NewUsersController(store),
		RentPropertyController:         NewRentPropertyController(store),
		RentPropertiesImagesController: NewRentPropertiesImagesController(store),
		OrderRentPropertyController:    NewOrderRentPropertyController(store),
		CartController:                 NewCartController(store),
	}
}
