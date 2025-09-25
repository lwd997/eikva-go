package testcasegroupcontroller

import (
	"fmt"
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type GetTestCaseGroupsResponse struct {
	Groups []models.TestCaseGroupFormatted `json:"groups"`
}

func GetTestCaseGroups(ctx *gin.Context) {
	var response GetTestCaseGroupsResponse
	cases := *database.GetTestCaseGroups()
	if cases != nil {
		response.Groups = cases
	} else {
		response.Groups = []models.TestCaseGroupFormatted{}
	}

	ctx.JSON(http.StatusOK, &response)
}

type AddTestCaseGroupPayload struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

func AddTestCaseGroup(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload AddTestCaseGroupPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	tcg, err := database.AddTestCaseGroup(payload.Name, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &tcg)
}

type GetTestCaseGroupContentsPayload struct {
	GroupUUID string `uri:"groupUUID" binding:"required,uuid"`
}

type GetTestCaseGroupContentResponse struct {
	TestCases []models.TestCaseFormatted `json:"test_cases"`
}

func GetTestCaseGroupContents(ctx *gin.Context) {
	var payload GetTestCaseGroupContentsPayload
	if err := ctx.ShouldBindUri(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	isGroupExists := database.IsTestGroupExisits(payload.GroupUUID)
	if !isGroupExists {
		ctx.JSON(http.StatusNotFound, &models.ServerErrorResponse{
			Error: fmt.Sprintf("Группы %s не существет", payload.GroupUUID),
		})

		return
	}

	var response GetTestCaseGroupContentResponse
	tc := *database.GetTestCaseGroupContents(payload.GroupUUID)
	if tc != nil {
		response.TestCases = tc
	} else {
		response.TestCases = []models.TestCaseFormatted{}
	}

	ctx.JSON(http.StatusOK, &response)
}
