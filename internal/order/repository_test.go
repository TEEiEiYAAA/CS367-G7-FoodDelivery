package order

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAssignRiderRepository(t *testing.T) {
	t.Run("Success - Should update rider_id and status", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock: %s", err)
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectExec("UPDATE orders SET rider_id = \\?, status = 'assigned' WHERE id = \\?").
			WithArgs(101, "1").
			WillReturnResult(sqlmock.NewResult(0, 1))


		err = repo.AssignRider("1", 101)


		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Failure - Database Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock: %s", err)
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectExec("UPDATE orders SET rider_id = \\?, status = 'assigned' WHERE id = \\?").
			WithArgs(101, "1").
			WillReturnError(errors.New("db connection lost"))

		err = repo.AssignRider("1", 101)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}