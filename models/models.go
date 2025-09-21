package models

import (
	"database/sql"

	"github.com/google/uuid"
)

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

type TestCaseGroupStatus int

const (
	TCGStatusNone TestCaseGroupStatus = iota
	TCGStatusLoading
)

func (tcgs TestCaseGroupStatus) Name() string {
	return [...]string{"None", "Loading"}[tcgs]
}

type TestCaseGroup struct {
	ID      int                 `db:"id" json:"id"`
	UUID    string              `db:"uuid" json:"uuid"`
	Status  TestCaseGroupStatus `db:"status" json:"status"`
	Name    string              `db:"name" json:"name"`
	Creator string              `db:"creator" json:"creator"`
}

type TestCaseGroupResponse struct {
	TestCaseGroup	`db:",inline" json:",inline"`
	Status  string  `db:"status" json:"status"`
	Creator  string `db:"creator" json:"creator"`
}

func (tcg TestCaseGroup) GetRequestPayloadPassedCreator(creator string) TestCaseGroupResponse {
	return TestCaseGroupResponse{
		TestCaseGroup: tcg,
		Creator: creator,
		Status: tcg.Status.Name(),
	}
}
/*
func (tcg TestCaseGroup) GetRequestPayloadSelectCrearor() (*TestCaseGroupResponse, error) {
	user, err := database.GetExistingUserByUUID(tcg.Creator)
	if err != nil {
		return nil, err
	}

	return TestCaseGroupResponse{
		TestCaseGroup: tcg,
		Creator: &user.Login,
		Status: tcg.Status.Name(),
	}
}
*/
