package restaurant

import "database/sql"

type Repository interface {
	CreateRestaurant()
	GetRestaurants()
	GetRestaurantByID()
	ConfirmOrder()
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRestaurant()  {}
func (r *repository) GetRestaurants()    {}
func (r *repository) GetRestaurantByID() {}
func (r *repository) ConfirmOrder()      {}
