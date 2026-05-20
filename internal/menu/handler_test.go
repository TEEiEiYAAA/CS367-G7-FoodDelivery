package menu

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockService struct {
	createMenuFn func(restaurantID int, req CreateMenuRequest) (Menu, error)
	getMenuFn    func(restaurantID int) ([]Menu, error)
}

func (m *mockService) CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error) {
	return m.createMenuFn(restaurantID, req)
}

func (m *mockService) GetMenu(restaurantID int) ([]Menu, error) {
	return m.getMenuFn(restaurantID)
}

func newTestRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/restaurant/:id/menu", h.CreateMenu)
	r.GET("/restaurant/:id/menu", h.GetMenu)
	return r
}

func TestCreateMenuHandler(t *testing.T) {

	t.Run("success_201", func(t *testing.T) {
		svc := &mockService{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return Menu{ID: 1, RestaurantID: restaurantID, Name: req.Name, Price: req.Price}, nil
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		body, _ := json.Marshal(CreateMenuRequest{Name: "Pad Thai", Price: 80})
		req := httptest.NewRequest(http.MethodPost, "/restaurant/10/menu", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}

		var got Menu
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if got.Name != "Pad Thai" || got.Price != 80 || got.RestaurantID != 10 {
			t.Errorf("unexpected response body: %+v", got)
		}
	})

	t.Run("invalid_restaurant_id_400", func(t *testing.T) {
		svc := &mockService{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return Menu{}, nil
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		body, _ := json.Marshal(CreateMenuRequest{Name: "Pad Thai", Price: 80})
		req := httptest.NewRequest(http.MethodPost, "/restaurant/abc/menu", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid_json_body_400", func(t *testing.T) {
		svc := &mockService{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return Menu{}, nil
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/restaurant/10/menu", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service_error_500", func(t *testing.T) {
		svc := &mockService{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return Menu{}, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		body, _ := json.Marshal(CreateMenuRequest{Name: "Pad Thai", Price: 80})
		req := httptest.NewRequest(http.MethodPost, "/restaurant/10/menu", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetMenuHandler(t *testing.T) {

	t.Run("success_200", func(t *testing.T) {
		svc := &mockService{
			getMenuFn: func(restaurantID int) ([]Menu, error) {
				return []Menu{
					{ID: 1, RestaurantID: restaurantID, Name: "Pad Thai", Price: 80},
				}, nil
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/restaurant/10/menu", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var got []Menu
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if len(got) != 1 || got[0].Name != "Pad Thai" {
			t.Errorf("unexpected response body: %+v", got)
		}
	})

	t.Run("invalid_restaurant_id_400", func(t *testing.T) {
		svc := &mockService{}
		h := NewHandler(svc)
		r := newTestRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/restaurant/abc/menu", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service_error_500", func(t *testing.T) {
		svc := &mockService{
			getMenuFn: func(restaurantID int) ([]Menu, error) {
				return nil, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		r := newTestRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/restaurant/10/menu", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}
