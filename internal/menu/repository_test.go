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

func TestGetMenu_Repository(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "restaurant_id", "name", "price"}).
			AddRow(1, 10, "Pad Thai", 80.0).
			AddRow(2, 10, "Khao Pad", 60.0)

		mock.ExpectQuery("SELECT id, restaurant_id, name, price FROM food_items").
			WithArgs(10).
			WillReturnRows(rows)

		repo := NewRepository(db)
		got, err := repo.GetMenu(10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(got) != 2 {
			t.Errorf("expected 2 items, got %d", len(got))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, restaurant_id, name, price FROM food_items").
			WithArgs(10).
			WillReturnError(errors.New("query error"))

		repo := NewRepository(db)
		_, err = repo.GetMenu(10)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "restaurant_id", "name", "price"}).
			AddRow("invalid_id", 10, "Pad Thai", 80.0)

		mock.ExpectQuery("SELECT id, restaurant_id, name, price FROM food_items").
			WithArgs(10).
			WillReturnRows(rows)

		repo := NewRepository(db)
		_, err = repo.GetMenu(10)
		if err == nil {
			t.Fatal("expected scan error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("rows_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "restaurant_id", "name", "price"}).
			AddRow(1, 10, "Pad Thai", 80.0).
			RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT id, restaurant_id, name, price FROM food_items").
			WithArgs(10).
			WillReturnRows(rows)

		repo := NewRepository(db)
		_, err = repo.GetMenu(10)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
