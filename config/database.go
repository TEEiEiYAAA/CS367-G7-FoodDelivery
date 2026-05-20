package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxDBRetries = 10
	retryDelay   = 3 * time.Second
)

var DB *sql.DB

func InitDB() {
	var err error

	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "food_delivery_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error configuring database: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute * 5)

	for i := 1; i <= maxDBRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Database connection failed (attempt %d/%d): %v", i, maxDBRetries, err)
		if i == maxDBRetries {
			log.Fatalf("Could not connect to database after %d attempts", maxDBRetries)
		}
		time.Sleep(retryDelay)
	}

	log.Println("Successfully connected to the database!")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}