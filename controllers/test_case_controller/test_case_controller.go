package testcasecontroller

import (
	"fmt"
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

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
}

type GetTestCaseStepsPayload struct {
	TestCaseUUID string `uri:"testCaseUUID" binding:"required,uuid"`
}

type GetTestCaseStepsResonse struct {
	Steps []models.TestCaseStepFormatted `json:"steps"`
}

func GetTestCaseSteps(ctx *gin.Context) {
	var payload GetTestCaseStepsPayload
	if err := ctx.ShouldBindUri(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	isTestCaseExisits := database.IsTestCaseExists(payload.TestCaseUUID)
	if !isTestCaseExisits {
		ctx.JSON(http.StatusNotFound, &models.ServerErrorResponse{
			Error: fmt.Sprintf("Тест-кейса %s не существет", payload.TestCaseUUID),
		})

		return
	}

	var response GetTestCaseStepsResonse
	tcs := *database.GetTestCaseSteps(payload.TestCaseUUID)
	if tcs != nil {
		response.Steps = tcs
	} else {
		response.Steps = []models.TestCaseStepFormatted{}
	}

	ctx.JSON(http.StatusOK, &response)
}
