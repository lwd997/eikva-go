package models

import "database/sql"

type ExportTestCase struct {
	Name          sql.NullString `db:"name" json:"name"`
	PreCondition  sql.NullString `db:"pre_condition" json:"pre_condition"`
	PostCondition sql.NullString `db:"post_condition" json:"post_condition"`
	Description   sql.NullString `db:"description" json:"description"`
}

type ExportTestCaseFormatted struct {
	Name          string `db:"name" json:"name"`
	PreCondition  string `db:"pre_condition" json:"pre_condition"`
	PostCondition string `db:"post_condition" json:"post_condition"`
	Description   string `db:"description" json:"description"`
}

type ExportTestCaseStep struct {
	Num            int            `db:"num" json:"num"`
	Description    sql.NullString `db:"description" json:"description"`
	Data           sql.NullString `db:"data" json:"data"`
	ExpectedResult sql.NullString `db:"expected_result" json:"expected_result"`
}

type ExportTestCaseStepFormatted struct {
	Num            int    `db:"num" json:"num"`
	Description    string `db:"description" json:"description"`
	Data           string `db:"data" json:"data"`
	ExpectedResult string `db:"expected_result" json:"expected_result"`
}

type ExportTestCaseWithSteps struct {
	ExportTestCaseFormatted `json:",inline"`
	Steps                       []ExportTestCaseStepFormatted `json:"steps"`
}
