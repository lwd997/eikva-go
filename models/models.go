package models

import (
	"database/sql"

	"github.com/google/uuid"
)

/* === User === */

type User struct {
	ID             int            `db:"id" json:"id"`
	UUID           string         `db:"uuid" json:"uuid"`
	Login          string         `db:"login" json:"login"`
	HashedPass     string         `db:"hashed_password" json:"hashed_password"`
	AccessTokenID  sql.NullString `db:"access_token_id" json:"access_token_id"`
	RefreshTokenID sql.NullString `db:"refresh_token_id" json:"refresh_token_id"`
}

func (u *User) UpdateTokenIDs() {
	u.AccessTokenID = sql.NullString{String: uuid.New().String(), Valid: true}
	u.RefreshTokenID = sql.NullString{String: uuid.New().String(), Valid: true}
}

/* === HTTP === */

type ServerErrorResponse struct {
	Error string `json:"error"`
}

type RequestError struct {
	Code    int
	Message string
}

func (e RequestError) Error() string {
	return e.Message
}

/* === Test Case Groups === */

type Status int

const (
	StatusNone Status = iota
	StatusLoading
)

func (tcgs Status) Name() string {
	switch tcgs {
	case StatusNone:
		return "none"
	case StatusLoading:
		return "loading"
	default:
		return string(tcgs)
	}
}

type TestCaseGroup struct {
	ID      int    `db:"id" json:"id"`
	UUID    string `db:"uuid" json:"uuid"`
	Status  Status `db:"status" json:"status"`
	Name    string `db:"name" json:"name"`
	Creator string `db:"creator" json:"creator"`
}

type TestCaseGroupFormatted struct {
	TestCaseGroup `db:",inline" json:",inline"`
	Status        string `db:"status" json:"status"`
	Creator       string `db:"creator" json:"creator"`
}

/* === Test Cases === */

type TestCase struct {
	ID            int            `db:"id" json:"id"`
	UUID          string         `db:"uuid" json:"uuid"`
	Status        Status         `db:"status" json:"status"`
	CreatedAt     string         `db:"created_at" json:"created_at"`
	Name          sql.NullString `db:"name" json:"name"`
	PreCondition  sql.NullString `db:"pre_condition" json:"pre_condition"`
	PostCondition sql.NullString `db:"post_condition" json:"post_condition"`
	Description   sql.NullString `db:"description" json:"description"`
	SorceRef      sql.NullString `db:"source_ref" json:"source_ref"`
	Creator       string         `db:"creator" json:"creator"`
	CreatorUUID   string         `db:"creator_uuid" json:"creator_uuid"`
	TestCaseGroup string         `db:"test_case_group" json:"test_case_group"`
}

type TestCaseFormatted struct {
	TestCase      `db:",inline" json:",inline"`
	Status        string `db:"status" json:"status"`
	Creator       string `db:"creator" json:"creator"`
	Name          string `db:"name" json:"name"`
	PreCondition  string `db:"pre_condition" json:"pre_condition"`
	PostCondition string `db:"post_condition" json:"post_condition"`
	Description   string `db:"description" json:"description"`
	SorceRef      string `db:"source_ref" json:"source_ref"`
}

func (tc *TestCase) UpdateUUID() {
	tc.UUID = uuid.New().String()
}

/* === Steps === */

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
	Creator        string `db:"creator" json:"creator"`
	Data           string `db:"data" json:"data"`
	Description    string `db:"description" json:"description"`
	ExpectedResult string `db:"expected_result" json:"expected_result"`
}

func (tcs *TestCaseStep) UpdateUUID() {
	tcs.UUID = uuid.New().String()
}
