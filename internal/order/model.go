package order

import "time"

// ---- Domain Models ----

type Order struct {
	ID              int         `json:"id"`
	CustomerUsername string     `json:"customer_username"`
	RestaurantID    int         `json:"restaurant_id"`
	RiderID         *int        `json:"rider_id"`        // pointer เพราะ nullable
	Status          string      `json:"status"`
	TotalPrice      int         `json:"total_price"`
	DeliveryAddress string      `json:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at"`
	GracePeriodEnd  time.Time   `json:"grace_period_end"` // deadline ที่ลูกค้ายังยกเลิกได้
	Items           []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID           int    `json:"id"`
	OrderID      int    `json:"order_id"`
	FoodItemID   int    `json:"food_item_id"`
	Quantity     int    `json:"quantity"`
	Subtotal     int    `json:"subtotal"`
}

// ---- Request / Response ----

// POST /order
type CreateOrderRequest struct {
	RestaurantID    int               `json:"restaurant_id"    binding:"required"`
	DeliveryAddress string            `json:"delivery_address" binding:"required"`
	Items           []OrderItemRequest `json:"items"            binding:"required,min=1"`
}

type OrderItemRequest struct {
	FoodItemID int `json:"food_item_id" binding:"required"`
	Quantity   int `json:"quantity"     binding:"required,min=1"`
}
