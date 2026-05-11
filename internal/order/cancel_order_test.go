package order

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// --- Repository Tests for CancelOrder ---

func TestCancelOrderRepository_AllCases(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	username := "testuser"
	orderID := 1

	t.Run("Success", func(t *testing.T) {
		gracePeriod := time.Now().Add(5 * time.Minute)
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow(username, "pending", gracePeriod)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders WHERE id = ?")).
			WithArgs(orderID).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status = 'cancelled' WHERE id = ?")).
			WithArgs(orderID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CancelOrder(username, orderID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnError(sql.ErrNoRows)

		err := repo.CancelOrder(username, orderID)
		if err == nil || err.Error() != fmt.Sprintf("order %d not found", orderID) {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow("otheruser", "pending", time.Now().Add(5*time.Minute))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnRows(rows)

		err := repo.CancelOrder(username, orderID)
		if err == nil || err.Error() != "forbidden: this order does not belong to you" {
			t.Errorf("expected forbidden error, got %v", err)
		}
	})

	t.Run("StatusNotPending", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow(username, "confirmed", time.Now().Add(5*time.Minute))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnRows(rows)

		err := repo.CancelOrder(username, orderID)
		if err == nil || err.Error() != "order cannot be cancelled: current status is 'confirmed'" {
			t.Errorf("expected status error, got %v", err)
		}
	})

	t.Run("GracePeriodExpired", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow(username, "pending", time.Now().Add(-1*time.Minute))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnRows(rows)

		err := repo.CancelOrder(username, orderID)
		if err == nil || err.Error() != "cancellation window has expired" {
			t.Errorf("expected expired error, got %v", err)
		}
	})

	t.Run("UpdateError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow(username, "pending", time.Now().Add(5*time.Minute))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE orders SET status = 'cancelled'")).
			WithArgs(orderID).
			WillReturnError(errors.New("db error"))

		err := repo.CancelOrder(username, orderID)
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		// ส่งข้อมูลผิดประเภทกลับไปเพื่อให้ Scan พัง
		rows := sqlmock.NewRows([]string{"customer_username", "status", "customer_grace_period_end"}).
			AddRow(123, nil, "invalid_time")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT customer_username, status, customer_grace_period_end FROM orders")).
			WithArgs(orderID).
			WillReturnRows(rows)

		err := repo.CancelOrder(username, orderID)
		if err == nil {
			t.Error("expected scan error, got nil")
		}
	})
}

// --- Service Tests for CancelOrder ---

func TestCancelOrderService_AllCases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockOrderRepository{err: nil}
		service := NewService(mockRepo)
		err := service.CancelOrder("testuser", 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockRepo := &mockOrderRepository{err: repoErr}
		service := NewService(mockRepo)
		err := service.CancelOrder("user", 1)
		if err != repoErr {
			t.Errorf("expected repo error, got %v", err)
		}
	})
}

// --- Handler Tests for CancelOrder ---

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
			c.Set("username", 123) // ส่งเป็นตัวเลขแทน String เพื่อให้ ok == false
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

// --- Mocks for Testing ---

type mockOrderRepository struct {
	Repository
	err error
}

func (m *mockOrderRepository) CancelOrder(username string, orderID int) error {
	return m.err
}

type mockOrderService struct {
	Service
	err error
}

func (m *mockOrderService) CancelOrder(username string, orderID int) error {
	return m.err
}
