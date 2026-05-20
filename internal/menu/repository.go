package menu

import "database/sql"

type Repository interface {
	CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error)
	GetMenu(restaurantID int) ([]Menu, error)
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
func (r *repository) GetMenu(restaurantID int) ([]Menu, error) {
	const q = `SELECT id, restaurant_id, name, price FROM food_items WHERE restaurant_id = ?`

	rows, err := r.db.Query(q, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menus := make([]Menu, 0)
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.RestaurantID, &menu.Name, &menu.Price); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}
