package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	// ปรับ Connection String ตามของทีมคุณ
	DB, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/food_db")
	if err != nil {
		panic(err)
	}
}
