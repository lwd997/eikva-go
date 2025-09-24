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
