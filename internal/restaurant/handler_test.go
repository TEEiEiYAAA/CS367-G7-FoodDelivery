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

// --- mock service ------------------------------------------------------------

// mockService implements Service. Each test sets the fields it needs.
type mockService struct {
	createFn func(Restaurant) (Restaurant, error)
	listFn   func() ([]Restaurant, error)
}

func (m *mockService) CreateRestaurant(r Restaurant) (Restaurant, error) {
	return m.createFn(r)
}
func (m *mockService) GetRestaurants() ([]Restaurant, error) {
	return m.listFn()
}
func (m *mockService) GetRestaurantByID() {}
func (m *mockService) ConfirmOrder()      {}

// --- helpers -----------------------------------------------------------------

// setupRouter wires the handler into a fresh Gin test engine.
// When username is non-empty, a tiny middleware injects it into the context —
// emulating what AuthMiddleware does in production.
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
	return r
}

func doRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- POST /restaurant --------------------------------------------------------

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

// --- GET /restaurant ---------------------------------------------------------

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
