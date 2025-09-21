package routes

import (
	authcontroller "eikva.ru/eikva/controllers/auth_controller"
	testcasegroupcontroller "eikva.ru/eikva/controllers/test_case_group_controller"
	"eikva.ru/eikva/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	auth := router.Group("/auth")
	{
		auth.POST("/register", authcontroller.RegisterNewUser)
		auth.POST("/login", authcontroller.Login)
		authProtect := auth.Group("/")
		authProtect.Use((middlewares.BearerAuth))
		{
			authProtect.POST("/logout", authcontroller.Logout)
		}
	}

	protected := router.Group("/")
	protected.Use(middlewares.BearerAuth)
	api := protected.Group("/groups")
	{
		api.GET("/get", testcasegroupcontroller.GetTestCaseGroups)
		api.POST("/add", testcasegroupcontroller.AddTestCaseGroup)
	}
}
