package routes

import (
	authcontroller "eikva.ru/eikva/controllers/auth_controller"
	testcasecontroller "eikva.ru/eikva/controllers/test_case_controller"
	testcasegroupcontroller "eikva.ru/eikva/controllers/test_case_group_controller"
	testcasestepscontroller "eikva.ru/eikva/controllers/test_case_steps_controller"
	"eikva.ru/eikva/middlewares"
	"eikva.ru/eikva/ws"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	router.Static("/static", "./test_client")
	router.GET("/ws", ws.HandleSubscribers)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authcontroller.RegisterNewUser)
		auth.POST("/login", authcontroller.Login)
		auth.POST("/update-tokens", authcontroller.UpdateTokens)
		authProtect := auth.Group("/")
		authProtect.Use((middlewares.BearerAuth))
		{
			authProtect.POST("/logout", authcontroller.Logout)
			authProtect.GET("/whoami", authcontroller.WhoAmI)
		}
	}

	protected := router.Group("/")
	protected.Use(middlewares.BearerAuth)

	groups := protected.Group("/groups")
	{
		groups.GET("/get", testcasegroupcontroller.GetTestCaseGroups)
		groups.POST("/add", testcasegroupcontroller.AddTestCaseGroup)
		groups.POST("/delete", testcasegroupcontroller.DeleteTestCaseGroup)
		groups.POST("/rename", testcasegroupcontroller.UpdateTestCaseName)
		groups.GET("/get-test-cases/:groupUUID", testcasegroupcontroller.GetTestCaseGroupContents)
	}

	testCases := protected.Group("/test-cases")
	{
		testCases.POST("/add", testcasecontroller.CreateTestCase)
		testCases.POST("/delete", testcasecontroller.DeleteTestCase)
		testCases.POST("/start-generation", testcasecontroller.StartTestCasesGeneration)
		testCases.POST("/update", testcasecontroller.UpdateTestCase)
		testCases.GET("/get-steps/:testCaseUUID", testcasecontroller.GetTestCaseSteps)
	}

	steps := protected.Group("/steps")
	{
		steps.POST("/add", testcasestepscontroller.CreateEmptyStep)
		steps.POST("/update", testcasestepscontroller.UpdateStep)
		steps.POST("/delete", testcasestepscontroller.DeleteStep)
		steps.POST("/swap", testcasestepscontroller.SwapSteps)
	}
}
