package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mamikost/config"
	"mamikost/controller"
	"mamikost/server"
	"mamikost/services"

	"github.com/spf13/viper"
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

	server.AutoMigrate(config)

	store := services.NewStoreManager(dbConnection)

	handlerCtrl := controller.NewControllerManager(store)

	router := server.CreateRouter(handlerCtrl, "dev")

	log.Println("Initializig HTTP sever")
	httpServer := server.NewHttpServer(&config, store, router)

	httpServer.MountSwaggerHandlers()
	httpServer.Start()

	// add graceful shutdown
	srv := &http.Server{
		Addr:    viper.GetString("http.server_address"),
		Handler: httpServer.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}

func getConfigFileName() string {
	env := os.Getenv("ENV")
	if env != "" {
		return "Mamikost-" + env
	}

	return "Mamikost"
}
