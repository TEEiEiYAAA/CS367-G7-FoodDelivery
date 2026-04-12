package main

import (
	"log"

	"CS367-G7-FoodDelivery/config"
	"CS367-G7-FoodDelivery/internal/auth"
	"CS367-G7-FoodDelivery/internal/menu"
	"CS367-G7-FoodDelivery/internal/middleware"
	"CS367-G7-FoodDelivery/internal/order"
	"CS367-G7-FoodDelivery/internal/restaurant"
	"CS367-G7-FoodDelivery/internal/rider"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	db := config.DB

	// Init Repositories
	restRepo := restaurant.NewRepository(db)
	menuRepo := menu.NewRepository(db)
	orderRepo := order.NewRepository(db)
	riderRepo := rider.NewRepository(db)

	// Init Services
	restSvc := restaurant.NewService(restRepo)
	menuSvc := menu.NewService(menuRepo)
	orderSvc := order.NewService(orderRepo)
	riderSvc := rider.NewService(riderRepo)

	// Init Handlers
	restHandler := restaurant.NewHandler(restSvc)
	menuHandler := menu.NewHandler(menuSvc)
	orderHandler := order.NewHandler(orderSvc)
	riderHandler := rider.NewHandler(riderSvc)

	r := gin.Default()

	// 11 Features List

	// 🏪 Restaurant
	r.POST("/restaurant", middleware.AuthMiddleware(), restHandler.CreateRestaurant)
	r.GET("/restaurant", restHandler.GetRestaurants)
	r.GET("/restaurant/:id", restHandler.GetRestaurantByID)
	r.PUT("/restaurant/order/confirm", middleware.AuthMiddleware(), restHandler.ConfirmOrder)

	// 🍽 Menu
	r.POST("/restaurant/:id/menu", middleware.AuthMiddleware(), menuHandler.CreateMenu)
	r.GET("/restaurant/:id/menu", menuHandler.GetMenu)

	// 🧾 Order
	r.POST("/order", middleware.AuthMiddleware(), orderHandler.CreateOrder)
	r.PUT("/order/cancel", middleware.AuthMiddleware(), orderHandler.CancelOrder)
	r.GET("/order/:id", orderHandler.GetOrderByID)
	r.PUT("/order/:id/status", middleware.AuthMiddleware(), orderHandler.UpdateOrderStatus)

	// 🛵 Rider
	r.POST("/order/:id/assign-rider", middleware.AuthMiddleware(), riderHandler.AssignRider)

	// Auth Route
	r.POST("/login", auth.LoginHandler)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
