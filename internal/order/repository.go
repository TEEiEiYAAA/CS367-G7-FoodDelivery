package order

import "database/sql"

type Repository interface {
	CreateOrder()
	CancelOrder()
	GetOrderByID()
	UpdateOrderStatus()
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
