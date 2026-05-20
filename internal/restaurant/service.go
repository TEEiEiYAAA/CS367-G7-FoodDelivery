package restaurant

import (
	"database/sql"
	"errors"
)

var ErrInvalidRestaurant = errors.New("invalid restaurant data")
var ErrNotFound = errors.New("restaurant not found")
var ErrOrderNotFound = errors.New("order not found or not authorized")

type Service interface {
	CreateRestaurant(r Restaurant) (Restaurant, error)
	GetRestaurants() ([]Restaurant, error)
	GetRestaurantByID(id int) (*Restaurant, error)
	ConfirmOrder(orderID int, ownerUsername string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRestaurant(r Restaurant) (Restaurant, error) {
	if !r.IsValid() {
		return Restaurant{}, ErrInvalidRestaurant
	}
	return s.repo.CreateRestaurant(r)
}

func (s *service) GetRestaurants() ([]Restaurant, error) {
	return s.repo.GetRestaurants()
}

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

func (s *service) ConfirmOrder(orderID int, ownerUsername string) error {
	return s.repo.ConfirmOrder(orderID, ownerUsername)
}
