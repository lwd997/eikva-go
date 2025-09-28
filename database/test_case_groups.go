package database

import (
	"errors"

	"eikva.ru/eikva/models"
	"github.com/google/uuid"
)

func AddTestCaseGroup(name string, user *models.User) (*models.TestCaseGroupFormatted, error) {
	if user.UUID == "" {
		return nil, errors.New("placeholer error: no uuid in passed user")
	}

	tcg := &models.TestCaseGroup{
		Name:        name,
		UUID:        uuid.New().String(),
		Status:      models.StatusNone,
		Creator:     user.Login,
		CreatorUUID: user.UUID,
	}

	res, err := GetDB().Exec(
		`INSERT INTO test_case_groups (uuid, status, name, creator)
		VALUES (?, ?, ?, ?)`,
		tcg.UUID, tcg.Status, tcg.Name, tcg.CreatorUUID,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	tcg.ID = int(id)

	formatted := &models.TestCaseGroupFormatted{
		TestCaseGroup: *tcg,
		Status:        tcg.Status.Name(),
	}

	return formatted, nil
}

func GetTestCaseGroups() *[]models.TestCaseGroupFormatted {
	var result []models.TestCaseGroupFormatted
	err := GetDB().Select(
		&result,
		`SELECT
			test_case_groups.id,
			test_case_groups.uuid,
			test_case_groups.name,
			test_case_groups.status,
			test_case_groups.creator as creator_uuid,
			users.login AS creator
		FROM test_case_groups
		JOIN users ON test_case_groups.creator = users.uuid`,
	)

	if err != nil {
		panic(err)
	}

	for i, entry := range result {
		result[i].Status = entry.TestCaseGroup.Status.Name()
	}

	return &result
}

func IsTestGroupExisits(groupUUID string) bool {
	var isGroupExists bool

	err := GetDB().Get(
		&isGroupExists,
		"SELECT EXISTS(SELECT 1 FROM test_case_groups where uuid = ?)",
		groupUUID,
	)

	if err != nil {
		panic(err)
	}

	return isGroupExists
}

func GetTestCaseGroupContents(groupUUID string) *[]models.TestCaseFormatted {
	var selectResult []models.TestCase

	err := GetDB().Select(
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
		JOIN users ON test_cases.creator = users.uuid where test_case_group = ?`,
		groupUUID,
	)

	if err != nil {
		panic(err)
	}

	result := make([]models.TestCaseFormatted, len(selectResult))

	for i, entry := range selectResult {
		result[i].TestCase = entry
		result[i].Status = entry.Status.Name()
		result[i].Name = entry.Name.String
		result[i].PreCondition = entry.PreCondition.String
		result[i].PostCondition = entry.PostCondition.String
		result[i].Description = entry.Description.String
		result[i].SourceRef = entry.SourceRef.String
	}

	return &result
}

func DeleteTestCaseGroup(uuid string, user *models.User) error {
	dbInst := GetDB()
	var tcgToDelete models.TestCaseGroup
	err := dbInst.Get(
		&tcgToDelete,
		`SELECT
			uuid,
			creator AS creator_uuid
		FROM test_case_groups where uuid = ?`,
		uuid,
	)

	if err != nil {
		return err
	}

	if tcgToDelete.CreatorUUID != user.UUID {
		return errors.New("Можно редактирвоать только свои записи")
	}

	_, errDel := dbInst.Exec("DELETE FROM test_case_groups WHERE uuid=?", tcgToDelete.UUID)
	if errDel != nil {
		return errDel
	}

	return nil
}

func RenameTestCaseGroup(
	uuid string,
	name string,
	user *models.User,
) (*models.TestCaseGroupFormatted, error) {
	dbInst := GetDB()
	var tcg models.TestCaseGroup
	err := dbInst.Get(
		&tcg,
		`SELECT
			test_case_groups.id,
			test_case_groups.uuid,
			test_case_groups.status,
			test_case_groups.name,
			users.login as creator,
			test_case_groups.creator AS creator_uuid
		FROM test_case_groups
		JOIN users ON test_case_groups.creator = users.uuid
		where test_case_groups.uuid = ?`,
		uuid,
	)

	if err != nil {
		return nil, err
	}

	if tcg.CreatorUUID != user.UUID {
		return nil, errors.New("Можно редактирвоать только свои записи")
	}

	_, errUpd := dbInst.Exec(
		`UPdate test_case_groups set name = ?
		where uuid = ?`,
		name,
		tcg.UUID,
	)

	if errUpd != nil {
		return nil, errUpd
	}

	tcg.Name = name
	formatted := &models.TestCaseGroupFormatted{
		TestCaseGroup: tcg,
		Status:        tcg.Status.Name(),
	}

	return formatted, nil
}
