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

/*
id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		status INTEGER NOT NULL,
		num INTEGER NOT NULL,
		description TEXT,
		data TEXT,
		expected_result TEXT,
		creator TEXT NOT NULL,
		test_case TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator) REFERENCES users(uuid),
		FOREIGN KEY (test_case) REFERENCES test_cases(uuid)

func CreateEmptyTestCase(groupUUID string, user *models.User) (*models.TestCaseFormatted, error) {
	tc := models.TestCase{
		Status:        models.StatusNone,
		Creator:       user.UUID,
		TestCaseGroup: groupUUID,
	}

	tc.UpdateUUID()

	res, err := GetDB().Exec(
		`INSERT INTO test_cases
			(uuid, status, creator, test_case_group)
		VALUES
			(?, ?, ?, ?)`,
		tc.UUID, tc.Status, tc.Creator, tc.TestCaseGroup,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	tc.ID = int(id)

	formatted := &models.TestCaseFormatted{
		TestCase: tc,
		Status:   tc.Status.Name(),
		Creator:  user.Login,
	}

	return formatted, nil
}
*/
