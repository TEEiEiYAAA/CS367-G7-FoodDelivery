package order

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"strings"
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