package tools

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"

	"eikva.ru/eikva/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

var (
	validate   *validator.Validate
	translator ut.Translator
	once       sync.Once
)

func Initvalidate() {
	validate = validator.New()

	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)
	var ok bool
	translator, ok = uni.GetTranslator("ru")
	if !ok {
		panic("")
	}

	if err := ruTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type FormErrorResponse struct {
	FormErrors []FieldError `json:"form_errors"`
}

func HadleRequestBodyValidation(ctx *gin.Context, payload interface{}) bool {
	once.Do(Initvalidate)

	err := validate.Struct(payload)

	if err != nil {
		response := FormErrorResponse{}
		response.FormErrors = []FieldError{}
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			response.FormErrors = append(response.FormErrors, FieldError{
				Field: e.Field(),
				Error: e.Translate(translator),
			})
		}

		ctx.JSON(http.StatusBadRequest, response)
		return false
	}

	return true
}

func CreateSha512Hash(input string) string {
	sum := sha512.Sum512([]byte(input))
	hash := hex.EncodeToString(sum[:])
	return hash
}

func HandleRequestBodyParsing(ctx *gin.Context, payload interface{}) bool {
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: "Некорректный формат запроса",
		})
		return false
	}

	return true
}

func HandleRequestError(ctx *gin.Context, err *models.RequestError) bool {
	if err != nil {
		ctx.JSON(err.Code, models.ServerErrorResponse{
			Error: err.Message,
		})
		return false
	}
	return true
}

func GetUserFromRequestCtx(ctx *gin.Context) (*models.User, error) {
	userField, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.ServerErrorResponse{
			Error: "placeholder error: no user in context",
		})
		return nil, errors.New("placeholder error: no user in context")
	}

	user, ok := userField.(*models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ServerErrorResponse{
			Error: "placeholder error: user in context is not type of User",
		})
		return nil, errors.New("placeholder error: user in context is not type of User")
	}

	return user, nil
}
