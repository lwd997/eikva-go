package main

import (
	"eikva.ru/eikva/database"
	"eikva.ru/eikva/middlewares"
	"eikva.ru/eikva/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(middlewares.PaincRecovery)
	routes.InitRoutes(router)

	database.Migrate()
	router.Run(":3000")
}
