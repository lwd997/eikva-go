package testcasecontroller

import (
	"fmt"
	"net/http"

	//"eikva.ru/eikva/ai"
	"eikva.ru/eikva/ai"
	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type CreateTestCasePayload struct {
	TestCaseGroup string `json:"test_case_group" validate:"required,uuid"`
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

	tc, err := database.CreateEmptyTestCase(payload.TestCaseGroup, models.StatusNone, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, tc)
}

type DeleteTestCasePayload struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

func DeleteTestCase(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload DeleteTestCasePayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	errDel := database.DeleteTestCase(payload.UUID, user)
	if errDel != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: errDel.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}

type StartTestCasesGenerationPayload struct {
	Amount        int    `json:"amount" validate:"required,min=1,max=10"`
	UserInput     string `json:"user_input" validate:"required"`
	TestCaseGroup string `json:"test_case_group" validate:"required"`
}

type StartTestCasesGenerationResponse struct {
	TestCases *[]models.TestCaseFormatted `json:"test_cases"`
}

func StartTestCasesGeneration(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload StartTestCasesGenerationPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	reslult, err := database.InitTestCasesGeneration(payload.TestCaseGroup, payload.Amount, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	go func() {
		generated, err := ai.StartTestCaseListGeneration(reslult.UUIDList, user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		database.UpdateTestCaseWithModelResponse(reslult.UUIDList, *generated, user)
	}()
	ctx.JSON(http.StatusOK, &StartTestCasesGenerationResponse{TestCases: reslult.TCList})
}

type UpdateTestCasePayload struct {
	UUID          string `json:"uuid" validate:"required,uuid"`
	Name          string `json:"name" validate:"max=200"`
	PreCondition  string `json:"pre_condition"`
	PostCondition string `json:"post_condition"`
	Description   string `json:"description"`
	SourceRef     string `json:"source_ref"`
}

func UpdateTestCase(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload UpdateTestCasePayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	tc := &models.TestCase{
		UUID:          payload.UUID,
		Name:          tools.MakeSqlNullString(payload.Name),
		PreCondition:  tools.MakeSqlNullString(payload.PreCondition),
		PostCondition: tools.MakeSqlNullString(payload.PostCondition),
		Description:   tools.MakeSqlNullString(payload.Description),
		SourceRef:     tools.MakeSqlNullString(payload.SourceRef),
	}

	result, err := database.UpdateTestCase(tc, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type GetTestCaseStepsPayload struct {
	TestCaseUUID string `uri:"testCaseUUID" validate:"required,uuid"`
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
