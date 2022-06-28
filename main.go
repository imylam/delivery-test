package main

import (
	"fmt"

	"github.com/imylam/delivery-test/configs"
	"github.com/imylam/delivery-test/db"
	"github.com/imylam/delivery-test/httpserver"
	"github.com/imylam/delivery-test/logger"

	"github.com/asaskevich/govalidator"
)

func main() {
	logger.Init()
	db.InitDBConn()
	govalidator.SetFieldsRequiredByDefault(true)

	router := httpserver.InitRoutes()

	port := configs.Get(configs.KeyAppPort)
	logger.Logger.Info(fmt.Sprintf("Starting server on port %s...", port))
	router.Run(":" + port)
}
