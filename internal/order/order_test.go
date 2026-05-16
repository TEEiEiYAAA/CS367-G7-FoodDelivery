package order

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// สร้าง mockRepository ขึ้นมาเพื่อจำลองพฤติกรรมของฐานข้อมูล
type mockRepository struct {
	Repository
	err error // เราจะใช้ตัวแปรนี้กำหนดว่าอยากให้ Repo คืนค่า error หรือไม่
}

// จำลองฟังก์ชัน AssignRider
func (m *mockRepository) AssignRider(orderID string, riderID int) error {
	return m.err
}

// ต้องประกาศฟังก์ชันอื่นๆ ให้ครบตาม Interface (แม้จะไม่ได้ใช้ในเทสนี้)
func (m *mockRepository) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	return 0, 0, m.err
}
func (m *mockRepository) CancelOrder(username string, orderID int) error { return nil }
func (m *mockRepository) GetOrderByID(orderID int) (*Order, []OrderItem, error) {
	if m.err != nil {
		return nil, nil, m.err
	}
	// จำลองข้อมูลสมมติส่งกลับไป
	mockOrder := &Order{ID: orderID, Status: "pending", TotalPrice: 500}
	mockItems := []OrderItem{{ID: 1, OrderID: orderID, FoodItemID: 10, Quantity: 2}}

	return mockOrder, mockItems, nil
}

func (m *mockRepository) UpdateOrderStatus(orderID int, newStatus string, role string) error {
	return m.err
}

// mockService สำหรับเทส Handler
type mockService struct {
	Service
	orderID    int64
	totalPrice int
	err        error
}

func (m *mockService) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	return m.orderID, m.totalPrice, m.err
}

func (m *mockService) AssignRider(orderID string, riderID int) error {
	return m.err
}

func (m *mockService) UpdateOrderStatus(orderID int, newStatus string, role string) error {
	return m.err
}

func TestAssignRider(t *testing.T) {
	// Case 1: มอบหมายไรเดอร์สำเร็จ (Happy Path)
	t.Run("Success - Should return nil when repo success", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil} // จำลองว่า DB ทำงานปกติ
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	// Case 2: เกิด Error จาก Database (Bad Path)
	t.Run("Failure - Should return error when repo fails", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("database connection failed")} // จำลอง DB พัง
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreateOrderRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	username := "testuser"
	req := CreateOrderRequest{
		RestaurantID:    1,
		DeliveryAddress: "Test Address",
		Items: []OrderItemRequest{
			{FoodItemID: 101, Quantity: 2},
		},
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()

		// 1. Mock Check Price & Available
		rows := sqlmock.NewRows([]string{"price", "is_available"}).AddRow(50, true)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT price, is_available FROM food_items WHERE id = ? AND restaurant_id = ?")).
			WithArgs(101, 1).
			WillReturnRows(rows)

		// 2. Mock Insert Orders
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO orders")).
			WithArgs(username, 1, 100, req.DeliveryAddress, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 3. Mock Insert Order Items
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO order_items")).
			WithArgs(int64(1), 101, 2, 100).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		orderID, totalPrice, err := repo.CreateOrder(username, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if orderID != 1 {
			t.Errorf("expected orderID 1, got %d", orderID)
		}
		if totalPrice != 100 {
			t.Errorf("expected total price 100, got %d", totalPrice)
		}
	})

	t.Run("FoodItemNotFound", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT price, is_available FROM food_items")).
			WithArgs(101, 1).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, _, err := repo.CreateOrder(username, req)

		if err == nil || err.Error() != "food item 101 not found in restaurant 1" {
			t.Errorf("expected specific not found error, got %v", err)
		}
	})

	t.Run("FoodItemNotAvailable", func(t *testing.T) {
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"price", "is_available"}).AddRow(50, false)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT price, is_available FROM food_items")).
			WithArgs(101, 1).
			WillReturnRows(rows)
		mock.ExpectRollback()

		_, _, err := repo.CreateOrder(username, req)

		if err == nil || err.Error() != "food item 101 is not available" {
			t.Errorf("expected not available error, got %v", err)
		}
	})
}

func TestCreateOrderService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		username := "testuser"
		req := CreateOrderRequest{
			RestaurantID: 1,
			Items:        []OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
		}

		_, _, err := service.CreateOrder(username, req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repoErr := errors.New("database error")
		mockRepo := &mockRepository{err: repoErr}
		service := NewService(mockRepo)

		_, _, err := service.CreateOrder("user", CreateOrderRequest{})
		if err != repoErr {
			t.Errorf("expected repo error, got %v", err)
		}
	})
}

func TestGetOrderByID(t *testing.T) {
	t.Run("Success - Should return order detail", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		order, items, err := service.GetOrderByID(1)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		// ตรวจสอบว่า order ไม่เป็น nil ก่อนเช็ก ID (กันโปรแกรมแครช)
		if order == nil {
			t.Fatal("Expected order object, got nil")
		}

		if order.ID != 1 {
			t.Errorf("Expected ID 1, got %d", order.ID)
		}

		// ตรวจสอบตัวแปร items เพื่อให้คอมไพเลอร์ยอมให้ผ่าน
		if len(items) == 0 {
			t.Error("Expected items, but got empty list")
		}
	})
}

func TestCreateOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := &mockService{orderID: 1, totalPrice: 100}
		handler := NewHandler(mockSvc)

		r := gin.Default()
		r.POST("/order", func(c *gin.Context) {
			c.Set("username", "testuser")
		}, handler.CreateOrder)

		reqBody, _ := json.Marshal(CreateOrderRequest{
			RestaurantID:    1,
			DeliveryAddress: "Bangkok",
			Items:           []OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
		})

		req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("NoUserInContext", func(t *testing.T) {
		mockSvc := &mockService{}
		handler := NewHandler(mockSvc)

		r := gin.Default()
		r.POST("/order", handler.CreateOrder) // ไม่ได้ set username

		reqBody, _ := json.Marshal(CreateOrderRequest{})
		req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := NewHandler(&mockService{})
		r := gin.Default()
		r.POST("/order", func(c *gin.Context) {
			c.Set("username", "testuser")
		}, handler.CreateOrder)

		req, _ := http.NewRequest("POST", "/order", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

