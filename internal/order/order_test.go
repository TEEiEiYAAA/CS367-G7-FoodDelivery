package order

import (
	"errors"
	"testing"
)

// mockRepository จำลองพฤติกรรมของฐานข้อมูล
type mockRepository struct {
	err error
}

func (m *mockRepository) CreateOrder(order Order, items []OrderItem) (Order, error) {
	if m.err != nil {
		return Order{}, m.err
	}
	// จำลองการดึงราคาจาก DB — ใช้ราคา mock = 100 ต่อ unit
	const mockPrice = 100
	total := 0
	for i := range items {
		subtotal := mockPrice * items[i].Quantity
		items[i].Subtotal = subtotal
		total += subtotal
	}
	order.ID = 1
	order.TotalPrice = total
	order.Items = items
	return order, nil
}

func (m *mockRepository) CancelOrder(orderID int, customerUsername string) error {
	return m.err
}

func (m *mockRepository) GetOrderByID(orderID int) (*Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &Order{ID: orderID, Status: "pending"}, nil
}

func (m *mockRepository) UpdateOrderStatus(orderID string, status string) error {
	return m.err
}

func (m *mockRepository) AssignRider(orderID string, riderID int) error {
	return m.err
}

// ---- Test: AssignRider ----

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

// ---- Test: CreateOrder ----

func TestCreateOrder(t *testing.T) {
	t.Run("Success - Should create order and return it", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		req := CreateOrderRequest{
			RestaurantID:    1,
			DeliveryAddress: "123 Test Road",
			Items: []OrderItemRequest{
				{FoodItemID: 10, Quantity: 2}, // mock price=100 → subtotal=200
				{FoodItemID: 11, Quantity: 1}, // mock price=100 → subtotal=100
			},
		}

		order, err := service.CreateOrder(req, "testuser")
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if order.ID == 0 {
			t.Error("Expected order ID to be set")
		}
		if order.CustomerUsername != "testuser" {
			t.Errorf("Expected customer_username 'testuser', got '%s'", order.CustomerUsername)
		}
		// ตรวจ total_price = 200 + 100 = 300
		if order.TotalPrice != 300 {
			t.Errorf("Expected total_price 300, got %d", order.TotalPrice)
		}
		// ตรวจ subtotal ของ item แรก
		if order.Items[0].Subtotal != 200 {
			t.Errorf("Expected item[0] subtotal 200, got %d", order.Items[0].Subtotal)
		}
	})

	t.Run("Failure - Should return error when repo fails", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("db error")}
		service := NewService(mockRepo)

		req := CreateOrderRequest{
			RestaurantID:    1,
			DeliveryAddress: "123 Test Road",
			Items:           []OrderItemRequest{{FoodItemID: 10, Quantity: 1}},
		}

		_, err := service.CreateOrder(req, "testuser")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// ---- Test: CancelOrder ----

func TestCancelOrder(t *testing.T) {
	t.Run("Success - Should cancel order", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil}
		service := NewService(mockRepo)

		err := service.CancelOrder(1, "testuser")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("Failure - Grace period expired", func(t *testing.T) {
		mockRepo := &mockRepository{err: ErrGracePeriodExpired}
		service := NewService(mockRepo)

		err := service.CancelOrder(1, "testuser")
		if !errors.Is(err, ErrGracePeriodExpired) {
			t.Errorf("Expected ErrGracePeriodExpired, got %v", err)
		}
	})
}