package order

import (
	"errors"
	"testing"
	"database/sql"
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


func TestGetOrderByIDRepository(t *testing.T) {
	t.Run("Success - Should return order and items", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock: %s", err)
		}
		defer db.Close()

		repo := NewRepository(db)
		orderRows := sqlmock.NewRows([]string{"id", "customer_username", "restaurant_id", "status", "total_price"}).
			AddRow(1, "testuser", 10, "pending", 500)
		
		mock.ExpectQuery("SELECT id, customer_username, restaurant_id, status, total_price FROM orders WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(orderRows)

		itemRows := sqlmock.NewRows([]string{"id", "order_id", "food_item_id", "quantity", "subtotal"}).
			AddRow(1, 1, 101, 2, 200).
			AddRow(2, 1, 102, 1, 300)

		mock.ExpectQuery("SELECT id, order_id, food_item_id, quantity, subtotal FROM order_items WHERE order_id = \\?").
			WithArgs(1).
			WillReturnRows(itemRows)

		order, items, err := repo.GetOrderByID(1)

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if order == nil || order.ID != 1 {
			t.Error("Expected order with ID 1")
		}
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Failure - Order Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock: %s", err)
		}
		defer db.Close()

		repo := NewRepository(db)

		mock.ExpectQuery("SELECT id, customer_username, restaurant_id, status, total_price FROM orders WHERE id = \\?").
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		_, _, err = repo.GetOrderByID(99)

		if err == nil {
			t.Error("Expected error for non-existent order, got nil")
		}
	})
	
	t.Run("Failure - Query Items Error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewRepository(db)

		// ให้ Query แรกผ่าน
		orderRows := sqlmock.NewRows([]string{"id", "customer_username", "restaurant_id", "status", "total_price"}).
			AddRow(1, "testuser", 10, "pending", 500)
		mock.ExpectQuery("SELECT id, customer_username, restaurant_id, status, total_price FROM orders WHERE id = \\?").
			WithArgs(1).WillReturnRows(orderRows)

		// แต่ให้ Query ที่สอง (Items) พัง
		mock.ExpectQuery("SELECT id, order_id, food_item_id, quantity, subtotal FROM order_items WHERE order_id = \\?").
			WithArgs(1).WillReturnError(errors.New("query items failed"))

		_, _, err := repo.GetOrderByID(1)
		if err == nil {
			t.Error("Expected error from query items, got nil")
		}
	})

	t.Run("Failure - Scan Items Error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewRepository(db)

		mock.ExpectQuery("SELECT id, customer_username, restaurant_id, status, total_price FROM orders WHERE id = \\?").
			WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "customer_username", "restaurant_id", "status", "total_price"}).AddRow(1, "testuser", 10, "pending", 500))

		// จำลองข้อมูลที่ส่งกลับมาไม่ตรงกับจำนวน Column ที่ Scan (ทำให้ Scan พัง)
		itemRows := sqlmock.NewRows([]string{"id"}).AddRow(1) // ส่งมาแค่ Column เดียว แต่โค้ดเรา Scan 5 ตัว
		mock.ExpectQuery("SELECT id, order_id, food_item_id, quantity, subtotal FROM order_items WHERE order_id = \\?").
			WithArgs(1).WillReturnRows(itemRows)

		_, _, err := repo.GetOrderByID(1)
		if err == nil {
			t.Error("Expected scan error, got nil")
		}
	})
}