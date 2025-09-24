package main

import (
	"log"
	"os"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/middlewares"
	"eikva.ru/eikva/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatalln("Не указан env JWT_SECRET")
	}

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()

	router.Use(middlewares.PaincRecovery)
	routes.InitRoutes(router)

	database.Migrate()
	router.Run(":3000")
}
