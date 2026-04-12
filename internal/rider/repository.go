package rider

import "database/sql"

type Repository interface {
	AssignRider()
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AssignRider() {}
