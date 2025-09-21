package testcasegroupcontroller

import (
	"fmt"
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

func GetTestCaseGroups(ctx *gin.Context) {
	cases := database.GetTestCaseGroups()
	ctx.JSON(http.StatusOK, &cases)
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

	fmt.Printf("user add test case: %+v\n", user)
	tcg, err := database.AddTestCaseGroup(payload.Name, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, tcg.GetRequestPayloadPassedCreator(user.Login))
}
