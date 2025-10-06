package uploadscontroller

import (
	"net/http"
	"regexp"

	"eikva.ru/eikva/ai"
	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"eikva.ru/eikva/ws"
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

	ctx.JSON(http.StatusOK, &models.FileFormatted{
		File:   *upload,
		Status: upload.Status.Name(),
	})
}

type UploadActionPayload struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

func CompressUpload(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload UploadActionPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	upload, err := database.GetFile(payload.UUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	tokenCount := tools.CountTokens(upload.Content)
	tokenCountTreshold := 20000

	if tokenCount >= tokenCountTreshold {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: "Превышен максимально допустимый размер запроса. Попробуйте разбить запрос на несколько частей.",
		})
		return
	}

	database.UpdateUploadStatus(upload.UUID, models.StatusLoading)
	ws.WSConntections.BroadCastUploadUpdate(upload.UUID)

	go func() {
		compressed, err := ai.StartTextCompression(&upload.Content)
		if err != nil {
			database.UpdateUploadStatus(upload.UUID, models.StatusError)
		} else {
			r := regexp.MustCompile(`<think>(.|\n|\r)*?</think>`)
			upload.Content = r.ReplaceAllString(*compressed, "")
			upload.TokenCount = tools.CountTokens(upload.Content)
			database.UpdateUpload(upload)
		}

		ws.WSConntections.BroadCastUploadUpdate(upload.UUID)
	}()

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}

func DeleteUpload(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload UploadActionPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	deleteErr := database.DeleteUpload(payload.UUID, user)
	if deleteErr != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: deleteErr.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}
