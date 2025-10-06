package testcasegroupcontroller

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xuri/excelize/v2"
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
	GroupUUID string `uri:"groupUUID" validate:"required,uuid"`
}

type GetTestCaseGroupContentResponse struct {
	TestCases []models.TestCaseFormatted `json:"test_cases"`
}

func GetTestCaseGroupContents(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

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

type DeleteGroupPayload struct {
	UUID string `json:"uuid" validate:"required,uuid"`
}

func DeleteTestCaseGroup(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload DeleteGroupPayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	errDel := database.DeleteTestCaseGroup(payload.UUID, user)
	if errDel != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: errDel.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.ServerBlankOk{Ok: true})
}

type UpdateTestCaseNamePayload struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
	UUID string `json:"uuid" validate:"required,uuid"`
}

func UpdateTestCaseName(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	var payload UpdateTestCaseNamePayload
	if !tools.HandleRequestBodyParsing(ctx, &payload) {
		return
	}

	if !tools.HadleRequestBodyValidation(ctx, &payload) {
		return
	}

	res, err := database.RenameTestCaseGroup(payload.UUID, payload.Name, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

type FilesResponse struct {
	Files []*models.FileFormatted `json:"files"`
}

func UploadFiles(ctx *gin.Context) {
	user, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

	payload, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	testCaseGroupUUID := payload.Value["group"]

	if len(testCaseGroupUUID) != 1 {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: "Должен присутствовать 1 параметр group",
		})
		return
	}

	if len(payload.File["files[]"]) < 1 {
		ctx.JSON(http.StatusBadRequest, &models.ServerErrorResponse{
			Error: "Должен присутствовать массив files[]",
		})
		return
	}

	fileList := []*models.File{}

	for _, file := range payload.File["files[]"] {
		f, err := file.Open()
		if err != nil {
			fileList = append(fileList, nil)
			break
		}

		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			fileList = append(fileList, nil)
			break
		}

		c := string(content)

		fileList = append(fileList, &models.File{
			Name:          file.Filename,
			Content:       c,
			TestCaseGroup: testCaseGroupUUID[0],
			CreatorUUID:   user.UUID,
			UUID:          uuid.New().String(),
			TokenCount:    tools.CountTokens(c),
			Status:        models.StatusNone,
		})
	}

	insertErr := database.SaveFiles(fileList)
	if insertErr != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: insertErr.Error(),
		})
		return
	}

	fileListResult := make([]*models.FileFormatted, len(fileList))

	for i, f := range fileList {
		fileListResult[i] = &models.FileFormatted{
			File:   *f,
			Status: f.Status.Name(),
		}
	}

	ctx.JSON(http.StatusOK, &FilesResponse{
		Files: fileListResult,
	})
}

func GetGroupUploads(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

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

	files, err := database.GetGroupFiles(payload.GroupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusOK, &FilesResponse{Files: []*models.FileFormatted{}})
		} else {
			ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
				Error: err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, &FilesResponse{Files: *files})
}

type ExportResponse struct {
	Content string `json:"content"`
}

func ExportExcel(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

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

	content, err := database.GetFullGroupContent(payload.GroupUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	xlsx := excelize.NewFile()
	sheet := "Sheet1"

	headers := []string{
		/* A */ "№",
		/* B */ "Название",
		/* C */ "Описание",
		/* D */ "Предусловиe",

		/* E */ "Номер шага",
		/* F */ "Описание действий",
		/* G */ "Данные",
		/* H */ "Ожидаемый результат",

		/* I */ "Постусловие",
	}

	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		xlsx.SetCellValue(sheet, cell, h)
	}

	pos := 2
	for n, tc := range content {
		if len(tc.Steps) < 1 {
			xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", pos), n+1)
			xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", pos), tc.Name)
			xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", pos), tc.Description)
			xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", pos), tc.PreCondition)
			xlsx.SetCellValue(sheet, fmt.Sprintf("I%d", pos), tc.PostCondition)
			pos++
		} else {
			for i, step := range tc.Steps {
				if i == 0 {
					xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", pos), n+1)
					xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", pos), tc.Name)
					xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", pos), tc.Description)
					xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", pos), tc.PreCondition)
					xlsx.SetCellValue(sheet, fmt.Sprintf("I%d", pos), tc.PostCondition)
				}

				xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", pos), step.Num)
				xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", pos), step.Description)
				xlsx.SetCellValue(sheet, fmt.Sprintf("G%d", pos), step.Data)
				xlsx.SetCellValue(sheet, fmt.Sprintf("H%d", pos), step.ExpectedResult)
				pos++
			}
		}
	}

	var buff bytes.Buffer
	if err := xlsx.Write(&buff); err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	xlsxB64 := base64.StdEncoding.EncodeToString(buff.Bytes())

	ctx.JSON(http.StatusOK, &ExportResponse{
		Content: xlsxB64,
	})
}

func ExportZephyr(ctx *gin.Context) {
	_, err := tools.GetUserFromRequestCtx(ctx)
	if err != nil {
		return
	}

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

	content, err := database.GetFullGroupContent(payload.GroupUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.ServerErrorResponse{
			Error: err.Error(),
		})

		return
	}

	jsonB64 := base64.StdEncoding.EncodeToString(bytes)

	ctx.JSON(http.StatusOK, &ExportResponse{
		Content: jsonB64,
	})
}
