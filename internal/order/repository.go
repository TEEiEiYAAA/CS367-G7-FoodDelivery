package order

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	GetFoodItem(foodItemID, restaurantID int) (*FoodItem, error)
	InsertOrderWithItems(username string, restaurantID, totalPrice int, deliveryAddress string, gracePeriodEnd time.Time, items []OrderItemRequest, itemPrices map[int]int) (int64, error)
	GetOrder(orderID int) (*Order, error)
	SetOrderStatus(orderID int, status string) error
	GetOrderByID(orderID int) (*Order, []OrderItem, error)
	AssignRider(orderID string, riderID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetFoodItem(foodItemID, restaurantID int) (*FoodItem, error) {
	var item FoodItem
	err := r.db.QueryRow(
		"SELECT price, is_available FROM food_items WHERE id = ? AND restaurant_id = ?",
		foodItemID, restaurantID,
	).Scan(&item.Price, &item.IsAvailable)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) InsertOrderWithItems(username string, restaurantID, totalPrice int, deliveryAddress string, gracePeriodEnd time.Time, items []OrderItemRequest, itemPrices map[int]int) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(
		`INSERT INTO orders
			(customer_username, restaurant_id, rider_id, status, total_price, delivery_address, created_at, customer_grace_period_end)
		 VALUES (?, ?, NULL, 'pending', ?, ?, NOW(), ?)`,
		username, restaurantID, totalPrice, deliveryAddress, gracePeriodEnd,
	)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		subtotal := itemPrices[item.FoodItemID] * item.Quantity
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, food_item_id, quantity, subtotal) VALUES (?, ?, ?, ?)",
			orderID, item.FoodItemID, item.Quantity, subtotal,
		)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	return orderID, err
}

func (r *repository) GetOrder(orderID int) (*Order, error) {
	var order Order
	err := r.db.QueryRow(
		`SELECT id, customer_username, restaurant_id, status, total_price, delivery_address, created_at, customer_grace_period_end
		 FROM orders WHERE id = ?`,
		orderID,
	).Scan(&order.ID, &order.CustomerUsername, &order.RestaurantID, &order.Status, &order.TotalPrice, &order.DeliveryAddress, &order.CreatedAt, &order.CustomerGracePeriodEnd)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order %d not found", orderID)
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *repository) SetOrderStatus(orderID int, status string) error {
	_, err := r.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderID)
	return err
}

func (r *repository) GetOrderByID(orderID int) (*Order, []OrderItem, error) {
	var order Order
	queryOrder := "SELECT id, customer_username, restaurant_id, status, total_price FROM orders WHERE id = ?"
	err := r.db.QueryRow(queryOrder, orderID).Scan(&order.ID, &order.CustomerUsername, &order.RestaurantID, &order.Status, &order.TotalPrice)
	if err != nil {
		return nil, nil, err
	}

	var items []OrderItem
	queryItems := "SELECT id, order_id, food_item_id, quantity, subtotal FROM order_items WHERE order_id = ?"
	rows, err := r.db.Query(queryItems, orderID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.FoodItemID, &item.Quantity, &item.Subtotal); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return &order, items, nil
}

func (r *repository) AssignRider(orderID string, riderID int) error {
	_, err := r.db.Exec("UPDATE orders SET rider_id = ?, status = 'assigned' WHERE id = ?", riderID, orderID)
	return err
}
