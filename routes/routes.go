package routes

import (
	authcontroller "eikva.ru/eikva/controllers/auth_controller"
	testcasecontroller "eikva.ru/eikva/controllers/test_case_controller"
	testcasegroupcontroller "eikva.ru/eikva/controllers/test_case_group_controller"
	testcasestepscontroller "eikva.ru/eikva/controllers/test_case_steps_controller"
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

	groups := protected.Group("/groups")
	{
		groups.GET("/get", testcasegroupcontroller.GetTestCaseGroups)
		groups.POST("/add", testcasegroupcontroller.AddTestCaseGroup)
		groups.GET("/get-test-cases/:groupUUID", testcasegroupcontroller.GetTestCaseGroupContents)
	}

	testCases := protected.Group("/test-cases")
	{
		testCases.POST("/add", testcasecontroller.CreateTestCase)
		testCases.GET("/get-steps/:testCaseUUID", testcasecontroller.GetTestCaseSteps)
	}


	steps := protected.Group("/steps")
	{
		steps.POST("/add", testcasestepscontroller.CreateEmptyStep)
	}
}
