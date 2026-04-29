package order

import "time"

// Order represents the orders table
type Order struct {
	ID                      int        `json:"id"`
	CustomerUsername        string     `json:"customer_username"`
	RestaurantID            int        `json:"restaurant_id"`
	RiderID                 *int       `json:"rider_id"`
	Status                  string     `json:"status"`
	TotalPrice              int        `json:"total_price"`
	DeliveryAddress         string     `json:"delivery_address"`
	CreatedAt               time.Time  `json:"created_at"`
	CustomerGracePeriodEnd  time.Time  `json:"customer_grace_period_end"`
}

// OrderItem represents the order_items table
type OrderItem struct {
	ID         int `json:"id"`
	OrderID    int `json:"order_id"`
	FoodItemID int `json:"food_item_id"`
	Quantity   int `json:"quantity"`
	Subtotal   int `json:"subtotal"`
}

// CreateOrderRequest คือ request body สำหรับ POST /order
type CreateOrderRequest struct {
	RestaurantID    int                  `json:"restaurant_id" binding:"required"`
	DeliveryAddress string               `json:"delivery_address" binding:"required"`
	Items           []OrderItemRequest   `json:"items" binding:"required,min=1"`
}

// OrderItemRequest คือแต่ละ item ใน request
type OrderItemRequest struct {
	FoodItemID int `json:"food_item_id" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,min=1"`
}

// CreateOrderResponse คือ response หลัง POST /order สำเร็จ
type CreateOrderResponse struct {
	OrderID    int64  `json:"order_id"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
}
