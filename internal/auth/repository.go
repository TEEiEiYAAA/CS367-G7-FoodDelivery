package auth

import (
	"database/sql"
)

type Repository interface {
	GetUserByUsername(username string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
