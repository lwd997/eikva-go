package database

import (
	"eikva.ru/eikva/models"
)

func CountStepsOfTestCase(testCaseUUID string) (*int, error) {
	var count int
	err := GetDB().Get(
		&count,
		`SELECT COUNT(*) FROM
			test_case_steps
		where test_case = ?`,
		testCaseUUID,
	)

	if err != nil {
		return nil, err
	}

	return &count, nil
}

func CreateEmptyStep(testCaseUUID string, user *models.User) (*models.TestCaseStepFormatted, error) {
	count, err := CountStepsOfTestCase(testCaseUUID)
	if err != nil {
		return nil, err
	}

	step := &models.TestCaseStep{
		Status:   models.StatusNone,
		Creator:  user.UUID,
		TestCase: testCaseUUID,
		Num:      *count + 1,
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
		TestCaseStep: *step,
		Status:       step.Status.Name(),
		Creator:      user.Login,
	}

	return formatted, nil
}
