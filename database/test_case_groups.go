package database

import (
	"database/sql"
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
		`UPDATE test_case_groups SET name = ?
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

func SaveFiles(fileList []*models.File) error {
	dbInst := GetDB()
	tx, err := dbInst.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(
		`INSERT INTO uploads (uuid, name, content, token_count, creator, test_case_group)
        VALUES (?, ?, ?, ?, ?, ?)`,
	)

	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, file := range fileList {
		_, err := stmt.Exec(
			file.UUID,
			file.Name,
			file.Content,
			file.TokenCount,
			file.CreatorUUID,
			file.TestCaseGroup,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetGroupFiles(testCaseGroupUUID string) (*[]*models.File, error) {
	var result []*models.File
	err := GetDB().Select(
		&result,
		`SELECT
			id,
			uuid,
			name,
			token_count,
			creator as creator
		FROM uploads
		WHERE test_case_group = ?`,
		testCaseGroupUUID,
	)

	if err != nil {
		return nil, err
	}

	if result == nil {
		result = []*models.File{}
	}

	return &result, nil
}

func GetFile(uuid string) (*models.File, error) {
	var result models.File
	err := GetDB().Get(
		&result,
		`SELECT
			id,
			uuid,
			name,
			content,
			token_count,
			creator as creator
		FROM uploads
		WHERE uuid = ?`,
		uuid,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetFullGroupContent(uuid string) ([]models.ExportTestCaseWithSteps, error) {
	type row struct {
		models.ExportTestCase
		StepNum            sql.NullInt64  `db:"step_num"`
		StepDesc           sql.NullString `db:"step_description"`
		StepData           sql.NullString `db:"step_data"`
		StepExpectedResult sql.NullString `db:"step_expected_result"`

		TCUUID string `db:"tc_uuid"`
	}

	var rows []row
	err := db.Select(&rows, `SELECT
		t.uuid as tc_uuid,
		t.name,
		t.description,
		t.pre_condition,
		t.post_condition,

		s.num AS step_num,
		s.description AS step_description,
		s.data as step_data,
		s.expected_result as step_expected_result

		FROM test_cases t
		LEFT JOIN test_case_steps s ON s.test_case = t.uuid
		WHERE t.test_case_group = ? AND t.status = ?
		ORDER BY t.id, s.num;`,
		uuid, models.StatusNone,
	)

	if err != nil {
		return nil, err
	}

	cases := map[string]*models.ExportTestCaseWithSteps{}
	for _, r := range rows {
		tc, ok := cases[r.TCUUID]
		if !ok {
			tc = &models.ExportTestCaseWithSteps{
				Steps: []models.ExportTestCaseStepFormatted{},
				ExportTestCaseFormatted: models.ExportTestCaseFormatted{
					Name:          r.Name.String,
					PreCondition:  r.PreCondition.String,
					PostCondition: r.PostCondition.String,
					Description:   r.Description.String,
				},
			}

			cases[r.TCUUID] = tc
		}

		if r.StepNum.Valid {
			tc.Steps = append(tc.Steps, models.ExportTestCaseStepFormatted{
				Num:            int(r.StepNum.Int64),
				Description:    r.StepDesc.String,
				ExpectedResult: r.StepExpectedResult.String,
				Data:           r.StepData.String,
			})
		}

	}

	result := []models.ExportTestCaseWithSteps{}
	for _, tc := range cases {
		result = append(result, *tc)
	}

	return result, nil
}
