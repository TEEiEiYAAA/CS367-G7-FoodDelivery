package auth

import (
	"CS367-G7-FoodDelivery/config"
	"database/sql"
)

func GetUserByUsername(username string) (*User, error) {
	var user User
	// แก้ไข SQL ตามชื่อ Table ในโปรเจกต์ (เช่น Customer หรือ Restaurant)
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	err := config.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
