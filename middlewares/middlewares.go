package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"eikva.ru/eikva/models"
	"eikva.ru/eikva/session"
	"github.com/gin-gonic/gin"
)

var PaincRecovery gin.HandlerFunc = gin.CustomRecovery(func(ctx *gin.Context, err any) {
	log.Printf("Panic: %v", err)

	ctx.JSON(http.StatusInternalServerError, models.ServerErrorResponse{
		Error: fmt.Sprintf("%v", err),
	})

	ctx.Abort()
})

var BearerAuth gin.HandlerFunc = func(ctx *gin.Context) {
	headerAuth := ctx.GetHeader("Authorization")
	bearerPrefix := "Bearer "
	if headerAuth == "" || !strings.HasPrefix(headerAuth, bearerPrefix) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ServerErrorResponse{
			Error: "Заголовок Authorization отстуствует или имеет не верное значение",
		})
		return
	}

	token := strings.TrimPrefix(headerAuth, bearerPrefix)
	user, err := session.ValidateSessionTokenAndGetUser(token, session.TokenTypeAccess)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
