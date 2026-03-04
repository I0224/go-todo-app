package view

import "todo-app/internal/model"

type IndexView struct {
	Todos []model.Todo
	Today string
}
