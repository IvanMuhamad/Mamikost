package server

import (
	"Mamikost2/config"
	"Mamikost2/docs"
	"Mamikost2/services"
	"fmt"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	Router *gin.Engine
	Store  services.Store
	Config *config.Config
}

func NewHttpServer(config *config.Config, store services.Store, router *gin.Engine) *HttpServer {
	return &HttpServer{
		Config: config,
		Store:  store,
		Router: router,
	}
}
func (hs HttpServer) Start() {
	httpAddr := viper.GetString("http.server_address")
	err := hs.Router.Run(httpAddr)
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}

func (hs HttpServer) MountSwaggerHandlers() {
	host := viper.GetString("http.host")
	httpAddr := viper.GetString("http.server_address")
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s%s", host, httpAddr)
	docs.SwaggerInfo.BasePath = "/api/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Title = "Mamikost"
	docs.SwaggerInfo.Description = "Mamikost API documentation"
	hs.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
