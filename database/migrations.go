package database

func Migrate() {
	db := GetDB()
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		login TEXT UNIQUE NOT NULL,
		hashed_password TEXT NOT NULL,
		access_token_id TEXT,
		refresh_token_id TEXT
	);
	`)

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS test_case_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		status INTEGER NOT NULL,
		name TEXT NOT NULL,
		creator TEXT NOT NULL,
		FOREIGN KEY (creator) REFERENCES users(uuid)
	);
	`)
}
