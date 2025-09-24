package testcasegroupcontroller

import (
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
		ctx.JSON(http.StatusInternalServerError, models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &tcg)
}
