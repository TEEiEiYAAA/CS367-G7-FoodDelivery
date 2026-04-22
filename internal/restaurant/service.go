package restaurant

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("restaurant not found")

type Service interface {
	CreateRestaurant()
	GetRestaurants()
	GetRestaurantByID(id int) (*Restaurant, error)
	ConfirmOrder()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRestaurant() {}
func (s *service) GetRestaurants()   {}
func (s *service) GetRestaurantByID(id int) (*Restaurant, error) {
	rest, err := s.repo.GetRestaurantByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return rest, nil
}
func (s *service) ConfirmOrder() {}
