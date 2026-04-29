package order

import (
	"database/sql"
	"errors"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrGracePeriodExpired = errors.New("grace period expired, cannot cancel")

type Repository interface {
	CreateOrder(order Order, items []OrderItem) (Order, error)
	CancelOrder(orderID int, customerUsername string) error
	GetOrderByID(orderID int) (*Order, error)
	UpdateOrderStatus(orderID string, status string) error
	AssignRider(orderID string, riderID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateOrder insert order + order_items ใน transaction เดียว
// โดยดึงราคาจาก food_items ภายใน transaction เพื่อ snapshot ราคา ณ เวลาสั่ง
func (r *repository) CreateOrder(order Order, items []OrderItem) (Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()

	// 1. ดึงราคาแต่ละ item และคำนวณ subtotal + total_price
	totalPrice := 0
	for i, item := range items {
		var price int
		err := tx.QueryRow(
			"SELECT price FROM food_items WHERE id = ? AND is_available = true",
			item.FoodItemID,
		).Scan(&price)
		if err != nil {
			// food_item ไม่เจอ หรือ ไม่ available
			return Order{}, errors.New("food item not found or unavailable")
		}
		subtotal := price * item.Quantity
		items[i].Subtotal = subtotal
		totalPrice += subtotal
	}
	order.TotalPrice = totalPrice

	// 2. Insert order พร้อม total_price ที่คำนวณแล้ว
	const orderQ = `
		INSERT INTO orders
			(customer_username, restaurant_id, status, total_price, delivery_address, created_at, grace_period_end)
		VALUES (?, ?, 'pending', ?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 3 MINUTE))`

	result, err := tx.Exec(orderQ,
		order.CustomerUsername,
		order.RestaurantID,
		order.TotalPrice,
		order.DeliveryAddress,
	)
	if err != nil {
		return Order{}, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return Order{}, err
	}
	order.ID = int(orderID)

	// 3. Insert order_items พร้อม subtotal ที่ snapshot ไว้แล้ว
	const itemQ = `INSERT INTO order_items (order_id, food_item_id, quantity, subtotal) VALUES (?, ?, ?, ?)`
	for i, item := range items {
		_, err := tx.Exec(itemQ, order.ID, item.FoodItemID, item.Quantity, item.Subtotal)
		if err != nil {
			return Order{}, err
		}
		items[i].OrderID = order.ID
	}

	if err := tx.Commit(); err != nil {
		return Order{}, err
	}

	order.Items = items
	return order, nil
}

// CancelOrder ยกเลิกออเดอร์ได้เฉพาะก่อน grace_period_end และเป็นเจ้าของออเดอร์
func (r *repository) CancelOrder(orderID int, customerUsername string) error {
	const q = `
		UPDATE orders
		SET status = 'cancelled'
		WHERE id = ?
		  AND customer_username = ?
		  AND status = 'pending'
		  AND NOW() < grace_period_end`

	result, err := r.db.Exec(q, orderID, customerUsername)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// เช็คว่า order มีอยู่จริงหรือเปล่า เพื่อ return error ที่ถูกต้อง
		var exists int
		r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ? AND customer_username = ?",
			orderID, customerUsername).Scan(&exists)
		if exists == 0 {
			return ErrOrderNotFound
		}
		return ErrGracePeriodExpired
	}
	return nil
}

// GetOrderByID ดึงออเดอร์พร้อม items
func (r *repository) GetOrderByID(orderID int) (*Order, error) {
	const q = `
		SELECT id, customer_username, restaurant_id, rider_id,
		       status, total_price, delivery_address, created_at, grace_period_end
		FROM orders WHERE id = ?`

	var o Order
	var riderID sql.NullInt64
	err := r.db.QueryRow(q, orderID).Scan(
		&o.ID, &o.CustomerUsername, &o.RestaurantID, &riderID,
		&o.Status, &o.TotalPrice, &o.DeliveryAddress, &o.CreatedAt, &o.GracePeriodEnd,
	)
	if err == sql.ErrNoRows {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if riderID.Valid {
		id := int(riderID.Int64)
		o.RiderID = &id
	}

	// ดึง order_items
	rows, err := r.db.Query(
		`SELECT id, order_id, food_item_id, quantity, subtotal FROM order_items WHERE order_id = ?`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.FoodItemID, &item.Quantity, &item.Subtotal); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *repository) UpdateOrderStatus(orderID string, status string) error {
	_, err := r.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderID)
	return err
}

func (r *repository) AssignRider(orderID string, riderID int) error {
	query := "UPDATE orders SET rider_id = ?, status = 'assigned' WHERE id = ?"
	_, err := r.db.Exec(query, riderID, orderID)
	return err
}