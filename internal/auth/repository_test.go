package auth

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestRepo(t *testing.T) (Repository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewRepository(db), mock
}

func TestRepo_GetUserByUsername_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role"}).
		AddRow(1, "alice", "password123", "customer")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, role FROM users WHERE username = ?")).
		WithArgs("alice").
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Username != "alice" || user.Password != "password123" || user.Role != "customer" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestRepo_GetUserByUsername_NotFound(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, role FROM users WHERE username = ?")).
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByUsername("unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}

func TestRepo_GetUserByUsername_DBError(t *testing.T) {
	repo, mock := newTestRepo(t)

	wantErr := errors.New("connection refused")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password, role FROM users WHERE username = ?")).
		WithArgs("alice").
		WillReturnError(wantErr)

	_, err := repo.GetUserByUsername("alice")
	if !errors.Is(err, wantErr) {
		t.Errorf("got %v, want %v", err, wantErr)
	}
}
