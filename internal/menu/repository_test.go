package menu

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestCreateMenu_Repository(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO food_items").
			WithArgs(10, "Pad Thai", 80.0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(db)
		got, err := repo.CreateMenu(10, CreateMenuRequest{Name: "Pad Thai", Price: 80.0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		want := Menu{ID: 1, RestaurantID: 10, Name: "Pad Thai", Price: 80.0}
		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("exec_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO food_items").
			WithArgs(10, "Pad Thai", 80.0).
			WillReturnError(errors.New("db exec error"))

		repo := NewRepository(db)
		_, err = repo.CreateMenu(10, CreateMenuRequest{Name: "Pad Thai", Price: 80.0})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("last_insert_id_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO food_items").
			WithArgs(10, "Pad Thai", 80.0).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		repo := NewRepository(db)
		_, err = repo.CreateMenu(10, CreateMenuRequest{Name: "Pad Thai", Price: 80.0})
		if err == nil {
			t.Fatal("expected error from LastInsertId, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
