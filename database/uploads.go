package database

import "eikva.ru/eikva/models"

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
