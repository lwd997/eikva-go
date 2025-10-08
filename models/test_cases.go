package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type TestCase struct {
	ID            int            `db:"id" json:"id"`
	UUID          string         `db:"uuid" json:"uuid"`
	Status        Status         `db:"status" json:"status"`
	CreatedAt     string         `db:"created_at" json:"created_at"`
	Name          sql.NullString `db:"name" json:"name"`
	PreCondition  sql.NullString `db:"pre_condition" json:"pre_condition"`
	PostCondition sql.NullString `db:"post_condition" json:"post_condition"`
	Description   sql.NullString `db:"description" json:"description"`
	SourceRef     sql.NullString `db:"source_ref" json:"source_ref"`
	Creator       string         `db:"creator" json:"creator"`
	CreatorUUID   string         `db:"creator_uuid" json:"creator_uuid"`
	TestCaseGroup string         `db:"test_case_group" json:"test_case_group"`
}

type TestCaseFormatted struct {
	TestCase      `db:",inline" json:",inline"`
	Status        string `db:"status" json:"status"`
	Name          string `db:"name" json:"name"`
	PreCondition  string `db:"pre_condition" json:"pre_condition"`
	PostCondition string `db:"post_condition" json:"post_condition"`
	Description   string `db:"description" json:"description"`
	SourceRef     string `db:"source_ref" json:"source_ref"`
}

func (tc *TestCase) UpdateUUID() {
	tc.UUID = uuid.New().String()
}
