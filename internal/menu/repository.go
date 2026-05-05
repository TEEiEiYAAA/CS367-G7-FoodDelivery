package menu

import "database/sql"

type Repository interface {
	CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error)
	GetMenu()
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error) {
	const q = `INSERT INTO food_items (restaurant_id, name, price, is_available) VALUES (?, ?, ?, true)`

	result, err := r.db.Exec(q, restaurantID, req.Name, req.Price)
	if err != nil {
		return Menu{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Menu{}, err
	}

	return Menu{
		ID:           int(id),
		RestaurantID: restaurantID,
		Name:         req.Name,
		Price:        req.Price,
	}, nil
}
func (r *repository) GetMenu()    {}
