package authcontroller

import (
	"fmt"
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/session"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
)

type AuthFormPayload struct {
	Login    string `json:"login" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=5,max=20"`
}

func RegisterNewUser(ctx *gin.Context) {
	var payload AuthFormPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	user, err := database.AddNewUser(payload.Login, payload.Password)
	if err != nil {
		var code int
		var message string
		if database.IsUniqueViolationError(err) {
			code = http.StatusBadRequest
			message = fmt.Sprintf("Пользователь %s уже зарегистрирован", payload.Login)
		} else {
			code = http.StatusBadRequest
			message = "Не удалось зарегистрировтаь пользователя"
		}

		ctx.JSON(code, models.ServerErrorResponse{
			Error: message,
		})

		return
	}

	tokens := session.CreateSessionTokens(user)
	ctx.JSON(http.StatusOK, tokens)
}

func Login(ctx *gin.Context) {
	var payload AuthFormPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	user, err := database.GetExistingUserByLogin(payload.Login)
	if err != nil {
		var code int
		var message string
		if database.IsErrNoRows(err) {
			code = http.StatusUnauthorized
			message = "Неверный логин или пароль"
		} else {
			code = http.StatusInternalServerError
			message = "Ошибка при обращении к базе данных"
		}

		ctx.JSON(code, models.ServerErrorResponse{
			Error: message,
		})
		return
	}

	if user.HashedPass != tools.CreateSha512Hash(payload.Password) {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: "Неверный логин или пароль",
		})
		return
	}

	user.UpdateTokenIDs()
	database.UpdateTokenIDs(user)
	tokens := session.CreateSessionTokens(user)
	ctx.JSON(http.StatusOK, tokens)
}

func Logout(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	deleteSessionErr := database.DeleteUserSessionInfo(user)
	if deleteSessionErr != nil {
		ctx.JSON(http.StatusInternalServerError, models.ServerErrorResponse{
			Error: deleteSessionErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true})
}
