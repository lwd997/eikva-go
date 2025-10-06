package database

import (
	"errors"

	"eikva.ru/eikva/models"
)

func UpdateUpload(upl *models.File) error {
	_, err := GetDB().Exec(
		`UPDATE uploads
		SET name = ?, content = ?, token_count = ?, status = ?
		WHERE uuid = ?`,
		upl.Name, upl.Content, upl.TokenCount, upl.Status,
		upl.UUID,
	)

	return err
}

func UpdateUploadStatus(uuid string, status models.Status) error {
	_, err := GetDB().Exec(
		`UPDATE uploads
		SET status = ?
		WHERE uuid = ?`,
		status,
		uuid,
	)

	return err
}


func DeleteUpload(uuid string, user *models.User) error {
	file, err := GetFile(uuid)

	if err != nil {
		return err
	}

	if file.CreatorUUID != user.UUID {
		return errors.New("Можно редактирвоать только свои записи")
	}

	_, deleteErr := GetDB().Exec("DELETE FROM uploads WHERE uuid=?", uuid)

	return deleteErr
}
