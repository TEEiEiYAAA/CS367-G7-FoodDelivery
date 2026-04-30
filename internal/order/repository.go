package order

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	CreateOrder(username string, req CreateOrderRequest) (int64, int, error)
	CancelOrder(username string, orderID int) error
	GetOrderByID()
	UpdateOrderStatus()
	AssignRider(orderID string, riderID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateOrder สร้าง order ใหม่ใน transaction
// คืน (orderID, totalPrice, error)
func (r *repository) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	// เริ่ม transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// คำนวณ total_price โดย query ราคาจาก food_items
	totalPrice := 0
	type itemDetail struct {
		price    int
		quantity int
	}
	details := make([]itemDetail, 0, len(req.Items))

	for _, item := range req.Items {
		var price int
		var isAvailable bool
		err = tx.QueryRow(
			"SELECT price, is_available FROM food_items WHERE id = ? AND restaurant_id = ?",
			item.FoodItemID, req.RestaurantID,
		).Scan(&price, &isAvailable)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, 0, fmt.Errorf("food item %d not found in restaurant %d", item.FoodItemID, req.RestaurantID)
			}
			return 0, 0, err
		}
		if !isAvailable {
			return 0, 0, fmt.Errorf("food item %d is not available", item.FoodItemID)
		}
		details = append(details, itemDetail{price: price, quantity: item.Quantity})
		totalPrice += price * item.Quantity
	}

	// grace period 5 นาที (ลูกค้ายกเลิกได้ใน 5 นาที)
	gracePeriodEnd := time.Now().Add(5 * time.Minute)

	// INSERT orders
	result, err := tx.Exec(
		`INSERT INTO orders
			(customer_username, restaurant_id, rider_id, status, total_price, delivery_address, created_at, customer_grace_period_end)
		 VALUES (?, ?, NULL, 'pending', ?, ?, NOW(), ?)`,
		username, req.RestaurantID, totalPrice, req.DeliveryAddress, gracePeriodEnd,
	)
	if err != nil {
		return 0, 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	// INSERT order_items
	for i, item := range req.Items {
		subtotal := details[i].price * details[i].quantity
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, food_item_id, quantity, subtotal) VALUES (?, ?, ?, ?)",
			orderID, item.FoodItemID, item.Quantity, subtotal,
		)
		if err != nil {
			return 0, 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return orderID, totalPrice, nil
}

// CancelOrder ตรวจสิทธิ์และ grace period ก่อนเปลี่ยน status เป็น cancelled
func (r *repository) CancelOrder(username string, orderID int) error {
	var dbUsername, status string
	var gracePeriodEnd time.Time

	err := r.db.QueryRow(
		"SELECT customer_username, status, customer_grace_period_end FROM orders WHERE id = ?",
		orderID,
	).Scan(&dbUsername, &status, &gracePeriodEnd)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order %d not found", orderID)
		}
		return err
	}

	// เช็คว่าเป็นเจ้าของ order นี้จริง
	if dbUsername != username {
		return fmt.Errorf("forbidden: this order does not belong to you")
	}

	// เช็คว่า status ยังเป็น pending อยู่
	if status != "pending" {
		return fmt.Errorf("order cannot be cancelled: current status is '%s'", status)
	}

	// เช็ค grace period
	if time.Now().After(gracePeriodEnd) {
		return fmt.Errorf("cancellation window has expired")
	}

	_, err = r.db.Exec(
		"UPDATE orders SET status = 'cancelled' WHERE id = ?",
		orderID,
	)
	return err
}
func (r *repository) GetOrderByID()      {}
func (r *repository) UpdateOrderStatus() {}
func (r *repository) AssignRider(orderID string, riderID int) error {
	query := "UPDATE orders SET rider_id = ?, status = 'assigned' WHERE id = ?"

	_, err := r.db.Exec(query, riderID, orderID)
	return err
}