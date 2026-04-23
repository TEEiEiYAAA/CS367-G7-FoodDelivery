package restaurant

import "database/sql"

type Repository interface {
	CreateRestaurant(r Restaurant) (Restaurant, error)
	GetRestaurants() ([]Restaurant, error)
	GetRestaurantByID(id int) (*Restaurant, error)
	ConfirmOrder(orderID int, ownerUsername string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRestaurant(rest Restaurant) (Restaurant, error) {
	const q = `INSERT INTO restaurants (name, address, owner_username) VALUES (?, ?, ?)`

	result, err := r.db.Exec(q, rest.Name, rest.Address, rest.OwnerUsername)
	if err != nil {
		return Restaurant{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Restaurant{}, err
	}

	rest.ID = int(id)
	return rest, nil
}

func (r *repository) GetRestaurants() ([]Restaurant, error) {
	const q = `SELECT id, name, address, owner_username FROM restaurants ORDER BY id`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]Restaurant, 0)
	for rows.Next() {
		var rest Restaurant
		if err := rows.Scan(&rest.ID, &rest.Name, &rest.Address, &rest.OwnerUsername); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, rest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return restaurants, nil
}

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

func (r *repository) ConfirmOrder(orderID int, ownerUsername string) error {
	const q = `
		UPDATE orders o
		JOIN restaurants res ON o.restaurant_id = res.id
		SET o.status = 'confirmed'
		WHERE o.id = ? AND res.owner_username = ?`

	result, err := r.db.Exec(q, orderID, ownerUsername)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrOrderNotFound
	}
	return nil
}
