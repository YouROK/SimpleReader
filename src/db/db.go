package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sqlite *sql.DB
}

func Open(name string) *DB {
	db := DB{}
	sqlite, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil
	}
	db.sqlite = sqlite
	return &db
}

func (db *DB) CreateTables() error {
	sql := `
	CREATE TABLE IF NOT EXISTS BookTable(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		FileName,
		FromHost,
		UserName text);
	`
	_, _ = db.sqlite.Exec(sql)
}
