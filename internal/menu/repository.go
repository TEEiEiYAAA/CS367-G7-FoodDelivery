package menu

import "database/sql"

type Repository interface {
	CreateMenu()
	GetMenu()
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMenu() {}
func (r *repository) GetMenu()    {}
