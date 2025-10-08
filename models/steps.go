package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type TestCaseStep struct {
	ID             int            `db:"id" json:"id"`
	UUID           string         `db:"uuid" json:"uuid"`
	Status         Status         `db:"status" json:"status"`
	Num            int            `db:"num" json:"num"`
	Description    sql.NullString `db:"description" json:"description"`
	Data           sql.NullString `db:"data" json:"data"`
	ExpectedResult sql.NullString `db:"expected_result" json:"expected_result"`
	Creator        string         `db:"creator" json:"creator"`
	CreatorUUID    string         `db:"creator_uuid" json:"creator_uuid"`
	TestCase       string         `db:"test_case" json:"test_case"`
	CreatedAt      string         `db:"created_at" json:"created_at"`
}

type TestCaseStepFormatted struct {
	TestCaseStep   `db:",inline" json:",inline"`
	Status         string `db:"status" json:"status"`
	Data           string `db:"data" json:"data"`
	Description    string `db:"description" json:"description"`
	ExpectedResult string `db:"expected_result" json:"expected_result"`
}

func (tcs *TestCaseStep) UpdateUUID() {
	tcs.UUID = uuid.New().String()
}
