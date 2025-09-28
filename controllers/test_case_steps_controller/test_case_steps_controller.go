package testcasestepscontroller

import (
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type CreateStepPayload struct {
	TestCase string `json:"test_case" validate:"required,uuid"`
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

type UpdateStepPayload struct {
	UUID           string `json:"uuid" validate:"required,uuid"`
	Data           string `json:"data"`
	ExpectedResult string `json:"expected_result"`
	Description    string `json:"description"`
}

func UpdateStep(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload UpdateStepPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	tcs := &models.TestCaseStep{
		UUID:           payload.UUID,
		Data:           tools.MakeSqlNullString(payload.Data),
		Description:    tools.MakeSqlNullString(payload.Description),
		ExpectedResult: tools.MakeSqlNullString(payload.ExpectedResult),
	}

	result, err := database.UpdateStep(tcs, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type SwapStepsPayload struct {
	First  string `json:"first" validate:"required,uuid"`
	Second string `json:"second" validate:"required,uuid"`
}

func SwapSteps(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload SwapStepsPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	errSwap := database.SwapSteps(payload.First, payload.Second, user)
	if errSwap != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: errSwap.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}

type DeleteStepPayload struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

func DeleteStep(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload DeleteStepPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	errDel := database.DeteteStep(payload.UUID, user)
	if errDel != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: errDel.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}

