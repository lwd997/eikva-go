package fakeresponse

import (
	"net/http"
	"os"
	"time"

	envvars "eikva.ru/eikva/env_vars"
	"github.com/gin-gonic/gin"
)

func FakeItForMePlease() {
	router := gin.Default()
	router.POST(string(envvars.Get(envvars.OpenAiCompletionsPathName)), func(ctx *gin.Context) {
		time.Sleep(time.Duration(time.Second * 7))
		data, _ := os.ReadFile("sample.json")
		ctx.Data(http.StatusOK, "application/json", data)
	})
	router.Run(":3001")
}
