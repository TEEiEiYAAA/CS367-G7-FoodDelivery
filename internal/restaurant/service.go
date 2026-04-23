package restaurant

import "errors"

// ErrInvalidRestaurant is returned when required fields are missing.
var ErrInvalidRestaurant = errors.New("invalid restaurant data")

type Service interface {
	CreateRestaurant(r Restaurant) (Restaurant, error)
	GetRestaurants() ([]Restaurant, error)
	GetRestaurantByID()
	ConfirmOrder()
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

func (s *service) GetRestaurantByID() {}
func (s *service) ConfirmOrder()      {}
