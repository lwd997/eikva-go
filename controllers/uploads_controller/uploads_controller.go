package uploadscontroller

import (
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type GetSingleUploadPayload struct {
	UUID string `uri:"uuid" validate:"required,uuid"`
}

func GetSingleUpload(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload GetSingleUploadPayload
	if err := ctx.ShouldBindUri(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	upload, err := database.GetFile(payload.UUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &upload)
}
