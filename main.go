package main

import (
	"log"

	"github.com/imylam/delivery-test/configs"
	"github.com/imylam/delivery-test/db"
	"github.com/imylam/delivery-test/httpserver"

	"github.com/asaskevich/govalidator"
)

func main() {
	configs.Init()
	db.InitDBConn()
	govalidator.SetFieldsRequiredByDefault(true)

	router := httpserver.InitRoutes()

	port := configs.Get(configs.KeyAppPort)
	log.Printf("Starting server on port %s...", port)
	router.Run(":" + port)
}
