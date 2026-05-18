package auth

import (
	"errors"

	"CS367-G7-FoodDelivery/pkg/jwt"
)

type Service interface {
	Login(username, password string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Login(username, password string) (string, error) {
	user, _ := s.repo.GetUserByUsername(username)

	if user == nil {
		return "", errors.New("user not found")
	}

	if user.Password != password {
		return "", errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
