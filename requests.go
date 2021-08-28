package main

import (
	"omsec.com/services/controllers"
	"omsec.com/services/middleware"
)

func registerControllers() {
	router.Use(middleware.CORSMiddleware())

	router.GET("/", controllers.SayHello)
}
