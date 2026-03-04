package controller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todo-app/internal/model"
)

type TodoServiceIF interface {
	GetTodos() ([]model.Todo, error)
	AddTodo(title, due string) (int64, error)
	DeleteTodo(id string) error
	ToggleTodo(id, completed string) error
	UpdateTitle(id, title string) error
	UpdateDate(id, due string) error
}

type TodoController struct {
	Service TodoServiceIF
	Tmpl    *template.Template
}

func (c *TodoController) Index(w http.ResponseWriter, r *http.Request) {
	// 一覧表示なのでGET以外は弾く(誤操作/直アクセス対策)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	todos, err := c.Service.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 期限切れ(完了済みでも赤表示にしたいのでCompletedは見ない)
	today := time.Now().Truncate(24 * time.Hour)
	for i := range todos {
		if todos[i].DueDate.Before(today) {
			todos[i].Expired = true
		}
	}

	if err := c.Tmpl.Execute(w, todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Add: 追加(POST/add)
func (c *TodoController) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// POST以外で来たら一覧に戻す
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	due := strings.TrimSpace(r.FormValue("due"))

	if title == "" || due == "" {
		// バリデーションは最小限(画面側requiredもある)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newID, err := c.Service.AddTodo(title, due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[ADD][OK] id=%d title=%q due=%s", newID, title, due)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete: 削除(POST/delete)
func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := c.Service.DeleteTodo(idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[DELETE][OK] id=%d", id)
	w.WriteHeader(http.StatusOK)
}

// Toggle: 完了/未完了切り替え(POST/toggle)
func (c *TodoController) Toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimSpace(r.FormValue("id"))
	completedStr := strings.TrimSpace(r.FormValue("completed"))

	if idStr == "" || completedStr == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	completedBool := completedStr == "1"

	if err := c.Service.ToggleTodo(idStr, completedStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[TOGGLE][OK] id=%d completed=%t", id, completedBool)
	w.WriteHeader(http.StatusOK)
}

// UpdateTitle: タイトル更新(POST/update-title)
func (c *TodoController) UpdateTitle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimSpace(r.FormValue("id"))
	title := strings.TrimSpace(r.FormValue("title"))

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || title == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := c.Service.UpdateTitle(idStr, title); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[UPDATE_TITLE][OK] id=%d title=%q", id, title)
	w.WriteHeader(http.StatusOK)
}

// UpdateDate: 期限日更新(POST/update-date)
func (c *TodoController) UpdateDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimSpace(r.FormValue("id"))
	due := strings.TrimSpace(r.FormValue("due"))

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || due == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := c.Service.UpdateDate(idStr, due); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[UPDATE_DATE][OK] id=%d due=%s", id, due)
	w.WriteHeader(http.StatusOK)
}
