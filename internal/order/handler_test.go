package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetOrderByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Should return 200 and order data", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		h.GetOrderByID(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["order"] == nil {
			t.Error("Expected order data in response, but got nil")
		}
	})

	t.Run("Failure - Invalid ID format should return 400", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "abc"}}

		h.GetOrderByID(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Failure - Order Not Found should return 404", func(t *testing.T) {

		mockRepo := &mockRepository{err: errors.New("order not found")}
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		h.GetOrderByID(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestAssignRiderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Should return 200", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		svc := NewService(mockRepo)
		h := NewHandler(svc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		body := `{"riderId": 101}`
		c.Request = httptest.NewRequest("POST", "/order/1/assign-rider", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AssignRider(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Failure - Invalid JSON body", func(t *testing.T) {
		h := NewHandler(&mockRepository{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"riderId": "not-a-number"}`))

		h.AssignRider(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Failure - Service Error", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("db error")}
		h := NewHandler(NewService(mockRepo))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"riderId": 101}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AssignRider(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestUpdateOrderStatusHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - restaurant changes status to preparing (200)", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		h := NewHandler(NewService(mockRepo))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Set("role", "restaurant")
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`{"status": "preparing"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Success - rider changes status to delivering (200)", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		h := NewHandler(NewService(mockRepo))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Set("role", "rider")
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`{"status": "delivering"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Failure - Invalid order ID (400)", func(t *testing.T) {
		h := NewHandler(NewService(&mockRepository{}))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Failure - No role in context (401)", func(t *testing.T) {
		h := NewHandler(NewService(&mockRepository{}))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`{"status": "preparing"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})

	t.Run("Failure - Invalid JSON body (400)", func(t *testing.T) {
		h := NewHandler(NewService(&mockRepository{}))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Set("role", "restaurant")
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`invalid`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Failure - Order not found (404)", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("order not found")}
		h := NewHandler(NewService(mockRepo))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}
		c.Set("role", "restaurant")
		c.Request = httptest.NewRequest("PUT", "/order/999/status", strings.NewReader(`{"status": "preparing"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Failure - Forbidden role (403)", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("forbidden: role 'customer' cannot update order status")}
		h := NewHandler(NewService(mockRepo))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Set("role", "customer")
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`{"status": "preparing"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", w.Code)
		}
	})

	t.Run("Failure - Invalid transition (422)", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("invalid transition: 'pending' → 'delivered' is not allowed for role 'rider'")}
		h := NewHandler(NewService(mockRepo))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Set("role", "rider")
		c.Request = httptest.NewRequest("PUT", "/order/1/status", strings.NewReader(`{"status": "delivered"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateOrderStatus(c)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected 422, got %d", w.Code)
		}
	})
}
