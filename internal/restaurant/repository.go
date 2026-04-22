package restaurant

import "database/sql"

type Repository interface {
	CreateRestaurant()
	GetRestaurants()
	GetRestaurantByID(id int) (*Restaurant, error)
	ConfirmOrder()
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRestaurant() {}
func (r *repository) GetRestaurants()   {}
func (r *repository) GetRestaurantByID(id int) (*Restaurant, error) {
	row := r.db.QueryRow(
		"SELECT id, name, address, owner_username FROM restaurants WHERE id = ?", id,
	)
	var rest Restaurant
	if err := row.Scan(&rest.ID, &rest.Name, &rest.Address, &rest.OwnerUsername); err != nil {
		return nil, err
	}
	return &rest, nil
}
func (r *repository) ConfirmOrder() {}
