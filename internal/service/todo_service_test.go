package service

import (
	"errors"
	"testing"
	"time"

	"todo-app/internal/model"
)

type fakeRepo struct {
	insertCalled bool
	insertTitle  string
	insertDue    string
	insertID     int64
	insertErr    error

	updateDateCalled bool
	updateDateID     string
	updateDateDue    string
	updateDateErr    error
}

func (f *fakeRepo) FindAll() ([]model.Todo, error)             { return nil, nil }
func (f *fakeRepo) Delete(id string) error                     { return nil }
func (f *fakeRepo) UpdateCompleted(id, completed string) error { return nil }
func (f *fakeRepo) UpdateTitle(id, title string) error         { return nil }
func (f *fakeRepo) UpdateDate(id, due string) error {
	f.updateDateCalled = true
	f.updateDateID = id
	f.updateDateDue = due
	return f.updateDateErr
}
func (f *fakeRepo) Insert(title, due string) (int64, error) {
	f.insertCalled = true
	f.insertTitle = title
	f.insertDue = due
	return f.insertID, f.insertErr
}

func TestTodoService_AddTodo_OK(t *testing.T) {
	// Arrange
	repo := &fakeRepo{insertID: 123}
	svc := &TodoService{Repo: repo}

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	// Act
	id, err := svc.AddTodo("  hello  ", tomorrow)

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if id != 123 {
		t.Fatalf("expected id=123, got %d", id)
	}
	if !repo.insertCalled {
		t.Fatal("expected repo.Insert called")
	}
	if repo.insertTitle != "hello" || repo.insertDue != tomorrow {
		t.Fatalf("unexpected args: title=%q due=%q", repo.insertTitle, repo.insertDue)
	}
}

func TestTodoService_AddTodo_RejectsEmptyTitle(t *testing.T) {
	// Arrange
	svc := &TodoService{Repo: &fakeRepo{}}
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	// Act
	_, err := svc.AddTodo("   ", tomorrow)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTodoService_AddTodo_RejectsInvalidDate(t *testing.T) {
	// Arrange
	svc := &TodoService{Repo: &fakeRepo{}}

	// Act
	_, err := svc.AddTodo("hello", "2099/01/01")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTodoService_AddTodo_RejectsPastDate(t *testing.T) {
	// Arrange
	svc := &TodoService{Repo: &fakeRepo{}}
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	// Act
	_, err := svc.AddTodo("hello", yesterday)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTodoService_UpdateDate_AllowsPastDate(t *testing.T) {
	// Arrange
	repo := &fakeRepo{}
	svc := &TodoService{Repo: repo}
	past := "1999-12-31"

	// Act
	err := svc.UpdateDate("1", past)

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.updateDateCalled {
		t.Fatal("expected repo.UpdateDate called")
	}
	if repo.updateDateID != "1" || repo.updateDateDue != past {
		t.Fatalf("unexpected args: id=%q due=%q", repo.updateDateID, repo.updateDateDue)
	}
}

func TestTodoService_AddTodo_PropagatesRepoError(t *testing.T) {
	// Arrange
	repoErr := errors.New("db down")
	repo := &fakeRepo{insertErr: repoErr}
	svc := &TodoService{Repo: repo}
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	// Act
	_, err := svc.AddTodo("hello", tomorrow)

	// Assert
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
