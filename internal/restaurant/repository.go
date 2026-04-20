package restaurant

import "database/sql"

// Repository defines the persistence contract for Restaurant entities.
//
// Only CreateRestaurant and GetRestaurants are implemented here (owner: Kanokporn).
// GetRestaurantByID and ConfirmOrder are reserved stubs for other teammates to
// fill in under their own feature branches.
type Repository interface {
	CreateRestaurant(r Restaurant) (Restaurant, error)
	GetRestaurants() ([]Restaurant, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateRestaurant inserts a new restaurant row and returns it populated with
// the generated ID.
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

// GetRestaurants returns every restaurant ordered by id.
// On an empty table it returns an empty slice (never nil) so handlers can
// always marshal it as a JSON array.
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

// Stubs owned by other teammates — keep empty signature so the package still
// compiles and main.go wiring stays intact.
func (r *repository) GetRestaurantByID() {}
func (r *repository) ConfirmOrder()      {}
