package database

func Migrate() {
	db := GetDB()
	db.MustExec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		login TEXT UNIQUE NOT NULL,
		hashed_password TEXT NOT NULL,
		access_token_id TEXT,
		refresh_token_id TEXT
	);`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS test_case_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		status INTEGER NOT NULL,
		name TEXT NOT NULL,
		creator TEXT NOT NULL,
		FOREIGN KEY (creator) REFERENCES users(uuid)
	);`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS test_cases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		status INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		name TEXT,
		pre_condition TEXT,
		post_condition TEXT,
		description TEXT,
		source_ref TEXT,
		creator TEXT NOT NULL,
		test_case_group TEXT NOT NULL,
		FOREIGN KEY (creator) REFERENCES users(uuid),
		FOREIGN KEY (test_case_group) REFERENCES test_case_groups(uuid) ON DELETE CASCADE
	);`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS test_case_steps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		status INTEGER NOT NULL,
		num INTEGER NOT NULL,
		description TEXT,
		data TEXT,
		expected_result TEXT,
		creator TEXT NOT NULL,
		test_case TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator) REFERENCES users(uuid),
		FOREIGN KEY (test_case) REFERENCES test_cases(uuid) ON DELETE CASCADE
	);`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		token_count INTEGER NOT NULL,
		creator TEXT NOT NULL,
		test_case_group TEXT NOT NULL,
		FOREIGN KEY (creator) REFERENCES users(uuid),
		FOREIGN KEY (test_case_group) REFERENCES test_case_groups(uuid) ON DELETE CASCADE
	);`)
}
