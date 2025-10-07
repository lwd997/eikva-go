package main

import (
	"eikva.ru/eikva/database"
	envvars "eikva.ru/eikva/env_vars"
	"eikva.ru/eikva/middlewares"
	"eikva.ru/eikva/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	envvars.Dotenv()
	envvars.Init()

	router := gin.Default()

	router.Use(middlewares.PaincRecovery)
	routes.InitRoutes(router)

	database.Migrate()

	router.Run(":3000")
}
