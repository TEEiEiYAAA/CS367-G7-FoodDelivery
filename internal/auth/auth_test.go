package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mock repository
type mockRepository struct {
	user *User
	err  error
}

func (m *mockRepository) GetUserByUsername(username string) (*User, error) {
	return m.user, m.err
}

// ---- Service tests ----

func TestLogin_Success(t *testing.T) {
	repo := &mockRepository{
		user: &User{ID: 1, Username: "alice", Password: "password123", Role: "customer"},
	}
	svc := NewService(repo)

	token, err := svc.Login("alice", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token, got empty string")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockRepository{user: nil}
	svc := NewService(repo)

	_, err := svc.Login("unknown", "password123")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got '%v'", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &mockRepository{
		user: &User{ID: 1, Username: "alice", Password: "password123", Role: "customer"},
	}
	svc := NewService(repo)

	_, err := svc.Login("alice", "wrongpassword")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != "invalid password" {
		t.Errorf("expected 'invalid password', got '%v'", err)
	}
}

func TestLogin_RepositoryError(t *testing.T) {
	repo := &mockRepository{err: errors.New("db error")}
	svc := NewService(repo)

	_, err := svc.Login("alice", "password123")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// ---- Handler tests ----

func TestLoginHandler_Success(t *testing.T) {
	repo := &mockRepository{
		user: &User{ID: 1, Username: "alice", Password: "password123", Role: "customer"},
	}
	h := NewHandler(NewService(repo))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"alice","password":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	h := NewHandler(NewService(&mockRepository{}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`invalid json`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	repo := &mockRepository{user: nil}
	h := NewHandler(NewService(repo))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"unknown","password":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	repo := &mockRepository{
		user: &User{ID: 1, Username: "alice", Password: "password123", Role: "customer"},
	}
	h := NewHandler(NewService(repo))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"alice","password":"wrongpassword"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
