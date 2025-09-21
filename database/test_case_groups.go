package database

import (
	"errors"

	"eikva.ru/eikva/models"
	"github.com/google/uuid"
)

func AddTestCaseGroup(name string, user *models.User) (*models.TestCaseGroup, error) {
	if (user.UUID == "") {
		return nil, errors.New("placeholer error: no uuid in passed user")
	}

	tcg := models.TestCaseGroup{
		Name:   name,
		UUID:   uuid.New().String(),
		Status: models.TCGStatusNone,
		Creator: user.UUID,
	}

	res, err := GetDB().Exec(
		`
		INSERT INTO test_case_groups (uuid, status, name, creator)
		VALUES (?, ?, ?, ?)
		`,
		tcg.UUID, tcg.Status, tcg.Name, tcg.Creator,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	tcg.ID = int(id)

	return &tcg, nil
}

func GetTestCaseGroups() *[]models.TestCaseGroupResponse  {
	var result []models.TestCaseGroupResponse
	err := GetDB().Select(
		&result,
		`
		SELECT
			test_case_groups.id,
			test_case_groups.uuid,
			test_case_groups.name,
			test_case_groups.status,
			users.login AS creator
		FROM test_case_groups
		JOIN users ON test_case_groups.creator = users.uuid
		`,
	)

	if err != nil {
		panic(err)
	}

	return &result
}
