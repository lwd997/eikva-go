package database

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
	sqliteDriver "modernc.org/sqlite" // for error types
	sqliteLib "modernc.org/sqlite/lib"
)

var db *sqlx.DB
var once sync.Once

var Code_ConstaintUnique = sqliteLib.SQLITE_CONSTRAINT_UNIQUE

func initDB(path string) {
	var err error
	db, err = sqlx.Connect("sqlite", path)
	if err != nil {
		panic("")
	}

	db.SetMaxOpenConns(1)
	db.MustExec("PRAGMA journal_mode=WAL;")
	db.MustExec("PRAGMA busy_timeout = 500;")
    db.MustExec("PRAGMA foreign_keys = ON;")
}

func GetDB() *sqlx.DB {
	once.Do(func() {
		initDB("eikva.db")
	})
	return db
}

func IsErrorType(err error, code int) bool {
	var sqliteError *sqliteDriver.Error
	if errors.As(err, &sqliteError) {
		errorCode := sqliteError.Code()
		return errorCode == code
	}

	return false
}

func IsUniqueViolationError(err error) bool {
	return IsErrorType(err, sqliteLib.SQLITE_CONSTRAINT_UNIQUE)
}

func IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}


