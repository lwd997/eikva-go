package testcasestepscontroller

import (
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type CreateStepPayload struct {
	TestCase string `json:"test_case" validate:"required"`
}

func CreateEmptyStep(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload CreateStepPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	step, err := database.CreateEmptyStep(payload.TestCase, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, step)
}

/*
type CreateTestCasePayload struct {
	UserInput     string `json:"user_input" validate:"required"`
	TestCaseGroup string `json:"test_case_group" validate:"required"`
}

func CreateTestCase(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload CreateTestCasePayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	tc, err := database.CreateEmptyTestCase(payload.TestCaseGroup, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, tc)
}*/
