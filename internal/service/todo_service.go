package service

import (
	"errors"
	"strings"
	"time"

	"todo-app/internal/model"
)

type TodoService struct {
	Repo TodoRepository
}

type TodoRepository interface {
	FindAll() ([]model.Todo, error)
	Insert(title, due string) (int64, error)
	Delete(id string) error
	UpdateCompleted(id, completed string) error
	UpdateTitle(id, title string) error
	UpdateDate(id, due string) error
}

// GetTodos: 一覧取得
func (s *TodoService) GetTodos() ([]model.Todo, error) {
	return s.Repo.FindAll()
}

// AddTodo: 追加
// 追加されたTodoのID(AUTOINCREMENT)を返す。
func (s *TodoService) AddTodo(title string, due string) (int64, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, errors.New("title is required")
	}

	// 画面の <input type="date"> はYYYY-MM-DDを返す。
	dueDate, err := time.Parse("2006-01-02", due)
	if err != nil {
		return 0, errors.New("invalid due date format")
	}

	// 追加時は「過去日を登録できない」ルール。
	// ※編集(UpdateDate)は過去日OKにしたいので、ここだけでチェックする。
	today := time.Now().Truncate(24 * time.Hour)
	if dueDate.Before(today) {
		return 0, errors.New("due date must be today or later")
	}

	return s.Repo.Insert(title, due)
}

// DeleteTodo: 削除
func (s *TodoService) DeleteTodo(id string) error {
	return s.Repo.Delete(id)
}

// ToggleTodo: 完了または未完了切り替え
func (s *TodoService) ToggleTodo(id string, completed string) error {
	return s.Repo.UpdateCompleted(id, completed)
}

// UpdateTitle: タイトル更新
func (s *TodoService) UpdateTitle(id, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("title is required")
	}
	return s.Repo.UpdateTitle(id, title)
}

// UpdateDate: 期限日更新
func (s *TodoService) UpdateDate(id, due string) error {
	if _, err := time.Parse("2006-01-02", due); err != nil {
		return errors.New("invalid due date format")
	}
	return s.Repo.UpdateDate(id, due)
}
