package httpserver

import (
	"github.com/imylam/delivery-test/db"
	"github.com/imylam/delivery-test/googlemap"
	_orderHandler "github.com/imylam/delivery-test/order/delivery/http"
	_orderRepo "github.com/imylam/delivery-test/order/repository/mysql"
	_orderUsecase "github.com/imylam/delivery-test/order/usecase"

	"github.com/gin-gonic/gin"
)

// InitRoutes creates routes to receive and respond to http requests
func InitRoutes() *gin.Engine {
	mysqlConn := db.GetDBConnection()
	mapClient := googlemap.NewMapClient()

	orderRepo := _orderRepo.NewOrderRepositoryMysql(mysqlConn)
	orderUC := _orderUsecase.NewOrderUsecase(orderRepo, mapClient)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	_orderHandler.NewOrderHandler(router, orderUC)

	return router
}
