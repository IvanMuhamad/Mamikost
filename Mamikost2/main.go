package main

import (
	"log"
	"os"

	"Mamikost2/config"
	"Mamikost2/controller"
	"Mamikost2/server"
	"Mamikost2/services"
)

// @title Mamikost
// @version 1.0
// @Description Lorem Ipsum
// @termsOfService http://swagger.io/terms/

// @contact.name ivanmuhammad977@gmail.com
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @basepath /api/

// @securityDefinitions.basic BasicAuth

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api

func main() {
	log.Println("Starting Mamikost App")
	log.Println("Initializing configuration")

	config := config.LoadConfig(getConfigFileName(), ".")
	log.Println("Initializing database")
	dbConnection := server.InitDatabase(&config)
	defer server.Close(dbConnection)

	store := services.NewStoreManager(dbConnection)

	handlerCtrl := controller.NewControllerManager(store)

	router := server.CreateRouter(handlerCtrl, "dev")

	log.Println("Initializig HTTP sever")
	httpServer := server.NewHttpServer(&config, store, router)

	//httpServer.MountSwaggerHandlers()
	httpServer.Start()

}

func getConfigFileName() string {
	env := os.Getenv("ENV")
	if env != "" {
		return "Mamikost-" + env
	}

	return "Mamikost"
}
