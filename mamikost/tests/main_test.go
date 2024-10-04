package tests

import (
	"mamikost/config"
	"mamikost/controller"
	"mamikost/server"
	"mamikost/services"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

/*
var (
	testGinContext   *gin.Context
	mockStore        *mocks.MockStore
	responseRecorder *httptest.ResponseRecorder
	apiController    *controller.ControllerManager
)

func setUpForTest(ctrl *gomock.Controller, method string, url string, data []byte) {

	//1.Create a new mock implementation of MockStore interface
	mockStore = mocks.NewMockStore(ctrl)

	//2.Create a new instance of APIHandler using the mock service
	apiController = controller.NewControllerManager(mockStore)

	//3.Create a new response recorder to capture HTTP responses
	responseRecorder = httptest.NewRecorder()

	//4.Create a new Gin test context with the response recorder
	testGinContext, _ = gin.CreateTestContext(responseRecorder)

	//5. Create a new HTTP request to use in the tests
	testGinContext.Request, _ = http.NewRequest(method, url, bytes.NewBuffer(data))
} */

func newTestHttpServer(t *testing.T, store services.Store) *server.HttpServer {
	config := config.LoadConfig(getConfigFileName(), "../")
	handlerCtrl := controller.NewControllerManager(store)
	router := server.CreateRouter(handlerCtrl, "dev")
	server := server.NewHttpServer(&config, store, router)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func getConfigFileName() string {
	env := os.Getenv("ENV")
	if env != "" {
		return "mamikost-" + env
	}

	return "mamikost"
}
