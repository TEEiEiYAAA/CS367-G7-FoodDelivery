package auth

import (
	"CS367-G7-FoodDelivery/pkg/jwt"
	"errors"
)

func Login(username, password string) (string, error) {
	user, _ := GetUserByUsername(username)

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
