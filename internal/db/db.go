package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// NewDB:
// - SQLite に接続
// - 必要なテーブルが無ければ作成
func NewDB() (*sql.DB, error) {
	// todo.db はプロジェクト直下に作成される
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		due_date TEXT,
		completed INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}
