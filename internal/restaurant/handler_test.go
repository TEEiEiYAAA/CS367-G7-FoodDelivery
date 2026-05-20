package restaurant

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// mock service
type mockService struct {
	createFn  func(Restaurant) (Restaurant, error)
	listFn    func() ([]Restaurant, error)
	getByIDFn func(int) (*Restaurant, error)
	confirmFn func(int, string) error
}

func (m *mockService) CreateRestaurant(r Restaurant) (Restaurant, error) {
	return m.createFn(r)
}
func (m *mockService) GetRestaurants() ([]Restaurant, error) {
	return m.listFn()
}
func (m *mockService) GetRestaurantByID(id int) (*Restaurant, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}
func (m *mockService) ConfirmOrder(orderID int, ownerUsername string) error {
	if m.confirmFn != nil {
		return m.confirmFn(orderID, ownerUsername)
	}
	return nil
}

// helpers
func setupRouter(svc Service, username string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if username != "" {
		r.Use(func(c *gin.Context) {
			c.Set("username", username)
			c.Next()
		})
	}
	h := NewHandler(svc)
	r.POST("/restaurant", h.CreateRestaurant)
	r.GET("/restaurant", h.GetRestaurants)
	r.GET("/restaurant/:id", h.GetRestaurantByID)
	r.PUT("/restaurant/order/confirm", h.ConfirmOrder)
	return r
}

func doRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// POST /restaurant
func TestHandler_CreateRestaurant_Success(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			r.ID = 1
			return r, nil
		},
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPost, "/restaurant", `{"name":"KFC","address":"Bangkok"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusCreated, w.Body.String())
	}

	var got Restaurant
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
	if got.Name != "KFC" || got.Address != "Bangkok" || got.OwnerUsername != "alice" {
		t.Errorf("unexpected body: %+v", got)
	}
}

func TestHandler_CreateRestaurant_BadJSON(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			t.Fatal("service should not be called")
			return Restaurant{}, nil
		},
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPost, "/restaurant", `{not-json}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateRestaurant_MissingRequiredField(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			t.Fatal("service should not be called")
			return Restaurant{}, nil
		},
	}
	r := setupRouter(svc, "alice")

	// missing "address"
	w := doRequest(r, http.MethodPost, "/restaurant", `{"name":"KFC"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateRestaurant_NoUsernameInContext(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			t.Fatal("service should not be called")
			return Restaurant{}, nil
		},
	}
	r := setupRouter(svc, "") // middleware does NOT inject username

	w := doRequest(r, http.MethodPost, "/restaurant", `{"name":"KFC","address":"Bangkok"}`)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandler_CreateRestaurant_ServiceValidationError(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			return Restaurant{}, ErrInvalidRestaurant
		},
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPost, "/restaurant", `{"name":"KFC","address":"Bangkok"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateRestaurant_ServiceInternalError(t *testing.T) {
	svc := &mockService{
		createFn: func(r Restaurant) (Restaurant, error) {
			return Restaurant{}, errors.New("db down")
		},
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPost, "/restaurant", `{"name":"KFC","address":"Bangkok"}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// GET /restaurant
func TestHandler_GetRestaurants_Success(t *testing.T) {
	want := []Restaurant{
		{ID: 1, Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"},
		{ID: 2, Name: "MK", Address: "Nonthaburi", OwnerUsername: "bob"},
	}
	svc := &mockService{
		listFn: func() ([]Restaurant, error) { return want, nil },
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got []Restaurant
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
}

func TestHandler_GetRestaurants_EmptyReturnsArrayNotNull(t *testing.T) {
	svc := &mockService{
		listFn: func() ([]Restaurant, error) { return []Restaurant{}, nil },
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := string(bytes.TrimSpace(w.Body.Bytes()))
	if body != "[]" {
		t.Errorf("body = %q, want []", body)
	}
}

func TestHandler_GetRestaurants_ServiceError(t *testing.T) {
	svc := &mockService{
		listFn: func() ([]Restaurant, error) {
			return nil, errors.New("db down")
		},
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// GET /restaurant/:id
func TestHandler_GetRestaurantByID_Success(t *testing.T) {
	want := &Restaurant{ID: 1, Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"}
	svc := &mockService{
		getByIDFn: func(id int) (*Restaurant, error) { return want, nil },
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant/1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var got Restaurant
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != 1 || got.Name != "KFC" || got.OwnerUsername != "alice" {
		t.Errorf("unexpected body: %+v", got)
	}
}

func TestHandler_GetRestaurantByID_NotFound(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(id int) (*Restaurant, error) { return nil, ErrNotFound },
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant/99", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_GetRestaurantByID_InvalidID(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_GetRestaurantByID_ServiceError(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(id int) (*Restaurant, error) { return nil, errors.New("db down") },
	}
	r := setupRouter(svc, "")

	w := doRequest(r, http.MethodGet, "/restaurant/1", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// PUT /restaurant/order/confirm
func TestHandler_ConfirmOrder_Success(t *testing.T) {
	svc := &mockService{
		confirmFn: func(orderID int, ownerUsername string) error { return nil },
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPut, "/restaurant/order/confirm", `{"order_id":1}`)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHandler_ConfirmOrder_BadJSON(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPut, "/restaurant/order/confirm", `{not-json}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_ConfirmOrder_NoUsernameInContext(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc, "") // no username

	w := doRequest(r, http.MethodPut, "/restaurant/order/confirm", `{"order_id":1}`)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandler_ConfirmOrder_OrderNotFound(t *testing.T) {
	svc := &mockService{
		confirmFn: func(orderID int, ownerUsername string) error { return ErrOrderNotFound },
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPut, "/restaurant/order/confirm", `{"order_id":99}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_ConfirmOrder_ServiceError(t *testing.T) {
	svc := &mockService{
		confirmFn: func(orderID int, ownerUsername string) error { return errors.New("db down") },
	}
	r := setupRouter(svc, "alice")

	w := doRequest(r, http.MethodPut, "/restaurant/order/confirm", `{"order_id":1}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
