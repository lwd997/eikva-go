package database

import (
	"errors"
	"sync"

	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
)

func CreateEmptyTestCase(groupUUID string, name string, status models.Status, user *models.User) (*models.TestCaseFormatted, error) {
	tc := &models.TestCase{
		Status:        status,
		Name:          tools.MakeSqlNullString(name),
		CreatorUUID:   user.UUID,
		Creator:       user.Login,
		TestCaseGroup: groupUUID,
	}

	tc.UpdateUUID()
	dbInst := GetDB()

	res, err := dbInst.Exec(
		`INSERT INTO test_cases
			(uuid, status, creator, test_case_group, name)
		VALUES
			(?, ?, ?, ?, ?)`,
		tc.UUID, tc.Status, tc.CreatorUUID, tc.TestCaseGroup, tc.Name,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	var inserted models.TestCase
	errSelect := dbInst.Get(
		&inserted,
		`SELECT
			test_cases.id,
			test_cases.uuid,
			test_cases.name,
			test_cases.status,
			test_cases.created_at,
			test_cases.pre_condition,
			test_cases.post_condition,
			test_cases.description,
			test_cases.source_ref,
			test_cases.creator as creator_uuid,
			users.login AS creator
		FROM test_cases
		JOIN users ON test_cases.creator = users.uuid where test_cases.id = ?`,
		id,
	)

	if errSelect != nil {
		return nil, err
	}

	formatted := &models.TestCaseFormatted{
		TestCase:      inserted,
		Status:        tc.Status.Name(),
		Name:          inserted.Name.String,
		Description:   inserted.Description.String,
		PreCondition:  inserted.PreCondition.String,
		PostCondition: inserted.PostCondition.String,
		SourceRef:     inserted.SourceRef.String,
	}

	return formatted, nil
}

type InitTestCasesGenerationResult struct {
	UUIDList *[]string
	TCList   *[]models.TestCaseFormatted
}

func InitTestCasesGeneration(groupUUID string, amount int, user *models.User) (*InitTestCasesGenerationResult, error) {
	// TODO: передалать на Beginx
	result := []models.TestCaseFormatted{}
	uuidList := []string{}

	for i := 0; i < amount; i++ {
		tc, err := CreateEmptyTestCase(groupUUID, "", models.StatusLoading, user)
		if err != nil {
			return nil, err
		}
		uuidList = append(uuidList, tc.UUID)
		result = append(result, *tc)
	}

	return &InitTestCasesGenerationResult{UUIDList: &uuidList, TCList: &result}, nil
}

func UpdateTestCase(tc *models.TestCase, user *models.User) (*models.TestCaseFormatted, error) {
	dbInst := GetDB()

	var currentTc models.TestCase
	err := dbInst.Get(
		&currentTc,
		`SELECT
			test_cases.id,
			test_cases.uuid,
			test_cases.name,
			test_cases.status,
			test_cases.created_at,
			test_cases.pre_condition,
			test_cases.post_condition,
			test_cases.description,
			test_cases.source_ref,
			test_cases.creator as creator_uuid,
			users.login AS creator
		FROM test_cases
		JOIN users ON test_cases.creator = users.uuid where test_cases.uuid = ?`,
		tc.UUID,
	)

	if err != nil {
		return nil, err
	}

	if currentTc.CreatorUUID != user.UUID {
		return nil, errors.New("Можно изменять только свои записи")
	}

	res, err := dbInst.Exec(`UPDATE test_cases
		SET name=?, pre_condition=?, post_condition=?, description=?, source_ref=?, status=?
		WHERE uuid=?`,
		tc.Name, tc.PreCondition, tc.PostCondition, tc.Description, tc.SourceRef,
		tc.Status, tc.UUID,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	currentTc.ID = int(id)
	currentTc.Name = tc.Name
	currentTc.PreCondition = tc.PreCondition
	currentTc.PostCondition = tc.PostCondition
	currentTc.Description = tc.Description
	currentTc.SourceRef = tc.SourceRef
	currentTc.Status = tc.Status

	result := &models.TestCaseFormatted{
		TestCase:      currentTc,
		Status:        currentTc.Status.Name(),
		Name:          currentTc.Name.String,
		PreCondition:  currentTc.PreCondition.String,
		PostCondition: currentTc.PostCondition.String,
		Description:   currentTc.Description.String,
		SourceRef:     currentTc.SourceRef.String,
	}

	return result, nil
}

func IsTestCaseExists(testCaseUUID string) bool {
	var isTestCaseExists bool

	err := GetDB().Get(
		&isTestCaseExists,
		"SELECT EXISTS(SELECT 1 FROM test_cases where uuid = ?)",
		testCaseUUID,
	)

	if err != nil {
		panic(err)
	}

	return isTestCaseExists
}

func GetTestCase(testCaseUUID string) (*models.TestCaseFormatted, error) {
	var selectResult models.TestCase
	err := GetDB().Get(
		&selectResult,
		`SELECT
			test_cases.id,
			test_cases.uuid,
			test_cases.name,
			test_cases.status,
			test_cases.created_at,
			test_cases.pre_condition,
			test_cases.post_condition,
			test_cases.description,
			test_cases.source_ref,
			test_cases.test_case_group,
			test_cases.creator as creator_uuid,
			users.login AS creator
		FROM test_cases
		JOIN users ON test_cases.creator = users.uuid where test_cases.uuid = ?`,
		testCaseUUID,
	)

	if err != nil {
		return nil, err
	}

	return &models.TestCaseFormatted{
		TestCase:      selectResult,
		Status:        selectResult.Status.Name(),
		Name:          selectResult.Name.String,
		PreCondition:  selectResult.PreCondition.String,
		PostCondition: selectResult.PostCondition.String,
		Description:   selectResult.Description.String,
		SourceRef:     selectResult.SourceRef.String,
	}, nil
}

func GetTestCaseSteps(testCaseUUID string) *[]models.TestCaseStepFormatted {
	var selectResult []models.TestCaseStep

	err := GetDB().Select(
		&selectResult,
		`SELECT
			test_case_steps.id,
			test_case_steps.uuid,
			test_case_steps.status,
			test_case_steps.num,
			test_case_steps.created_at,
			test_case_steps.description,
			test_case_steps.data,
			test_case_steps.expected_result,
			test_case_steps.test_case,
			test_case_steps.creator as creator_uuid,
			users.login AS creator
		FROM test_case_steps
		JOIN users ON test_case_steps.creator = users.uuid where test_case = ?
		ORDER BY num ASC`,
		testCaseUUID,
	)

	if err != nil {
		panic(err)
	}

	result := make([]models.TestCaseStepFormatted, len(selectResult))

	for i, entry := range selectResult {
		result[i].TestCaseStep = entry
		result[i].Status = entry.Status.Name()
		result[i].Data = entry.Data.String
		result[i].Description = entry.Description.String
		result[i].ExpectedResult = entry.ExpectedResult.String
	}

	return &result
}

func SetTestCaseErrorStatus(uuidList *[]string) error {
	dbInst := GetDB()
	tx, err := dbInst.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(
		`UPDATE test_cases SET status = ?
		WHERE uuid = ?`,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	for _, uuid := range *uuidList {
		_, err := stmt.Exec(
			models.StatusError,
			uuid,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func UpdateTestCaseWithModelResponse(
	uuidList *[]string,
	data []*models.CreateTestCaseOutputEntry,
	user *models.User,
) {
	uuidListLen := len(*uuidList)
	dataLen := len(data)

	if dataLen < uuidListLen {
		empty := make([]*models.CreateTestCaseOutputEntry, uuidListLen-dataLen)
		data = append(data, empty...)
	}

	tcList := []*models.TestCase{}
	stepsList := []*models.TestCaseStep{}

	for i, uuid := range *uuidList {
		entry := data[i]
		tc := &models.TestCase{}
		tc.UUID = uuid

		if entry == nil {
			tc.Status = models.StatusError
		} else {
			tc.Name = tools.MakeSqlNullString(entry.Name)
			tc.Description = tools.MakeSqlNullString(entry.Name)
			tc.PreCondition = tools.MakeSqlNullString(entry.PreCondition)
			tc.PostCondition = tools.MakeSqlNullString(entry.PostCondition)
			tc.SourceRef = tools.MakeSqlNullString(entry.SourceRef)
			tc.Status = models.StatusNone

			for _, step := range entry.Steps {
				s := &models.TestCaseStep{
					CreatorUUID:    user.UUID,
					Data:           tools.MakeSqlNullString(step.Data),
					Description:    tools.MakeSqlNullString(step.Description),
					ExpectedResult: tools.MakeSqlNullString(step.ExpectedResult),
					TestCase:       tc.UUID,
				}
				stepsList = append(stepsList, s)
			}
		}
		tcList = append(tcList, tc)
	}

	var wg sync.WaitGroup

	for _, tc := range tcList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			UpdateTestCase(tc, user)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, step := range stepsList {
			r, err := CreateEmptyStep(step.TestCase, user)
			if err == nil {
				step.UUID = r.UUID
			}

			UpdateStep(step, user)
		}
	}()

	wg.Wait()
}

func DeleteTestCase(uuid string, user *models.User) error {
	dbInst := GetDB()
	var tcToDelete models.TestCase
	err := dbInst.Get(
		&tcToDelete,
		`SELECT
			uuid,
			creator AS creator_uuid
		FROM test_cases where uuid = ?`,
		uuid,
	)

	if err != nil {
		return err
	}

	if tcToDelete.CreatorUUID != user.UUID {
		return errors.New("Можно редактирвоать только свои записи")
	}

	_, errDel := dbInst.Exec("DELETE FROM test_cases WHERE uuid=?", tcToDelete.UUID)
	if errDel != nil {
		return errDel
	}

	return nil
}
