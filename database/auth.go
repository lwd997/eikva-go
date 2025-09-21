package database

import (
	"database/sql"

	"eikva.ru/eikva/models"
	"eikva.ru/eikva/tools"
	"github.com/google/uuid"
)

func AddNewUser(login string, password string) (*models.User, error) {
	user := models.User{
		Login:      login,
		HashedPass: tools.CreateSha512Hash(password),
		UUID:       uuid.New().String(),
	}
	user.UpdateTokenIDs()

	res, err := GetDB().Exec(
		`
		INSERT INTO users (uuid, login, hashed_password, access_token_id, refresh_token_id)
		VALUES (?, ?, ?, ?, ?)
		`,
		user.UUID, user.Login, user.HashedPass, user.AccessTokenID, user.RefreshTokenID,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	user.ID = int(id)
	return &user, nil
}

func GetExistingUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := GetDB().Get(&user, "SELECT * FROM users WHERE login=?", login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetExistingUserByID(id int) (*models.User, error) {
	var user models.User
	err := GetDB().Get(&user, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUserSessionInfo(user *models.User) error {
	_, err := GetDB().Exec(
		`
		UPDATE users SET access_token_id=NULL, refresh_token_id=NULL
		where uuid=?
		`,
		user.UUID,
	)
	return err
}

func GetExistingUserByUUID(uuid string) (*models.User, error) {
	var user models.User
	safeUUID := sql.NullString{String: uuid, Valid: true}
	err := GetDB().Get(&user, "SELECT * FROM users WHERE uuid=?", safeUUID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateTokenIDs(user *models.User) {
	GetDB().MustExec(`
		UPDATE users SET access_token_id=?,refresh_token_id=?
		WHERE id=?
	`, user.AccessTokenID, user.RefreshTokenID, user.ID)
}
