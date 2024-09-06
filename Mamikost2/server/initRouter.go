package server

import (
	"Mamikost2/controller"
	"Mamikost2/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CreateRouter(handlers *controller.ControllerManager, mode string) *gin.Engine {
	var router *gin.Engine
	if mode == "test" {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	} else {
		router = gin.Default()
	}

	//router := gin.Default()
	//set a lower memory limit for multipart forms
	router.MaxMultipartMemory = 8 << 20 //8 Mib
	router.Static("/static", "./public")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://google.com"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	// config.AllowAllOrigins = true

	router.Use(cors.New(config))

	api := router.Group("/api")

	api.GET("/home", func(ctx *gin.Context) {
		ctx.String(200, "Hello Gin FB")
	})

	categoryRoute := api.Group("/category")
	{
		categoryRoute.Use(middleware.AuthMiddleware())
		categoryRoute.POST("/", handlers.CreateCategory)
		categoryRoute.GET("/:id", handlers.GetCategoryById)
		categoryRoute.PUT("/:id", handlers.UpdateCategory)

	}

	reproRoute := api.Group("/repro")
	{
		reproRoute.Use(middleware.AuthMiddleware())
		reproRoute.POST("/", handlers.CreateRentProperty)
		reproRoute.GET("/", handlers.GetAllRentProperties)
		reproRoute.GET("/:id", handlers.GetRentPropertyByID)
		reproRoute.PUT("/:id", handlers.UpdateRentProperty)
		reproRoute.DELETE("/:id", handlers.DeleteCategory)
	}

	imageRoute := api.Group("/repro/image")
	{
		imageRoute.Use(middleware.AuthMiddleware())
		imageRoute.POST("/", handlers.UploadImage)
		imageRoute.GET("/", handlers.GetAllImages)
		imageRoute.DELETE("/:id", handlers.DeleteImage)
	}

	orderRoute := api.Group("/order")
	{
		orderRoute.Use(middleware.AuthMiddleware())
		orderRoute.POST("/", handlers.CreateOrder)
		orderRoute.GET("/:id", handlers.GetOrderByID)
		orderRoute.DELETE("/:id", handlers.DeleteOrder)

	}

	detailRoute := api.Group("/order/detail")
	{
		detailRoute.Use(middleware.AuthMiddleware())
		detailRoute.POST("/", handlers.AddOrderDetail)
		detailRoute.DELETE("/:id", handlers.RemoveItemFromOrder)
	}

	cartRoute := api.Group("/cart")
	{
		cartRoute.Use(middleware.AuthMiddleware())
		cartRoute.GET("/:id", handlers.GetCartByUserID)
		cartRoute.DELETE("/:id", handlers.DeleteCart)
		cartRoute.POST("/", handlers.AddToCart)
	}

	userRoute := api.Group("/user")
	{
		userRoute.POST("/signup", handlers.Signup)
		userRoute.POST("/signin", handlers.Sigin)
		userRoute.POST("/signout", handlers.Signout)
		userRoute.GET("/profile", handlers.GetProfile)
	}

	return router
}
