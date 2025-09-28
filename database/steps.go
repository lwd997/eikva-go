package database

import (
	"errors"

	"eikva.ru/eikva/models"
)

func CountStepsOfTestCase(testCaseUUID string) (*int, error) {
	var count int
	err := GetDB().Get(
		&count,
		`SELECT COUNT(*) FROM
			test_case_steps
		WHERE test_case = ?`,
		testCaseUUID,
	)

	if err != nil {
		return nil, err
	}

	return &count, nil
}

func CreateEmptyStep(testCaseUUID string, user *models.User) (*models.TestCaseStepFormatted, error) {
	dbInst := GetDB()

	var uuids models.GenericOwnershipCheckFields
	err := dbInst.Get(
		&uuids,
		`SELECT uuid, creator
		FROM test_cases where uuid = ?`,
		testCaseUUID,
	)

	if err != nil {
		return nil, err
	}

	if uuids.Creator != user.UUID {
		return nil, errors.New("Можно изменять только свои записи")
	}

	count, err := CountStepsOfTestCase(testCaseUUID)
	if err != nil {
		return nil, err
	}

	step := &models.TestCaseStep{
		Status:      models.StatusNone,
		Creator:     user.Login,
		CreatorUUID: user.UUID,
		TestCase:    testCaseUUID,
		Num:         *count + 1,
	}

	step.UpdateUUID()

	res, err := GetDB().Exec(
		`INSERT INTO test_case_steps
			(uuid, status, num, creator, test_case)
		 VALUES
			(?, ?, ?, ?, ?)`,
		step.UUID, step.Status, step.Num, user.UUID, step.TestCase,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	step.ID = int(id)

	formatted := &models.TestCaseStepFormatted{
		TestCaseStep:   *step,
		Status:         step.Status.Name(),
		Data:           step.Data.String,
		ExpectedResult: step.ExpectedResult.String,
		Description:    step.Description.String,
	}

	return formatted, nil
}

func UpdateStep(tcs *models.TestCaseStep, user *models.User) (*models.TestCaseStepFormatted, error) {
	dbInst := GetDB()

	var currentTcs models.TestCaseStep
	err := dbInst.Get(
		&currentTcs,
		`SELECT
			test_case_steps.id,
			test_case_steps.uuid,
			test_case_steps.status,
			test_case_steps.num,
			test_case_steps.created_at,
			test_case_steps.description,
			test_case_steps.data,
			test_case_steps.expected_result,
			test_case_steps.creator as creator_uuid,
			users.login AS creator
		FROM test_case_steps
		JOIN users ON test_case_steps.creator = users.uuid where test_case_steps.uuid = ?`,
		tcs.UUID,
	)

	if err != nil {
		return nil, err
	}

	if currentTcs.CreatorUUID != user.UUID {
		return nil, errors.New("Можно изменять только свои записи")
	}

	res, err := dbInst.Exec(`UPDATE test_case_steps
		SET description=?, expected_result=?, data=?
		WHERE uuid=?`,
		tcs.Description, tcs.ExpectedResult, tcs.Data, tcs.UUID,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	currentTcs.ID = int(id)
	currentTcs.Description = tcs.Description
	currentTcs.ExpectedResult = tcs.ExpectedResult
	currentTcs.Data = tcs.Data

	result := &models.TestCaseStepFormatted{
		TestCaseStep:   currentTcs,
		Status:         currentTcs.Status.Name(),
		Data:           currentTcs.Data.String,
		ExpectedResult: currentTcs.ExpectedResult.String,
		Description:    currentTcs.Description.String,
	}

	return result, nil
}

func SwapSteps(uuidFirst string, uuidSecond string, user *models.User) error {
	dbInst := GetDB()
	var stepsToSwap []models.TestCaseStep
	err := dbInst.Select(
		&stepsToSwap,
		`SELECT
			test_case_steps.id,
			test_case_steps.uuid,
			test_case_steps.num,
			test_case_steps.creator as creator_uuid
		FROM test_case_steps
		where uuid = ? OR uuid = ?`,
		uuidFirst, uuidSecond,
	)

	if err != nil {
		return err
	}

	if len(stepsToSwap) != 2 {
		return errors.New("Не найдены шаги с переданными uuid")
	}

	stepFirst := stepsToSwap[0]
	stepsSecond := stepsToSwap[1]

	if user.UUID != stepFirst.CreatorUUID || user.UUID != stepsSecond.CreatorUUID {
		return errors.New("Можно редактирвоать только свои записи")
	}

	_, errExecOne := dbInst.Exec(
		"UPDATE test_case_steps SET num=? WHERE uuid=?",
		stepFirst.Num, stepsSecond.UUID,
	)

	if errExecOne != nil {
		return errExecOne
	}

	_, errExecTwo := dbInst.Exec(
		"UPDATE test_case_steps SET num=? WHERE uuid=?",
		stepsSecond.Num, stepFirst.UUID,
	)

	if errExecTwo != nil {
		return errExecTwo
	}

	return nil
}

func DeteteStep(uuid string, user *models.User) error {
	dbInst := GetDB()
	var stepToDelete models.TestCaseStep
	err := dbInst.Get(
		&stepToDelete,
		`SELECT
			num,
			test_case,
			creator AS creator_uuid
		FROM test_case_steps where uuid = ?`,
		uuid,
	)

	if err != nil {
		return err
	}

	if stepToDelete.CreatorUUID != user.UUID {
		return errors.New("Можно редактирвоать только свои записи")
	}

	_, errDel := dbInst.Exec("DELETE FROM test_case_steps WHERE uuid=?", uuid)

	if errDel != nil {
		return errDel
	}

	_, numAdjustErr := db.Exec(
		`UPDATE test_case_steps
		SET num = num - 1
		WHERE num > ? AND test_case = ?`,
		stepToDelete.Num,
		stepToDelete.TestCase,
	)

	if numAdjustErr != nil {
		return numAdjustErr
	}

	return nil
}
