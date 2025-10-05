package models

type File struct {
	ID            int    `db:"id" json:"id"`
	UUID          string `db:"uuid" json:"uuid"`
	Name          string `db:"name" json:"name"`
	Content       string `db:"content" json:"content"`
	TokenCount    int    `db:"token_count" json:"token_count"`
	CreatorUUID   string `db:"creator" json:"creator"`
	TestCaseGroup string `db:"test_case_group" json:"test_case_group"`
}
