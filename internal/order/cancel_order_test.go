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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// --- Repository Tests ---

func TestGetOrderRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	t.Run("Success", func(t *testing.T) {
		gracePeriod := time.Now().Add(5 * time.Minute)
		rows := sqlmock.NewRows([]string{"id", "customer_username", "restaurant_id", "status", "total_price", "delivery_address", "created_at", "customer_grace_period_end"}).
			AddRow(1, "testuser", 1, "pending", 100, "Bangkok", time.Now(), gracePeriod)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, customer_username, restaurant_id, status, total_price, delivery_address, created_at, customer_grace_period_end")).
			WithArgs(1).
			WillReturnRows(rows)

		order, err := repo.GetOrder(1)
		if err != nil || order == nil || order.ID != 1 {
			t.Errorf("expected order 1, got %v %v", order, err)
		}
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, customer_username, restaurant_id, status, total_price, delivery_address, created_at, customer_grace_period_end")).
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetOrder(99)
		if err == nil || err.Error() != "order 99 not found" {
			t.Errorf("expected 'order 99 not found', got %v", err)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "customer_username", "restaurant_id", "status", "total_price", "delivery_address", "created_at", "customer_grace_period_end"}).
			AddRow(1, nil, nil, nil, nil, nil, "invalid_time", "invalid_time")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, customer_username, restaurant_id, status, total_price, delivery_address, created_at, customer_grace_period_end")).
			WithArgs(1).
			WillReturnRows(rows)

		_, err := repo.GetOrder(1)
		if err == nil {
			t.Error("expected scan error, got nil")
		}
	})
}

func TestSetOrderStatusRepository(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewRepository(db)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status = ? WHERE id = ?")).
			WithArgs("cancelled", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SetOrderStatus(1, "cancelled")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewRepository(db)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status = ? WHERE id = ?")).
			WithArgs("cancelled", 1).
			WillReturnError(errors.New("db error"))

		err := repo.SetOrderStatus(1, "cancelled")
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

// --- Service Tests ---

func TestCancelOrderService_AllCases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockOrderRepository{}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockRepo := &mockOrderRepository{err: errors.New("order 1 not found")}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			order: &Order{
				ID:                     1,
				CustomerUsername:       "otheruser",
				Status:                 "pending",
				CustomerGracePeriodEnd: time.Now().Add(5 * time.Minute),
			},
		}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err == nil || err.Error() != "forbidden: this order does not belong to you" {
			t.Errorf("expected forbidden error, got %v", err)
		}
	})

	t.Run("StatusNotPending", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			order: &Order{
				ID:                     1,
				CustomerUsername:       "testuser",
				Status:                 "confirmed",
				CustomerGracePeriodEnd: time.Now().Add(5 * time.Minute),
			},
		}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err == nil || err.Error() != "order cannot be cancelled: current status is 'confirmed'" {
			t.Errorf("expected status error, got %v", err)
		}
	})

	t.Run("GracePeriodExpired", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			order: &Order{
				ID:                     1,
				CustomerUsername:       "testuser",
				Status:                 "pending",
				CustomerGracePeriodEnd: time.Now().Add(-1 * time.Minute),
			},
		}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err == nil || err.Error() != "cancellation window has expired" {
			t.Errorf("expected expired error, got %v", err)
		}
	})

	t.Run("UpdateError", func(t *testing.T) {
		mockRepo := &mockOrderRepository{setStatusErr: errors.New("db error")}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

// --- Handler Tests ---

func TestCancelOrderHandler_AllCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := &mockOrderService{err: nil}
		handler := NewHandler(mockSvc)
		r := gin.Default()
		r.PUT("/order/cancel", func(c *gin.Context) {
			c.Set("username", "testuser")
		}, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("PUT", "/order/cancel", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		handler := NewHandler(&mockOrderService{})
		r := gin.Default()
		r.PUT("/order/cancel", handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("PUT", "/order/cancel", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("ServiceError_NotFound", func(t *testing.T) {
		mockSvc := &mockOrderService{err: errors.New("order 1 not found")}
		handler := NewHandler(mockSvc)
		r := gin.Default()
		r.PUT("/order/cancel", func(c *gin.Context) {
			c.Set("username", "testuser")
		}, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("PUT", "/order/cancel", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("ServiceError_Forbidden", func(t *testing.T) {
		mockSvc := &mockOrderService{err: errors.New("forbidden: error")}
		handler := NewHandler(mockSvc)
		r := gin.Default()
		r.POST("/test", func(c *gin.Context) { c.Set("username", "u") }, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("ServiceError_Unprocessable", func(t *testing.T) {
		mockSvc := &mockOrderService{err: errors.New("order cannot be cancelled: error")}
		handler := NewHandler(mockSvc)
		r := gin.Default()
		r.POST("/test", func(c *gin.Context) { c.Set("username", "u") }, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		mockSvc := &mockOrderService{err: errors.New("something went wrong")}
		handler := NewHandler(mockSvc)
		r := gin.Default()
		r.POST("/test", func(c *gin.Context) { c.Set("username", "u") }, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})

	t.Run("InvalidTokenClaims", func(t *testing.T) {
		handler := NewHandler(&mockOrderService{})
		r := gin.Default()
		r.POST("/test", func(c *gin.Context) {
			c.Set("username", 123)
		}, handler.CancelOrder)

		reqBody, _ := json.Marshal(CancelOrderRequest{OrderID: 1})
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		handler := NewHandler(&mockOrderService{})
		r := gin.Default()
		r.POST("/test", func(c *gin.Context) {
			c.Set("username", "u")
		}, handler.CancelOrder)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString("{invalid json}"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

// --- Mocks ---

// mockOrderRepository จำลอง Repository สำหรับ CancelOrder service tests
type mockOrderRepository struct {
	Repository
	err          error
	order        *Order
	setStatusErr error
}

func (m *mockOrderRepository) GetOrder(orderID int) (*Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.order != nil {
		return m.order, nil
	}
	return &Order{
		ID:                     orderID,
		CustomerUsername:       "testuser",
		Status:                 "pending",
		CustomerGracePeriodEnd: time.Now().Add(5 * time.Minute),
	}, nil
}

func (m *mockOrderRepository) SetOrderStatus(orderID int, status string) error {
	if m.setStatusErr != nil {
		return m.setStatusErr
	}
	return nil
}

// mockOrderService จำลอง Service สำหรับ CancelOrder handler tests
type mockOrderService struct {
	Service
	err error
}

func (m *mockOrderService) CancelOrder(username string, orderID int) error {
	return m.err
}
