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

// mockRepository จำลอง Repository interface สำหรับ service/handler tests
type mockRepository struct {
	Repository
	err         error
	foodItem    *FoodItem
	noFoodItem  bool
	order       *Order
}

func (m *mockRepository) AssignRider(orderID string, riderID int) error {
	return m.err
}

func (m *mockRepository) GetFoodItem(foodItemID, restaurantID int) (*FoodItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.noFoodItem {
		return nil, nil
	}
	if m.foodItem != nil {
		return m.foodItem, nil
	}
	return &FoodItem{Price: 50, IsAvailable: true}, nil
}

func (m *mockRepository) InsertOrderWithItems(username string, restaurantID, totalPrice int, deliveryAddress string, gracePeriodEnd time.Time, items []OrderItemRequest, itemPrices map[int]int) (int64, error) {
	return 1, m.err
}

func (m *mockRepository) GetOrder(orderID int) (*Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.order != nil {
		return m.order, nil
	}
	return &Order{
		ID:                     orderID,
		Status:                 "confirmed",
		CustomerUsername:       "testuser",
		CustomerGracePeriodEnd: time.Now().Add(5 * time.Minute),
	}, nil
}

func (m *mockRepository) GetOrderByID(orderID int) (*Order, []OrderItem, error) {
	if m.err != nil {
		return nil, nil, m.err
	}
	mockOrder := &Order{ID: orderID, Status: "pending", TotalPrice: 500}
	mockItems := []OrderItem{{ID: 1, OrderID: orderID, FoodItemID: 10, Quantity: 2}}
	return mockOrder, mockItems, nil
}

func (m *mockRepository) SetOrderStatus(orderID int, status string) error {
	return m.err
}

// mockService สำหรับเทส Handler โดยตรง
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

// --- Service Tests ---

func TestAssignRider(t *testing.T) {
	t.Run("Success - Should return nil when repo success", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("Failure - Should return error when repo fails", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("database connection failed")}
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestCreateOrderService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		_, _, err := service.CreateOrder("testuser", CreateOrderRequest{
			RestaurantID:    1,
			DeliveryAddress: "Bangkok",
			Items:           []OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("FoodItemNotFound", func(t *testing.T) {
		mockRepo := &mockRepository{noFoodItem: true}
		service := NewService(mockRepo)

		_, _, err := service.CreateOrder("testuser", CreateOrderRequest{
			RestaurantID: 1,
			Items:        []OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
		})
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("FoodItemNotAvailable", func(t *testing.T) {
		mockRepo := &mockRepository{foodItem: &FoodItem{Price: 50, IsAvailable: false}}
		service := NewService(mockRepo)

		_, _, err := service.CreateOrder("testuser", CreateOrderRequest{
			RestaurantID: 1,
			Items:        []OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
		})
		if err == nil || !strings.Contains(err.Error(), "not available") {
			t.Errorf("expected not available error, got %v", err)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repoErr := errors.New("database error")
		mockRepo := &mockRepository{err: repoErr}
		service := NewService(mockRepo)

		_, _, err := service.CreateOrder("user", CreateOrderRequest{
			RestaurantID: 1,
			Items:        []OrderItemRequest{{FoodItemID: 101, Quantity: 1}},
		})
		if err == nil {
			t.Errorf("expected error, got nil")
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
		if order == nil {
			t.Fatal("Expected order object, got nil")
		}
		if order.ID != 1 {
			t.Errorf("Expected ID 1, got %d", order.ID)
		}
		if len(items) == 0 {
			t.Error("Expected items, but got empty list")
		}
	})
}

func TestUpdateOrderStatusService(t *testing.T) {
	t.Run("Success - restaurant: confirmed → preparing", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		svc := NewService(mockRepo)

		err := svc.UpdateOrderStatus(1, "preparing", "restaurant")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Success - rider: assigned → delivering", func(t *testing.T) {
		mockRepo := &mockRepository{order: &Order{ID: 2, Status: "assigned"}}
		svc := NewService(mockRepo)

		err := svc.UpdateOrderStatus(2, "delivering", "rider")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Failure - Order not found", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("order not found")}
		svc := NewService(mockRepo)

		err := svc.UpdateOrderStatus(999, "preparing", "restaurant")
		if err == nil || err.Error() != "order not found" {
			t.Errorf("expected 'order not found', got %v", err)
		}
	})

	t.Run("Failure - Invalid role (customer)", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		svc := NewService(mockRepo)

		err := svc.UpdateOrderStatus(1, "preparing", "customer")
		if err == nil || !strings.Contains(err.Error(), "forbidden") {
			t.Errorf("expected forbidden error, got %v", err)
		}
	})

	t.Run("Failure - Invalid transition (pending → delivered)", func(t *testing.T) {
		mockRepo := &mockRepository{order: &Order{ID: 1, Status: "pending"}}
		svc := NewService(mockRepo)

		err := svc.UpdateOrderStatus(1, "delivered", "rider")
		if err == nil || !strings.Contains(err.Error(), "invalid transition") {
			t.Errorf("expected invalid transition error, got %v", err)
		}
	})
}

// --- Repository Tests ---

func TestGetFoodItemRepository(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock: %s", err)
		}
		defer db.Close()
		repo := NewRepository(db)

		rows := sqlmock.NewRows([]string{"price", "is_available"}).AddRow(50, true)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT price, is_available FROM food_items WHERE id = ? AND restaurant_id = ?")).
			WithArgs(101, 1).
			WillReturnRows(rows)

		item, err := repo.GetFoodItem(101, 1)
		if err != nil || item == nil || item.Price != 50 || !item.IsAvailable {
			t.Errorf("expected item{50, true}, got %v %v", item, err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT price, is_available FROM food_items WHERE id = ? AND restaurant_id = ?")).
			WithArgs(999, 1).
			WillReturnError(sql.ErrNoRows)

		item, err := repo.GetFoodItem(999, 1)
		if err != nil || item != nil {
			t.Errorf("expected nil item and nil error, got %v %v", item, err)
		}
	})
}

func TestInsertOrderWithItemsRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()
	repo := NewRepository(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO orders")).
			WithArgs("testuser", 1, 100, "Bangkok", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO order_items")).
			WithArgs(int64(1), 101, 2, 100).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		orderID, err := repo.InsertOrderWithItems(
			"testuser", 1, 100, "Bangkok",
			time.Now().Add(5*time.Minute),
			[]OrderItemRequest{{FoodItemID: 101, Quantity: 2}},
			map[int]int{101: 50},
		)
		if err != nil || orderID != 1 {
			t.Errorf("expected orderID 1, got %d %v", orderID, err)
		}
	})
}


// --- Handler Tests ---

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
		r.POST("/order", handler.CreateOrder)

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

