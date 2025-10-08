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
