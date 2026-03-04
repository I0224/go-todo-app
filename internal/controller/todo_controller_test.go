package controller

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"todo-app/internal/model"
)

type fakeSvc struct {
	// GetTodos
	todos []model.Todo
	err   error

	// AddTodo
	addCalled bool
	addTitle  string
	addDue    string
	addID     int64
	addErr    error
}

func (f *fakeSvc) GetTodos() ([]model.Todo, error) { return f.todos, f.err }
func (f *fakeSvc) AddTodo(title, due string) (int64, error) {
	f.addCalled = true
	f.addTitle = title
	f.addDue = due
	return f.addID, f.addErr
}
func (f *fakeSvc) DeleteTodo(id string) error            { return nil }
func (f *fakeSvc) ToggleTodo(id, completed string) error { return nil }
func (f *fakeSvc) UpdateTitle(id, title string) error    { return nil }
func (f *fakeSvc) UpdateDate(id, due string) error       { return nil }

func TestIndex_OK_RendersTitles(t *testing.T) {
	// Arrange
	tmpl := template.Must(template.New("index").Parse(`{{range .}}{{.Title}}{{end}}`))

	svc := &fakeSvc{
		todos: []model.Todo{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}},
	}
	ctrl := &TodoController{Service: svc, Tmpl: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	ctrl.Index(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "AB" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestIndex_SetsExpiredFlag(t *testing.T) {
	// Arrange
	tmpl := template.Must(template.New("index").Parse(`{{range .}}{{if .Expired}}E{{else}}N{{end}}{{end}}`))

	today := time.Now().Truncate(24 * time.Hour)
	svc := &fakeSvc{
		todos: []model.Todo{
			{ID: 1, Title: "Past", DueDate: today.Add(-24 * time.Hour)},
			{ID: 2, Title: "Future", DueDate: today.Add(24 * time.Hour)},
		},
	}
	ctrl := &TodoController{Service: svc, Tmpl: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	ctrl.Index(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	// 1件目は期限切れ(E)、2件目は期限内(N)
	if rr.Body.String() != "EN" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestIndex_MethodNotAllowed(t *testing.T) {
	// Arrange
	ctrl := &TodoController{
		Service: &fakeSvc{},
		Tmpl:    template.Must(template.New("index").Parse(`ok`)),
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	ctrl.Index(rr, req)

	// Assert
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestAdd_OK_RedirectsAndCallsService(t *testing.T) {
	// Arrange
	svc := &fakeSvc{addID: 123}
	ctrl := &TodoController{Service: svc, Tmpl: template.Must(template.New("index").Parse(`ok`))}

	form := "title=hello&due=2099-01-01"
	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	// Act
	ctrl.Add(rr, req)

	// Assert
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/" {
		t.Fatalf("expected redirect to /, got %q", rr.Header().Get("Location"))
	}
	if !svc.addCalled {
		t.Fatal("expected AddTodo called")
	}
	if svc.addTitle != "hello" || svc.addDue != "2099-01-01" {
		t.Fatalf("unexpected args: title=%q due=%q", svc.addTitle, svc.addDue)
	}
}
