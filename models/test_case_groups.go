package models

type TestCaseGroup struct {
	ID          int    `db:"id" json:"id"`
	UUID        string `db:"uuid" json:"uuid"`
	Status      Status `db:"status" json:"status"`
	Name        string `db:"name" json:"name"`
	Creator     string `db:"creator" json:"creator"`
	CreatorUUID string `db:"creator_uuid" json:"creator_uuid"`
}

type TestCaseGroupFormatted struct {
	TestCaseGroup `db:",inline" json:",inline"`
	Status        string `db:"status" json:"status"`
}
