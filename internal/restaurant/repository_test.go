package restaurant

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

func TestRepo_CreateRestaurant_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO restaurants`)).
		WithArgs("KFC", "Bangkok", "alice").
		WillReturnResult(sqlmock.NewResult(7, 1))

	in := Restaurant{Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"}
	got, err := repo.CreateRestaurant(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("ID = %d, want 7", got.ID)
	}
	if got.Name != "KFC" || got.Address != "Bangkok" || got.OwnerUsername != "alice" {
		t.Errorf("unexpected result: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestRepo_CreateRestaurant_InsertError(t *testing.T) {
	repo, mock := newTestRepo(t)

	wantErr := errors.New("connection refused")
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO restaurants`)).
		WithArgs("KFC", "Bangkok", "alice").
		WillReturnError(wantErr)

	_, err := repo.CreateRestaurant(Restaurant{Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"})
	if !errors.Is(err, wantErr) {
		t.Errorf("got %v, want %v", err, wantErr)
	}
}

func TestRepo_CreateRestaurant_LastInsertIdError(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO restaurants`)).
		WithArgs("KFC", "Bangkok", "alice").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("boom")))

	_, err := repo.CreateRestaurant(Restaurant{Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// GetRestaurants

func TestRepo_GetRestaurants_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "owner_username"}).
		AddRow(1, "KFC", "Bangkok", "alice").
		AddRow(2, "MK", "Nonthaburi", "bob")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, address, owner_username FROM restaurants`)).
		WillReturnRows(rows)

	got, err := repo.GetRestaurants()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "KFC" || got[1].Name != "MK" {
		t.Errorf("unexpected rows: %+v", got)
	}
}

func TestRepo_GetRestaurants_EmptyIsNotNil(t *testing.T) {
	repo, mock := newTestRepo(t)

	emptyRows := sqlmock.NewRows([]string{"id", "name", "address", "owner_username"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(emptyRows)

	got, err := repo.GetRestaurants()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("got nil, expected empty non-nil slice")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestRepo_GetRestaurants_QueryError(t *testing.T) {
	repo, mock := newTestRepo(t)

	wantErr := errors.New("connection refused")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(wantErr)

	_, err := repo.GetRestaurants()
	if !errors.Is(err, wantErr) {
		t.Errorf("got %v, want %v", err, wantErr)
	}
}

func TestRepo_GetRestaurants_ScanError(t *testing.T) {
	repo, mock := newTestRepo(t)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "owner_username"}).
		AddRow("not-an-int", "KFC", "Bangkok", "alice")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)

	_, err := repo.GetRestaurants()
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
}

func TestRepo_GetRestaurants_RowsError(t *testing.T) {
	repo, mock := newTestRepo(t)

	wantErr := errors.New("row iteration failed")
	rows := sqlmock.NewRows([]string{"id", "name", "address", "owner_username"}).
		AddRow(1, "KFC", "Bangkok", "alice").
		RowError(0, wantErr)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)

	_, err := repo.GetRestaurants()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// GetRestaurantByID
func TestRepo_GetRestaurantByID_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	rows := sqlmock.NewRows([]string{"id", "name", "address", "owner_username"}).
		AddRow(1, "KFC", "Bangkok", "alice")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, address, owner_username FROM restaurants WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(rows)

	got, err := repo.GetRestaurantByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 || got.Name != "KFC" || got.Address != "Bangkok" || got.OwnerUsername != "alice" {
		t.Errorf("unexpected: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestRepo_GetRestaurantByID_NotFound(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, address, owner_username FROM restaurants WHERE id = ?")).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetRestaurantByID(99)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("got %v, want sql.ErrNoRows", err)
	}
}

// ConfirmOrder
func TestRepo_ConfirmOrder_Success(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(`UPDATE orders`).
		WithArgs(1, "alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.ConfirmOrder(1, "alice"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}

func TestRepo_ConfirmOrder_OrderNotFound(t *testing.T) {
	repo, mock := newTestRepo(t)

	mock.ExpectExec(`UPDATE orders`).
		WithArgs(99, "alice").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.ConfirmOrder(99, "alice")
	if !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("got %v, want ErrOrderNotFound", err)
	}
}

func TestRepo_ConfirmOrder_ExecError(t *testing.T) {
	repo, mock := newTestRepo(t)

	wantErr := errors.New("connection refused")
	mock.ExpectExec(`UPDATE orders`).
		WithArgs(1, "alice").
		WillReturnError(wantErr)

	err := repo.ConfirmOrder(1, "alice")
	if !errors.Is(err, wantErr) {
		t.Errorf("got %v, want %v", err, wantErr)
	}
}
