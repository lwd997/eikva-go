package database

import "eikva.ru/eikva/models"

func CreateEmptyTestCase(groupUUID string, user *models.User) (*models.TestCaseFormatted, error) {
	tc := &models.TestCase{
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
		TestCase: *tc,
		Status:   tc.Status.Name(),
		Creator:  user.Login,
	}

	return formatted, nil
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
			test_case_steps.creator as creator_uuid,
			users.login AS creator
		FROM test_case_steps
		JOIN users ON test_case_steps.creator = users.uuid where test_case = ?`,
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
