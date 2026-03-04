package repository

import (
	"database/sql"
	"time"

	"todo-app/internal/model"
)

// TodoRepositoryはDB(SQLite)とのやりとりを担当する層。
// SQLをここに集約することで、Service/ControllerからDBの詳細を隠す。
type TodoRepository struct {
	DB *sql.DB
}

// FindAll: 全件取得(期限日順)
func (r *TodoRepository) FindAll() ([]model.Todo, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, due_date, completed
		FROM todos
		ORDER BY due_date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]model.Todo, 0)

	for rows.Next() {
		var t model.Todo
		var due sql.NullString
		var completed int

		if err := rows.Scan(&t.ID, &t.Title, &due, &completed); err != nil {
			return nil, err
		}

		t.Completed = completed == 1

		// due_dateがNULLまたは空のときはゼロ値(0001-01-01)にならないようにする
		if due.Valid && due.String != "" {
			if d, err := time.Parse("2006-01-02", due.String); err == nil {
				t.DueDate = d
			}
		}

		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// Insert: 追加
func (r *TodoRepository) Insert(title string, due string) (int64, error) {
	res, err := r.DB.Exec(`
		INSERT INTO todos (title, due_date, completed)
		VALUES (?, ?, 0)
	`, title, due)
	if err != nil {
		return 0, err
	}

	// SQLiteのAUTOINCREMENT IDを取得(go-sqlite3で利用可能)
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Delete: 削除
func (r *TodoRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM todos WHERE id = ?`, id)
	return err
}

// UpdateCompleted: 完了または未完了更新(0/1)
func (r *TodoRepository) UpdateCompleted(id string, completed string) error {
	_, err := r.DB.Exec(
		`UPDATE todos SET completed = ? WHERE id = ?`,
		completed,
		id,
	)
	return err
}

// UpdateTitle: タイトル更新
func (r *TodoRepository) UpdateTitle(id, title string) error {
	_, err := r.DB.Exec(
		`UPDATE todos SET title = ? WHERE id = ?`,
		title,
		id,
	)
	return err
}

// UpdateDate: 期限日更新
func (r *TodoRepository) UpdateDate(id, due string) error {
	_, err := r.DB.Exec(
		`UPDATE todos SET due_date = ? WHERE id = ?`,
		due,
		id,
	)
	return err
}
