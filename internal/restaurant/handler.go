package restaurant

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /restaurant
func (h *Handler) CreateRestaurant(c *gin.Context) {
	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	rest := Restaurant{
		Name:          req.Name,
		Address:       req.Address,
		OwnerUsername: username,
	}

	created, err := h.service.CreateRestaurant(rest)
	if err != nil {
		if errors.Is(err, ErrInvalidRestaurant) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create restaurant"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GET /restaurant
func (h *Handler) GetRestaurants(c *gin.Context) {
	list, err := h.service.GetRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch restaurants"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /restaurant/{id} (ดูข้อมูลของร้านอาหาร)
func (h *Handler) GetRestaurantByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rest, err := h.service.GetRestaurantByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, rest)
}

// PUT /restaurant/order/confirm (ยืนยันออเดอร์)
func (h *Handler) ConfirmOrder(c *gin.Context) {
	var req ConfirmOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ownerUsername := c.GetString("username")
	if ownerUsername == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	if err := h.service.ConfirmOrder(req.OrderID, ownerUsername); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or not authorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order confirmed"})
}
