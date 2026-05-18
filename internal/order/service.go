package order

import (
	"fmt"
	"time"
)

type Service interface {
	CreateOrder(username string, req CreateOrderRequest) (int64, int, error)
	CancelOrder(username string, orderID int) error
	GetOrderByID(orderID int) (*Order, []OrderItem, error)
	UpdateOrderStatus(orderID int, newStatus string, role string) error
	AssignRider(orderID string, riderID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	totalPrice := 0
	itemPrices := make(map[int]int)

	for _, item := range req.Items {
		foodItem, err := s.repo.GetFoodItem(item.FoodItemID, req.RestaurantID)
		if err != nil {
			return 0, 0, err
		}
		if foodItem == nil {
			return 0, 0, fmt.Errorf("food item %d not found in restaurant %d", item.FoodItemID, req.RestaurantID)
		}
		if !foodItem.IsAvailable {
			return 0, 0, fmt.Errorf("food item %d is not available", item.FoodItemID)
		}
		itemPrices[item.FoodItemID] = foodItem.Price
		totalPrice += foodItem.Price * item.Quantity
	}

	gracePeriodEnd := time.Now().Add(5 * time.Minute)
	orderID, err := s.repo.InsertOrderWithItems(username, req.RestaurantID, totalPrice, req.DeliveryAddress, gracePeriodEnd, req.Items, itemPrices)
	return orderID, totalPrice, err
}

func (s *service) CancelOrder(username string, orderID int) error {
	order, err := s.repo.GetOrder(orderID)
	if err != nil {
		return err
	}
	if order.CustomerUsername != username {
		return fmt.Errorf("forbidden: this order does not belong to you")
	}
	if order.Status != "pending" {
		return fmt.Errorf("order cannot be cancelled: current status is '%s'", order.Status)
	}
	if time.Now().After(order.CustomerGracePeriodEnd) {
		return fmt.Errorf("cancellation window has expired")
	}
	return s.repo.SetOrderStatus(orderID, "cancelled")
}

func (s *service) GetOrderByID(orderID int) (*Order, []OrderItem, error) {
	return s.repo.GetOrderByID(orderID)
}

func (s *service) UpdateOrderStatus(orderID int, newStatus string, role string) error {
	order, err := s.repo.GetOrder(orderID)
	if err != nil {
		return err
	}

	allowed := map[string]map[string]string{
		"restaurant": {
			"confirmed": "preparing",
			"preparing": "ready",
		},
		"rider": {
			"confirmed":  "delivering",
			"ready":      "delivering",
			"assigned":   "delivering",
			"delivering": "delivered",
		},
	}

	transitions, roleExists := allowed[role]
	if !roleExists {
		return fmt.Errorf("forbidden: role '%s' cannot update order status", role)
	}

	expectedNext, ok := transitions[order.Status]
	if !ok || expectedNext != newStatus {
		return fmt.Errorf("invalid transition: '%s' → '%s' is not allowed for role '%s'", order.Status, newStatus, role)
	}

	return s.repo.SetOrderStatus(orderID, newStatus)
}

func (s *service) AssignRider(orderID string, riderID int) error {
	return s.repo.AssignRider(orderID, riderID)
}
