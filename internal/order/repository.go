package order

import "database/sql"

type Repository interface {
	CreateOrder()
	CancelOrder()
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

func (r *repository) CreateOrder()       {}
func (r *repository) CancelOrder()       {}
func (r *repository) GetOrderByID()      {}
func (r *repository) UpdateOrderStatus() {}
func (r *repository) AssignRider(orderID string, riderID int) error {
	query := "UPDATE orders SET rider_id = ?, status = 'assigned' WHERE id = ?"

	_, err := r.db.Exec(query, riderID, orderID)
	return err
}